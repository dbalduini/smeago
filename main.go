package main

import (
	"flag"
	"github.com/dbalduini/smeago/src"
	"log"
	"os"
	"path"
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
	origin := host + ":" + port

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
	if err := s.WriteToFile(true); err != nil {
		log.Println(err)
	}
	log.Println("Finished in", time.Since(start))
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
