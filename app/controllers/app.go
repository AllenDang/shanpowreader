package controllers

import (
  "fmt"
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

  dal := models.NewDal(session)

  logicFunc(dal, &r)

  return &r
}

func (c *App) Index() revel.Result {

  logicFunc := func(dal *models.Dal, r *models.AjaxResult) {
    r.Result = false
    r.ErrorMsg = "It works"
  }

  r := ajaxWrapper(c.Controller, c.MongoSession, logicFunc)

  return c.RenderJson(*r)
}
