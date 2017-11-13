package main

import (
	"time"

	mgo "gopkg.in/mgo.v2"

	"log"

	"github.com/xiaoping378/blockchain-on-sql/parser"
	"github.com/xiaoping378/blockchain-on-sql/sql"
)

func main() {

	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	s := sql.SQL{}
	s.SetCollection(session.DB("Eth").C("Block"))

	sync := make(chan int, 1)
	go s.Sync(s.GetSyncedBlockCount(), parser.GetLatestBlockNumber(), sync)

	// 周期同步
	for {
		select {
		case <-sync:
			log.Println("syncing task is completed.")
			time.Sleep(7 * time.Second) // TODO: using event listen
			s.Sync(s.GetSyncedBlockCount(), parser.GetLatestBlockNumber(), sync)
		}
	}

}
