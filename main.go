package main

import (
	"errors"
	"fmt" // TODO: log stuff instead of print.
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"log"
	"manga_server/mongoutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var users = map[string]string{
	"fenxy": "fenxy",
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

type IndexTemplateData struct {
	Username  string
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
		fmt.Fprintf(w, "Wrong passowrd or username")
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

	md, err := mongoutil.GetMangaList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = t.Execute(w, &IndexTemplateData{
		Username:  username,
		MangaData: md,
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

func StaticFileHandler(w http.ResponseWriter, r *http.Request) {
	_, err := VerifyAndGetUsername(r)

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
	_, err := VerifyAndGetUsername(r)

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
	_, err := VerifyAndGetUsername(r)

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
	http.HandleFunc("/logout", LogOutHandler)
	http.HandleFunc("/static/", StaticFileHandler)

	fmt.Println("running server on localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
