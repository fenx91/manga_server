package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"log"
	"net/http"
)

var users = map[string]string{
	"fenxy": "password1",
}

var jwtSecret = []byte("my_secret_key")

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type TemplateData struct {
	Username string
}

func setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
	})
}

func deleteTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "token",
		MaxAge: 0,
	})
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("entered signin handler")
	credentials := &Credentials{
		r.FormValue("password"),
		r.FormValue("username"),
	}
	fmt.Println(credentials)

	expectedPassword, ok := users[credentials.Username]

	if !ok || expectedPassword != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Wwong passowrd or username")
		return
	}
	claims := &Claims{
		Username:       credentials.Username,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	setTokenCookie(w, tokenString)
	//fmt.Fprintf(w, "login succeed. <a href=\"/\"> Go to homepage</a>")
	http.Redirect(w, r, "/", http.StatusFound)
}

func VerifyAndGetUsername(r *http.Request) (Username string, err error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", err
	}

	tokenString := c.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

	if err != nil {
		return "", err
	}
	if token.Valid == false {
		return "", errors.New("token not valid")
	}

	return claims.Username, nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	username, err := VerifyAndGetUsername(r)

	if err != nil {
		http.Redirect(w, r, "/loginpage", http.StatusFound)
		return
	}

	t, _ := template.ParseFiles("html/index.html")

	_ = t.Execute(w, &TemplateData{
		Username: username,
	})
}

func LogInPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("html/login.html")

	_ = t.Execute(w, nil)
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	deleteTokenCookie(w)
	fmt.Fprintf(w, "Logout succeed")
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/signinaction", SignInHandler)
	http.HandleFunc("/loginpage", LogInPageHandler)
	http.HandleFunc("/logout", LogOutHandler)
	fmt.Println("running server on :10001")
	log.Fatal(http.ListenAndServe(":10001", nil))
}
