package sql

import (
	"fmt"
	"log"
	"strconv"

	"github.com/xiaoping378/blockchain-on-sql/common"
	"github.com/xiaoping378/blockchain-on-sql/parser"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var server = "127.0.0.1:27017"

type SQL struct {
	C *mgo.Collection
}

// func (s *SQL) Connect(addr string) error {
// 	session, err := mgo.Dial(addr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// }

func (s *SQL) SetCollection(c *mgo.Collection) *SQL {
	s.C = c
	return s
}

// InsertOne mongo full-block with tx details
func (s *SQL) InsertOne(block interface{}) error {
	if err := s.C.Insert(block); err != nil {
		return err
	}
	return nil

}

func (s *SQL) GetSyncedBlockCount() uint64 {

	result := common.MBlock{}
	s.C.Find(bson.M{}).Sort("-number").Limit(1).One(&result)
	syncedNumber, _ := strconv.ParseUint(result.Number.String(), 10, 64)
	return syncedNumber
}

func (s *SQL) Sync(syncedNumber, latestBlock uint64, c chan int) {
	block := common.Block{}
	if syncedNumber > 0 {
		// 从下一个块开始同步
		syncedNumber++
	}
	for i := syncedNumber; i <= latestBlock; i++ {

		number := fmt.Sprintf("0x%s", strconv.FormatUint(uint64(i), 16))
		resp, err := parser.Call("eth_getBlockByNumber", []interface{}{number, true})
		if err != nil {
			log.Fatal(err)
		}

		if err := parser.MapToObject(resp.Result, &block); err != nil {
			log.Fatalln(err)
		}

		mBlock := block.ToMBlock()

		if err := s.InsertOne(mBlock); err != nil {
			log.Fatal(err)
		}
	}

	c <- 1
}
