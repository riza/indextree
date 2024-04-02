package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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
	version = `v1.0.2`

	treeBranch = `├── `
	treeEnd    = `└── `
	treeSuffix = `│   `
)

type (
	Options struct {
		URL           string
		Banner        bool
		Tree          bool
		ShowOnlyFiles bool

		Extensions []string
		Matchers   []string
	}
)

func main() {
	options := parseOptions()

	if options.Banner {
		showBanner()
	}

	wg := &sync.WaitGroup{}

	line := make(chan string)

	go func() {
		for l := range line {
			if !options.Tree {
				l = strings.ReplaceAll(l, treeBranch, "")
				l = strings.ReplaceAll(l, treeEnd, "")
				l = strings.ReplaceAll(l, treeSuffix, "")
				l = strings.ReplaceAll(l, "    ", "")
			}

			if options.ShowOnlyFiles && strings.HasSuffix(l, "/") {
				continue
			}

			fmt.Println(l)
		}
	}()

	err := crawl(
		line,
		wg,
		options.URL,
		options.Extensions,
		options.Matchers,
		"")

	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()
	time.Sleep(1 * time.Second)

	close(line)
}

func showBanner() {
	fmt.Printf(banner, version)
}

func parseOptions() *Options {
	options := &Options{}

	extensions := ""
	matchers := ""

	flag.StringVar(&options.URL, "u", "", "url to parse index")
	flag.StringVar(&extensions, "e", "", "extensions to filter, example: -e jpg,png,gif")
	flag.StringVar(&matchers, "m", "", "match in url, example: -mu admin,login")
	flag.BoolVar(&options.Banner, "b", true, "show banner")
	flag.BoolVar(&options.Tree, "t", true, "show tree")
	flag.BoolVar(&options.ShowOnlyFiles, "of", false, "show only files")

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

	if matchers != "" {
		options.Matchers = strings.Split(matchers, ",")
	}

	return options
}

func crawl(line chan string, wg *sync.WaitGroup, url string, extensions, matchers []string, prefix string) error {
	wg.Add(1)
	defer wg.Done()

	body, err := get(url)
	if err != nil {
		return fmt.Errorf("error getting index: %w", err)
	}
	defer body.Close()

	urls, err := parseIndex(body)
	if err != nil {
		return fmt.Errorf("error parsing index: %w", err)
	}

	for i, u := range urls {
		u = url + u

		suffix := treeSuffix
		if i == len(urls)-1 {
			suffix = "    "
		}

		//is dir
		if strings.HasSuffix(u, "/") {
			line <- prefix + treeBranch + u

			err = crawl(line, wg, u, extensions, matchers, prefix+suffix)
			if err != nil {
				log.Print(err)
			}
			continue
		}

		//is match with matchers
		if len(matchers) > 0 {
			for _, m := range matchers {
				if strings.Contains(strings.ToLower(u), strings.ToLower(m)) {
					line <- prefix + treeBranch + u
					continue
				}
			}
			continue
		}

		//is not dir
		if !strings.HasSuffix(u, "/") && len(extensions) > 0 {
			ext := filepath.Ext(u)
			if len(ext) > 0 {
				ext = ext[1:]
				for _, e := range extensions {
					if e == ext {
						line <- prefix + treeEnd + u
						continue
					}
				}
			}
			continue
		}

		line <- prefix + treeBranch + u
	}

	return nil
}

func get(url string) (io.ReadCloser, error) {
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error getting index: %w", err)
	}

	return resp.Body, nil
}

func parseIndex(body io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing index: %w", err)
	}

	urls := make([]string, 0)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		//skip parent directory
		href, _ := s.Attr("href")
		if strings.HasPrefix(href, "?") || href == "../" || href == "/" || strings.Contains(strings.ToLower(s.Text()), "parent") { // parent directory
			return
		}

		urls = append(urls, href)
	})

	return urls, nil
}
