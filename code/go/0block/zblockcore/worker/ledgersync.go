package worker

import (
	"context"

	"github.com/0chain/gosdk/core/block"
)

func LedgerSync(ctx context.Context, roundToProcess int64) {
	blockChan := make(chan *block.Block)
	go fetchBlock(ctx, blockChan, roundToProcess, forward)
	for {
		select {
		case <-ctx.Done():
			return
		case blockToProcess := <-blockChan:
			insertBlock(ctx, blockToProcess)
		}
	}
}
