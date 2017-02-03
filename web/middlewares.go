package web

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/btshopng/btshopng/config"

	"context"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/btshopng/btshopng/messages"
)

// Middlewares

func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				messages.WriteError(w, messages.ErrInternalServer)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func AcceptHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			messages.WriteError(w, messages.ErrNotAcceptable)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func ContentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			messages.WriteError(w, messages.ErrUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

//BodyHandler copies the body of the request object into the context, in case multiple middlewares want to read the content, so that the content doesnt disappear before the final handler
func BodyHandler(v interface{}) func(http.Handler) http.Handler {
	t := reflect.TypeOf(v)

	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			val := reflect.New(t).Interface()
			err := json.NewDecoder(r.Body).Decode(val)

			if err != nil {
				log.Println(err)
				messages.WriteError(w, messages.ErrBadRequest)
				return
			}

			if next != nil {
				//context.Set(r, "body", val)
				ctx := context.WithValue(r.Context(), "body", val)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			}
		}

		return http.HandlerFunc(fn)
	}

	return m
}

//FrontAuthHandler checks if authenticated, and moves the user details from the JWT token into the request context
func FrontAuthHandler(next http.Handler) http.Handler {
	ac := config.GetConf()
	fn := func(w http.ResponseWriter, r *http.Request) {

		// check if we have a cookie with out tokenName
		//tokenValue := r.Header.Get("X-AUTH-TOKEN")
		var cookie, err = r.Cookie("AuthToken")
		if err != nil {
			log.Println(err)

			messages.WriteError(w, messages.ErrBadToken)
			return
		}
		tokenValue := cookie.Value
		log.Println(tokenValue)

		// validate the token
		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(ac.Encryption.Public)

			if err != nil {
				return publicKey, err
			}
			return publicKey, nil
		})

		// branch out into the possible error from signing
		switch err.(type) {

		case nil: // no error
			if !token.Valid { // but may still be invalid
				log.Println(err)

				messages.WriteError(w, messages.ErrBadToken)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				ctx := context.WithValue(r.Context(), "User", claims["User"])
				r = r.WithContext(ctx)
				//context.Set(r, "User", claims["User"])
				//context.Set(r, "UserID", claims["UserID"])
			} else {
				log.Println(err)
				return
			}

			next.ServeHTTP(w, r)

		case *jwt.ValidationError: // something was wrong during the validation
			vErr := err.(*jwt.ValidationError)

			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				messages.WriteError(w, messages.ErrBadToken)
				return
			default:
				messages.WriteError(w, messages.ErrBadToken)
				log.Printf("ValidationError error: %+v\n", vErr.Errors)
				return
			}

		default: // something else went wrong
			messages.WriteError(w, messages.ErrBadToken)
			return
		}
	}
	return http.HandlerFunc(fn)

}
