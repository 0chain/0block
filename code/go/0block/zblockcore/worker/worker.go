package worker

import (
	"0block/core/config"
	"0block/zblockcore/models"
	"context"
	"runtime"
	"time"

	. "0block/core/logging"

	"github.com/0chain/gosdk/core/block"
	"github.com/0chain/gosdk/zcncore"
	"go.uber.org/zap"
)

type FetchMode string

const (
	forward  FetchMode = "forward"
	backward FetchMode = "backward"
)

func SetupWorkers(ctx context.Context) {
	latestBlock, err := zcncore.GetLatestFinalized(ctx, 1)
	if err != nil {
		Logger.Error("Failed to get latest finalized block from blockchain", zap.Error(err))
		panic("GetLatestFinalized failed")
	}
	roundToProcess := latestBlock.Round
	go LedgerSync(ctx, roundToProcess)
	go Scanner(ctx, roundToProcess-1)
	go MemUsageLogger(ctx)
}

func fetchBlock(ctx context.Context, blockChan chan *block.Block, round int64, mode FetchMode) {
	for round > 0 {
		select {
		case <-ctx.Done():
			return
		default:
			if models.CheckBlockPresentInDB(ctx, round) {
				Logger.Info("Block already present in DB", zap.Any("round", round), zap.Any("mode", mode))
			} else {
				retries := config.Configuration.RoundFetchRetries
				Logger.Info("Fetching block by round from blockchain", zap.Any("round", round))
				for retries > 0 {
					block, err := zcncore.GetBlockByRound(ctx, zcncore.GetMinShardersVerify(), round)
					if err != nil {
						retries--
						Logger.Info("Unable to get block by round from blockchain", zap.Error(err), zap.Any("round", round),
							zap.Any("Attempts left", retries), zap.Any("mode", mode))
						time.Sleep(time.Duration(config.Configuration.RoundFetchDelayInMilliSeconds) * time.Millisecond)
						if retries == 0 && mode == forward {
							panic("A block is missed, Network probably stuck, Killing block worker...")
						}
						continue
					}

					Logger.Info("Got block by round from blockchain", zap.Any("round", round), zap.Any("mode", mode))
					blockChan <- block
					break
				}
			}

			if mode == forward {
				round++
			} else {
				round--
			}
		}
	}
}

func insertBlock(ctx context.Context, blockToProcess *block.Block) {
	go func(ctx context.Context, b *block.Block) {
		retries := config.Configuration.RoundFetchRetries
		for retries > 0 {
			err := models.InsertBlock(ctx, b)
			if err != nil {
				Logger.Error("Failed to insert block to the DB", zap.Any("round", b.Round), zap.Error(err))
				retries--
				continue
			}
			Logger.Info("Insert block successfully to the DB", zap.Any("round", b.Round))
			return
		}
	}(ctx, blockToProcess)
}

func MemUsageLogger(ctx context.Context) {
	const MB = 1024 * 1024
	ticker := time.NewTicker(2 * time.Minute)
	for true {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			MemUsage.Info("runtime", zap.Int("goroutines", runtime.NumGoroutine()), zap.Uint64("heap_objects", mem.HeapObjects), zap.Uint32("gc", mem.NumGC), zap.Uint64("gc_pause", mem.PauseNs[(mem.NumGC+255)%256]))
			MemUsage.Info("runtime", zap.Uint64("total_alloc", mem.TotalAlloc/MB), zap.Uint64("sys", mem.Sys/MB), zap.Uint64("heap_sys", mem.HeapSys/MB), zap.Uint64("heap_alloc", mem.HeapAlloc/MB))
		}
	}
}
