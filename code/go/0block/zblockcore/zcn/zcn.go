package zcn

import (
	"encoding/json"
	"sync"

	"0block/core/config"
	. "0block/core/logging"

	"github.com/0chain/gosdk/zcncore"
	"go.uber.org/zap"
)

const (
	ZCNStatusSuccess int = 0
	ZCNStatusError   int = 1
)

type ZCNStatus struct {
	walletString string
	wg           *sync.WaitGroup
	success      bool
	errMsg       string
	balance      int64
	wallets      []string
	clientID     string
}

func (zcn *ZCNStatus) OnWalletCreateComplete(status int, wallet string, err string) {
	defer zcn.wg.Done()
	if status == ZCNStatusError {
		zcn.success = false
		zcn.errMsg = err
		zcn.walletString = ""
		return
	}
	zcn.success = true
	zcn.errMsg = ""
	zcn.walletString = wallet
	return
}

// func (zcn *ZCNStatus) OnAuthComplete(t *zcncore.Transaction, status int) {
// 	fmt.Println("Authorization complete on zauth.", status)
// }

// func (zcn *ZCNStatus) OnTransactionComplete(t *zcncore.Transaction, status int) {
// 	defer zcn.wg.Done()
// 	if status == zcncore.StatusSuccess {
// 		zcn.success = true
// 	} else {
// 		zcn.errMsg = t.GetTransactionError()
// 	}
// }

// func (zcn *ZCNStatus) OnVerifyComplete(t *zcncore.Transaction, status int) {
// 	defer zcn.wg.Done()
// 	if status == zcncore.StatusSuccess {
// 		zcn.success = true
// 	} else {
// 		zcn.errMsg = t.GetVerifyError()
// 	}
// }

func InitZCN() {
	// No logs from SDK
	zcncore.SetLogLevel(0)
	err := zcncore.InitZCNSDK(config.Configuration.DNSURL, config.Configuration.SignatureScheme)
	if err != nil {
		Logger.Error("Failed io initialize SDK", zap.Error(err))
		panic("Error: Unable to init SDK")
	}

	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	wg.Add(1)
	err = zcncore.RegisterToMiners(config.Configuration.Wallet, statusBar)
	if err != nil {
		Logger.Error("Wallet registration to miners failed with err", zap.Error(err))
		panic("Error trying to register wallet to miners")
	}
	wg.Wait()
	if statusBar.success {
		Logger.Info("Wallet registered")
		walletBytes, _ := json.Marshal(config.Configuration.Wallet)
		err := zcncore.SetWalletInfo(string(walletBytes), false)
		if err != nil {
			panic("Error in setting wallet info")
		}
	} else {
		Logger.Error("Wallet registration failed. " + statusBar.errMsg)
		panic("Wallet registration failed")
	}
}
