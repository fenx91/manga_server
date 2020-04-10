package main

import (
	"encoding/json"
	"fmt"
	"github.com/mssola/user_agent"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"manga_server/jwtutil"
	"manga_server/mongoutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	_, nickname, err := jwtutil.VerifyTokenAndGetUsername(r)

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
	signedToken, err := jwtutil.CreateJwt(credentials.Email, ud.Nickname)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setTokenCookie(w, signedToken)
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
	_, _, err := jwtutil.VerifyTokenAndGetUsername(r)

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
	_, _, err := jwtutil.VerifyTokenAndGetUsername(r)

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
	if user_agent.New(r.UserAgent()).Mobile() {
		http.Redirect(w, r, "/m"+r.URL.String(), http.StatusFound)
	}

	_, _, err := jwtutil.VerifyTokenAndGetUsername(r)

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

func MobileChapterReaderHandler(w http.ResponseWriter, r *http.Request) {
	_, _, err := jwtutil.VerifyTokenAndGetUsername(r)

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

	t, _ := template.ParseFiles("html/mobilechapterreader.html")
	_ = t.Execute(w, crtd)
}

func MangaListHandler(w http.ResponseWriter, r *http.Request) {
	mangaDataList, err := mongoutil.GetMangaList()
	if err != nil {
		fmt.Println("get manga list failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(mangaDataList)
	if err != nil {
		fmt.Println("js marshal manga list failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
