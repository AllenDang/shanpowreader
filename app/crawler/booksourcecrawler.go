package crawler

import (
  "fmt"
  "github.com/AllenDang/shanpowreader/app/models"
  "github.com/AllenDang/shanpowreader/app/util"
  "github.com/PuerkitoBio/goquery"
  iconv "github.com/djimenez/iconv-go"
  "net/http"
  "net/url"
  "regexp"
  "strings"
)

//
// 通过搜索引擎搜索书籍
// 得到书籍来源网站及最新章节
//
// 实现新 Crawler
// 1 实现 models.BookSourcesCrawler 接口
// 2 在 BookSourcesCrawl 中调用
//
func BookSourcesCrawl(name, bookTitle, bookAuthor string) ([]models.BookSource, error) {
  var c models.SearchCrawlContext

  c.BookTitle = bookTitle
  c.BookAuthor = bookAuthor

  switch name {
  case "sodu":
    sodu := new(SoDuSearch)
    return bookSourcesCrawl(sodu, &c)
  case "easou":
    easou := new(EASOUSearch)
    return bookSourcesCrawl(easou, &c)
  }

  return nil, util.ErrSearchEngineNotSurpport
}

func bookSourcesCrawl(crawler models.BookSourcesCrawler,
  c *models.SearchCrawlContext) ([]models.BookSource, error) {
  url, err := crawler.Search(c)
  if err != nil {
    return nil, err
  }

  return crawler.Crawl(url, c)
}

//
// SoDu 小说搜索
type SoDuSearch struct {
}

var SoDuClient = http.Client{
  CheckRedirect: redirectPolicyFunc,
}

var SoDuSearcher SoDuSearch

// 重定向处理
// 禁止自动重定向 以便于得到重定向网址
func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
  return util.ErrAutoRedirectForbidden
}

// 获取 SoDu 小说搜索结果列表页面链接
func (s *SoDuSearch) Search(sc *models.SearchCrawlContext) (string, error) {
  title, err := iconv.ConvertString(sc.BookTitle, "utf-8", "gbk")
  if err != nil {
    return "", err
  }

  searchUrl := fmt.Sprintf("http://www.sodu.so/search/index.aspx?key=%s", title)

  html, _, _, err := util.GetHtmlFromUrl(searchUrl, "gbk")
  if err != nil {
    return "", err
  }

  pattern := fmt.Sprintf(`<a href="([^"]+)[^>]+><b>%s`, sc.BookTitle)
  rx := regexp.MustCompile(pattern)

  matches := rx.FindStringSubmatch(html)
  if len(matches) < 2 || strings.TrimSpace(matches[1]) == "" {
    return "", util.ErrRegexCannotMatch
  }

  return fmt.Sprintf("http://www.sodu.so%s", matches[1]), nil
}

func (s *SoDuSearch) Crawl(sourcesUrl string, sc *models.SearchCrawlContext) ([]models.BookSource, error) {
  html, _, _, err := util.GetHtmlFromUrl(sourcesUrl, "gbk")
  if err != nil {
    return nil, err
  }

  pattern := `<div[^<]*<div[^<]*<a[^>]+>%s_[^>]+[^<]*</div[^<]*<div[^<]*<a[^<]*</a[^<]*</div[^<]*<div[^<]*</div[^<]*</div>`
  pattern = fmt.Sprintf(pattern, sc.BookTitle)

  rx := regexp.MustCompile(pattern)

  matches := rx.FindAllString(html, -1)

  var bookSources []models.BookSource
  existNameMap := map[string]int{}

  for _, m := range matches {
    chapterUrl, chapter, err := exactChapterAndUrl(m, sc.BookTitle)
    if err != nil {
      continue
    }

    updateTime, err := exactUpdateTime(m)
    if err != nil {
      continue
    }

    bs := models.BookSource{
      ChapterUrl: fmt.Sprintf("http://www.sodu.so%s", chapterUrl),
      Chapter:    chapter,
      UpdateTime: updateTime,
    }

    // 将链接替换为目录链接
    resp, err := SoDuClient.Get(bs.ChapterUrl)
    if err != nil && resp.StatusCode != 302 {
      continue
    }

    redirectUrl, err := resp.Location()
    if err != nil {
      continue
    }

    bs.ChapterUrl = redirectUrl.String()
    bs.Host = redirectUrl.Host

    // 去掉来源重复的
    if _, ok := existNameMap[bs.Host]; ok { // 搜索引擎已经按照更新时间排好序 相同网站只取排序靠前的
      continue
    } else {
      existNameMap[bs.Host] = 1
    }

    bookSources = append(bookSources, bs)
  }

  return bookSources, nil
}

// 章节 url 名称
// <a[^'"]+['"]([^'"]+)[^>]+>裁决_[^<]*(第[^>]+章[^>]+)</a>
func exactChapterAndUrl(s, bookTitle string) (string, string, error) {
  pattern := `<a[^'"]+['"]([^'"]+)[^>]+>%s_([^<]+)</a>`
  pattern = fmt.Sprintf(pattern, bookTitle)

  rx := regexp.MustCompile(pattern)

  matches := rx.FindStringSubmatch(s)
  if len(matches) < 3 || strings.TrimSpace(matches[1]) == "" || strings.TrimSpace(matches[2]) == "" {
    return "", "", util.ErrRegexCannotMatch
  }

  return matches[1], matches[2], nil
}

// 网站
// <a[^>]+class=["']tl["']>([^<]+)</a>
// func exactHostName(s string) (string, error) {
//   rx := regexp.MustCompile(`<a[^>]+class=["']tl["']>([^<]+)</a>`)

//   matches := rx.FindStringSubmatch(s)
//   if len(matches) < 2 || strings.TrimSpace(matches[1]) == "" {
//     return "", util.ErrRegexCannotMatch
//   }

//   return matches[1], nil
// }

// 更新时间
// \d{4}-\d{1,2}-\d{1,2}\s*\d{1,2}:\d{1,2}:\d{1,2}
func exactUpdateTime(s string) (string, error) {
  rx := regexp.MustCompile(`\d{4}-\d{1,2}-\d{1,2}\s*\d{1,2}:\d{1,2}:\d{1,2}`)

  match := rx.FindString(s)
  if match == "" {
    return "", util.ErrRegexCannotMatch
  }

  return match, nil
}

//
// easou 小说搜索
type EASOUSearch struct {
}

var EASOUSearcher EASOUSearch

func (s *EASOUSearch) Search(sc *models.SearchCrawlContext) (string, error) {

  searchUrl := fmt.Sprintf("http://book.easou.com/c/s.m?q=%s", sc.BookTitle)

  doc, err := goquery.NewDocument(searchUrl)
  if err != nil {
    return "", err
  }

  selection := doc.Find(".easou_box").EachWithBreak(func(n int, gs *goquery.Selection) bool {
    if gs.Find("p>a").Eq(0).Text() == sc.BookTitle &&
      gs.Find("p>a").Eq(1).Text() == sc.BookAuthor {
      return false
    }

    return true
  })

  sourceUrl, isExist := selection.Find(".link_green").Attr("href")
  if isExist {
    sourceUrl = "http://book.easou.com" + sourceUrl
  } else {
    return "", util.ErrCanNotFindBook
  }

  return sourceUrl, nil
}

func (s *EASOUSearch) Crawl(sourceUrl string, sc *models.SearchCrawlContext) ([]models.BookSource, error) {
  doc, err := goquery.NewDocument(sourceUrl)
  if err != nil {
    return nil, err
  }

  var bookSources []models.BookSource

  foreach := func(n int, gs *goquery.Selection) {
    cs := gs.Find("p>a").Eq(0)
    chapter := cs.Text()
    chapterUrl, isExist := cs.Attr("href")
    if !isExist {
      return
    }

    u, err := url.Parse(chapterUrl)
    if err != nil {
      return
    }

    values, err := url.ParseQuery(u.RawQuery)
    if err != nil {
      return
    }

    if cu, ok := values["cu"]; ok {
      chapterUrl = cu[0]
    } else {
      return
    }

    // 得到 host
    u, err = url.Parse(chapterUrl)
    if err != nil {
      return
    }

    updateTime, err := util.ExtendTimeLayoutWithYear(gs.Find("p>span").Eq(0).Text())
    if err != nil {
      return
    }

    bookSources = append(bookSources, models.BookSource{
      Host:       u.Host,
      ChapterUrl: chapterUrl,
      Chapter:    chapter,
      UpdateTime: updateTime,
    })
  }

  doc.Find(".easou_pdb4").Each(foreach)

  return bookSources, nil
}
