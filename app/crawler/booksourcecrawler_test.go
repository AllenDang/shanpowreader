package crawler

import (
  "github.com/AllenDang/shanpowreader/app/models"
  "log"
  "testing"
)

func TestSoDuSearch(t *testing.T) {
  var context models.SearchCrawlContext

  bookTitles := []string{"裁决", "道士下山"}

  sodu := SoDuSearch{}

  for _, title := range bookTitles {
    context.BookTitle = title
    _, err := sodu.Search(&context)
    if title == "裁决" && err != nil {
      t.Fatal(err)
    }
  }
}

func TestSoDuBookSourcesCrawl(t *testing.T) {

  var context models.SearchCrawlContext
  context.BookTitle = "裁决"

  sodu := SoDuSearch{}

  listUrl, err := sodu.Search(&context)
  if err != nil {
    t.Fatal(err)
  }

  sources, _, err := sodu.Crawl(listUrl, &context)
  if err != nil {
    t.Fatal(err)
  }

  for _, s := range sources {
    log.Println(s)
  }
}

func TestEASOUSearch(t *testing.T) {
  var context models.SearchCrawlContext
  context.BookTitle = "仙路争锋"
  context.BookAuthor = "缘分0"

  easou := EASOUSearch{}
  _, err := easou.Search(&context)
  if err != nil {
    t.Fatal(err)
  }
}

func TestEASOUBookSourceCrawl(t *testing.T) {
  var context models.SearchCrawlContext
  context.BookTitle = "仙路争锋"
  context.BookAuthor = "缘分0"

  easou := EASOUSearch{}
  listUrl, err := easou.Search(&context)
  if err != nil {
    t.Fatal(err)
  }

  var bookSources []models.BookSource

  sources, nextPageUrl, err := easou.Crawl(listUrl, &context)
  if err != nil {
    t.Fatal(err)
  }

  bookSources = append(bookSources, sources...)

  for nextPageUrl != "" && len(bookSources) < BookSourcesLimitCount {
    log.Println(nextPageUrl)
    sources, nextPageUrl, err = easou.Crawl(nextPageUrl, &context)
    if err != nil {
      t.Fatal(err)
    }

    if len(sources) == 0 {
      break
    }

    bookSources = append(bookSources, sources...)
  }

  for _, s := range bookSources {
    log.Println(s)
  }
}
