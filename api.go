package main

import "github.com/lunny/tango"

// APIAction 提供了API访问的能力
type APIAction struct {
	tango.Json
}

func (a *APIAction) Get() interface{} {
	news, err := getNews()
	if err != nil {
		return err
	}

	return news
}
