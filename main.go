package main

import (
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/lunny/log"
	"github.com/lunny/nodb"
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

type HomeAction struct {
	renders.Renderer
	tango.Ctx
}

type set map[string]interface{}

func (s *set) Get(name string) interface{} {
	return (*s)[name]
}

func (s *set) Set(name string, value interface{}) string {
	(*s)[name] = value
	return ""
}

func (h *HomeAction) Get() error {
	last_read_time := h.Cookies().Get("last_read_time")
	var lastSeconds int64
	if last_read_time != nil {
		lastSeconds, _ = strconv.ParseInt(last_read_time.Value, 10, 64)
	}

	news, err := getNews()
	if err != nil {
		return err
	}

	h.Cookies().Set(tango.NewCookie("last_read_time", fmt.Sprintf("%d", time.Now().Unix())))
	vars := make(set)
	err = h.Render("home.html", renders.T{
		"news":        news,
		"lastSeconds": lastSeconds,
		"vars":        &vars,
	})
	if err != nil {
		return err
	}

	return nil
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

	w := log.NewFileWriter()
	log.Std.SetOutput(w)
	log.Std.SetOutputLevel(log.Lall)

	t := tango.Classic(log.Std)
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
	t.Get("/api/v1/news", new(APIAction))
	t.Run(":8980")
}
