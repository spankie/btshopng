package web

import (
	"net/http"
	"time"

	"github.com/btshopng/btshopng/config"
	"github.com/btshopng/btshopng/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

//Userget reads the json web token(JWT) content from context and marshals it ito a user struct,
func Userget(r *http.Request) (models.User, error) {
	//id := context.Get(r, "UserID")
	//u := context.Get(r, "User")

	u := r.Context().Value("User")
	// log.Println(u)
	// if !ok {
	// 	return models.User{}, errors.New("not a value ")
	// }
	// return *u, nil
	//
	user := models.User{}
	if u != nil {
		err := mapstructure.Decode(u, &user)

		if err != nil {
			return user, err
		}
		return user, nil
	}
	return user, nil
}

//Turn user details into a hasked token that can be used to recognize the user in the future.
func GenerateJWT(user models.User) (resp LoginResponse, err error) {
	claims := jwt.MapClaims{}

	// set our claims
	claims["User"] = user
	claims["Name"] = user.Name

	// set the expire time

	claims["exp"] = time.Now().Add(time.Hour * 24 * 30 * 12).Unix() //24 hours inn a day, in 30 days * 12 months = 1 year in milliseconds

	// create a signer for rsa 256
	t := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)

	pub, err := jwt.ParseRSAPrivateKeyFromPEM(config.GetConf().Encryption.Private)
	if err != nil {
		return
	}
	tokenString, err := t.SignedString(pub)

	if err != nil {
		return
	}

	resp = LoginResponse{
		User:    user,
		Message: "Token succesfully generated",
		Token:   tokenString,
	}

	return

}
