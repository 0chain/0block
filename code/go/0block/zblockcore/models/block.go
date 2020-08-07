package models

import (
	"0block/core/datastore"
	"context"
	"time"

	. "0block/core/logging"

	"github.com/0chain/gosdk/core/block"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2/bson"
)

// type VerificationTicket struct {
// 	VerifierID string `bson:"verifier_id;index" json:"verifier_id"`
// 	Signature  string `bson:"signature" json:"signature"`
// }

type Block struct {
	Hash                  string    `bson:"hash,unique_index,not null"`
	Version               string    `bson:"version,not null"`
	CreationDate          int64     `bson:"creation_date,not null"`
	Round                 int64     `bson:"round,unique_index,not null"`
	MinerID               string    `bson:"miner_id,index,not null"`
	RoundRandomSeed       int64     `bson:"round_random_seed,not null"`
	MerkleTreeRoot        string    `bson:"merkle_tree_root"`
	StateHash             string    `bson:"state_hash,not null"`
	ReceiptMerkleTreeRoot string    `bson:"receipt_merkle_tree_root"`
	NumTxns               int64     `bson:"num_txns"`
	MagicBlockHash        string    `bson:"magic_block_hash"`
	PrevHash              string    `bson:"prev_hash"`
	Signature             string    `bson:"signature"`
	ChainID               string    `bson:"chain_id"`
	ChainWeight           float64   `bson:"chain_weight"`
	RunningTxnCount       int64     `bson:"running_txn_count"`
	RoundTimeoutCount     int       `bson:"round_timeout_count"`
	CreatedAt             time.Time `bson:"created_at,not null"`
	UpdatedAt             time.Time `bson:"updated_at,not null"`
	// PrevVerificationTickets []VerificationTicket `bson:"prev_verification_tickets" json:"prev_verification_tickets"`
	// VerificationTickets     []VerificationTicket `bson:"verification_tickets" json:"verification_tickets"`
}

func (Block) GetCollection() *mongo.Collection {
	return datastore.GetStore().GetDB().Collection("blocks")
}

func InsertBlock(ctx context.Context, block *block.Block) error {
	blocksCollection := (&Block{}).GetCollection()
	dbBlock := transferBlockData(block)
	_, err := blocksCollection.InsertOne(ctx, dbBlock)
	if err != nil {
		Logger.Error("Failed to insert block data", zap.Error(err))
		return err
	}

	err = InsertTransactions(ctx, dbBlock.Hash, block.Txns)
	if err != nil {
		Logger.Error("Failed to insert transaction data", zap.Error(err))
		return err
	}
	return nil
}

func GetLatestBlockInDB(ctx context.Context) (*Block, error) {
	var block Block
	opts := options.FindOne()
	opts.SetSort(bson.M{"round": -1})
	if err := (&Block{}).GetCollection().FindOne(ctx, bson.M{}, opts).Decode(&block); err != nil {
		Logger.Error("Failed to get latest block from DB", zap.Error(err))
		return nil, err
	}
	return &block, nil
}

func CheckBlockPresentInDB(ctx context.Context, round int64) bool {
	var block Block
	filter := bson.M{"round": round}
	if err := (&Block{}).GetCollection().FindOne(ctx, filter).Decode(&block); err != nil {
		return false
	}
	return true
}

func transferBlockData(block *block.Block) *Block {
	dbBlock := new(Block)
	dbBlock.Hash = string(block.Hash)
	dbBlock.Version = block.Version
	dbBlock.CreationDate = int64(block.CreationDate)
	dbBlock.Round = block.Round
	dbBlock.MinerID = string(block.MinerID)
	dbBlock.RoundRandomSeed = block.RoundRandomSeed
	dbBlock.MerkleTreeRoot = block.Header.MerkleTreeRoot
	dbBlock.StateHash = block.Header.StateHash
	dbBlock.ReceiptMerkleTreeRoot = block.Header.ReceiptMerkleTreeRoot
	dbBlock.NumTxns = block.Header.NumTxns
	dbBlock.MagicBlockHash = block.MagicBlockHash
	dbBlock.PrevHash = block.PrevHash
	dbBlock.Signature = block.Signature
	dbBlock.ChainID = string(block.ChainID)
	dbBlock.ChainWeight = block.ChainWeight
	dbBlock.RunningTxnCount = block.RunningTxnCount
	dbBlock.RoundTimeoutCount = block.RoundTimeoutCount
	dbBlock.CreatedAt = time.Now()
	dbBlock.UpdatedAt = time.Now()
	return dbBlock
}
