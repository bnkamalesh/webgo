package webgo

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// This file contains a wrapper to use over Mgo driver.
// Only basic operations are implemented as of now.
// Need to add Delete function

type DBConfig struct {
	Name          string `json:"name"`
	Host          string `json:"host"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	AuthSource    string `json:"authSource"`
	MgoDialString string `json:"mgoDialString"`
}

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
	return ds.Session.Copy()
}

// ===

// Get appropriate MongoDB collection
func (ds *DataStore) GetSessionCollection(dbName, collection string) (*mgo.Session, *mgo.Collection) {
	s := ds.getSession()
	c := s.DB(dbName).C(collection)

	return s, c
}

// ===

// Do a MongoDB Get
func (ds *DataStore) Get(dbName, collection string, conditions interface{}, resultStruct interface{}) ([]bson.M, error) {

	s, c := ds.GetSessionCollection(dbName, collection)
	defer s.Close()

	if resultStruct != nil {
		err := c.Find(conditions).All(resultStruct)
		if err != nil {
			if err == mgo.ErrNotFound {
				return nil, nil
			}
			return nil, err
		}
		return nil, nil
	}

	var data []bson.M
	err := c.Find(conditions).All(&data)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// ===

// Do a MongoDB GetAll
func (ds *DataStore) GetAll(dbName, collection string, resultStruct interface{}) ([]bson.M, error) {

	s, c := ds.GetSessionCollection(dbName, collection)
	defer s.Close()

	if resultStruct != nil {
		err := c.Find(nil).All(resultStruct)
		if err != nil {
			if err == mgo.ErrNotFound {
				return nil, nil
			}
			return nil, err
		}
		return nil, nil
	}

	var data []bson.M
	err := c.Find(nil).All(&data)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// ===

// Do a MongoDB GetOne
func (ds *DataStore) GetOne(dbName, collection string, conditions interface{}, resultStruct interface{}) (bson.M, error) {

	s, c := ds.GetSessionCollection(dbName, collection)
	defer s.Close()

	if resultStruct != nil {
		err := c.Find(conditions).One(resultStruct)
		if err != nil {
			if err == mgo.ErrNotFound {
				return nil, nil
			}
			return nil, err
		}
		return nil, nil
	}

	var data bson.M
	err := c.Find(conditions).One(&data)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

// ===

// Do a MongoDB Save
func (ds *DataStore) Save(dbName, collection string, data interface{}) error {
	s, c := ds.GetSessionCollection(dbName, collection)
	defer s.Close()

	err := c.Insert(data)
	if err != nil {
		return err
	}
	return nil
}

// ===

// Do a MongoDB Update - multiple records
func (ds *DataStore) Update(dbName, collection string, condition, updateData interface{}) error {
	s, c := ds.GetSessionCollection(dbName, collection)
	defer s.Close()

	err := c.Update(condition, updateData)
	if err != nil {
		return err
	}

	return nil
}

// ===

// Do a MongoDB update - single record, by MongoID
func (ds *DataStore) UpdateId(dbName, collection string, _id, data interface{}) error {
	s, c := ds.GetSessionCollection(dbName, collection)
	defer s.Close()

	err := c.UpdateId(_id, data)
	if err != nil {
		return err
	}

	return nil
}

// ===

// Create a new data store
func newDataStore(user, pass, host, port, name, authSource, mgoDialString string) (*DataStore, error) {
	dialString := "mongodb://"

	if len(mgoDialString) > 0 {
		dialString = mgoDialString
	} else {
		if len(user) > 0 && len(pass) > 0 {
			dialString += (user + ":" + pass)
		}

		if len(host) > 0 {
			dialString += ("@" + host)
			if len(port) > 0 {
				dialString += (":" + port)
			}
		}

		if len(name) > 0 {
			dialString += ("/" + name)
		}

		if len(authSource) > 0 {
			dialString += "?authSource=" + authSource
		}
	}

	session, err := mgo.Dial(dialString)
	if err != nil {
		return nil, err
	}
	session.SetSafe(&mgo.Safe{})

	return &DataStore{DbName: name, Session: session}, nil
}

// ===

// Initializing Mongo DB
func InitDB(dbc DBConfig) *DataStore {
	dStore, err := newDataStore(
		dbc.Username,
		dbc.Password,
		dbc.Host,
		dbc.Port,
		dbc.Name,
		dbc.AuthSource,
		dbc.MgoDialString,
	)
	if err != nil {
		Err.Fatal("db.go", "InitDB()", err)
		return nil
	}
	return dStore
}

// ===
