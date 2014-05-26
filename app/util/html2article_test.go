package util

import (
  "log"
  "testing"
)

func TestHtml2Article(t *testing.T) {
  htmlStr, _, _, err := GetHtmlFromUrl("http://www.piaotian.net/html/4/4940/3257334.html", "gbk") // "http://www.78xs.com/article/13/19226/7565228.shtml", "gbk")
  if err != nil {
    t.Fatal(err)
  }

  contentWithTags := Html2Article(htmlStr)
  if len(contentWithTags) == 0 {
    t.Fatal("没有找到正文")
  }

  log.Println(contentWithTags)
}
