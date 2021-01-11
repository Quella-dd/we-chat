package test

import (
	"fmt"
	"testing"
	"we-chat/message"
)

func TestGetSession(t *testing.T) {
	scope := &message.Scope{
		SourceID: 1,
		DestinationID: 2,
	}
	fmt.Println(scope.GetSession())
	t.Log(scope.GetSession())
}