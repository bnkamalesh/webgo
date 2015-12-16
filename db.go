package webgo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// This file contains a wrapper to use over Mgo driver.
// Only basic operations are implemented as of now.
// Need to add Delete function

/*
 DB Session creation is done inside a struct
 Developed based on the answer in
 http://stackoverflow.com/questions/26574594/best-practice-to-maintain-a-mgo-session
*/
type DataStore struct {
	DbName  string
	Session *mgo.Session
}

// ===

// Clone the master session and return
func (ds *DataStore) getSession() *mgo.Session {
	return ds.Session.Clone()
}

// ===

// Get appropriate MongoDB collection
func (ds *DataStore) getSessionCollection(collection string) (*mgo.Session, *mgo.Collection) {
	s := ds.getSession()
	c := s.DB(ds.DbName).C(collection)

	return s, c
}

// ===

// Do a MongoDB Get
func (ds *DataStore) Get(collection string, conditions interface{}) ([]bson.M, error) {
	var data []bson.M

	s, c := ds.getSessionCollection(collection)
	defer s.Close()

	err := c.Find(conditions).All(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ===

// Do a MongoDB GetAll
func (ds *DataStore) GetAll(collection string) ([]bson.M, error) {
	var data []bson.M

	s, c := ds.getSessionCollection(collection)
	defer s.Close()

	err := c.Find(nil).All(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ===

// Do a MongoDB GetOne
func (ds *DataStore) GetOne(collection string, conditions interface{}) (bson.M, error) {
	var data bson.M

	s, c := ds.getSessionCollection(collection)
	defer s.Close()

	err := c.Find(conditions).One(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ===

// Do a MongoDB Save
func (ds *DataStore) Save(collection string, data interface{}) error {
	s, c := ds.getSessionCollection(collection)
	defer s.Close()

	err := c.Insert(data)
	if err != nil {
		Err.Log("db.go", "Save()", err)
		return err
	}
	return nil
}

// ===

// Do a MongoDB Update - multiple records
func (ds *DataStore) Update(collection string, condition, updateData interface{}) error {
	s, c := ds.getSessionCollection(collection)
	defer s.Close()

	err := c.Update(condition, updateData)
	if err != nil {
		return err
	}

	return nil
}

// ===

// Do a MongoDB update - single record, by MongoID
func (ds *DataStore) UpdateId(collection string, _id, data interface{}) error {
	s, c := ds.getSessionCollection(collection)
	defer s.Close()

	err := c.UpdateId(_id, data)
	if err != nil {
		return err
	}

	return nil
}

// ===

// Create a new data store
func newDataStore(user, pass, host, port, name string) (*DataStore, error) {
	session, err := mgo.Dial("mongodb://" + user + ":" + pass + "@" + host + ":" + port + "/" + name)
	if err != nil {
		return nil, err
	}
	session.SetSafe(&mgo.Safe{})

	return &DataStore{DbName: name, Session: session}, nil
}

// ===

// Initializing Mongo DB
func InitDB(user, pass, host, port, name string) *DataStore {
	dStore, err := newDataStore(user, pass, host, port, name)
	if err != nil {
		Err.Fatal("db.go", "InitDB()", err)
		return nil
	}
	return dStore
}

// ===
