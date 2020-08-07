package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"0block/core/common"
	"0block/core/config"
	"0block/core/datastore"
	"0block/core/encryption"
	"0block/core/logging"
	. "0block/core/logging"

	"0block/zblockcore/worker"
	"0block/zblockcore/zcn"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func initializeConfig() {
	config.Configuration.ChainID = viper.GetString("server_chain.id")
	config.Configuration.SignatureScheme = viper.GetString("server_chain.signature_scheme")
	config.Configuration.Port = viper.GetInt("port")

	config.Configuration.MongoURL = viper.GetString("mongo.url")
	config.Configuration.DBName = viper.GetString("mongo.db_name")
	config.Configuration.MongoPoolSize = viper.GetInt64("mongo.pool_size")

	config.Configuration.RoundFetchDelayInMilliSeconds = viper.GetInt64("worker.round_fetch_delay")
	config.Configuration.RoundFetchRetries = viper.GetInt64("worker.round_fetch_retires")
	config.Configuration.DNSURL = viper.GetString("dns_url")
}

func initHandlers(r *mux.Router) {
	r.HandleFunc("/", HomePageHandler)
	r.HandleFunc("/logs", logging.LogWriter)
	r.HandleFunc("/mem_logs", logging.MemLogWriter)
}

var startTime time.Time

func main() {
	deploymentMode := flag.Int("deployment_mode", 2, "deployment_mode")
	keysFile := flag.String("keys_file", "", "keys_file")
	flag.Parse()

	config.Configuration.DeploymentMode = byte(*deploymentMode)
	config.SetupDefaultConfig()
	config.SetupConfig()

	if config.Development() {
		logging.InitLogging("development")
	} else {
		logging.InitLogging("production")
	}
	initializeConfig()

	reader, err := os.Open(*keysFile)
	if err != nil {
		panic(err)
	}

	publicKey, privateKey := encryption.ReadKeys(reader)
	config.Configuration.SetWallet(publicKey, privateKey)

	common.SetupRootContext(context.Background())

	checkForDBConnection(context.Background())

	zcn.InitZCN()

	address := fmt.Sprintf(":%v", config.Configuration.Port)

	var server *http.Server
	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET"})
	rHandler := handlers.CORS(originsOk, headersOk, methodsOk)(r)
	if config.Development() {
		server = &http.Server{
			Addr:           address,
			ReadTimeout:    30 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        rHandler,
		}
	} else {
		server = &http.Server{
			Addr:           address,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        rHandler,
		}
	}
	common.HandleShutdown(server)

	initHandlers(r)
	go worker.SetupWorkers(context.Background())

	startTime = time.Now().UTC()
	Logger.Info("Ready to listen to the requests on ", zap.Any("port", config.Configuration.Port))
	log.Fatal(server.ListenAndServe())
}

// HomePageHandler for 0block
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<div>Running since %v ...\n", startTime)
	fmt.Fprintf(w, "<div>Working on the chain: %v</div>\n", config.Configuration.ChainID)
	fmt.Fprintf(w, "<div>I am 0block with <ul><li>DNS:%v</li></ul></div>\n", config.Configuration.DNSURL)
	fmt.Fprintf(w, "<div>To check network details <a href='%v'>Click here</a>", config.Configuration.DNSURL+"/network")
}

func checkForDBConnection(ctx context.Context) {
	retries := 0
	var err error
	for retries < 600 {
		Logger.Info("Trying to connect to mongoDB ...")
		err = datastore.GetStore().Open(ctx)
		if err != nil {
			time.Sleep(1 * time.Second)
			retries++
			continue
		}
		Logger.Info("DB Connection done.")
		break
	}

	if err != nil {
		Logger.Error("Error in opening the database. Shutting the server down")
		panic(err)
	}
}
