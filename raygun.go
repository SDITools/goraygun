package goraygun

import (
  "encoding/json"
  "errors"
  "log"
  "net/http"
  "strings"
)

const ENDPOINT = "https://api.raygun.io/entries"

const (
  ClientName    = "Go-Raygun"
  ClientVersion = "0.0.1"
  ClientRepo    = "http://github.com/sditools/go-raygun"
)

type Settings struct {
  ApiKey   string
  Enabled  bool
  Endpoint string
}

type Client struct {
  settings Settings
  Entry    Entry
}

func Init(s Settings, e Entry) (c Client) {
  c.settings = s
  c.Entry = e

  // provide user the option to override (for testing)
  if s.Endpoint == "" {
    s.Endpoint = ENDPOINT
  }
  return
}

func (c *Client) Recover() {
  if err := recover(); err != nil {
    c.Report(getError(err), c.Entry)
  }
}

func (c *Client) Report(err error, entry Entry) {
  if !c.settings.Enabled {
    return
  }
  st, stErr := GetStackTrace(2)
  if stErr != nil {
    // handle stErr
    return
  }
  entry.populate(err, st)
  c.post(entry, c.settings.Endpoint)
}

func (c *Client) post(e Entry, uri string) {
  data, err := json.Marshal(e)
  if err != nil {
    log.Printf("Error Marshalling RayGun Message: %v:", err)
  }

  req, err := http.NewRequest("POST", uri, strings.NewReader(string(data)))
  if err != nil {
    log.Printf("Error creating POST request: %v", err)
  }

  req.Header.Set("X-ApiKey", c.settings.ApiKey)
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    log.Printf("Error sending Request: %v", err)
  }

  if resp.StatusCode != http.StatusAccepted {
    log.Printf("Error status sent back: %v", resp.StatusCode)
  }

  defer resp.Body.Close()
}

func getError(err interface{}) error {
  switch err := err.(type) {
  case error:
    return err
  case string:
    return errors.New(err)
  default:
    return errors.New("")
  }
}