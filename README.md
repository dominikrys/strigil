# IMDB Web Crawler

[![Build Status](https://img.shields.io/github/workflow/status/dominikrys/web-crawler/Continuous%20Integration?style=flat-square)](https://github.com/dominikrys/web-crawler/actions)

Web crawler for fetching information of people born on a specified day from [IMDB](https://www.imdb.com/), which is written to a MongoDB database. Information from the most popular x profiles if fetched. The crawler part is based off [Michael Okoko's blog post](https://blog.logrocket.com/web-scraping-with-go-and-colly/).

Note that there is rate limiting in place as the client may be blocked if too many requests are sent.

The aim of this project was to learn about Go and web scraping/crawling.

## Build and Run Instructions

Make sure [Go](https://golang.org/) is installed.

To compile, run:

```bash
go build ./crawler.go
```

Before running the program, run a [MongoDB](https://www.mongodb.com/) instance on port `27017`. This can be easily done using [Docker](https://www.docker.com/):

```bash
docker run --name mongo -p 27017:27017 -d mongo:4.4.6
```

Note that if MongoDB is not running the crawler will still work, but writing to MongoDB will be disabled.

Then, run the crawler:

```bash
./crawler.go --day <day of birthday> --month <month of birthday> [--profileNo <number of profiles to fetch>] [--mongoUri <MongoDB URI>]
```

Alternatively, for development, `go run` can be used:

```bash
go run . --day <day of birthday> --month <month of birthday>
```

To get more help on how to run the program and to check the program defaults, run:

```bash
./crawler --help
```

## TODO

- Specify MongoDB collection and database + mention in README
- Merge db branch
- refactor stuff so less comments are needed
- add a test?
