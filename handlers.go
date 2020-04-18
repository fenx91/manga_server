package main

import (
	"encoding/json"
	"fmt"
	"log"
	"manga_server/commonutil"
	"manga_server/mongoutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func GetStaticFileHandler(path string, dir string, allowRoot bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Reject requests to directory if "allowRoot" is set to false.
		if !allowRoot && strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		// simulates net delay.
		/*if path == "/static/" {
			num := 3000 + rand.Intn(2000)
			time.Sleep(time.Duration(num) * time.Millisecond)
		}*/
		http.StripPrefix(path, http.FileServer(http.Dir(dir))).ServeHTTP(w, r)
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix(r.URL.Path, http.FileServer(http.Dir("./mangaserver_frontend/dist"))).
		ServeHTTP(w, r)
}

func ApiMangaListHandler(w http.ResponseWriter, r *http.Request) {
	mangaDataList, err := mongoutil.GetMangaList()
	if err != nil {
		fmt.Println("get manga list failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type ApiMangaList struct {
		MangaDataList []mongoutil.MangaData
	}
	mangalist := ApiMangaList{MangaDataList: mangaDataList}
	js, err := json.Marshal(mangalist)
	if err != nil {
		fmt.Println("js marshal manga list failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ApiMangaInfoHandler(w http.ResponseWriter, r *http.Request) {
	values, _ := url.ParseQuery(r.URL.RawQuery)
	value, ok := values["mangaid"]
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
		log.Println("failed to get manga info for " + string(mangaId) + " in mongodb.")
		http.NotFound(w, r)
		return
	}

	js, err := json.Marshal(mangaData)
	if err != nil {
		fmt.Println("js marshal manga info failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ApiChapterPageCount(w http.ResponseWriter, r *http.Request) {
	// Response Struct
	type MangaChapterInfo struct {
		ChapterNo int
		PageCount int
	}
	type ApiChapterPageCountResponse struct {
		MangaId              int
		MangaTitle           string
		MangaChapterInfoList []MangaChapterInfo
	}
	// Start extracting params.
	values, _ := url.ParseQuery(r.URL.RawQuery)
	// extract mangaId.
	mangaIds, ok := values["mangaid"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mangaId, err := strconv.Atoi(mangaIds[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mangaData, err := mongoutil.GetMandaData(mangaId)
	if err != nil {
		log.Println("failed to get manga info for mangaId=" + string(mangaId) + " in mongodb.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mangaTitle := mangaData.MangaTitle
	totalChapterCount := mangaData.ChapterCount

	response := ApiChapterPageCountResponse{
		MangaId:              mangaId,
		MangaTitle:           mangaTitle,
		MangaChapterInfoList: []MangaChapterInfo{},
	}
	// extract chapter no. if present.
	chapterNos, chapterNoPresent := values["chapterno"]
	if chapterNoPresent {
		chapterNo, err := strconv.Atoi(chapterNos[0])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		pageCount, err := commonutil.GetChapterPageCount(mangaTitle, chapterNo)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		chapterInfo := MangaChapterInfo{
			ChapterNo: chapterNo,
			PageCount: pageCount,
		}
		response.MangaChapterInfoList = append(response.MangaChapterInfoList, chapterInfo)
	} else {
		// if chapter no. not present, return info for all chapters.
		for i := 1; i <= totalChapterCount; i++ {
			pageCount, err := commonutil.GetChapterPageCount(mangaTitle, i)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			chapterInfo := MangaChapterInfo{
				ChapterNo: i,
				PageCount: pageCount,
			}
			response.MangaChapterInfoList = append(response.MangaChapterInfoList, chapterInfo)
		}
	}

	js, err := json.Marshal(response)
	if err != nil {
		fmt.Println("js marshal chapter info response failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*   LOGIN/SIGNUP related functionality disabled at the moment.
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
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
*/
