package models

import (
  "github.com/revel/revel"
  "labix.org/v2/mgo"
)

const (
  DbName                 = "shanpowreader"
  BookContentsCollection = "bookcontents"
  BookChapterCollection  = "bookchapter"
)

//Data access layer
type Dal struct {
  session *mgo.Session
}

//创建新的Dal
func NewDal(session *mgo.Session) *Dal {
  if session == nil {
    panic("session cannot be nil.")
  }
  return &Dal{session}
}

func NewDalByDial() *Dal {
  var dial string
  var found bool
  if dial, found = revel.Config.String("revmgo.dial"); !found {
    // Default to 'localhost'
    dial = "localhost"
  }

  // Read configuration.
  var session *mgo.Session
  var err error
  if session, err = mgo.Dial(dial); err != nil {
    revel.ERROR.Panic(err)
  }

  return &Dal{session}
}

func (d *Dal) Close() {
  d.session.LogoutAll()
  d.session.Close()
}

//仅供测试使用
func newDalForTest(ip string) *Dal {
  session, err := mgo.Dial(ip)
  if err != nil {
    panic(err)
  }
  session.SetMode(mgo.Monotonic, true)

  return &Dal{session}
}
