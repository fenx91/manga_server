package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

type MangaData struct {
	MangaTitle string
	//  TODO() add: MangaId int  // Id in database
	ChapterCount int
}

type ChapterData struct {
	ChapterNo string // string as easier to use
}

type IndexTemplateData struct {
	Username  string
	MangaData []MangaData
}

type MangaPageTemplateData struct {
	MangaData   MangaData
	ChapterData []ChapterData
}

type ChapterReaderTemplateData struct {
	MangaData    MangaData
	ChapterData  ChapterData
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

	var md []MangaData
	mangaRootDir := "static/manga/"
	filepath.Walk(mangaRootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != mangaRootDir {
			md = append(md, MangaData{
				MangaTitle:   info.Name(),
				ChapterCount: 0, // TODO() add logic if necessary
			})
			return filepath.SkipDir
		}
		return nil
	})

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

	mangaName := value[0]
	// TODO: check if mangaName is valid

	var cd []ChapterData
	mangaRootDir := "static/manga/" + mangaName + "/"
	filepath.Walk(mangaRootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != mangaRootDir {
			cd = append(cd, ChapterData{ChapterNo: info.Name()})
			return filepath.SkipDir
		}
		return nil
	})

	md := MangaData{
		MangaTitle:   mangaName,
		ChapterCount: len(cd),
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

	mangaName := value[0]
	chapterNo := chapterNos[0]
	// TODO: check if mangaName and chapter no is valid

	var picFileNames []string
	mangaChapterRootDir := "static/manga/" + mangaName + "/" + chapterNo + "/"
	filepath.Walk(mangaChapterRootDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			picFileNames = append(picFileNames, info.Name())
		}
		return nil
	})

	crtd := ChapterReaderTemplateData{
		MangaData:    MangaData{MangaTitle: mangaName},
		ChapterData:  ChapterData{ChapterNo: chapterNo},
		PicFileNames: picFileNames,
	}

	t, _ := template.ParseFiles("html/chapterreader.html")
	_ = t.Execute(w, crtd)
}

func main() {
	// Handlers for different paths
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/mangapage", MangaPageHandler)
	http.HandleFunc("/chapterreader", ChapterReaderHandler)
	http.HandleFunc("/signinaction", SignInHandler)
	http.HandleFunc("/loginpage", LogInPageHandler)
	http.HandleFunc("/logout", LogOutHandler)
	http.HandleFunc("/static/", StaticFileHandler)

	fmt.Println("running server on localhost:10001")
	log.Fatal(http.ListenAndServe(":10001", nil))
}
