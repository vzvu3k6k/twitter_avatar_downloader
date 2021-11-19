package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/chromedp/chromedp"
)

func getAvatarURL(ctx context.Context, twitterId string) string {
	url := "https://twitter.com/" + twitterId
	log.Printf("get %s\n", url)

	var src string
	var ok bool
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.AttributeValue(`img[src^="https://pbs.twimg.com/profile_images/"]`, "src", &src, &ok, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("node not found")
	}

	src = strings.TrimSpace(src)
	return deleteStrings(src, "_200x200", "_400x400")
}

func deleteStrings(s string, dels ...string) string {
	for _, d := range dels {
		s = strings.Replace(s, d, "", 1)
	}
	return s
}

func downloadFile(filename, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP status: %s", resp.Status)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	twitterIds := os.Args[1:]
	for _, id := range twitterIds {
		url := getAvatarURL(ctx, id)
		filename := id + filepath.Ext(url)
		downloadFile(filename, url)
	}
}
