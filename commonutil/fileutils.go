package commonutil

import (
	"io/ioutil"
)

func GetChapterPageCount(mangaTitle string, chapterNo int) (int, error) {
	files, err := ioutil.ReadDir("static/manga/" + mangaTitle + "/" + TwoDigitInt(chapterNo))
	if err != nil {
		return 0, err
	}
	return len(files), nil
}
