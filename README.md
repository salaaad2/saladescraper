# SaladeScraper

This is a very cool web scraper that I started writing in go, then finished in python

It is simple and for now only supports a single url and no custom options;
the number of posts etched is hardcoded as 25, but you can ask it for more if you want

# Usage

```bash
# use the python version
cd python/
python3 ./saladescraper.py `url` [--limit n] [--offset n]
# OR, use the go version
cd go/
go build
./saladescraper `url` 
```

The `limit` argument determines the number of articles to scrape
The `offset` argument determines the post index to start scraping from
It is useful because the likelihood of running into a `too many requests` error is likely.
It is possible to see how many posts you got using your previous request, and to start another with that number as your offset argument

