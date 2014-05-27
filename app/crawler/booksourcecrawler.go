package crawler

import (
  "code.google.com/p/go.net/html"
  "fmt"
  "github.com/AllenDang/shanpowreader/app/models"
  "github.com/AllenDang/shanpowreader/app/util"
  "github.com/PuerkitoBio/goquery"
  iconv "github.com/djimenez/iconv-go"
  "net/http"
  "net/url"
  "regexp"
  "strings"
  "time"
)

var (
  BookSourcesLimitCount = 15
)

//
// 通过搜索引擎搜索书籍
// 得到书籍来源网站及最新章节
//
// 实现新 Crawler
// 1 实现 models.BookSourcesCrawler 接口
// 2 在 BookSourcesCrawl 中调用
//
func BookSourcesCrawl(name, bookTitle, bookAuthor string,
  isCrawlableFunc func(string) bool) ([]models.BookSource, error) {
  var c models.SearchCrawlContext

  c.BookTitle = bookTitle
  c.BookAuthor = bookAuthor
  c.IsCrawlableFunc = isCrawlableFunc
  c.HostExists = map[string]int{}

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

func check2AppendBookSources(sources []models.BookSource, c *models.SearchCrawlContext) (bookSources []models.BookSource) {
  if c.IsCrawlableFunc != nil {
    for _, s := range sources {
      // 去掉来源重复的
      if _, ok := c.HostExists[s.Host]; ok { // 搜索引擎已经按照更新时间排好序 相同网站只取排序靠前的
        continue
      } else {
        c.HostExists[s.Host] = 1
      }

      if c.IsCrawlableFunc(s.Host) {
        bookSources = append(bookSources, s)
      }
    }
  }

  return
}

func bookSourcesCrawl(crawler models.BookSourcesCrawler,
  c *models.SearchCrawlContext) ([]models.BookSource, error) {
  url, err := crawler.Search(c)
  if err != nil {
    return nil, err
  }

  var bookSources []models.BookSource

  sources, nextPageUrl, err := crawler.Crawl(url, c)
  if err != nil {
    return nil, err
  }

  bookSources = append(bookSources, check2AppendBookSources(sources, c)...)

  for nextPageUrl != "" && len(bookSources) < BookSourcesLimitCount {
    sources, nextPageUrl, err = crawler.Crawl(nextPageUrl, c)
    if err != nil {
      return bookSources, err
    }

    if len(sources) == 0 {
      break
    }

    bookSources = append(bookSources, check2AppendBookSources(sources, c)...)
  }

  return bookSources, nil
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

  htmlStr, _, _, err := util.GetHtmlFromUrl(searchUrl, "gbk")
  if err != nil {
    return "", err
  }

  pattern := fmt.Sprintf(`<a href="([^"]+)[^>]+><b>%s`, sc.BookTitle)
  rx := regexp.MustCompile(pattern)

  matches := rx.FindStringSubmatch(htmlStr)
  if len(matches) < 2 || strings.TrimSpace(matches[1]) == "" {
    return "", util.ErrRegexCannotMatch
  }

  return fmt.Sprintf("http://www.sodu.so%s", matches[1]), nil
}

func (s *SoDuSearch) Crawl(sourcesUrl string, sc *models.SearchCrawlContext) ([]models.BookSource, string, error) {
  htmlStr, _, _, err := util.GetHtmlFromUrl(sourcesUrl, "gbk")
  if err != nil {
    return nil, "", err
  }

  pattern := `<div[^<]*<div[^<]*<a[^>]+>%s_[^>]+[^<]*</div[^<]*<div[^<]*<a[^<]*</a[^<]*</div[^<]*<div[^<]*</div[^<]*</div>`
  pattern = fmt.Sprintf(pattern, sc.BookTitle)

  rx := regexp.MustCompile(pattern)

  matches := rx.FindAllString(htmlStr, -1)

  var bookSources []models.BookSource

  for _, m := range matches {
    chapterUrl, chapter, err := exactChapterAndUrl(m, sc.BookTitle)
    if err != nil {
      continue
    }

    updateTime, err := exactUpdateTime(m)
    if err != nil {
      continue
    }

    updateTime, err = s.tranSoDuUpdateTime2Standard(updateTime)
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

    bookSources = append(bookSources, bs)
  }

  return bookSources, "", nil
}

func (s *SoDuSearch) tranSoDuUpdateTime2Standard(dateTime string) (string, error) {
  t, err := time.Parse(util.CKSoDuYearDateTimeLayout, dateTime)
  if err != nil {
    return "", err
  }

  return t.Format(util.CKYearDateTimeLayout), nil
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

func (s *EASOUSearch) Crawl(sourceUrl string, sc *models.SearchCrawlContext) ([]models.BookSource, string, error) {

  htmlStr, _, _, err := util.GetHtmlFromUrl(sourceUrl, "")
  if err != nil {
    return nil, "", err
  }

  nodes, err := html.Parse(strings.NewReader(htmlStr))
  if err != nil {
    return nil, "", err
  }

  doc := goquery.NewDocumentFromNode(nodes)

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

  nextPageRegex := regexp.MustCompile(`<a[^‘“]+['"]([^'"]+)['"][^>]*>下页</a>`)
  matches := nextPageRegex.FindStringSubmatch(htmlStr)
  if len(matches) < 2 || strings.TrimSpace(matches[1]) == "" { // 无匹配认为无下页
    return bookSources, "", nil
  }

  // 得到的url html编码了
  nextPageUrl := html.UnescapeString("http://book.easou.com" + matches[1])

  return bookSources, nextPageUrl, nil
}
