package util

import (
  "html"
  "regexp"
  "strings"
)

var (
  BodyRegex         = `(?s)<body.*?</body>` // 多行匹配
  PreDepthLine      = 6                     // 预先统计行数
  StartLimitCount   = 180
  HeadEmptyLines    = 2
  EndLimitCharCount = 20
)

func formatTags(tag string) string {
  result := strings.Replace(tag, "\r", "", -1)
  return strings.Replace(result, "\n", "", -1)
}

// 从 body 标签文本中分析正文内容
//
// 剔除标签
// 预先统计 PreDepthLine 行的字符个数 PreTextLen
// PreTextLen > StartLimitCount && 后续还有字符时 认为正文开始
// 连续空行数满足 HeadEmptyLines 时人为找到文章头 没找到则以以上找到的正文开始为起始
// 赋值内容直到满足结束条件
func getContent(body string) (string, string) {

  var contentLines, contentWithTagsLines []string

  orgLines := strings.Split(body, "\n")
  lines := make([]string, len(orgLines))

  // 去除每行空白字符 剔除标签
  for k, lineInfo := range orgLines {
    lineInfo = ReplaceNewlineTags(lineInfo, "[crlf]")
    lineInfo = FilterAllTags(lineInfo)
    lineInfo = strings.TrimSpace(lineInfo)

    lines[k] = lineInfo
  }

  // 提取正文文本
  var preTextLen int    // 上次统计字符数量
  var startPos int = -1 // 正文起始位置
  var isExistEnd bool

  for i := 0; i < len(lines)-PreDepthLine; i++ {
    depthTextLen := 0
    for j := 0; j < PreDepthLine; j++ { // 从第i行开始 PreDepthLine 行字符个数
      depthTextLen += len(lines[i+j])
    }

    if startPos == -1 { // 还没有找到文章起始
      if preTextLen > StartLimitCount && depthTextLen > 0 {
        // 查找文章起始位置
        emptyCount := 0
        for j := i - 1; j > 0; j-- { // 向上查找 根据连续空行个数确定头部
          if len(lines[j]) == 0 {
            emptyCount += 1
          } else {
            emptyCount = 0
          }

          if emptyCount == HeadEmptyLines {
            startPos = j + HeadEmptyLines
            break
          }
        }

        // 如果没有定位到文章头，则以当前查找位置作为文章头
        if startPos == -1 {
          startPos = i
        }

        // 赋值发现的文章起始部分
        for j := startPos; j < i; j++ {
          contentLines = append(contentLines, lines[j])
          contentWithTagsLines = append(contentWithTagsLines, orgLines[j])
        }
      }
    } else { // 已找到文章起始
      if depthTextLen <= EndLimitCharCount && preTextLen < EndLimitCharCount {
        isExistEnd = true
        break
      }

      contentLines = append(contentLines, lines[i])
      contentWithTagsLines = append(contentWithTagsLines, orgLines[i])
    }

    preTextLen = depthTextLen
  }

  if !isExistEnd { // 不添加此段代码 统计时会直接扔到最后 PreDepthLine 行 有可能导致丢失内容 加了 会导致有时内容冗余
    for i := len(lines) - PreDepthLine; i < len(lines); i++ {
      contentLines = append(contentLines, lines[i])
    }
  }

  content := strings.Join(contentLines, "")
  content = strings.Replace(content, "\r", "", -1)
  content = strings.Replace(content, "[crlf]", "\r\n", -1)
  content = html.UnescapeString(content)

  contentWithTags := strings.Join(contentWithTagsLines, "")

  return content, contentWithTags
}

func UnCompressHtml(htmlStr string) string {
  if strings.Count(htmlStr, "\n") < 10 { // 换行符小于10个人为 htmlStr 为压缩过的
    htmlStr = strings.Replace(htmlStr, ">", ">\n", -1)
  }

  return htmlStr
}

func TranHtmlTagToLower(htmlStr string) string {
  // 将所有标签处理为小写
  toLowerRegex := regexp.MustCompile(`<[^!]*?[^>]+>`)
  return toLowerRegex.ReplaceAllStringFunc(htmlStr, strings.ToLower)
}

func FilterAllTags(htmlStr string) string {
  tagsRegex := regexp.MustCompile(`<[^>]+>`)
  return tagsRegex.ReplaceAllString(htmlStr, "")
}

func ReplaceNewlineTags(htmlStr, repl string) string {
  crlfRegex := regexp.MustCompile(`<p>|</p>|<br.*?/?>`)

  return crlfRegex.ReplaceAllString(htmlStr, repl)
}

// htmlStr utf-8 编码
func Content2Article(htmlStr string) string {

  // 过滤掉样式、脚本等标签
  filterRegex := regexp.MustCompile(`(?s)<script.*?>.*?</script>`) // 过滤脚本
  body := filterRegex.ReplaceAllString(htmlStr, "")

  filterRegex = regexp.MustCompile(`(?s)<style.*?>.*?</style>`) // 过滤样式
  body = filterRegex.ReplaceAllString(body, "")

  filterRegex = regexp.MustCompile(`<!--.*?-->`) // 过滤注释 仅能匹配单行 注释的匹配比较麻烦  例如 // --> 时
  body = filterRegex.ReplaceAllString(body, "")

  filterRegex = regexp.MustCompile(`(?s)<a[^>]+>.*?</a>`) // 过滤超链接
  body = filterRegex.ReplaceAllString(body, "")

  // 标签规整化处理，将标签属性格式化处理到同一行
  // 处理形如以下的标签：
  //  <a
  //   href='http://www.baidu.com'
  //   class='test'
  // 处理后为
  //  <a href='http://www.baidu.com' class='test'>
  formatRegex := regexp.MustCompile(`(<[^<>]+)\s*\n\s*`)
  body = formatRegex.ReplaceAllStringFunc(body, formatTags)

  content, _ := getContent(body)

  return content
}

// htmlStr utf-8 编码
// 仅处理 <body> 与 </body>间内容
func Html2Article(htmlStr string) string {

  // 提取 body 内容
  bodyRegex := regexp.MustCompile(BodyRegex)
  body := bodyRegex.FindString(htmlStr)

  return Content2Article(body)
}
