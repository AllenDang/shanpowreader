package crawler

import (
  "log"
  "regexp"
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

  chapterUrls := []string{
    "http://www.xiaoshuoan.com/73/73571/7848348.html",
    "http://www.hao662.com/haoshu/1/1113/511564.html",
    "http://www.78xs.com/article/215/11354/3732084.shtml",
    "http://www.173wx.com/xiaoshuo/8/8916/7574300.html",
    "http://www.du7.com/html/20/20982/5834907.html",
    "http://www.bingdi.cc/book/47_2136228.html",
    "http://www.epzw.com/files/article/html/89/89282/10085689.html",
    "http://www.siluke.com/0/106/106745/15926288.html",
    "http://www.aszw.com/book/41/41989/9337024.html",
    "http://read.qidian.com/BookReader/3144241,53936102.aspx",
    "http://book108.com/a/25351/7351145.html",
  }

  tagsRegex := regexp.MustCompile(`<[^>]+>`)

  for _, chapterUrl := range chapterUrls {

    content, err := cm.ChapterContentCrawl(chapterUrl)
    if err != nil {
      t.Errorf("Chapter content crawl (%s) failed. Error: %s", chapterUrl, err.Error())
    }

    // 不能为空
    if len(content) == 0 {
      t.Errorf("Chapter content crawl (%s) failed. Error: %s", chapterUrl, "内容为空")
    }

    // 不能有 标签
    tag := tagsRegex.FindString(content)
    if tag != "" {
      t.Errorf("Chapter content crawl (%s) failed. Error: %s", chapterUrl, "还存在 Html 标签")
    }
  }
}
