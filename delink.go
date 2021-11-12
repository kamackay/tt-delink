package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/alecthomas/kong"
	"github.com/atotto/clipboard"
	"github.com/dustin/go-humanize"
)

type Opts struct {
	NoCopy bool   `short:"n" help:"Don't add link to Clipboard"`
	Link   string `arg:"d" help:"Link" default:"."`
}

func main() {
	var opts Opts
	ctx := kong.Parse(&opts)

	start := time.Now()
	defer func() {
		if time.Now().Sub(start) > 100*time.Millisecond {
			fmt.Printf("Done in %s\n", humanize.RelTime(start, time.Now(), "", ""))
		}
	}()

	req, err := http.NewRequest("GET", opts.Link, nil)
	if err != nil {
		panic(err)
	}
	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		link := req.URL
		u, err := url.Parse(link.String())
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		u.RawQuery = ""
		result := u.String()
		fmt.Println(result)
		if !opts.NoCopy {
			fmt.Println("Adding to clipboard")
			clipboard.WriteAll(result)
		}
		return nil
	}

	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	ctx.Exit(0)
}
