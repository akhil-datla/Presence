package dbmanager

import (
	"github.com/asdine/storm"
)

//DBManager defines the type of database
type DBManager struct {
	db *storm.DB
}

var dbm *DBManager = &DBManager{}


// Start opens the database
func Start(name string) error {
	var err error
	dbm.db, err = storm.Open(name)
	return err
}

//CreateBucket creates a table for a given data type
func CreateBucket(data interface{}) error {
	err := dbm.db.Init(data)
	return err
}

//Save saves the entry into the database
func Save(entry interface{}) error {
	err := dbm.db.Save(entry)
	return err
}

//Delete an entry from the database
func Delete(entry interface{}) error {
	err := dbm.db.DeleteStruct(entry)
	return err
}

//Query queries the database for a specific entry
func Query(group string, key string, store interface{}) error {
	err := dbm.db.One(group, key, store)
	return err

}


//GroupQuery queries the database for a group of entries
func GroupQuery(group string, key string, store interface{}) error {
	err := dbm.db.Find(group, key, store)
	return err
}

// Close closes the database
func Close(name string) error {
	err := dbm.db.Close()
	return err
}

//Update updates an entry
func Update(data interface{}) error {
	err := dbm.db.Update(data)
	return err
}
