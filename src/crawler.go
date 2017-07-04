package smeago

import (
	"log"
	"net/http"
)

type Job struct {
	ID         int
	Path       string
	Links      []string
	Completed  bool
	RetryCount int
}

// NewJob creates a new Job
func NewJob(id int, path string) *Job {
	return &Job{
		ID:   id,
		Path: path,
	}
}

type Crawler struct {
	Domain  string
	Results chan Job
	Retries chan Job
}

// NewCrawler creates a crawler for the given domain
func NewCrawler(d string) *Crawler {
	c := &Crawler{}
	c.Domain = d
	c.Results = make(chan Job)
	c.Retries = make(chan Job)
	return c
}

// Crawl the job path and retries in case of failures
func (c *Crawler) Crawl(j Job) {
	link := c.Domain + j.Path

	if j.RetryCount > 0 {
		log.Printf("Retrying (%d): %s\n", j.RetryCount, link)
	} else {
		log.Println("Visiting:", link)
	}

	resp, err := http.Get(link)
	if err != nil {
		log.Println(err)
		c.Retries <- j
		return
	}
	defer resp.Body.Close()

	r, _ := ReadStringSize(resp.Body, int(resp.ContentLength))
	j.Links = r.Links
	j.Completed = true
	c.Results <- j
}
