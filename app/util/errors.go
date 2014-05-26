package util

import (
  "errors"
)

var (
  ErrRegexCannotMatch               = errors.New("正则表达式没有匹配到内容")
  ErrAutoRedirectForbidden          = errors.New("禁用自动重定向")
  ErrCanNotConstructBookContentsUrl = errors.New("无法得到书籍目录 url")
  ErrCrawBookContentsFailed         = errors.New("抓取书籍目录失败")
  ErrCanNotFindBook                 = errors.New("没有搜索到书籍")
  ErrSearchEngineNotSurpport        = errors.New("搜索引擎不支持")
  ErrCanNotFindContentsUrlCfg       = errors.New("找不到书籍目录url配置")
  ErrCanNotFindChapterUrlCfg        = errors.New("找不到章节目录url配置")
  ErrCanNotFindChapterCfg           = errors.New("找不到章节抓取配置")
  ErrOnlySupportRegexMethodForNow   = errors.New("目前仅支持 regexp 方法")
  ErrOnlySupportGoQueryMethodForNow = errors.New("目前仅支持 goquery 方法")
  ErrNoValidRegexPattern            = errors.New("正则表达式不可用")
  ErrNotSupportCrawl                = errors.New("该网页不支持抓取")
  ErrHtml2ArticleFailed             = errors.New("没有找到正文")
  ErrConfigParasError               = errors.New("配置文件参数有误")
)
