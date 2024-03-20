package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	banner = `
    _           __          __               
   (_)___  ____/ /__  _  __/ /_________  ___ 
  / / __ \/ __  / _ \| |/_/ __/ ___/ _ \/ _ \
 / / / / / /_/ /  __/>  </ /_/ /  /  __/  __/
/_/_/ /_/\__,_/\___/_/|_|\__/_/   \___/\___/ %s

`
	version = `v1.0.0`
)

type (
	Options struct {
		URL    string
		Banner bool

		Extensions []string
		Matchers   []string
	}
)

func main() {
	options := parseOptions()

	if options.Banner {
		showBanner()
	}

	err := crawl(
		options.URL,
		options.Extensions,
		"")
	if err != nil {
		log.Error(err)
	}
}

func showBanner() {
	log.Printf(banner, version)
}

func parseOptions() *Options {
	options := &Options{}
	extensions := ""

	flag.StringVar(&options.URL, "u", "", "url to parse index")
	flag.StringVar(&extensions, "e", "", "extensions to filter, example: -e jpg,png,gif")
	flag.BoolVar(&options.Banner, "b", true, "show banner")
	flag.Parse()

	if options.URL == "" {
		flag.Usage()
		os.Exit(1)
	}

	if extensions != "" {
		options.Extensions = strings.Split(extensions, ",")
	}

	return options
}

func crawl(url string, extensions []string, prefix string) error {
	body, err := get(url)
	if err != nil {
		return fmt.Errorf("error getting index: %w", err)
	}

	urls, err := parseIndex(url, prefix, extensions, body)
	if err != nil {
		return fmt.Errorf("error parsing index: %w", err)
	}

	for i, u := range urls {
		if i == len(urls)-1 {
			fmt.Println(prefix + "└── " + u)
			if strings.HasSuffix(u, "/") {
				err = crawl(u, extensions, prefix+"    ")
				if err != nil {
					log.Error(err)
				}
			}
		} else {
			fmt.Println(prefix + "├── " + u)
			if strings.HasSuffix(u, "/") {
				err = crawl(u, extensions, prefix+"│   ")
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	return nil
}

func get(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting index: %w", err)
	}

	return resp.Body, nil
}

func parseIndex(url, prefix string, extensions []string, body io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing index: %w", err)
	}

	urls := make([]string, 0)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		//skip parent directory
		if strings.Contains(s.Text(), " Parent Directory") {
			return
		}

		if url[len(url)-1:] != "/" {
			url += "/"
		}

		href, _ := s.Attr("href")

		//handle files
		if s.Text()[len(s.Text())-1:] != "/" {
			//filter by extensions
			if len(extensions) > 0 {
				pos := strings.LastIndex(url+href, ".")
				if pos == -1 {
					return
				}

				ext := (url + href)[pos+1 : len(url+href)]
				if len(ext) > 0 {
					for _, e := range extensions {
						if e == ext {
							fmt.Println(prefix + "└── " + url + href)
							return
						}
					}
				}
				return
			}

			fmt.Println(prefix + "└── " + url + href)
			return
		}

		urls = append(urls, url+href)
	})

	return urls, nil
}
