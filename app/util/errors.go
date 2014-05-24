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
)
