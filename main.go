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

const jsSnippet = `(function() {
	var body = document.querySelector("body");
	var div = document.createElement("div");
	div.id="div2";
	div.innerText = "div2";
	%s;
	return {};
})()
`

func insert(f string) chromedp.Action {
	var res interface{}
	return chromedp.Evaluate(fmt.Sprintf(jsSnippet, f), &res)
}

func assertXPath(sel, needle string) chromedp.Action {
	return chromedp.QueryAfter(sel, func(ctx context.Context, nodes ...*cdp.Node) error {
		if len(nodes) < 1 {
			return errors.New("nodes not found")
		}
		fullPath := nodes[0].FullXPath()
		if fullPath != needle {
			return fmt.Errorf("fullPath mismatch: got %s, want %s", fullPath, needle)
		}
		return nil
	}, chromedp.ByQuery)
}

func main() {
	absPath, err := filepath.Abs("template.html")
	if err != nil {
		panic(err)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Tasks{
		// Test 1
		chromedp.Navigate("file://" + absPath),
		assertXPath("#div1", "/html[1]/body[1]/div[1]"), // this line is required to build DOM before js.
		insert("body.appendChild(div)"),
		assertXPath("#div1", "/html[1]/body[1]/div[1]"),
		assertXPath("#div2", "/html[1]/body[1]/div[2]"),
		// Test 2
		chromedp.Navigate("file://" + absPath),
		assertXPath("#div1", "/html[1]/body[1]/div[1]"), // this line is required to build DOM before js.
		insert("body.insertBefore(div,body.childNodes[0])"),
		assertXPath("#div1", "/html[1]/body[1]/div[2]"), // this check will fail, without fix
		assertXPath("#div2", "/html[1]/body[1]/div[1]"),
	}); err != nil {
		log.Fatal(err)
	}
}
