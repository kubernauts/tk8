package test

import (
	"testing"
	"github.com/kubernauts/tk8/cmd"
)

func TestVersion(t *testing.T) {
	
	if cmd.VERSION != "dev-build" {
		t.Errorf("Version was not passed correctly")
	} else {
		t.Log("Version passed correctly")
	}
}