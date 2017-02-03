package models

import (
	// mgo "gopkg.in/mgo.v2"
	"log"
	"time"

	"github.com/btshopng/btshopng/config"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
	FBPicture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
	Link                 string `json:"link"`
	DateCreated          time.Time
	FormattedDateCreated string
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
	}, user)
	if err != nil {
		log.Println(err)
		return err
	}
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
	return user, nil
}
