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

func TestExtractVariablesFromStringWithMultipleVariablesInSingleLine(t *testing.T) {
	variables, err := ExtractVariablesFromString("export NAME=\"__FNAME__ __LNAME__\"", "__")
	if err != nil {
		t.Error(err)
	}
	if len(variables) != 2 {
		t.Errorf("Expected 2 variables, but got %d instead\n", len(variables))
	}
	if variables[0] != "__FNAME__" {
		t.Error("Expected variable __FNAME__")
	}
	if variables[1] != "__LNAME__" {
		t.Error("Expected variable __LNAME__")
	}
}

func TestExtractVariablesFromStringWithEmptyName(t *testing.T) {
	variables, err := ExtractVariablesFromString("export EMPTY=\"____\"", "__")
	if err != nil {
		t.Error(err)
	}
	if len(variables) != 0 {
		t.Errorf("Expected 0 variables, but got %d instead\n", len(variables))
	}
}

func TestExtractVariablesFromStringWithVariablesAtExtremities(t *testing.T) {
	variables, err := ExtractVariablesFromString("__GREETING__\nMy name is __NAME__\n__COLOR__ is my favorite color", "__")
	if err != nil {
		t.Error(err)
	}
	if len(variables) != 3 {
		t.Errorf("Expected 3 variables, but got %d instead\n", len(variables))
	}
	if variables[0] != "__GREETING__" {
		t.Error("Expected variable __GREETING__")
	}
	if variables[1] != "__NAME__" {
		t.Error("Expected variable __NAME__")
	}
	if variables[2] != "__COLOR__" {
		t.Error("Expected variable __COLOR__")
	}
}

func TestExtractVariablesFromStringWithBadMultilineVariable(t *testing.T) {
	variables, err := ExtractVariablesFromString("__HEL\nLO__", "__")
	if err != nil {
		t.Error(err)
	}
	if len(variables) != 0 {
		t.Errorf("Expected 0 variables, but got %d instead\n", len(variables))
	}
}
