package main

import (
	"testing"

	"github.com/brotherlogic/keystore/client"
)

func InitTest() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient("./testing")
	return s
}

func Test(t *testing.T) {
	doNothing()
}
