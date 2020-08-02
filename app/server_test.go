package app

import "testing"

func TestNewStatusCheck(t *testing.T) {
	check := StatusCheck{}
	check.Name = "Database"
	if check.Name != "Database" {
		t.Errorf("Expected %v, received %v\n", "Database", check.Name)
	}
}
