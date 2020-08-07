package config

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/0chain/gosdk/core/zcncrypto"
	"github.com/spf13/viper"
	"golang.org/x/crypto/sha3"
)

//SetupDefaultConfig - setup the default config options that can be overridden via the config file
func SetupDefaultConfig() {
	viper.SetDefault("logging.level", "info")
}

/*SetupConfig - setup the configuration system */
func SetupConfig() {
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetConfigName("0block")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

const (
	DeploymentDevelopment = 0
	DeploymentTestNet     = 1
	DeploymentMainNet     = 2
	HTTPProtocol          = "http://"
	HTTPSProtocol         = "https://"
)

/*Config - all the config options passed from the command line*/
type Config struct {
	Port            int
	ChainID         string
	DeploymentMode  byte
	SignatureScheme string

	MongoURL      string
	DBName        string
	MongoPoolSize int64

	RoundFetchDelayInMilliSeconds int64
	RoundFetchRetries             int64

	Wallet *zcncrypto.Wallet
	DNSURL string
}

/*Configuration of the system */
var Configuration Config

/*TestNet is the program running in TestNet mode? */
func TestNet() bool {
	return Configuration.DeploymentMode == DeploymentTestNet
}

/*Development - is the programming running in development mode? */
func Development() bool {
	return Configuration.DeploymentMode == DeploymentDevelopment
}
func (c *Config) SetWallet(publicKey string, privateKey string) {
	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		panic(err)
	}
	c.Wallet = &zcncrypto.Wallet{}
	c.Wallet.ClientID = Hash(publicKeyBytes)
	c.Wallet.ClientKey = publicKey
	c.Wallet.Keys = make([]zcncrypto.KeyPair, 1)
	c.Wallet.Keys[0].PublicKey = publicKey
	c.Wallet.Keys[0].PrivateKey = privateKey
	c.Wallet.Version = zcncrypto.CryptoVersion
}

const HASH_LENGTH = 32

type HashBytes [HASH_LENGTH]byte

/*Hash - hash the given data and return the hash as hex string */
func Hash(data interface{}) string {
	return hex.EncodeToString(RawHash(data))
}

/*RawHash - Logic to hash the text and return the hash bytes */
func RawHash(data interface{}) []byte {
	var databuf []byte
	switch dataImpl := data.(type) {
	case []byte:
		databuf = dataImpl
	case HashBytes:
		databuf = dataImpl[:]
	case string:
		databuf = []byte(dataImpl)
	default:
		panic("unknown type")
	}
	hash := sha3.New256()
	hash.Write(databuf)
	var buf []byte
	return hash.Sum(buf)
}
