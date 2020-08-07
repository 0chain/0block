package models

import (
	"0block/core/datastore"
	. "0block/core/logging"
	"context"
	"encoding/json"
	"time"

	"github.com/0chain/gosdk/core/transaction"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Transaction struct {
	Hash                string    `bson:"hash,unique,not null"`
	BlockHash           string    `bson:"block_hash,not null"`
	Version             string    `bson:"version,not null"`
	ClientID            string    `bson:"client_id"`
	ToClientID          string    `bson:"to_client_id"`
	ChainID             string    `bson:"chain_id,not null"`
	TransactionData     string    `bson:"transaction_data"`
	Value               int64     `bson:"transaction_value"`
	Signature           string    `bson:"signature"`
	CreationDate        int64     `bson:"creation_date"`
	Fee                 int64     `bson:"transaction_fee"`
	TransactionType     int       `bson:"transaction_type"`
	TransactionOutput   string    `bson:"transaction_output"`
	OutputHash          string    `bson:"txn_output_hash"`
	Status              string    `bson:"transaction_status"`
	ConfirmationFetched bool      `bson:"confirmation_fetched"`
	LookUpHash          string    `bson:"lookup_hash"`
	Name                string    `bson:"name"`
	ContentHash         string    `bson:"content_hash"`
	AllocationID        string    `bson:"allocation_id"`
	CreatedAt           time.Time `bson:"created_at,not null"`
	UpdatedAt           time.Time `bson:"updated_at,not null"`
}

func (Transaction) GetCollection() *mongo.Collection {
	return datastore.GetStore().GetDB().Collection("transactions")
}

func InsertTransactions(ctx context.Context, blockHash string, transactions []*transaction.Transaction) error {
	var ITransactions []interface{}
	transactionCollection := (&Transaction{}).GetCollection()
	for _, transaction := range transactions {
		dbTransaction := transferTransactionData(blockHash, transaction)

		if isJSONString(dbTransaction.TransactionData) {
			var TxnData interface{}
			json.Unmarshal([]byte(dbTransaction.TransactionData), &TxnData)
			if txnData, ok := TxnData.(map[string]interface{}); ok {

				if metaData, ok := txnData["MetaData"].(map[string]interface{}); ok {

					if lookUphash, ok := metaData["LookupHash"].(string); ok {
						dbTransaction.LookUpHash = lookUphash
					}

					if name, ok := metaData["Name"].(string); ok {
						dbTransaction.Name = name
					}

					if contentHash, ok := metaData["Hash"].(string); ok {
						dbTransaction.ContentHash = contentHash
					}
				}
			}
		}

		if isJSONString(dbTransaction.TransactionOutput) {
			var TxnOutput interface{}
			json.Unmarshal([]byte(dbTransaction.TransactionOutput), &TxnOutput)

			if parsedOutput, ok := TxnOutput.(map[string]interface{}); ok {
				if allocationID, ok := parsedOutput["allocation_id"].(string); ok {
					dbTransaction.AllocationID = allocationID
				}
			}
		}

		ITransactions = append(ITransactions, dbTransaction)
	}
	_, err := transactionCollection.InsertMany(ctx, ITransactions)
	if err != nil {
		Logger.Error("transaction_insert_many_failed", zap.Error(err), zap.Any("block_hash", blockHash))
		return err
	}
	Logger.Info("Transactions stored successfully", zap.Any("block_hash", blockHash))
	return nil
}

func isJSONString(data string) bool {
	var v interface{}
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		return true
	}
	return false
}

func transferTransactionData(blockHash string, transaction *transaction.Transaction) *Transaction {
	dbTransaction := new(Transaction)
	dbTransaction.Hash = transaction.Hash
	dbTransaction.BlockHash = blockHash
	dbTransaction.Version = transaction.Version
	dbTransaction.ClientID = transaction.ClientID
	dbTransaction.ToClientID = transaction.ToClientID
	dbTransaction.ChainID = transaction.ChainID
	dbTransaction.TransactionData = transaction.TransactionData
	dbTransaction.Value = transaction.Value
	dbTransaction.Signature = transaction.Signature
	dbTransaction.CreationDate = transaction.CreationDate
	dbTransaction.Fee = transaction.TransactionFee
	dbTransaction.TransactionType = transaction.TransactionType
	dbTransaction.TransactionOutput = transaction.TransactionOutput
	dbTransaction.OutputHash = transaction.OutputHash
	dbTransaction.ConfirmationFetched = false
	dbTransaction.CreatedAt = time.Now()
	dbTransaction.UpdatedAt = time.Now()
	return dbTransaction
}
