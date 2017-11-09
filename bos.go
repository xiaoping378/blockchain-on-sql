package main

import (
	"fmt"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"log"

	"github.com/xiaoping378/blockchain-on-sql/common"

	"github.com/xiaoping378/blockchain-on-sql/parser"
)

func main() {

	fullTx := true

	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("Eth").C("Block")

	result := common.MBlock{}
	c.Find(bson.M{}).Sort("-number").Limit(1).One(&result)
	syncedNumber, _ := strconv.ParseUint(result.Number.String(), 10, 64)

	block := common.Block{}
	latest, _ := parser.Call("eth_getBlockByNumber", []interface{}{"latest", fullTx})
	parser.MapToObject(latest.Result, &block)
	latestBlock, _ := strconv.ParseUint(block.Number[2:], 16, 64)

	for i := syncedNumber; i < latestBlock; i++ {
		number := fmt.Sprintf("0x%s", strconv.FormatInt(int64(i), 16))
		resp, err := parser.Call("eth_getBlockByNumber", []interface{}{number, fullTx})
		if err != nil {
			log.Fatal(err)
		}

		if err := parser.MapToObject(resp.Result, &block); err != nil {
			log.Fatalln(err)
		}

		mBlock := block.ToMBlock()

		if err := c.Insert(mBlock); err != nil {
			log.Fatal(err)
		}
		// err = sql.InsertOne(mBlock)
		if err != nil {
			log.Fatalln(err)
		}
	}

}
