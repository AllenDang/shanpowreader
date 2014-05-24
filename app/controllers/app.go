package controllers

import (
  "fmt"
  "github.com/AllenDang/shanpowreader/app/crawler"
  "github.com/AllenDang/shanpowreader/app/models"
  "github.com/jgraham909/revmgo"
  "github.com/revel/revel"
  "labix.org/v2/mgo"
)

type App struct {
  *revel.Controller
  revmgo.MongoController
}

func validationErrorString(validation *revel.Validation) string {
  errMsg := ""
  if validation.HasErrors() {
    for _, e := range validation.Errors {
      errMsg += fmt.Sprintf("%s:%s", e.Key, e.Message)
    }
  }
  return errMsg
}

func ajaxWrapper(c *revel.Controller, session *mgo.Session, logicFunc func(dal *models.Dal, r *models.AjaxResult)) *models.AjaxResult {
  var r models.AjaxResult
  r.Result = true

  if c.Validation.HasErrors() {
    r.Result = false
    r.ErrorMsg = validationErrorString(c.Validation)
    return &r
  }

  //dal := models.NewDal(session)

  logicFunc(nil, &r)

  return &r
}

// se 搜索引擎
// title 书籍名称
// author 书籍作者
// id 书籍 Id
func (c *App) GetBookSources(se, title, author, id string) revel.Result {

  logicFunc := func(dal *models.Dal, r *models.AjaxResult) {
    sources, err := crawler.BookSourcesCrawl("easou", title, author)
    if err != nil {
      r.Result = false
      r.ErrorMsg = err.Error()
      return
    }

    r.Data = sources
  }

  r := ajaxWrapper(c.Controller, c.MongoSession, logicFunc)

  return c.RenderJson(*r)
}
