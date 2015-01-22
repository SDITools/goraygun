package goraygun_test

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "testing"
)

func TestGoRaygun(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Raygun Suite")
}
