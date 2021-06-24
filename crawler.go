package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Movie struct {
	Title string
	Year  string
}

type Profile struct {
	Name      string
	Photo     string
	JobTitle  string
	BirthDate string
	Bio       string
	TopMovies []Movie
}

func main() {
	month := flag.Int("month", 0, "Month to fetch birthdays for")
	day := flag.Int("day", 0, "Day to fetch birthdays for")
	flag.Parse()

	if *month == 0 || *day == 0 {
		fmt.Println("Not enough arguments provided. Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Fetching birthdays for Day: %d, Month: %d\n", *month, *day)

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
		profile := Profile{}
		profile.Name = element.ChildText("h1.header > span.itemprop")
		profile.Photo = element.ChildAttr("#name-poster", "src")
		profile.JobTitle = element.ChildText("#name-job-categories > a > span.itemprop")
		profile.BirthDate = element.ChildAttr("#name-born-info time", "datetime")
		profile.Bio = strings.TrimSpace(element.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))

		element.ForEach("div.knownfor-title", func(_ int, knownForElem *colly.HTMLElement) {
			movie := Movie{}
			movie.Title = knownForElem.ChildText("div.knownfor-title-role > a.knownfor-ellipsis")
			movie.Year = knownForElem.ChildText("div.knownfor-year > span.knownfor-ellipsis")

			profile.TopMovies = append(profile.TopMovies, movie)
		})

		profileJson, err := json.MarshalIndent(profile, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(profileJson))
	})

	birthdayUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	collector.Visit(birthdayUrl)

	// TODO: add toggle for async mode
	//collector.Wait()
}
