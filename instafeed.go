package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Masterminds/goutils"
	"github.com/ahmdrz/goinsta/v2"
	sentry "github.com/getsentry/raven-go"
	"github.com/gorilla/feeds"
)

var (
	sendtoSentry bool
	configFile   string
	feedMaxItems int
)

func init() {
	if os.Getenv("SENTRY_DSN") != "" {
		sendtoSentry = true
	}

	flag.StringVar(&configFile, "f", path.Join(os.Getenv("HOME"), ".instafeed"),
		"Path to file where to store profile configuration")
	flag.IntVar(&feedMaxItems, "n", 20, "Number of user feed items")
	flag.Parse()
}

func main() {
	var (
		insta *goinsta.Instagram
		err   error
	)

	if len(os.Args) != 2 {
		dieOnError("missing Instagram username argument")
	}

	igUser := os.Args[1]
	igLogin := os.Getenv("IG_LOGIN")
	igPassword := os.Getenv("IG_PASSWORD")

	if igLogin == "" || igPassword == "" {
		dieOnError("IG_LOGIN or IG_PASSWORD environment variable unset")
	}

	if insta, err = goinsta.Import(configFile); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to import Instagram configuration: %s\n", err)
		fmt.Fprintln(os.Stderr, "Attempting new login")

		insta = goinsta.New(igLogin, igPassword)
		if err = insta.Login(); err != nil {
			dieOnError("unable to initialize Instagram client: %s", err)
		}
	}

	if err := insta.Export(configFile); err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to export Instagram client configuration: %s\n", err)
	}

	user, err := insta.Profiles.ByName(igUser)
	if err != nil {
		dieOnError("unable to get user information: %s", err)
	}

	latest := user.Feed()
	if err != nil {
		dieOnError("unable to get user latest feed: %s", err)
	}

	feed := &feeds.Feed{
		Title: fmt.Sprintf("Instagram Feed: %s", user.FullName),
		Link:  &feeds.Link{Href: fmt.Sprintf("https://www.instagram.com/%s", igUser)},
		Author: &feeds.Author{
			Name:  user.FullName,
			Email: user.Username + "@instagram.com",
		},
		Description: user.Biography,
		Created:     time.Now(),
	}

	for latest.Next(false) {
		for _, item := range latest.Items {
			fi, err := formatFeedContent(&item)
			if err != nil {
				dieOnError("unable to format feed item: %s", err)
			}

			feed.Add(fi)
			if len(feed.Items) >= feedMaxItems {
				goto done
			}
		}

		if err := latest.Error(); err != nil {
			if err := latest.Error(); err == goinsta.ErrNoMore {
				break
			}
			dieOnError("unable to retrieve user feed: %s", err)
		}
	}

done:
	rss, err := feed.ToRss()
	if err != nil {
		dieOnError("unable to render RSS feed: %s", err)
	}

	fmt.Println(rss)
}

func formatFeedContent(item *goinsta.Item) (*feeds.Item, error) {
	shortDesc, _ := goutils.Abbreviate(item.Caption.Text, 50)

	content := fmt.Sprintf("<p>%s</p>", item.Caption.Text)

	if len(item.Images.Versions) > 0 {
		content += fmt.Sprintf("<p><img src=%q></p>", item.Images.Versions[0].URL)
	}

	if len(item.CarouselMedia) > 0 {
		for _, i := range item.CarouselMedia {
			content += fmt.Sprintf("<p><img src=%q></p>", i.Images.Versions[0].URL)
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
