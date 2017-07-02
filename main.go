package main

import (
	"bufio"
	"flag"
	"github.com/dbalduini/smeago/src"
	"log"
	"os"
	"path"
	"sort"
	"time"
)

var (
	host        string
	port        string
	urlLoc      string
	outputDir   string
	sitemapFile string
)

func main() {
	flag.Parse()

	origin := host + ":" + port
	sitemapFile = path.Join(outputDir, "sitemap.xml")

	log.Println("Crawling Host:", origin)
	log.Println("Urlset Loc:", urlLoc)
	log.Println("Sitemap File:", sitemapFile)

	// Control Variables
	start := time.Now()
	done := make(chan bool, 1)
	pending := make(map[int]int)
	visited := make(map[string]int)
	stop := false

	// Start crawling on the home page
	crawler := smeago.NewCrawler(origin)
	jobID := 1
	pending[jobID] = 1
	visited["/"] = 1
	j := *smeago.NewJob(jobID, "/")
	go crawler.Crawl(j)

	for !stop {
		select {
		case j := <-crawler.Results:
			delete(pending, j.ID)

			for _, l := range j.Links {
				_, ok := visited[l]
				if !ok {
					// new link
					jobID++
					pending[jobID] = 1
					visited[l] = 1
					job := *smeago.NewJob(jobID, l)
					go crawler.Crawl(job)
				} else {
					// already visited
					visited[l]++
				}
			}

			if len(pending) == 0 {
				done <- true
			}
		case j := <-crawler.Retries:
			go func() {
				time.Sleep(time.Second * 3)
				j.IsRetry = true
				crawler.Crawl(j)
			}()
		case stop = <-done:
		}
	}

	err := writeSitemap(visited, true)
	if err != nil {
		log.Println(err)
	}
	log.Println("Finished in", time.Since(start))
}

func writeSitemap(visited map[string]int, sortLinks bool) error {
	links := make([]string, 0)
	for k := range visited {
		links = append(links, urlLoc+k)
	}

	if sortLinks {
		sort.Strings(links)
	}

	log.Println("Writing Sitemap:", len(links))
	f, err := os.Create(sitemapFile)
	if err != nil {
		return err
	}
	defer f.Close()

	wd := bufio.NewWriter(f)

	err = smeago.WriteSitemap(wd, links)
	if err != nil {
		return err
	}

	wd.Flush()
	return nil
}

func init() {
	const (
		defaultHost = "http://localhost"
		defaultPort = "8080"
	)

	wordDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		wordDir = ""
	}

	flag.StringVar(&host, "h", defaultHost, "the host name")
	flag.StringVar(&port, "p", defaultPort, "the host port")
	flag.StringVar(&urlLoc, "loc", defaultHost, "the prefix of sitemap loc tags")
	flag.StringVar(&outputDir, "o", wordDir, "the sitemap output dir")
}
