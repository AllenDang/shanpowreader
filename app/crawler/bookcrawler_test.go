package crawler

import (
  "log"
  "testing"
)

// 章节 http://www.hao662.com/haoshu/0/168/6380396.html
// 章节目录 http://www.hao662.com/haoshu/0/168
func TestHSWBookContentsCrawl(t *testing.T) {
  chapters, err := HSWBookContentsCrawl("http://www.hao662.com/haoshu/1/1426/")
  if err != nil {
    t.Fatal(err)
  }

  for _, c := range chapters {
    log.Println(c)
  }
}

func TestHSWChapterContentCrawl(t *testing.T) {
  content, err := HSWChapterContentCrawl("http://www.hao662.com/haoshu/0/168/53683.html")
  if err != nil {
    t.Fatal(err)
  }

  log.Println(content)
}
