package mongoutil

type MangaData struct {
	//ObjectId     primitive.ObjectID `bson:"_id"`
	MangaId      int    `bson:"id"`
	MangaTitle   string `bson:"name"`
	ChapterCount int    `bson:"ChapterNo"`
}

type ChapterData struct {
	ChapterNo string // string as easier to use
}

type UserRegistrationData struct {
	Email    string
	Nickname string
	Password string
}
