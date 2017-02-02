package config

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

// Conf contains application wide configurations
type Conf struct {
	MongoDB     string
	MongoServer string
	Database    *mgo.Database
	//DBPassword  string
	//DBUser      string
}

var (
	config Conf
)

const (
	USERCOLLECTION = "users"
)

// Init initialize the configurations
func Init() {

	// get mongo server address from the system variables
	MONGOSERVER := os.Getenv("MONGO_URL")

	// get name of mongo DB from system variables
	MONGODB := os.Getenv("MONGODB")

	// set mongo server to default if it is not set as a system variable
	if MONGOSERVER == "" {
		//MONGOSERVER = "mongodb://spankie:506dad@ds163738.mlab.com:63738/btshopng"
		MONGOSERVER = "127.0.0.1:27017"
		log.Println("No mongo server address set, Using default address:", MONGOSERVER)
	}

	// set name of mongo DB to default if it is not set as a system variable
	if MONGODB == "" {
		MONGODB = "btshopng"
	}

	// obtain session from connecting to the mongo server
	session, err := mgo.Dial(MONGOSERVER)
	// log error if available from connecting to DB
	if err != nil {
		log.Println("Error connecting to DB:", err, "shutting down...")
		panic(err)
	}

	// Set safeMode of the session
	// session.SetSafe(&mgo.Safe{})

	// Set configurations
	config = Conf{

		MongoDB:     MONGODB,
		MongoServer: MONGOSERVER,
		Database:    session.DB(MONGODB),
	}

	// log the database in use
	log.Println("Mongo server:", MONGOSERVER)

	// return the Configuration
	// return &config
}

// GetConf returns the App configurations
func GetConf() *Conf {
	return &config
}
