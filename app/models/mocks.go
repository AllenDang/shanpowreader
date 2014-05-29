// Copyright 2013 The Shanpow Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// mocks包含所有不需要存储到数据的实体定义

package models

type AjaxResult struct {
  Result   bool
  ErrorMsg string
  Data     interface{}
}

type BookSource struct {
  Host       string // 来源网站
  ChapterUrl string // 最新章节链接
  Chapter    string // 最新章节名称 第一百零一章 以牙还牙
  UpdateTime string
}

type SearchCrawlContext struct {
  BookTitle       string
  BookAuthor      string
  IsCrawlableFunc func(string) bool
  HostExists      map[string]int // 本次抓取中书籍源是否已存在
}

type BookSourcesCrawler interface {
  Search(*SearchCrawlContext) (string, error)                      // 搜索书籍来源
  Crawl(string, *SearchCrawlContext) ([]BookSource, string, error) // 抓取书籍最新章节链接 并返回下一页Url
}
