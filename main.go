package main

import (
	"flag"
	"github.com/dbalduini/smeago/src"
	"log"
	"os"
	"path"
	"sort"
	"time"
)

var (
	host      string
	port      string
	urlLoc    string
	outputDir string
)

func main() {
	flag.Parse()

	start := time.Now()

	origin := host
	if port != "80" {
		origin += ":" + port
	}
	if urlLoc == "" {
		urlLoc = origin
	}

	s := &smeago.Sitemap{}
	s.Filename = path.Join(outputDir, "sitemap.xml")
	s.Path = urlLoc

	log.Println("Crawling Host:", origin)
	log.Println("Urlset Loc:", s.Path)
	log.Println("Sitemap File:", s.Filename)

	// Start crawling on the home page
	c := smeago.NewCrawler(origin)
	cs := smeago.NewCrawlerSupervisor(c)
	cs.AddJobToBuffer("/")

	// Block main until the crawler is done
	done := make(chan bool, 1)
	cs.Start(done)
	<-done
	close(done)

	s.Links = cs.GetVisitedLinks()
	sort.Strings(s.Links)
	if err := s.WriteToFile(); err != nil {
		log.Println(err)
	}
	log.Println("Finished in", time.Since(start))
}

func init() {
	wordDir, err := os.Getwd()
	if err != nil {
		log.Println(err)
		wordDir = ""
	}

	flag.StringVar(&host, "h", "http://localhost", "the host name")
	flag.StringVar(&port, "p", "80", "the host port")
	flag.StringVar(&urlLoc, "loc", "", "the prefix of sitemap loc tags")
	flag.StringVar(&outputDir, "o", wordDir, "the sitemap output dir")
}
