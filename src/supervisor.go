package smeago

import "time"

const (
	maxBufferSize    = 5
	bufferingTimeout = 200 * time.Millisecond
	retryTimeout     = 1 * time.Second
)

// CrawlerSupervisor control the execution of the crawler
type CrawlerSupervisor struct {
	jobID   int
	pending map[int]int
	visited map[string]int
	buff    []*Job
	crawler *Crawler
}

// NewCrawlerSupervisor returns a new CrawlerSupervisor
func NewCrawlerSupervisor(c *Crawler) *CrawlerSupervisor {
	return &CrawlerSupervisor{
		jobID:   1,
		pending: make(map[int]int),
		visited: make(map[string]int),
		buff:    make([]*Job, 0),
		crawler: c,
	}
}

// BuffSize returns the len of the buffer
func (cs *CrawlerSupervisor) BuffSize() int {
	return len(cs.buff)
}

// HasPending returns true if there are jobs in the pending list
func (cs *CrawlerSupervisor) HasPending() bool {
	return len(cs.pending) > 0
}

//CompleteJob Removes the job from pending list
func (cs *CrawlerSupervisor) CompleteJob(j Job) {
	delete(cs.pending, j.ID)
}

// GetVisitedLinks returns a set of all visited links
func (cs *CrawlerSupervisor) GetVisitedLinks() []string {
	links := make([]string, 0)
	for k := range cs.visited {
		links = append(links, k)
	}
	return links
}

// AddJobToBuffer creates a new job for the given path and adds it to the buffer
func (cs *CrawlerSupervisor) AddJobToBuffer(path string) {
	_, ok := cs.visited[path]
	if ok {
		// already visited
		cs.visited[path]++
		return
	}

	cs.visited[path] = 1
	cs.pending[cs.jobID] = 1
	j := NewJob(cs.jobID, path)
	cs.jobID++
	cs.buff = append(cs.buff, j)
}

// CrawlJobs crawls all jobs in the buffer concurrently
func (cs *CrawlerSupervisor) CrawlJobs() {
	for _, job := range cs.buff {
		go cs.crawler.Crawl(*job)
	}
	// Clear the buffer
	cs.buff = cs.buff[:0]
}

// Start crawls buffered jobs until pending list is empty
func (cs *CrawlerSupervisor) Start(done chan bool) {
	go func() {
		for {
			select {
			case j := <-cs.crawler.Results:
				cs.CompleteJob(j)

				for _, l := range j.Links {
					cs.AddJobToBuffer(l)
				}

				// Circular Buffer
				if cs.BuffSize() >= maxBufferSize {
					cs.CrawlJobs()
				}

				if !cs.HasPending() {
					done <- true
					close(cs.crawler.Results)
					close(cs.crawler.Retries)
					return
				}
			case j := <-cs.crawler.Retries:
				go func() {
					time.Sleep(retryTimeout)
					j.RetryCount++
					cs.crawler.Crawl(j)
				}()
			case <-time.After(bufferingTimeout):
				cs.CrawlJobs()
			}
		}
	}()
}
