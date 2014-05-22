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
  Index      uint // 索引 1、2、3...
  Name       string
  Url        string // 章节链接
  UpdateTime string // 更新时间
}

// 书籍目录
// 单独 Collection 存储
type BookContents struct {
  BookId        bson.ObjectId
  Host          string  // 网站
  Url           string  // 书籍目录链接
  LatestChapter Chapter // 最新更新章节
  Chapters      []Chapter
}

// 章节内容
// 单独 Collection 存储
// 使用 Url 作为索引
type ChapterContent struct {
  Url           string // 章节链接
  Content       string // 章节内容
  ContentLength uint
}
