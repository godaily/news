package main

import (
	"fmt"
	"testing"
)

func TestGouYouTuan(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}

	err = catchGouYouTuan()
	if err != nil {
		t.Fatal(err)
	}

	news, err := getNews()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(news)
}

func TestGolangTC(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}

	err = catchGolangTC()
	if err != nil {
		t.Fatal(err)
	}

	news, err := getNews()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(news)
}

func TestStudyGolang(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}

	err = catchStudyGolang()
	if err != nil {
		t.Fatal(err)
	}

	news, err := getNews()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(news)
}
