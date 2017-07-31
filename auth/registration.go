package auth

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Auth struct {
	db *sql.DB
}

type Login_data struct {
	Username string `json:"un"`
	Password string `json:"pwd"`
}

func NewAuth(db *sql.DB) Auth {
	return Auth{db: db}
}

func (a Auth) Register(res http.ResponseWriter, req *http.Request) {
	var login *Login_data

	err := json.NewDecoder(req.Body).Decode(login)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	log.Printf("register data %v", *login)

	err = a.db.QueryRow("SELECT username FROM users WHERE username=?", login.Username).Scan()

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Couldn't generate password!\n Error:%s\n", err)
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = a.db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", login.Username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			log.Printf("Couldn't insert user!\n Error:%s\n", err)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		log.Printf("Couldn't select user!\n Error:%s\n", err)
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}
