package lexer

import (
	"strings"
	"testing"
)

var samples = map[string]string{
	"feature": `Feature: gherkin lexer
  in order to run features
  as gherkin lexer
  I need to be able to parse a feature`,

	"background": `Background:`,

	"scenario": "Scenario: tokenize feature file",

	"step_given": `Given a feature file`,

	"step_when": `When I try to read it`,

	"comment": `# an important comment`,

	"step_then": `Then it should give me tokens`,

	"step_given_table": `Given there are users:
      | name | lastname | num |
      | Jack | Sparrow  | 4   |
      | John | Doe      | 79  |`,
}

func indent(n int, s string) string {
	return strings.Repeat(" ", n) + s
}

func Test_feature_read(t *testing.T) {
	l := New(strings.NewReader(samples["feature"]))
	tok := l.Next()
	if tok.Type != FEATURE {
		t.Fatalf("Expected a 'feature' type, but got: '%s'", tok.Type)
	}
	val := "gherkin lexer"
	if tok.Value != val {
		t.Fatalf("Expected a token value to be '%s', but got: '%s'", val, tok.Value)
	}
	if tok.Line != 0 {
		t.Fatalf("Expected a token line to be '0', but got: '%d'", tok.Line)
	}
	if tok.Indent != 0 {
		t.Fatalf("Expected a token identation to be '0', but got: '%d'", tok.Indent)
	}

	tok = l.Next()
	if tok.Type != TEXT {
		t.Fatalf("Expected a 'text' type, but got: '%s'", tok.Type)
	}
	val = "in order to run features"
	if tok.Value != val {
		t.Fatalf("Expected a token value to be '%s', but got: '%s'", val, tok.Value)
	}
	if tok.Line != 1 {
		t.Fatalf("Expected a token line to be '1', but got: '%d'", tok.Line)
	}
	if tok.Indent != 2 {
		t.Fatalf("Expected a token identation to be '2', but got: '%d'", tok.Indent)
	}

	tok = l.Next()
	if tok.Type != TEXT {
		t.Fatalf("Expected a 'text' type, but got: '%s'", tok.Type)
	}
	val = "as gherkin lexer"
	if tok.Value != val {
		t.Fatalf("Expected a token value to be '%s', but got: '%s'", val, tok.Value)
	}
	if tok.Line != 2 {
		t.Fatalf("Expected a token line to be '2', but got: '%d'", tok.Line)
	}
	if tok.Indent != 2 {
		t.Fatalf("Expected a token identation to be '2', but got: '%d'", tok.Indent)
	}

	tok = l.Next()
	if tok.Type != TEXT {
		t.Fatalf("Expected a 'text' type, but got: '%s'", tok.Type)
	}
	val = "I need to be able to parse a feature"
	if tok.Value != val {
		t.Fatalf("Expected a token value to be '%s', but got: '%s'", val, tok.Value)
	}
	if tok.Line != 3 {
		t.Fatalf("Expected a token line to be '3', but got: '%d'", tok.Line)
	}
	if tok.Indent != 2 {
		t.Fatalf("Expected a token identation to be '2', but got: '%d'", tok.Indent)
	}

	tok = l.Next()
	if tok.Type != EOF {
		t.Fatalf("Expected an 'eof' type, but got: '%s'", tok.Type)
	}
}

func Test_minimal_feature(t *testing.T) {
	file := strings.Join([]string{
		samples["feature"] + "\n",

		indent(2, samples["background"]),
		indent(4, samples["step_given"]) + "\n",

		indent(2, samples["comment"]),
		indent(2, samples["scenario"]),
		indent(4, samples["step_given"]),
		indent(4, samples["step_when"]),
		indent(4, samples["step_then"]),
	}, "\n")
	l := New(strings.NewReader(file))

	var tokens []TokenType
	for tok := l.Next(); tok.Type != EOF; tok = l.Next() {
		tokens = append(tokens, tok.Type)
	}
	expected := []TokenType{
		FEATURE,
		TEXT,
		TEXT,
		TEXT,
		NEW_LINE,

		BACKGROUND,
		GIVEN,
		NEW_LINE,

		COMMENT,
		SCENARIO,
		GIVEN,
		WHEN,
		THEN,
	}
	for i := 0; i < len(expected); i++ {
		if expected[i] != tokens[i] {
			t.Fatalf("expected token '%s' at position: %d, is not the same as actual token: '%s'", expected[i], i, tokens[i])
		}
	}
}

func Test_table_row_reading(t *testing.T) {
	file := strings.Join([]string{
		indent(2, samples["background"]),
		indent(4, samples["step_given_table"]),
		indent(4, samples["step_given"]),
	}, "\n")
	l := New(strings.NewReader(file))

	var types []TokenType
	var values []string
	var indents []int
	for tok := l.Next(); tok.Type != EOF; tok = l.Next() {
		types = append(types, tok.Type)
		values = append(values, tok.Value)
		indents = append(indents, tok.Indent)
	}
	expectedTypes := []TokenType{
		BACKGROUND,
		GIVEN,
		TABLE_ROW,
		TABLE_ROW,
		TABLE_ROW,
		GIVEN,
	}
	expectedIndents := []int{2, 4, 6, 6, 6, 4}
	for i := 0; i < len(expectedTypes); i++ {
		if expectedTypes[i] != types[i] {
			t.Fatalf("expected token type '%s' at position: %d, is not the same as actual: '%s'", expectedTypes[i], i, types[i])
		}
	}
	for i := 0; i < len(expectedIndents); i++ {
		if expectedIndents[i] != indents[i] {
			t.Fatalf("expected token indentation '%d' at position: %d, is not the same as actual: '%d'", expectedIndents[i], i, indents[i])
		}
	}
	if values[2] != "name | lastname | num |" {
		t.Fatalf("table row value '%s' was not expected", values[2])
	}
}
