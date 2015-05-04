package main

import (
	"crypto/md5"
	"fmt"
	"math"
	"time"

	"github.com/lunny/nodb"
	"github.com/lunny/nodb/config"
)

type Site int

const (
	GoYouTuan Site = iota + 1
	GolangTC
	StudyGoLang
)

var (
	sites = map[Site]string{
		GoYouTuan:   "Go友团",
		GolangTC:    "Golang中国",
		StudyGoLang: "Study Golang",
	}

	links = map[Site]string{
		GoYouTuan:   "http://golanghome.com",
		GolangTC:    "http://golangtc.com",
		StudyGoLang: "http://studygolang.com",
	}
)

func SiteName(site Site) string {
	return sites[site]
}

func SiteLink(site Site) string {
	return links[site]
}

func gmd5(ori string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(ori)))
}

type News struct {
	Id    int64
	Title string
	Image string
	Url   string
	Site
	Author      string
	AuthorLink  string
	Author2     string
	Author2Link string
	Updated     time.Time
}

func Time2Str(t time.Time) string {
	delta := time.Now().Sub(t)
	if delta < time.Minute {
		return fmt.Sprintf("%d 秒前", int(delta/time.Second))
	} else if delta < time.Hour {
		return fmt.Sprintf("%d 分钟前", int(delta/time.Minute))
	} else if delta < time.Hour*24 {
		return fmt.Sprintf("%d 小时前", int(delta/time.Hour))
	}
	return fmt.Sprintf("%d 天前", int(delta/(24*time.Hour)))
}

var (
	db *nodb.DB
)

func Init() error {
	cfg := config.NewConfigDefault()
	cfg.DataDir = "./db"

	var err error
	// init nosql
	ndb, err := nodb.Open(cfg)
	if err != nil {
		return err
	}

	// select db
	db, err = ndb.Select(0)
	return err
}

var (
	updatedKey = []byte("updated")
)

// TODO: transaction
func saveNews(site Site, imgUrl, articleUrl, title, author1, author1Link,
	author2, author2Link string, postTime time.Time) error {
	id, err := nodb.StrInt64(db.Get([]byte("urlkey:" + gmd5(articleUrl))))
	if err != nil {
		return err
	}

	if id > 0 {
		member := nodb.StrPutInt64(id)
		score, err := db.ZScore(updatedKey, member)
		if err != nil {
			return err
		}

		db.Set([]byte(fmt.Sprintf("author2:%d", id)), []byte(author2))
		db.Set([]byte(fmt.Sprintf("author2Link:%d", id)), []byte(author2Link))
		delta := postTime.Unix() - score
		if delta == 0 {
			delta = 1
		}
		_, err = db.ZIncrBy(updatedKey, delta, member)
		if err != nil {
			return err
		}
	} else {
		id, err = db.Incr([]byte("index"))
		if err != nil {
			return err
		}
		member := nodb.StrPutInt64(id)
		db.Set([]byte(fmt.Sprintf("site:%d", id)), nodb.StrPutInt64(int64(site)))
		db.Set([]byte(fmt.Sprintf("image:%d", id)), []byte(imgUrl))
		db.Set([]byte(fmt.Sprintf("url:%d", id)), []byte(articleUrl))
		db.Set([]byte(fmt.Sprintf("title:%d", id)), []byte(title))
		db.Set([]byte(fmt.Sprintf("author:%d", id)), []byte(author1))
		db.Set([]byte(fmt.Sprintf("authorLink:%d", id)), []byte(author1Link))
		db.Set([]byte(fmt.Sprintf("author2:%d", id)), []byte(author2))
		db.Set([]byte(fmt.Sprintf("author2Link:%d", id)), []byte(author2Link))
		db.Set([]byte("urlkey:"+gmd5(articleUrl)), member)
		db.ZAdd(updatedKey, nodb.ScorePair{postTime.Unix(), member})
	}
	return nil
}

func getNews() ([]News, error) {
	scores, err := db.ZRevRangeByScore([]byte("updated"), 0, math.MaxInt64, 0, 20)
	if err != nil {
		return nil, err
	}

	var news = make([]News, len(scores))
	for i, scorepair := range scores {
		id, _ := nodb.StrInt64(scorepair.Member, nil)

		bsite, err := db.Get([]byte(fmt.Sprintf("site:%d", id)))
		if err != nil {
			return nil, err
		}
		site, _ := nodb.StrInt64(bsite, nil)

		bImage, err := db.Get([]byte(fmt.Sprintf("image:%d", id)))
		if err != nil {
			return nil, err
		}

		bUrl, err := db.Get([]byte(fmt.Sprintf("url:%d", id)))
		if err != nil {
			return nil, err
		}

		bTitle, err := db.Get([]byte(fmt.Sprintf("title:%d", id)))
		if err != nil {
			return nil, err
		}

		bAuthor, err := db.Get([]byte(fmt.Sprintf("author:%d", id)))
		if err != nil {
			return nil, err
		}

		bAuthorLink, err := db.Get([]byte(fmt.Sprintf("authorLink:%d", id)))
		if err != nil {
			return nil, err
		}

		bAuthor2, err := db.Get([]byte(fmt.Sprintf("author2:%d", id)))
		if err != nil {
			return nil, err
		}

		bAuthor2Link, err := db.Get([]byte(fmt.Sprintf("author2Link:%d", id)))
		if err != nil {
			return nil, err
		}

		news[i] = News{
			Id:          id,
			Site:        Site(site),
			Image:       string(bImage),
			Title:       string(bTitle),
			Url:         string(bUrl),
			Author:      string(bAuthor),
			AuthorLink:  string(bAuthorLink),
			Author2:     string(bAuthor2),
			Author2Link: string(bAuthor2Link),
			Updated:     time.Unix(scorepair.Score, 0),
		}
	}
	return news, nil
}
