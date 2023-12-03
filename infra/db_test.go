package infra

import (
	"testing"
)

func TestDBConnect(t *testing.T) {
	client := &DBClient{}

	err := client.DBConnect()
	if err != nil {
		t.Errorf("DBConnect failed: %v", err)
	}
}
