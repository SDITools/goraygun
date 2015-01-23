package goraygun_test

import (
  "encoding/json"
  "errors"
  // "fmt"
  // "github.com/kr/pretty"
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "github.com/sditools/goraygun"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
)

var _ = Describe("Raygun", func() {
  Context("Receiving a valid, basic stack trace", func() {
    expected := []goraygun.StackTraceElement{
      goraygun.StackTraceElement{
        LineNumber: 31,
        ClassName:  "github.com/sditools/go-raygun",
        FileName:   "/Users/dedalus/Go/src/github.com/sditools/go-raygun/stacktrace.go",
        MethodName: "getRawStackTrace",
      },
      goraygun.StackTraceElement{
        LineNumber: 20,
        ClassName:  "github.com/sditools/go-raygun",
        FileName:   "/Users/dedalus/Go/src/github.com/sditools/go-raygun/stacktrace.go",
        MethodName: "GetStackTrace",
      },
      goraygun.StackTraceElement{
        LineNumber: 51,
        ClassName:  "github.com/sditools/go-raygun",
        FileName:   "/Users/dedalus/Go/src/github.com/sditools/go-raygun/raygun.go",
        MethodName: "(*Client).Report",
      },
      goraygun.StackTraceElement{
        LineNumber: 88,
        ClassName:  "github.com/sditools/go-raygun_test",
        FileName:   "/Users/dedalus/Go/src/github.com/sditools/go-raygun/raygun_test.go",
        MethodName: "func·009",
      },
    }
    rawStackTrace, _ := ioutil.ReadFile("test/stacktrace1")
    stackTrace, err := goraygun.ParseStackTrace(rawStackTrace)
    It("Should correctly serialize the stack trace data", func() {
      Expect(err).To(BeNil())
      Expect(expected).To(Equal(stackTrace))
    })
  })

  Context("Receiving an invalid stack trace", func() {
    stackTrace, err := goraygun.ParseStackTrace([]byte("BAD STACK TRACE"))
    It("Should return an error", func() {
      Expect(err).NotTo(BeNil())
      Expect(len(stackTrace)).To(Equal(0))
    })
  })

  Context("Requesting a stack trace", func() {
    stackTrace, err := goraygun.GetStackTrace(2)
    It("Should return a slice of stack trace elements", func() {
      Expect(err).To(BeNil())
      Expect(len(stackTrace)).NotTo(Equal(0))
    })
  })

  Context("Reporting an error", func() {
    It("Should send an property formatted POST request to the specified endpoint", func() {
      ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        decoder := json.NewDecoder(r.Body)
        var entry goraygun.Entry
        err := decoder.Decode(&entry)
        Expect(err).To(BeNil())
        Expect(len(entry.Details.Error.StackTrace)).To(BeNumerically(">", 0))
        Expect(entry.Details.Error.Message).To(Equal("Test Error"))
        w.WriteHeader(http.StatusAccepted)
      }))
      defer ts.Close()

      settings := goraygun.Settings{
        ApiKey:   "123",
        Enabled:  true,
        Endpoint: ts.URL,
      }

      raygunClient := goraygun.Init(settings, goraygun.Entry{})
      raygunClient.Report(errors.New("Test Error"), raygunClient.Entry)
    })
  })
})
