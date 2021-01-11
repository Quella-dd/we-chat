package test

import (
	"fmt"
	"testing"
	"we-chat/message"
)

func TestGetSession(t *testing.T) {
	scope := &message.Scope{
		SourceID: "2",
		DestinationID: "1",
	}
	fmt.Println(scope.GetSession())
	t.Log(scope.GetSession())
	t.Log(scope.GetSession())
}