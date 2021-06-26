package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	// Parse arguments
	month := flag.Int("month", 0, "Month to fetch birthdays for")
	day := flag.Int("day", 0, "Day to fetch birthdays for")
	profileNo := flag.Int("profileNo", 5, "(Optional) Amount of profiles to fetch")
	mongoUri := flag.String("mongoUri", "mongodb://localhost:27017", "(Optional) MongoDB URI")
	flag.Parse()

	if *month == 0 || *day == 0 {
		fmt.Println("Not enough arguments provided. Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Fetching %d profiles for people born on Day: %d, Month: %d\n", *profileNo, *day, *month)

	// Connect to MongoDB database
	usingMongo := true
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(*mongoUri))
	if err != nil {
		log.Println("MongoDB support disabled - couldn't connect:", err)
		usingMongo = false
	} else {
		defer func() {
			if err = mongoClient.Disconnect(ctx); err != nil {
				panic(err)
			}
		}()

		// Check if the database is connected
		err = mongoClient.Ping(ctx, readpref.Primary())
		if err != nil {
			log.Println("MongoDB support disabled - couldn't ping database:", err)
			usingMongo = false
		} else {
			fmt.Println("Successfully connected to MongoDB on", *mongoUri)
		}
	}

	crawl(*month, *day, *profileNo, *mongoClient, usingMongo)
}

func crawl(month int, day int, profileNo int, mongoClient mongo.Client, usingMongo bool) {
	profilesCrawled := 0

	// infoCollector panics when enough profiles have been fetched.
	// We recover from the panic here to stop crawling.
	// More info: https://github.com/gocolly/colly/issues/109#issuecomment-506995291
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Successfully crawled %d profiles. Exiting.\n", profilesCrawled)
		}
	}()

	// Set up the collector for the birthday list
	mainCollector := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)
	mainCollector.Limit(&colly.LimitRule{
		Delay:       50 * time.Millisecond,
		RandomDelay: 25 * time.Millisecond,
	})

	// Create a copy of the collector to initialise later
	profileCollector := mainCollector.Clone()

	mainCollector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting birthday page:", request.URL.String())
	})
	mainCollector.OnHTML(".mode-detail", func(element *colly.HTMLElement) {
		profileHref := element.ChildAttr("div.lister-item-image > a", "href")
		profileUrl := element.Request.AbsoluteURL(profileHref)
		profileCollector.Visit(profileUrl)
	})
	mainCollector.OnHTML("a.lister-page-next", func(element *colly.HTMLElement) {
		nextPage := element.Request.AbsoluteURL(element.Attr("href"))
		mainCollector.Visit(nextPage)
	})

	// Set up the collector for the profiles
	profileCollector.OnRequest(func(request *colly.Request) {
		fmt.Printf("Fetching profile %d/%d: %v\n", profilesCrawled+1, profileNo, request.URL.String())
	})
	profileCollector.OnHTML("#content-2-wide", func(element *colly.HTMLElement) {
		defer func() {
			// Check if enough profiles have been fetched
			profilesCrawled++
			if profilesCrawled >= profileNo {
				panic("Exit")
			}
		}()

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

		// Print JSON to console
		profileJson, err := json.MarshalIndent(profile, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(profileJson))

		if !usingMongo {
			return
		}

		// Create BSON for MongoDB
		profileBson, err := bson.Marshal(profile)
		if err != nil {
			log.Fatal(err)
		}

		// Write profile to database
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		collection := mongoClient.Database("crawler").Collection("profiles")
		countRes, err := collection.CountDocuments(ctx, profileBson)
		if err != nil {
			log.Fatal(err)
		} else if countRes > 0 {
			fmt.Println("Profile already exists in the database, skipping")
		} else {
			insertRes, err := collection.InsertOne(ctx, profileBson)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Wrote profile to database: %v\n", insertRes.InsertedID)
		}
	})

	// Run the crawler
	birthdayUrl := fmt.Sprintf("https://www.imdb.com/search/name/?birth_monthday=%d-%d", month, day)
	mainCollector.Visit(birthdayUrl)
}
