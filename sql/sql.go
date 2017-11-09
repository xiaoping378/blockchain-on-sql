package sql

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var server = "127.0.0.1:27017"

type SQL struct {
	session mgo.Session
}

// func (s *SQL) Connect(addr string) error {
// 	session, err := mgo.Dial(addr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()
// 	session.SetMode(mgo.Monotonic, true)
// }

// InsertOne mongo full-block with tx details
func (s *SQL) InsertOne(block interface{}) error {
	c := s.session.DB("Eth").C("Block")
	if err := c.Insert(block); err != nil {
		return err
	}
	return nil

}

func (s *SQL) GetSyncedBlockCount() uint64 {

	result := struct {
		number uint64
	}{}

	c := s.session.DB("Eth").C("Block")
	c.Find(bson.M{}).Sort("-number").Limit(1).One(&result)
	return result.number
}
