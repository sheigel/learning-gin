package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/dgrijalva/jwt-go"

	"net/http"
	"database/sql"

	"time"
	"fmt"
	"context"
	"log"
	"os"

	"github.com/sheigel/learning-gin/auth"
)

var db *sql.DB

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func main() {
	db = connectToDb()
	a := auth.NewAuth(db)
	http.HandleFunc("/register", a.Register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/profile", validate(protectedProfile))
	http.HandleFunc("/logout", validate(logout))

	log.Println("Listenning!!!")
	http.ListenAndServe(":9000", nil)
}

func connectToDb() *sql.DB {
	db, err := sql.Open("mysql", "gotest:1234@tcp(db:3306)/demo")

	if err != nil {
		log.Printf("Couldn't connect to database! Error: %s", err)
		os.Exit(-1)
	}

	err = db.Ping()

	if err != nil {
		log.Printf("Couldn't connect to database! Error: %s", err)
		os.Exit(-1)
	}

	return db
}

func login(res http.ResponseWriter, req *http.Request) {

	// Expires the token and cookie in 1 hour
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	expireCookie := time.Now().Add(time.Hour * 1)

	// We'll manually assign the claims but in production you'd insert values from a database
	claims := Claims{
		"myusername",
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "localhost:9000",
		},
	}

	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	signedToken, _ := token.SignedString([]byte("secret"))

	// Place the token in the client's cookie
	cookie := http.Cookie{Name: "Auth", Value: signedToken, Expires: expireCookie, HttpOnly: true}
	http.SetCookie(res, &cookie)

	// Redirect the user to his profile
	http.Redirect(res, req, "/profile", 307)
}

func validate(protectedPage http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Auth")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {

			ctx := context.WithValue(r.Context(), "MyKey", *claims)

			protectedPage(w, r.WithContext(ctx))
			return
		} else {
			http.NotFound(w, r)
			return
		}
	})
}

func protectedProfile(res http.ResponseWriter, req *http.Request) {
	claims, ok := req.Context().Value("MyKey").(Claims)
	if !ok {
		http.NotFound(res, req)
		return
	}

	fmt.Fprintf(res, "Hello %s", claims.Username)
}

func logout(res http.ResponseWriter, req *http.Request) {
	deleteCookie := http.Cookie{Name: "Auth", Value: "none", Expires: time.Now()}
	http.SetCookie(res, &deleteCookie)
	return
}
