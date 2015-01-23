package goraygun

import (
  "os"
  "os/exec"
  // "strconv"
  // "fmt"
  // "github.com/kr/pretty"
  "net/http"
  "runtime"
  "strings"
  "time"
)

func (e *Entry) populate(err error, st []StackTraceElement) {
  e.OccurredOn = time.Now().Format("2006-01-02T15:04:05Z")

  hn, _ := os.Hostname()
  e.Details.MachineName = hn

  e.Details.Client.Name = ClientName
  e.Details.Client.Version = ClientVersion
  e.Details.Client.ClientUrl = ClientRepo

  e.Details.Error.StackTrace = st
  e.Details.Error.ClassName = st[0].ClassName
  e.Details.Error.Message = err.Error()

  e.Details.Context.Identifier = uuid()

  return
}

func getMemStats() (m runtime.MemStats) {
  runtime.ReadMemStats(&m)
  return
}

type Entry struct {
  OccurredOn string  `json:"occurredOn"`
  Details    details `json:"details"`
}

type details struct {
  MachineName string             `json:"machineName"`
  Version     string             `json:"version"`
  Client      clientDetails      `json:"client"`
  Error       errorDetails       `json:"error"`
  Context     contextDetails     `json:"context"`
  Environment environmentDetails `json:"environment"`
  Request     requestDetails     `json:"request"`
}

type clientDetails struct {
  Name      string `json:"name"`
  Version   string `json:"version"`
  ClientUrl string `json:"clientUrl"`
}

type errorDetails struct {
  Data       string              `json:"data"`
  ClassName  string              `json:"className"`
  Message    string              `json:"message"`
  StackTrace []StackTraceElement `json:"stackTrace"`
}

func uuid() string {
  out, _ := exec.Command("uuidgen").Output()
  return strings.TrimSpace(string(out))
}

type environmentDetails struct {
  ProcessorCount          int    `json:"processorCount"`
  OsVersion               string `json:"osVersion"`
  Cpu                     string `json:"cpu"`
  PackageVersion          string `json:"packageVersion"`
  Architecture            string `json:"architecture"`
  TotalPhysicalMemory     uint64 `json:"totalPhysicalMemory"`
  AvailablePhysicalMemory uint64 `json:"availablePhysicalMemory"`
  Locale                  string `json:"locale"`
}

func (ed *environmentDetails) populate() {
  memstats := getMemStats()
  ed.ProcessorCount = runtime.NumCPU()
  ed.OsVersion = runtime.GOOS
  ed.TotalPhysicalMemory = memstats.Sys
  ed.AvailablePhysicalMemory = memstats.Sys - memstats.Alloc
}

type requestDetails struct {
  HostName    string            `json:"hostName"`
  Url         string            `json:"url"`
  HttpMethod  string            `json:"httpMethod"`
  IpAddress   string            `json:"iPAddress"`
  Querystring string            `json:"querystring"`
  Form        map[string]string `json:"form"`
  Headers     map[string]string `json:"headers"`
}

func (rd *requestDetails) Populate(r http.Request) {
  rd.HostName = r.Host
  rd.Url = r.URL.String()
  rd.HttpMethod = r.Method
  rd.IpAddress = r.RemoteAddr
  rd.Querystring = r.URL.RawQuery
  rd.Headers = joinChild(r.Header, ", ")
  rd.Form = joinChild(r.Form, ", ")
}

func joinChild(m map[string][]string, sep string) map[string]string {
  newMap := make(map[string]string)
  for k, v := range m {
    newMap[k] = strings.Join(v, sep)
  }
  return newMap
}

type contextDetails struct {
  Identifier string `json:"identifier"`
}
