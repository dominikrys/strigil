# Web Crawler

Web crawler for getting information on people born on a specified day from [IMDB](https://www.imdb.com/). The crawler returns information from the profiles of the people born on a specified date, where the people are sorted by popularity. Based off [Michael Okoko's blog post](https://blog.logrocket.com/web-scraping-with-go-and-colly/).

The aim of this project was to learn about Go and web scraping/crawling.

## Build and Run Instructions

Make sure [Go](https://golang.org/) is installed.

To compile, run:

```bash
go build ./crawler.go
```

Then, to run the program:

```bash
./crawler.go --day <day of birthday> --month <month of birthday>
```

Alternatively, for development, `go run` can be used:

```bash
go run . --day <day of birthday> --month <month of birthday>
```

To get more help on how to run the program, run:

```bash
./crawler --help
```

## TODO

- build status - CI/CD
- You can as well go ahead and append it to an array of celebrities or store it in a database
- add rate limiting: http://go-colly.org/docs/examples/rate_limit/
- go through todos
- Make repo public
