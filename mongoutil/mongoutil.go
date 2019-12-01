package mongoutil

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbUsername = "fenxy"
const dbPassword = "4REQlkmb"
const dbName = "manga_server"
const collectionMangas = "mangas"

var client *mongo.Client
var mangasCollection *mongo.Collection

func Init() (err error) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{Username: dbUsername, Password: dbPassword})

	// Connect to MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return err
	}

	mangasCollection = client.Database(dbName).Collection(collectionMangas)
	fmt.Println("Connected to MongoDB as fenxy!")
	return nil
}

func GetMangaList() (MangaData []MangaData, err error) {
	// Find every entry in the DB.
	cursor, err := mangasCollection.Find(context.TODO(), bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.TODO(), &MangaData)
	if err != nil {
		return nil, err
	}

	return MangaData, nil
}

func GetMandaData(mangaId int) (md MangaData, err error) {
	err = mangasCollection.FindOne(context.TODO(), bson.D{{"id", mangaId}}).Decode(&md)
	if err != nil {
		return MangaData{}, nil
	}
	return md, nil
}

func GetChapterData(mangaId int) (md MangaData, cd []ChapterData, err error) {
	err = mangasCollection.FindOne(context.TODO(), bson.D{{"id", mangaId}}).Decode(&md)
	if err != nil {
		return MangaData{}, nil, err
	}

	for i := 1; i <= md.ChapterCount; i++ {
		cd = append(cd, ChapterData{ChapterNo: fmt.Sprintf("%02d", i)})
	}
	return md, cd, nil
}
