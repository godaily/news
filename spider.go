package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/mahonia"
	"github.com/PuerkitoBio/goquery"
)

var (
	gbkDecoder = mahonia.NewDecoder("gbk")
)

func gouYouTuanS2T(tm string) time.Time {
	if strings.HasSuffix(tm, "分钟前") {
		min, _ := strconv.Atoi(strings.TrimSpace(strings.TrimRight(tm, "分钟前")))
		return time.Now().Add(-time.Duration(min) * time.Minute)
	} else if strings.HasSuffix(tm, "小时前") {
		hour, _ := strconv.Atoi(strings.TrimSpace(strings.TrimRight(tm, "小时前")))
		return time.Now().Add(-time.Duration(hour) * time.Hour)
	} else if strings.HasSuffix(tm, "天前") {
		min, _ := strconv.Atoi(strings.TrimSpace(strings.TrimRight(tm, "天前")))
		return time.Now().Add(-time.Duration(min) * time.Hour * 24)
	}
	panic("unknow time")
}

func catchGouYouTuan() error {
	resp, err := http.Get("http://golanghome.com/")
	if err != nil {
		return err
	}

	d, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	d.Find("div.post-list").Find("div.post").Each(func(idx int, s *goquery.Selection) {
		imgsrc, _ := s.Find("img").Attr("src")

		a := s.Find("h3.title").Find("a")
		title := a.Text()
		href, _ := a.Attr("href")

		t := s.Find("div.meta").Find("span.time").Last()

		var author string
		s.Find("div.meta").Find("a").Each(func(idx int, s *goquery.Selection) {
			if idx == 2 {
				author = s.Text()
			}
		})

		err := saveNews(GoYouTuan, imgsrc, href, title, author, gouYouTuanS2T(t.Text()))
		if err != nil {
			fmt.Println(err)
		}
	})

	return nil
}

func golangTCS2T(tm string) time.Time {
	if strings.HasSuffix(tm, "分钟前") {
		min, _ := strconv.Atoi(strings.TrimSpace(strings.TrimRight(tm, "分钟前")))
		return time.Now().Add(-time.Duration(min) * time.Minute)
	} else if strings.HasSuffix(tm, "小时前") {
		hour, _ := strconv.Atoi(strings.TrimSpace(strings.TrimRight(tm, "小时前")))
		return time.Now().Add(-time.Duration(hour) * time.Hour)
	} else if strings.HasSuffix(tm, "天前") {
		min, _ := strconv.Atoi(strings.TrimSpace(strings.TrimRight(tm, "天前")))
		return time.Now().Add(-time.Duration(min) * time.Hour * 24)
	}
	t, err := time.Parse("2006-01-02 15:04", tm)
	if err != nil {
		panic(fmt.Sprintf("unknow time format: %v", err))
	}
	return t
}

/*
<dl class="topics">
	<dd>
      <a href="/member/freej" class="pull-left" style="margin-right: 10px;">
		<img class="img-rounded" src="http://gopher.qiniudn.com/avatar/625cba24aee811e2bc7f4e508e16aa57.jpg-middle" alt="freej">
	  </a>
	  <a class="badge pull-right" href="/t/55389a20421aa95094000064#.LatestReplyId.Hex">2</a>
	  <a href="/t/55389a20421aa95094000064" class="title">【北京】【2015-6-6】Golang &amp; Docker Hackathon <span class="glyphicon glyphicon-pushpin"></span></a>
	  <div class="space"></div>
	  <div class="info" style="margin-left:60px">
		<a class="label label-info" href="/go/activity">活动</a> •
		<a href="/member/freej"><strong>freej</strong></a> •
		2015-04-25 09:08 • 最后回复来自 <a href="/member/freej">freej</a>
	  </div>
	  <div class="clear"></div>
	</dd>
*/

var golangTCLink = "http://golangtc.com"

func catchGolangTC() error {
	resp, err := http.Get(golangTCLink)
	if err != nil {
		return err
	}

	d, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	d.Find("dl.topics").Find("dd").Each(func(idx int, s *goquery.Selection) {
		imgsrc, _ := s.Find("img").Attr("src")
		var href string
		var title string

		s.Find("a").Each(func(idx int, g *goquery.Selection) {
			if idx == 2 {
				href, _ = g.Attr("href")
				href = golangTCLink + href
				title = g.Text()
			}
		})

		info := strings.Replace(strings.TrimSpace(s.Find("div.info").Text()), "\n", "", -1)
		info = strings.Replace(info, "\t", "", -1)
		infos := strings.Split(info, "•")
		author := strings.TrimSpace(infos[1])
		t := strings.TrimSpace(infos[2])

		err := saveNews(GolangTC, imgsrc, href, title, author, golangTCS2T(t))
		if err != nil {
			fmt.Println(err)
		}
	})
	return nil
}

var (
	stduyGolangLink = "http://studygolang.com"
)

func catchStudyGolang() error {
	resp, err := http.Get(stduyGolangLink + "/topics")
	if err != nil {
		return err
	}

	d, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	d.Find("div.topic").Each(func(idx int, s *goquery.Selection) {
		imgsrc, _ := s.Find("img.img-rounded").Attr("src")

		a := s.Find("div.title").Find("a")
		href, _ := a.Attr("href")
		href = stduyGolangLink + href
		title := a.Text()
		lastauthor := s.Find("a.author").Last()
		author := lastauthor.Text()
		t, _ := lastauthor.Parent().Find("abbr").Attr("title")

		tm, _ := time.Parse("2006-01-02 15:04:05", t)

		err := saveNews(StudyGoLang, imgsrc, href, title, author, tm)
		if err != nil {
			fmt.Println(err)
		}
	})
	return nil
}

func spiders() {
	for {
		catchGouYouTuan()

		catchGolangTC()

		catchStudyGolang()

		time.Sleep(time.Minute)
	}
}
