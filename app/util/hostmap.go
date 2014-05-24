package util

var HostMap = map[string]string{
  "www.hao662.com": "好书网",
}

func IsHostCanBeCrawled(hostUrl string) bool {
  if _, ok := HostMap[hostUrl]; ok {
    return true
  }

  return false
}

func GetHostName(hostUrl string) (string, bool) {
  if name, ok := HostMap[hostUrl]; ok {
    return name, true
  }

  return "", false
}
