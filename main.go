package main

import (
	"flag"
	"fmt" // TODO: log stuff instead of print.
	"log"
	"manga_server/jwtutil"
	"manga_server/mongoutil"
	"net/http"
)

func main() {
	// Get cmdline flags.
	jwtSecret := flag.String("jwtSecret", "default", "Secret used in JWT token.")
	dbUsername := flag.String("dbUsername", "default", "Database username.")
	dbPassword := flag.String("dbPassword", "default", "Database password.")
	flag.Parse()
	if *jwtSecret == "default" || *dbUsername == "default" || *dbPassword == "default" {
		log.Fatal("Did you set cmd line flags?")
	}
	jwtutil.Secret = []byte(*jwtSecret)

	// Initialization
	dbErr := mongoutil.Init(*dbUsername, *dbPassword)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	// Handlers for different paths
	http.HandleFunc("/", RootHandler) // GetStaticFileHandler("/", "./mangaserver_frontend/dist", true))
	http.HandleFunc("/main.js", GetStaticFileHandler("/", "./mangaserver_frontend/dist", false))
	//http.HandleFunc("/mangapage/", GetStaticFileHandler("/mangapage/", "./mangaserver_frontend/dist", true))
	//http.HandleFunc("/mangapage", MangaPageHandler)
	http.HandleFunc("/chapterreader", ChapterReaderHandler)
	http.HandleFunc("/m/chapterreader", MobileChapterReaderHandler)
	http.HandleFunc("/signinaction", SignInHandler)
	http.HandleFunc("/loginpage", LogInPageHandler)
	http.HandleFunc("/signuppage", SignUpPageHandler)
	http.HandleFunc("/signupaction", SignUpActionHandler)
	http.HandleFunc("/logout", LogOutHandler)
	http.HandleFunc("/static/", GetStaticFileHandler("/static/", "./static", false))
	http.HandleFunc("/images/", GetStaticFileHandler("/images/", "./mangaserver_frontend/images", false))
	http.HandleFunc("/api/mangalist", ApiMangaListHandler)
	http.HandleFunc("/api/mangainfo", ApiMangaInfoHandler)

	fmt.Println("running server on localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
