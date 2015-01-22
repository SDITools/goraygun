package goraygun

import (
  "errors"
  "regexp"
  "runtime"
  "strconv"
  "strings"
)

type StackTraceElement struct {
  LineNumber int    `json:"lineNumber"`
  ClassName  string `json:"className"`
  FileName   string `json:"fileName"`
  MethodName string `json:"methodName"`
}

func GetStackTrace(offset int) ([]StackTraceElement, error) {
  rst := getRawStackTrace()
  st, err := ParseStackTrace(rst)
  if err != nil {
    return st, err
  }
  // Omit calling functions
  return st[offset:], err
}

func getRawStackTrace() []byte {
  st := make([]byte, 1<<16)
  return st[:runtime.Stack(st, false)]
}

func ParseStackTrace(rst []byte) (sts []StackTraceElement, err error) {
  str := string(rst)
  lines := strings.Split(str, "\n")

  re := regexp.MustCompile("(.+?)\\.([^/]+)\\([\\w\\d\\s,.]*\\)\\s+(.+\\.go):(\\d+)")

  for i := 1; i < len(lines); i = i + 2 {
    elem := StackTraceElement{}

    if len(lines) <= i+1 {
      break
    }
    submatches := re.FindStringSubmatch(lines[i] + lines[i+1])
    if len(submatches) < 5 {
      continue
    }
    elem.ClassName = submatches[1]
    elem.MethodName = submatches[2]
    elem.FileName = submatches[3]
    elem.LineNumber, _ = strconv.Atoi(submatches[4])
    sts = append(sts, elem)
  }
  if len(sts) > 0 {
    return
  }
  return sts, errors.New("Invalid stack trace input")
}
