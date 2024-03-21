package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	// Track visited URLs to avoid loops
	visited := make(map[string]bool)

	err := crawl(
		options.URL,
		options.Extensions,
		"", visited)
	if err != nil {
		log.Fatal(err)
	}
}

func showBanner() {
	fmt.Printf(banner, version)
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

	// Add https:// if not present, to avoid errors
	if !strings.HasPrefix(options.URL, "http") {
		options.URL = "https://" + options.URL
	}

	if extensions != "" {
		options.Extensions = strings.Split(extensions, ",")
	}

	return options
}

func crawl(url string, extensions []string, prefix string, visited map[string]bool) error {
	if visited[url] {
		return nil // Skip already visited URL
	}

	body, err := get(url)
	if err != nil {
		return fmt.Errorf("error getting index: %w", err)
	}

	urls, err := parseIndex(url, prefix, extensions, body)
	if err != nil {
		return fmt.Errorf("error parsing index: %w", err)
	}

	for i, u := range urls {
		visited[u] = true
		if i == len(urls)-1 {
			fmt.Println(prefix + "└── " + u)
			if strings.HasSuffix(u, "/") {
				err = crawl(u, extensions, prefix+"    ", visited)
				if err != nil {
					log.Print(err)
				}
			}
		} else {
			fmt.Println(prefix + "├── " + u)
			if strings.HasSuffix(u, "/") {
				err = crawl(u, extensions, prefix+"│   ", visited)
				if err != nil {
					log.Print(err)
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
		if len(s.Text()) > 0 && s.Text()[len(s.Text())-1:] != "/" {
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
