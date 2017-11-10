package parser

import (
	"strconv"

	"github.com/xiaoping378/blockchain-on-sql/common"
)

func GetLatestBlockNumber() uint64 {
	block := common.Block{}
	latest, _ := Call("eth_getBlockByNumber", []interface{}{"latest", false})
	MapToObject(latest.Result, &block)
	latestBlock, _ := strconv.ParseUint(block.Number[2:], 16, 64)
	return latestBlock
}
