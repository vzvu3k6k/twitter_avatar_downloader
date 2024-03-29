package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func getAvatarURL(ctx context.Context, twitterId string) (string, error) {
	url := "https://twitter.com/" + twitterId

	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		return "", fmt.Errorf("could not navigate profile page: %w", err)
	}

	var src string
	var ok bool
	sel := `img[src^="https://pbs.twimg.com/profile_images/"], div[data-testid="emptyState"]`
	if err := chromedp.Run(ctx,
		chromedp.AttributeValue(sel, "src", &src, &ok, chromedp.ByQuery),
	); err != nil {
		return "", fmt.Errorf("could not found profile icon: %w", err)
	}
	if !ok {
		return "", errors.New("user does not exist")
	}

	src = strings.TrimSpace(src)
	return deleteStrings(src, "_200x200", "_400x400"), nil
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
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx, cancel := chromedp.NewContext(timeoutCtx)
	defer cancel()

	twitterIds := os.Args[1:]
	for _, id := range twitterIds {
		log.Printf("get %s\n", id)

		url, err := getAvatarURL(ctx, id)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)

		filename := id + filepath.Ext(url)
		downloadFile(filename, url)
		time.Sleep(1 * time.Second)
	}
}
