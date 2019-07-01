package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	absPath, err := filepath.Abs("template.html")
	if err != nil {
		panic(err)
	}
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate("file://" + absPath),
		chromedp.Query("a.toggle-nav", chromedp.ByQuery),
		chromedp.QueryAfter("a[href='#mm-2']", func(ctx context.Context, nodes ...*cdp.Node) error {
			fullPath := nodes[0].FullXPath()
			fmt.Println("[[FullXPath]]", fullPath)
			if fullPath != "/html[1]/body[1]/div[1]/div[1]/div[1]/ul[1]/li[1]/a[1]" {
				return errors.New("wrong path")
			}
			return nil
		}, chromedp.ByQuery),
	}); err != nil {
		log.Fatal(err)
	}
}
