package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/goutils"
	"github.com/ahmdrz/goinsta"
	instarep "github.com/ahmdrz/goinsta/response"
	sentry "github.com/getsentry/raven-go"
	"github.com/gorilla/feeds"
)

var sendtoSentry bool

func init() {
	if os.Getenv("SENTRY_DSN") != "" {
		sendtoSentry = true
	}
}

func main() {
	if len(os.Args) != 2 {
		dieOnError("missing Instagram username argument")
	}

	igLogin := os.Getenv("IG_LOGIN")
	igPassword := os.Getenv("IG_PASSWORD")

	if igLogin == "" || igPassword == "" {
		dieOnError("IG_LOGIN or IG_PASSWORD environment variable unset")
	}

	igUser := os.Args[1]

	insta := goinsta.New(igLogin, igPassword)
	if err := insta.Login(); err != nil {
		dieOnError("unable to initialize Instagram client: %s", err)
	}
	defer insta.Logout()

	user, err := insta.GetUserByUsername(igUser)
	if err != nil {
		dieOnError("unable to get user information: %s", err)
	}

	latest, err := insta.LatestUserFeed(user.User.ID)
	if err != nil {
		dieOnError("unable to get user latest feed: %s", err)
	}

	feed := &feeds.Feed{
		Title: fmt.Sprintf("Instagram Feed: %s", user.User.FullName),
		Link:  &feeds.Link{Href: fmt.Sprintf("https://www.instagram.com/%s", igUser)},
		Author: &feeds.Author{
			Name:  user.User.FullName,
			Email: user.User.Username + "@instagram.com",
		},
		Description: user.User.Biography,
		Created:     time.Now(),
	}

	for _, item := range latest.Items {
		fi, err := formatFeedContent(&item)
		if err != nil {
			dieOnError("unable to format feed item: %s", err)
		}

		feed.Add(fi)
	}

	rss, err := feed.ToRss()
	if err != nil {
		dieOnError("unable to render RSS feed: %s", err)
	}

	fmt.Println(rss)
}

func formatFeedContent(item *instarep.Item) (*feeds.Item, error) {
	shortDesc, _ := goutils.Abbreviate(item.Caption.Text, 50)

	content := fmt.Sprintf("<p>%s</p>", item.Caption.Text)

	if len(item.ImageVersions2.Candidates) > 0 {
		content += fmt.Sprintf("<p><img src=%q></p>", item.ImageVersions2.Candidates[0].URL)
	}

	if len(item.CarouselMedia) > 0 {
		for _, i := range item.CarouselMedia {
			content += fmt.Sprintf("<p><img src=%q></p>", i.ImageVersions.Candidates[0].URL)
		}
	}

	return &feeds.Item{
		Title:       shortDesc,
		Created:     time.Unix(item.TakenAt, 0),
		Description: shortDesc,
		Content:     content,
		Link:        &feeds.Link{Href: fmt.Sprintf("https://www.instagram.com/p/%s", item.Code)},
	}, nil
}

func dieOnError(format string, a ...interface{}) {
	if sendtoSentry {
		sentry.CaptureErrorAndWait(fmt.Errorf(format, a...), nil)
	}

	fmt.Fprintf(os.Stderr, fmt.Sprintf("error: %s\n", format), a...)
	os.Exit(1)
}
