package main

import (
	"fmt"
	"html/template"
	"time"

	"github.com/lunny/nodb"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

type HomeAction struct {
	renders.Renderer
}

func (h *HomeAction) Get() error {
	news, err := getNews()
	if err != nil {
		return err
	}
	return h.Render("home.html", renders.T{
		"news": news,
	})
}

var (
	local, _ = time.LoadLocation("Asia/Chongqing")
)

func main() {
	err := Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	go spiders()

	t := tango.Classic()
	t.Use(renders.New(renders.Options{
		Reload: true,
		Funcs: template.FuncMap{
			"SiteName": SiteName,
			"SiteLink": SiteLink,
			"Time2Str": Time2Str,
		},
		Vars: renders.T{
			"TangoVer": tango.Version(),
			"NodbVer":  nodb.Version,
		},
	}))
	t.Get("/", new(HomeAction))
	t.Run(":8980")
}
