package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"log"

	"github.com/btshopng/btshopng/config"
)

// Barter holds data for individual barter.
type Barter struct {
	ID           string
	UserID       string
	Have         string
	HaveCategory string
	Need         string
	NeedCategory string
	Location     string
	DateCreated  time.Time
	Status       string   // status of the barter: available or not.
	Images       []string // array of string to the path where the images will be stored.

}

func (barter Barter) Get(c *config.Conf) (Barter, error) {
	mgoSession := c.Database.Session.Copy()
	defer mgoSession.Close()

	collection := c.Database.C(config.BARTERCOLLECTION).With(mgoSession)

	result := Barter{}

	err := collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{"id": barter.ID},
			bson.M{"userid": barter.UserID},
		},
	}).One(&result)

	if err != nil {
		log.Println("Get barter error:", err)
		return barter, err
	}

	return result, nil
}

// Upsert inserts or update a barter
func (barter Barter) Upsert(c *config.Conf) error {
	mgoSession := c.Database.Session.Copy()
	defer mgoSession.Close()

	collection := c.Database.C(config.BARTERCOLLECTION).With(mgoSession)

	_, err := collection.Upsert(bson.M{"id": barter.ID}, bson.M{"$set": barter})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// GetAllBarters Barter by a particular user
func GetAllBarters(c *config.Conf, id string) ([]Barter, error) {
	mgoSession := c.Database.Session.Copy()
	defer mgoSession.Close()

	collection := c.Database.C(config.BARTERCOLLECTION).With(mgoSession)
	result := []Barter{}
	err := collection.Find(bson.M{"userid": id}).All(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
