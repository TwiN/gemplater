package install

import "testing"

func TestExtractVariablesFromString(t *testing.T) {
	variables, err := ExtractVariablesFromString("export EDITOR=\"__EDITOR__\"\nexport GOPATH=\"__GOPATH__\"", "__")
	if err != nil {
		t.Error(err)
	}
	if len(variables) != 2 {
		t.Errorf("Expected 2 variables, but got %d instead\n", len(variables))
	}
	if variables[0] != "__EDITOR__" {
		t.Error("Expected variable __EDITOR__")
	}
	if variables[1] != "__GOPATH__" {
		t.Error("Expected variable __GOPATH__")
	}
}
