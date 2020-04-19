package main

import (
	"flag"
	"fmt" // TODO: log stuff instead of print.
	"log"
	"manga_server/mongoutil"
	"net/http"
)

func main() {
	// Get cmdline flags.
	dbUsername := flag.String("dbUsername", "default", "Database username.")
	dbPassword := flag.String("dbPassword", "default", "Database password.")
	// jwt secret not being used right now due to the disabling of login/signup functionality.
	// jwtSecret := flag.String("jwtSecret", "default", "Secret used in JWT token.")
	flag.Parse()
	if /* *jwtSecret == "default" ||*/ *dbUsername == "default" || *dbPassword == "default" {
		log.Fatal("Did you set cmd line flags?")
	}
	//jwtutil.Secret = []byte(*jwtSecret)

	// Initialization
	dbErr := mongoutil.Init(*dbUsername, *dbPassword)
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	// Handlers for different paths
	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/main.js", GetStaticFileHandler("/", "./mangaserver_frontend/dist", false))
	http.HandleFunc("/static/", GetStaticFileHandler("/static/", "./static", false))
	http.HandleFunc("/images/", GetStaticFileHandler("/images/", "./mangaserver_frontend/images", false))
	http.HandleFunc("/api/mangalist", ApiMangaListHandler)
	http.HandleFunc("/api/mangainfo", ApiMangaInfoHandler)
	http.HandleFunc("/api/chapterpagecount", ApiChapterPageCount)

	fmt.Println("running server on localhost:80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
