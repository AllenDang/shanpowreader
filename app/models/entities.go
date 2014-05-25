//Copyright 2014 The Shanpow Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// entities包含所有需要存储到数据的实体定义
package models

import (
  "labix.org/v2/mgo/bson"
)

// 书籍章节
type Chapter struct {
  Index uint   // 索引 1、2、3...
  Name  string // 章节名称
  Url   string // 章节链接
}

// 书籍目录
// 单独 Collection 存储
// 通过 BookId 聚合书籍在各个 Host 的最新更新章节
type BookContents struct {
  BookId           bson.ObjectId // 与网站存储书籍Id 一致
  Host             string        // 网站 url
  ContentsUrl      string        // 书籍目录链接 后续爬取时使用
  LatestChapter    Chapter       // 最新更新章节
  LatestUpdateTime string        // 最近更新时间
  Chapters         []Chapter
}

// 章节内容
// 单独 Collection 存储
// 使用 Url 作为索引
type ChapterContent struct {
  Url           string // 章节链接
  Content       string // 章节内容
  ContentLength uint
}
