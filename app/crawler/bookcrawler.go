package crawler

import (
  "code.google.com/p/go.net/html"
  "encoding/json"
  "github.com/AllenDang/shanpowreader/app/models"
  "github.com/AllenDang/shanpowreader/app/util"
  "github.com/PuerkitoBio/goquery"
  "github.com/howeyc/fsnotify"
  "github.com/revel/revel"
  "io/ioutil"
  "net/url"
  "path"
  "regexp"
  "strings"
)

type CrawlParam struct {
  Method  string // regex、goquery  目前 目录仅支持 regex，章节内容仅支持 goquery
  Pattern string
}

type CrawlerConfig struct {
  ForUrl   string                // 网址
  Encoding string                // 编码方式 空 为 utf-8
  Params   map[string]CrawlParam // 书籍目录、章节内容抓取配置
}

type Crawler struct {
  config CrawlerConfig
}

func NewCrawler(config CrawlerConfig) *Crawler {
  var c Crawler
  c.config = config
  return &c
}

// 好书网 www.hao662.com
// 抓取章节列表
// contentsUrl 可以为 章节url 或 章节目录 url
// 章节 http://www.hao662.com/haoshu/0/168/6380396.html
// 章节目录 http://www.hao662.com/haoshu/0/168

func (c *Crawler) BookContentsCrawl(contentsUrl string) ([]models.Chapter, error) {

  contentsRegex, err := regexp.Compile(c.config.Params["ContentsUrl"].Pattern)
  if err != nil {
    return nil, err
  }

  chapterRegex, err := regexp.Compile(c.config.Params["ChapterUrl"].Pattern)
  if err != nil {
    return nil, err
  }

  contentsUrl = contentsRegex.FindString(contentsUrl)
  if contentsUrl == "" {
    return nil, util.ErrCanNotConstructBookContentsUrl
  }

  html, _, _, err := util.GetHtmlFromUrl(contentsUrl, c.config.Encoding)
  if err != nil {
    return nil, err
  }

  // 章节链接
  contents := chapterRegex.FindAllStringSubmatch(html, -1)

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

func (c *Crawler) ChapterContentCrawl(chapterUrl string) (string, error) {

  htmlStr, _, _, err := util.GetHtmlFromUrl(chapterUrl, c.config.Encoding)
  if err != nil {
    return "", err
  }

  node, err := html.Parse(strings.NewReader(htmlStr))
  if err != nil {
    return "", err
  }

  doc := goquery.NewDocumentFromNode(node)

  return doc.Find(c.config.Params["Chapter"].Pattern).Html()
}

type CrawlerManager struct {
  crawlerMap map[string]*Crawler
  configPath string
}

func NewCrawlerManager() *CrawlerManager {
  cm := new(CrawlerManager)

  configFilename := path.Join(revel.ConfPaths[0], "bookcrawler.conf")
  cm.configPath = configFilename
  cm.Init(configFilename)

  return cm
}

func (c *CrawlerManager) Init(configFilename string) {
  b, err := ioutil.ReadFile(configFilename)
  if err != nil {
    revel.ERROR.Panicf("Cannot load cralwer config file (%s)", configFilename)
  }

  var configs []CrawlerConfig
  err = json.Unmarshal(b, &configs)
  if err != nil {
    revel.ERROR.Panicf("Cannot unmarshal crawler config file (%s)", err)
  }

  c.crawlerMap = make(map[string]*Crawler)
  for _, config := range configs {
    c.crawlerMap[config.ForUrl] = NewCrawler(config)
  }
}

//检测Crawler.conf的更改，注意，整个应用程序周期内只能调用一次
func (c *CrawlerManager) MonitorConfigChange() {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    revel.ERROR.Panicf("Cannot create fswatcher (%s)", err)
  }

  go func() {
    for {
      select {
      case ev := <-watcher.Event:
        if ev.IsModify() {
          c.Init(c.configPath)
          revel.INFO.Println("Crawler config file is reloaded")
        }
      case err := <-watcher.Error:
        revel.ERROR.Println("fswatcher error:", err)
      }
    }
  }()

  err = watcher.Watch(c.configPath)
  if err != nil {
    revel.ERROR.Panicf("Failed to watch crawler.conf (%s)", err)
  }
}

//判断域名是否可以抓取
func (c *CrawlerManager) IsCrawlable(host string) bool {
  for _, v := range c.crawlerMap {
    if v.config.ForUrl == host {
      return true
    }
  }

  return false
}

func (c *CrawlerManager) BookContentsCrawl(contentsUrl string) ([]models.Chapter, error) {

  u, err := url.Parse(contentsUrl)
  if err != nil {
    return nil, err
  }

  crawler, ok := c.crawlerMap[u.Host]
  if !ok {
    return nil, util.ErrNotSupportCrawl
  }

  // 配置参数检查
  contentsUrlCfg, ok := crawler.config.Params["ContentsUrl"]
  if !ok {
    return nil, util.ErrCanNotFindContentsUrlCfg
  }

  chapterUrlCfg, ok := crawler.config.Params["ChapterUrl"]
  if !ok {
    return nil, util.ErrCanNotFindChapterUrlCfg
  }

  if contentsUrlCfg.Method != "" || chapterUrlCfg.Method != "" {
    return nil, util.ErrOnlySupportRegexMethodForNow
  }

  return crawler.BookContentsCrawl(contentsUrl)
}

func (c *CrawlerManager) ChapterContentCrawl(chapterUrl string) (string, error) {

  u, err := url.Parse(chapterUrl)
  if err != nil {
    return "", err
  }

  crawler, ok := c.crawlerMap[u.Host]
  if !ok {
    return "", util.ErrNotSupportCrawl
  }

  // 配置参数检查
  chapterCfg, ok := crawler.config.Params["Chapter"]
  if !ok {
    return "", util.ErrCanNotFindChapterCfg
  }

  if chapterCfg.Method != "goquery" {
    return "", util.ErrOnlySupportGoQueryMethodForNow
  }

  return crawler.ChapterContentCrawl(chapterUrl)
}
