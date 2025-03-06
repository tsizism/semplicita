package main

import (
	"authentication/data"
	"os"
	"testing"
)

var cfgTestApp config 

func TestMain(m *testing.M) {
	println("_______________________TestMain_________________")

	testRepo := data.NewPostgressRepositoryTest(nil)
	cfgTestApp.repo = testRepo


	os.Exit(m.Run())
}