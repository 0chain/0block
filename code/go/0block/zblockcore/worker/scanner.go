package worker

import (
	"context"

	"github.com/0chain/gosdk/core/block"
)

func Scanner(ctx context.Context, roundToProcess int64) {
	blockChan := make(chan *block.Block)
	go fetchBlock(ctx, blockChan, roundToProcess, backward)
	for {
		select {
		case <-ctx.Done():
			return
		case blockToProcess := <-blockChan:
			insertBlock(ctx, blockToProcess)
			if blockToProcess.Round == 1 {
				return
			}
		}
	}
}
