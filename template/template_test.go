package template

import "testing"

func TestTemplate_Replace(t *testing.T) {
	template := &Template{
		variables: map[string]string{"__IS__": "is_not"},
	}
	output := template.Replace("export this=\"__IS__/a/test\"")
	const ExpectedOutput = "export this=\"is_not/a/test\""
	if output != ExpectedOutput {
		t.Errorf("Expected '%s', but got '%s' instead", ExpectedOutput, output)
	}
}

func TestTemplate_ReplaceWithMultipleOccurrenceOfSameVariable(t *testing.T) {
	template := &Template{
		variables: map[string]string{"__NUMBER__": "1"},
	}
	output := template.Replace("__NUMBER__ = __NUMBER__")
	const ExpectedOutput = "1 = 1"
	if output != ExpectedOutput {
		t.Errorf("Expected '%s', but got '%s' instead", ExpectedOutput, output)
	}
}

func TestTemplate_ReplaceWithMultipleVariables(t *testing.T) {
	template := &Template{
		variables: map[string]string{"__FIRST_NAME__": "John", "__LAST_NAME__": "Doe"},
	}
	output := template.Replace("__FIRST_NAME__ __LAST_NAME__")
	const ExpectedOutput = "John Doe"
	if output != ExpectedOutput {
		t.Errorf("Expected '%s', but got '%s' instead", ExpectedOutput, output)
	}
}
