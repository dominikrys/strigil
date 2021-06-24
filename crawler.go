package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type Movie struct {
	Title string
	Year  string
}

type Celebrity struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []Movie
}

func main() {
	// TODO: add something saying which args are chosen
	month := flag.Int("month", 1, "Month to fetch birthdays for")
	day := flag.Int("day", 1, "Day to fetch birthdays for")
	flag.Parse()

	crawl(*month, *day)
}

func crawl(month int, day int) {
	// TODO: add result limit

	collector := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
		//colly.Async(true),
	)

	infoCollector := collector.Clone()

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting: ", request.URL.String())
	})

	infoCollector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting Profile URL: ", request.URL.String())
	})

	collector.OnHTML("a.lister-page-next", func(element *colly.HTMLElement) {
		nextPage := element.Request.AbsoluteURL(element.Attr("href"))
		collector.Visit(nextPage)
	})

	collector.OnHTML(".mode-detail", func(element *colly.HTMLElement) {
		profileHref := element.ChildAttr("div.lister-item-image > a", "href")
		profileUrl := element.Request.AbsoluteURL(profileHref)
		infoCollector.Visit(profileUrl)
	})

	infoCollector.OnHTML("#content-2-wide", func(element *colly.HTMLElement) {
		celeb := Celebrity{}
		celeb.Name = element.ChildText("h1.header > span.itemprop")
		celeb.Photo = element.ChildAttr("#name-poster", "src")
		celeb.JobTitle = element.ChildText("#name-job-categories > a > span.itemprop")
		celeb.BirthDate = element.ChildAttr("#name-born-info time", "datetime")
		celeb.Bio = strings.TrimSpace(element.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))

		element.ForEach("div.knownfor-title", func(_ int, knownForElem *colly.HTMLElement) {
			movie := Movie{}
			movie.Title = knownForElem.ChildText("div.knownfor-title-role > a.knownfor-ellipsis")
			movie.Year = knownForElem.ChildText("div.knownfor-year > span.knownfor-ellipsis")

			celeb.TopMovies = append(celeb.TopMovies, movie)
		})

		celebJson, err := json.MarshalIndent(celeb, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(celebJson))
	})

	birthdayUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	collector.Visit(birthdayUrl)

	// TODO: add toggle for async mode
	//collector.Wait()
}
