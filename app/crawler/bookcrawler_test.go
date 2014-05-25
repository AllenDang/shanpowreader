package crawler

import (
  "log"
  "testing"
)

func TestNewCrawlerManager(t *testing.T) {

  var cm CrawlerManager

  cm.Init("../../conf/bookcrawler.conf")
}

// 章节 http://www.hao662.com/haoshu/0/168/6380396.html
// 章节目录 http://www.hao662.com/haoshu/0/168
func TestBookContentsCrawl(t *testing.T) {
  var cm CrawlerManager

  cm.Init("../../conf/bookcrawler.conf")

  chapters, err := cm.BookContentsCrawl("http://www.hao662.com/haoshu/1/1426/")
  if err != nil {
    t.Fatal(err)
  }

  for _, c := range chapters {
    log.Println(c)
  }
}

func TestChapterContentCrawl(t *testing.T) {
  var cm CrawlerManager

  cm.Init("../../conf/bookcrawler.conf")

  content, err := cm.ChapterContentCrawl("http://www.hao662.com/haoshu/0/168/53683.html")
  if err != nil {
    t.Fatal(err)
  }

  log.Println(content)
}
