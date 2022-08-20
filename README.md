# SaladeScraper

This is a very cool web scraper that I started writing in go, then finished in python

It is simple and for now only supports a single url and no custom options;
the number of posts etched is hardcoded as 25, but you can ask it for more if you want

# Usage

```bash
# use the python version
cd python/
python3 ./saladescraper.py `url`
# OR, use the go version
cd go/
go build
./saladescraper `url`
```

# Future
- For now, the output is still html. I think it would be cool if it outputted markdown
- include title
- make big bucks
