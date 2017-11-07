package sql

import (
	"gopkg.in/mgo.v2"
)

// Insert mongo full-block with tx details
func InsertOne(block interface{}) error {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("Eth").C("Block")
	if err := c.Insert(block); err != nil {
		return err
	}
	return nil

}
