package goraygun

import (
  "os"
  "os/exec"
  // "strconv"
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
  //UserCustomData userCustomerDataDetail `json:"userCustomData"`
  //Response       detailResponse         `json:"response"`
  //User           detailUser             `json:"user"`
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
  TotalPhysicalMemory     int    `json:"totalPhysicalMemory"`
  AvailablePhysicalMemory int    `json:"availablePhysicalMemory"`
  TotalVirtualMemory      int    `json:"totalVirtualMemory"`
  AvailableVirtualMemory  int    `json:"availableVirtualMemory"`
  DiskSpaceFree           []int  `json:"diskSpaceFree"`
  Locale                  string `json:"locale"`
}

type requestDetails struct {
  HostName    string `json:"hostName"`
  Url         string `json:"url"`
  HttpMethod  string `json:"httpMethod"`
  IpAddress   string `json:"ipAddress"`
  Querystring string `json:"querystring"`
  Form        string `json:"form"`
  Headers     string `json:"headers"`
  RawData     string `json:"rawData"`
  // Querystring requestQuerystring `json:"querystring"`
  // Form        requestForm        `json:"form"`
  // Headers     requestHeaders     `json:"headers"`
  // RawData     requestRawData     `json:"rawData"`
}

type requestQuerystring struct {
  data string
}

type requestForm struct{}

type requestHeaders struct{}

type requestRawData struct{}

type detailResponse struct {
  StatusCode int `json:"statusCode"`
}

type detailUser struct {
  Identifier string `json:"identifier"`
}

type contextDetails struct {
  Identifier string `json:"identifier"`
}
