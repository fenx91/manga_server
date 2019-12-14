package main

import (
	"errors"
	"fmt" // TODO: log stuff instead of print.
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"manga_server/mongoutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var jwtSecret = []byte("my_secret_key_fenxy")
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Claims struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	jwt.StandardClaims
}

type IndexTemplateData struct {
	Nickname  string
	MangaData []mongoutil.MangaData
}

type MangaPageTemplateData struct {
	MangaData   mongoutil.MangaData
	ChapterData []mongoutil.ChapterData
}

type ChapterReaderTemplateData struct {
	MangaData    mongoutil.MangaData
	ChapterData  mongoutil.ChapterData
	PicFileNames []string
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

func VerifyTokenAndGetUsername(r *http.Request) (UserEmail string, Nickname string, err error) {
	c, err := r.Cookie("token")
	if err != nil {
		return "", "", err
	}

	tokenString := c.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

	if err != nil {
		return "", "", err
	}
	if token.Valid == false {
		return "", "", errors.New("token not valid")
	}

	return claims.Email, claims.Nickname, nil
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, nickname, err := VerifyTokenAndGetUsername(r)

	if err != nil {
		http.Redirect(w, r, "/loginpage", http.StatusFound)
		return
	}

	t, _ := template.ParseFiles("html/index.html")

	md, err := mongoutil.GetMangaList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = t.Execute(w, &IndexTemplateData{
		Nickname:  nickname,
		MangaData: md,
	})
}

func LogInPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("html/login.html")

	_ = t.Execute(w, nil)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("entered signin handler")
	credentials := &Credentials{
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}
	fmt.Println(credentials)
	if len(credentials.Email) > 254 || !rxEmail.MatchString(credentials.Email) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s is not a valid email address", credentials.Email)
		return
	}

	// Check password match
	expectedHashedPassword, err := mongoutil.GetExpectedPassword(credentials.Email)
	err2 := bcrypt.CompareHashAndPassword([]byte(expectedHashedPassword), []byte(credentials.Password))
	if err != nil || err2 != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Wrong passowrd or user email")
		return
	}

	// Password matched. Creating JWT.
	ud, err := mongoutil.GetUserRegistrationData(credentials.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	claims := &Claims{
		Email:          credentials.Email,
		Nickname:       ud.Nickname,
		StandardClaims: jwt.StandardClaims{},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setTokenCookie(w, tokenString)
	http.Redirect(w, r, "/", http.StatusFound)
}

func SignUpPageHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("html/register.html")

	_ = t.Execute(w, nil)
}

func SignUpActionHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("entered signup handler")
	ud := &mongoutil.UserRegistrationData{
		Email:    r.FormValue("email"),
		Nickname: r.FormValue("nickname"),
		Password: r.FormValue("password"),
	}
	fmt.Println(ud)
	// Check email address is valid
	if len(ud.Email) > 254 || !rxEmail.MatchString(ud.Email) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s is not a valid email address", ud.Email)
		return
	}
	// Check if user registered
	flag, err := mongoutil.DoesUserExist(ud.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if flag {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "User %s has already registered.", ud.Email)
		return
	}

	// Hash and store password
	bytes, err := bcrypt.GenerateFromPassword([]byte(ud.Password), 6)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	ud.Password = string(bytes)
	err = mongoutil.SaveUserRegistrationInfo(*ud)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Register succeeded.")
}

func LogOutHandler(w http.ResponseWriter, r *http.Request) {
	deleteTokenCookie(w)
	fmt.Fprintf(w, "Logout succeed")
}

func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	_, _, err := VerifyTokenAndGetUsername(r)

	if err != nil {
		http.Redirect(w, r, "/loginpage", http.StatusFound)
		return
	}
	// Reject requests to directory.
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}

func MangaPageHandler(w http.ResponseWriter, r *http.Request) {
	_, _, err := VerifyTokenAndGetUsername(r)

	if err != nil {
		http.Redirect(w, r, "/loginpage", http.StatusFound)
		return
	}

	values, _ := url.ParseQuery(r.URL.RawQuery)
	value, ok := values["book"]
	if !ok {
		http.NotFound(w, r)
		return
	}

	mangaId, err := strconv.Atoi(value[0])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	md, cd, err := mongoutil.GetChapterData(mangaId)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t, _ := template.ParseFiles("html/mangapage.html")
	_ = t.Execute(w, &MangaPageTemplateData{
		MangaData:   md,
		ChapterData: cd,
	})
}

func ChapterReaderHandler(w http.ResponseWriter, r *http.Request) {
	_, _, err := VerifyTokenAndGetUsername(r)

	if err != nil {
		http.Redirect(w, r, "/loginpage", http.StatusFound)
		return
	}

	values, _ := url.ParseQuery(r.URL.RawQuery)
	value, ok := values["book"]
	if !ok {
		http.NotFound(w, r)
		return
	}

	chapterNos, ok := values["chapterno"]
	if !ok {
		http.NotFound(w, r)
		return
	}

	if !ok {
		http.NotFound(w, r)
		return
	}

	mangaId, err := strconv.Atoi(value[0])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mangaData, err := mongoutil.GetMandaData(mangaId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	mangaName := mangaData.MangaTitle
	chapterNo := chapterNos[0]
	chapterNoInt, err := strconv.Atoi(chapterNo)
	if err != nil || chapterNoInt > mangaData.ChapterCount {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var picFileNames []string
	mangaChapterRootDir := "static/manga/" + mangaName + "/" + chapterNo + "/"
	filepath.Walk(mangaChapterRootDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			picFileNames = append(picFileNames, info.Name())
		}
		return nil
	})

	crtd := ChapterReaderTemplateData{
		MangaData:    mongoutil.MangaData{MangaTitle: mangaName},
		ChapterData:  mongoutil.ChapterData{ChapterNo: chapterNo},
		PicFileNames: picFileNames,
	}

	t, _ := template.ParseFiles("html/chapterreader.html")
	_ = t.Execute(w, crtd)
}

func main() {
	// Initialization
	dbErr := mongoutil.Init()
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	// Handlers for different paths
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/mangapage", MangaPageHandler)
	http.HandleFunc("/chapterreader", ChapterReaderHandler)
	http.HandleFunc("/signinaction", SignInHandler)
	http.HandleFunc("/loginpage", LogInPageHandler)
	http.HandleFunc("/signuppage", SignUpPageHandler)
	http.HandleFunc("/signupaction", SignUpActionHandler)
	http.HandleFunc("/logout", LogOutHandler)
	http.HandleFunc("/static/", StaticFileHandler)

	fmt.Println("running server on localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
