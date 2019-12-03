package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"net/http"
	"os"
	"strings"
	// "time"
)

var TotalBooks int

func main() {
	c := colly.NewCollector(colly.AllowedDomains("mirrors.pglaf.org"))
	initLinks := []string{"http://mirrors.pglaf.org/gutenberg/0/",
		"http://mirrors.pglaf.org/gutenberg/1/",
		"http://mirrors.pglaf.org/gutenberg/2/",
		"http://mirrors.pglaf.org/gutenberg/3/",
		"http://mirrors.pglaf.org/gutenberg/4/",
		"http://mirrors.pglaf.org/gutenberg/5/",
		"http://mirrors.pglaf.org/gutenberg/6/",
		"http://mirrors.pglaf.org/gutenberg/7/",
		"http://mirrors.pglaf.org/gutenberg/8/",
		"http://mirrors.pglaf.org/gutenberg/9/",
		"http://mirrors.pglaf.org/gutenberg/0/1/",
	}

	c.Limit(&colly.LimitRule{
		// Filter domains affected by this rule
		DomainRegexp: "http://mirrors.pglaf.org/gutenberg/*",
		// Set a delay between requests to these domains
		// Delay: 1 * time.Second,
		// Add an additional random delay
		// RandomDelay: 1 * time.Second,
	})

	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Printf("visited %s\n", r.Request.URL)
	// })

	ignoreFiles := []string{".htm", ".png", ".zip", "jpg", ".jpeg", ".lit", ".prc"}
	ignoreTexts := []string{"Name", "Last modified", "Parent Directory", "Size", "Description"}

	c.OnHTML("a", func(e *colly.HTMLElement) {
		// fmt.Printf("text: %v", e.Text)
		if strings.Contains(e.Text, ".txt") {
			filesplit := strings.Split(e.Text, "/")
			filename := "harvests/" + filesplit[len(filesplit)-1]
			DownloadFile(filename, e.Request.AbsoluteURL(e.Attr("href")))

		} else if !(IgnoreTexts(e, ignoreTexts) || IgnoreFiles(e, ignoreFiles)) {
			e.Request.Visit(e.Attr("href"))
		}
	})

	for i := 0; i < len(initLinks); i++ {
		c.Visit(initLinks[i])
	}
}

// IgnoreTexts true if text of element is excluded
func IgnoreTexts(e *colly.HTMLElement, texts []string) bool {
	ignore := false
	for i := 0; i < len(texts); i++ {
		if strings.Contains(e.Text, texts[i]) {
			ignore = true
		}
	}
	return ignore
}

// IgnoreFiles true if text of element is excluded
func IgnoreFiles(e *colly.HTMLElement, files []string) bool {
	ignore := false
	for i := 0; i < len(files); i++ {
		if strings.Contains(e.Attr("href"), files[i]) {
			ignore = true
		}
	}
	return ignore
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
		// return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	TotalBooks++
	fmt.Printf("Total Books download: %d\n", TotalBooks)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
