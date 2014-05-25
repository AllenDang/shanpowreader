package util

import (
  "fmt"
  iconv "github.com/djimenez/iconv-go"
  "io/ioutil"
  "net/http"
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
