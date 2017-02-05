package models

import (
	// mgo "gopkg.in/mgo.v2"
	"log"
	"time"

	"github.com/btshopng/btshopng/config"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID    string `json:"id" bson:",omitempty"`
	Name  string `json:"name" bson:",omitempty"`
	Email string `json:"email" bson:",omitempty"`
	Image struct {
		URL string `json:"url" bson:",omitempty"`
	} `json:"image" bson:",omitempty"`
	FBPicture struct {
		Data struct {
			URL string `json:"url" bson:",omitempty"`
		} `json:"data" bson:",omitempty"`
	} `json:"picture" bson:",omitempty"`
	Link                 string    `json:"link" bson:",omitempty"`
	DateCreated          time.Time `bson:",omitempty"`
	FormattedDateCreated string    `bson:",omitempty"`
	Password             []byte    `bson:",omitempty"`
}

func (user User) Upsert(c *config.Conf) error {

	mgoSession := c.Database.Session.Copy()
	defer mgoSession.Close()

	collection := c.Database.C(config.USERCOLLECTION).With(mgoSession)

	_, err := collection.Upsert(bson.M{
		"$or": []bson.M{
			bson.M{
				"id": user.ID,
			},
			bson.M{
				"email": user.Email,
			},
		},
	}, bson.M{"$set": user})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("upsert user:", user)
	return nil
}

func (user User) Get(c *config.Conf) (User, error) {

	mgoSession := c.Database.Session.Copy()
	defer mgoSession.Close()

	collection := c.Database.C(config.USERCOLLECTION).With(mgoSession)

	result := User{}
	err := collection.Find(bson.M{
		"$or": []bson.M{
			bson.M{
				"id": user.ID,
			},
			bson.M{
				"email": user.Email,
			},
		},
	}).One(&result)

	if err != nil {
		log.Println(err)
		return user, err
	}
	return result, nil
}
