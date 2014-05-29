package util

import (
  "fmt"
  iconv "github.com/djimenez/iconv-go"
  "io/ioutil"
  "log"
  "net/http"
  "regexp"
  "strconv"
  "strings"
  "time"
)

var (
  CKYearDateTimeLayout     = "2006-01-02 15:04:05"
  CKSoDuYearDateTimeLayout = "2006-1-02 15:04:05"
  CKDateTimeLayout         = "01-02 15:04"
)

// httpUrl 网址
// encoding 编码格式 == "" 为 utf-8
func GetHtmlFromUrl(httpUrl, encoding string) (html, host, actualUrl string, err error) {
  var resp *http.Response
  var body []byte

  if resp, err = http.Get(httpUrl); err != nil {
    return
  }
  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    return html, host, actualUrl, fmt.Errorf("http.Get(%s) failed with status code %d", httpUrl, resp.StatusCode)
  }

  if body, err = ioutil.ReadAll(resp.Body); err != nil {
    return
  }

  html = string(body)

  if encoding != "" {
    if html, err = iconv.ConvertString(html, encoding, "utf-8"); err != nil {
      return
    }
  }

  host = resp.Request.URL.Host
  actualUrl = resp.Request.URL.String()

  return
}

// 01-02 15:04 -> 2006-01-02 15:04
func ExtendTimeLayoutWithYear(timeStr string) (string, error) {
  t, err := time.Parse(CKDateTimeLayout, timeStr)
  if err != nil {
    return "", err
  }

  t = t.AddDate(time.Now().Year(), 0, 0)

  dateTime := t.Format(CKYearDateTimeLayout)

  return dateTime, nil
}

//返回给定时间与当前时间之间的差距
func GetDurationSubNow(datetimeStr string) string {

  t, err := time.ParseInLocation(CKYearDateTimeLayout, datetimeStr, time.Local)
  if err != nil {
    log.Fatalln(datetimeStr, err)
    return "1分钟"
  }

  result := time.Now().Sub(t).String()
  result = strings.Replace(result, "-", "", 1)
  result = strings.Replace(result, "ms", "", 1)

  //如果有小时，则忽略分钟和秒
  if strings.Contains(result, "h") {
    result, _ = ExtractDataByRegex(result, `(\d+)h`, nil)

    if i, er := strconv.Atoi(result); er == nil {
      switch {
      case i < 24:
        result = fmt.Sprintf("%d小时", i)
      case i >= 24 && i < 168:
        result = fmt.Sprintf("%d天", i/24)
      case i >= 168 && i < 672:
        result = fmt.Sprintf("%d周", i/168)
      case i >= 672 && i < 8064:
        result = fmt.Sprintf("%d个月", i/672)
      case i >= 8064:
        result = fmt.Sprintf("%d年", i/8064)
      }
    }
  } else {
    result, _ = ExtractDataByRegex(result, `(\d+)m`, nil)
    if result == "" {
      result = "1"
    }
    result += "分钟"
  }

  return result
}

func ExtractDataByRegex(html, query string, option map[string]interface{}) (string, error) {
  rx := regexp.MustCompile(query)
  value := rx.FindStringSubmatch(html)

  var v interface{}
  ok := false
  if len(value) == 0 {
    if v, ok = option["or"]; ok {
      var op string
      if op, ok = v.(string); ok {
        return ExtractDataByRegex(html, op, nil)
      }
    }

    return "", fmt.Errorf("正则表达式没有匹配到内容:(%s)", query)
  }

  if strings.TrimSpace(value[1]) == "" {
    return "", fmt.Errorf("正则表达式没有匹配到内容:(%s)", query)
  }

  return strings.TrimSpace(value[1]), nil
}
