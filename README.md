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

Then, to run the program:

```bash
./crawler.go --day <day of birthday> --month <month of birthday> [--profileNo <number of profiles to fetch>]
```

Alternatively, for development, `go run` can be used:

```bash
go run . --day <day of birthday> --month <month of birthday> [--profileNo <number of profiles to fetch>]
```

To get more help on how to run the program, run:

```bash
./crawler --help
```

## TODO

- Mention how to run MongoDB instance
- Tweak CI job to make this work. May possibly need some tweaking?
  - Optional writing to mongodb? Specify mongodb address? Optional JSON output?
  - Mention what collection and database the data is written to
- Reflect repo headline with mongodb
- refactor stuff so no comments are needed
- add a test?
