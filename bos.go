package main

import (
	"fmt"
	"math/big"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"log"

	"github.com/xiaoping378/blockchain-on-sql/common"
	"github.com/xiaoping378/blockchain-on-sql/sql"

	"github.com/xiaoping378/blockchain-on-sql/parser"
)

func hexToDecimal(s string) bson.Decimal128 {
	bigInt := new(big.Int)
	bigInt.SetString(s, 0)
	bigIntByte, _ := bigInt.MarshalText()
	decimal, _ := bson.ParseDecimal128(string(bigIntByte))
	return decimal
	// return bigInt.Int64()
}

func converToMTransaction(txs []common.Transaction) []common.MTransaction {
	var mts []common.MTransaction

	for _, t := range txs {
		var mt common.MTransaction

		mt.BlockHash = t.BlockHash
		mt.BlockNumber = hexToDecimal(t.BlockNumber)
		mt.From = t.From
		mt.Gas = hexToDecimal(t.Gas)
		mt.GasPrice = hexToDecimal(t.GasPrice)
		mt.Hash = t.Hash
		mt.Input = t.Input
		mt.Nonce = t.Nonce
		mt.R = t.R
		mt.S = t.S
		mt.To = t.To
		mt.TransactionIndex = hexToDecimal(t.TransactionIndex)
		mt.V = t.V
		mt.Value = hexToDecimal(t.Value)

		mts = append(mts, mt)
	}

	return mts
}

func converToMBlock(r *common.Block) *common.MBlock {
	var m = common.MBlock{}
	// var err error
	m.Bloom = r.Bloom
	m.Coinbase = r.Coinbase
	m.Difficulty = hexToDecimal(r.Difficulty)
	m.Extra = r.Extra
	m.GasLimit = hexToDecimal(r.GasLimit)
	m.GasUsed = hexToDecimal(r.GasUsed)
	m.Hash = r.Hash
	m.MixDigest = r.MixDigest
	m.Nonce = r.Nonce
	m.Number = hexToDecimal(r.Number)
	m.ParentHash = r.ParentHash
	m.ReceiptHash = r.ReceiptHash
	m.Root = r.Root
	m.Size = hexToDecimal(r.Size)
	m.Time = hexToDecimal(r.Time)
	m.TotalDifficulty = hexToDecimal(r.TotalDifficulty)
	m.TxHash = r.TxHash
	m.TXs = converToMTransaction(r.TXs)
	m.UncleHash = r.UncleHash
	m.Uncles = r.Uncles

	return &m
}

func intern(n int64) []struct{} {
	return make([]struct{}, n)
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func main() {

	fullTx := true

	for i := range intern(20000) {
		number := fmt.Sprintf("0x%s", strconv.FormatInt(int64(i), 16))
		resp, err := parser.Call("eth_getBlockByNumber", []interface{}{number, fullTx})
		if err != nil {
			log.Fatal(err)
		}

		// var respData = resp.Result
		// fmt.Println(respData)
		block := new(common.Block)

		if err := parser.MapToObject(resp.Result, block); err != nil {
			log.Fatalln(err)
		}

		mBlock := converToMBlock(block)

		err = sql.InsertOne(mBlock)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
