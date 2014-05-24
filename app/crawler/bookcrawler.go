package crawler

import (
  "code.google.com/p/go.net/html"
  "github.com/AllenDang/shanpowreader/app/models"
  "github.com/AllenDang/shanpowreader/app/util"
  "github.com/PuerkitoBio/goquery"
  iconv "github.com/djimenez/iconv-go"
  "io/ioutil"
  "net/http"
  "regexp"
  "strings"
)

// 好书网 www.hao662.com
// 抓取章节列表
// refUrl 可以为 章节url 或 章节目录 url
// 章节 http://www.hao662.com/haoshu/0/168/6380396.html
// 章节目录 http://www.hao662.com/haoshu/0/168
func HSWBookContentsCrawl(refUrl string) ([]models.Chapter, error) {
  contentsRegex := regexp.MustCompile(`http://www.hao662.com/haoshu/\d+/\d+`)

  contentsUrl := contentsRegex.FindString(refUrl)
  if contentsUrl == "" {
    return nil, util.ErrCanNotConstructBookContentsUrl
  }

  resp, err := http.Get(contentsUrl)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  htmlStr, err := iconv.ConvertString(string(body), "gbk", "utf-8")
  if err != nil {
    return nil, err
  }

  // 章节链接
  // <a[^>]+["'](\d+\.html)['"]>([^<]+)</a>
  rx := regexp.MustCompile(`<a[^>]+["'](\d+\.html)['"]>([^<]+)</a>`)
  contents := rx.FindAllStringSubmatch(htmlStr, -1)

  var chapters []models.Chapter
  var index uint
  for _, c := range contents {
    if len(c) < 3 || strings.TrimSpace(c[1]) == "" || strings.TrimSpace(c[2]) == "" {
      continue
    }

    chapters = append(chapters, models.Chapter{
      Index: index,
      Url:   contentsUrl + "/" + c[1],
      Name:  c[2],
    })

    index += 1
  }

  return chapters, nil
}

func HSWChapterContentCrawl(chapterUrl string) (string, error) {

  resp, err := http.Get(chapterUrl)
  if err != nil {
    return "", nil
  }
  defer resp.Body.Close()

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }

  htmlStr, err := iconv.ConvertString(string(body), "gbk", "utf-8")
  if err != nil {
    return "", err
  }

  node, err := html.Parse(strings.NewReader(htmlStr))
  if err != nil {
    return "", err
  }

  doc := goquery.NewDocumentFromNode(node)

  return doc.Find("#content").Html()
}
