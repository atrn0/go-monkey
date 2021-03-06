package parser

import (
	"fmt"
	"github.com/atrn0/go-monkey/ast"
	"github.com/atrn0/go-monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{input: `let x = 5;`, expectedIdentifier: "x", expectedValue: 5},
		{input: `let y = true;`, expectedIdentifier: "y", expectedValue: true},
		{input: `let foobar = y;`, expectedIdentifier: "foobar", expectedValue: "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)
		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral() not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf(
			"letStmt.Name.TokenLiteral() not '%s'. got=%s",
			name,
			letStmt.Name.TokenLiteral(),
		)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{input: "return 5;", expectedValue: 5},
		{input: "return x;", expectedValue: "x"},
		{input: "return true;", expectedValue: true},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatalf("ParseProgram() returned nil")
		}
		if len(program.Statements) != 1 {
			t.Fatalf("not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("program.Statements[0] is not ast.ReturnStatement. got %T",
				program.Statements[0])
		}

		testLiteralExpression(t, stmt.ReturnValue, tt.expectedValue)
	}

}

func TestIdentifierExpressions(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. got %d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionsStatement. got %T",
			program.Statements[0])
	}
	id, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got %s",
			stmt.Expression)
	}
	if id.Value != "foobar" {
		t.Errorf("id.Value not %s. got %s",
			"foobar", id.Value)
	}
	if id.TokenLiteral() != "foobar" {
		t.Errorf("id.TokenLiteral not %s. got %s", "foobar",
			id.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. got %d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionsStatement. got %T",
			program.Statements[0])
	}
	i, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got %T",
			stmt.Expression)
	}
	if i.Value != 5 {
		t.Errorf("i.Value not %d. got %d",
			5, i.Value)
	}
	if i.TokenLiteral() != "5" {
		t.Errorf("i.TokenLiteral not %s. got %s", "5",
			i.TokenLiteral())
	}
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement. got %d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionsStatement. got %T",
				program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp not *ast.PrefixExpression. got %T",
				stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got %s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	InfixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"1 - 5;", 1, "-", 5},
		{"5 * 4;", 5, "*", 4},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 8;", 5, "==", 8},
		{"1 != 5;", 1, "!=", 5},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}

	for _, tt := range InfixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement. got %d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionsStatement. got %T",
				program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got %T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp Operator is not '%s'. got %q", operator, opExp.Operator)
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestOperationPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a + b", "((-a) + b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b * c", "(a + (b * c))"},
		{"a - b / c + d - f * -g", "(((a - (b / c)) + d) - (f * (-g)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(2 + 7) * 2", "((2 + 7) * 2)"},
		{"1 / (6 + 6)", "(1 / (6 + 6))"},
		{"-(4 + 7)", "(-(4 + 7))"},
		{"!(!true != true)", "(!((!true) != true))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if actual := program.String(); actual != tt.expected {
			t.Errorf("expected %q. got %q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. got %d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionsStatement. got %T",
			program.Statements[0])
	}

	b, ok := stmt.Expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("exp not *ast.Boolean. got %s", stmt.Expression)
	}

	if b.Value != true {
		t.Errorf("b.Value not %t. got %t", true, b.Value)
	}

	if b.TokenLiteral() != "true" {
		t.Errorf("b.TokenLiteral not %s. got %s", "true", b.TokenLiteral())
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. got %d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionsStatement. got %T",
			program.Statements[0])
	}

	ie, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *ast.IfExpression. got %s", stmt.Expression)
	}

	if !testInfixExpression(t, ie.Condition, "x", "<", "y") {
		return
	}

	if len(ie.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got %d",
			len(ie.Consequence.Statements))
	}

	cons, ok := ie.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionsStatement. got %T",
			program.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	if ie.Alternative != nil {
		t.Errorf("ie.Alternative is not nil. got %+v", ie.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. got %d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionsStatement. got %T",
			program.Statements[0])
	}

	ie, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("exp not *ast.IfExpression. got %s", stmt.Expression)
	}

	if !testInfixExpression(t, ie.Condition, "x", "<", "y") {
		return
	}

	if len(ie.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got %d",
			len(ie.Consequence.Statements))
	}

	cons, ok := ie.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"ie.Consequence.Statements[0] is not ast.ExpressionsStatement. got %T",
			ie.Consequence.Statements[0])
	}

	if !testIdentifier(t, cons.Expression, "x") {
		return
	}

	if len(ie.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statements. got %d",
			len(ie.Alternative.Statements))
	}

	alt, ok := ie.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"ie.Alternative.Statements[0] is not ast.ExpressionsStatement. got %T",
			ie.Alternative.Statements[0])
	}

	if !testIdentifier(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := `fn (x ,y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. got %d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionsStatement. got %T",
			program.Statements[0])
	}

	fl, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("exp not *ast.FunctionLiteral. got %T", stmt.Expression)
	}

	if len(fl.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got %d", len(fl.Parameters))
	}

	testLiteralExpression(t, fl.Parameters[0], "x")
	testLiteralExpression(t, fl.Parameters[1], "y")

	if len(fl.Body.Statements) != 1 {
		t.Fatalf("fl.Body.Statements has not 1 statements. got %d", len(fl.Body.Statements))
	}

	body, ok := fl.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("fl.Body.Statements[0] is not ast.ExpressionStatement")
	}

	testInfixExpression(t, body.Expression, "x", "+", "y")
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input              string
		expectedParameters []string
	}{
		{input: `fn () { };`, expectedParameters: []string{}},
		{input: `fn (x) { };`, expectedParameters: []string{"x"}},
		{input: `fn (x, y, z) { };`, expectedParameters: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParameters) {
			t.Errorf("length of parameters is not %d. got %d",
				len(tt.expectedParameters), len(function.Parameters))
		}

		for i, ident := range tt.expectedParameters {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2, 3 + 5, 3 / 5)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement. got %d",
			len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionsStatement. got %T",
			program.Statements[0])
	}

	ce, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("exp not *ast.CallExpression. got %T", stmt.Expression)
	}

	if !testIdentifier(t, ce.Function, "add") {
		return
	}

	if len(ce.Arguments) != 4 {
		t.Fatalf("length of ce.Arguments is not 4. got %d.", len(ce.Arguments))
	}

	testLiteralExpression(t, ce.Arguments[0], 1)
	testLiteralExpression(t, ce.Arguments[1], 2)
	testInfixExpression(t, ce.Arguments[2], 3, "+", 5)
	testInfixExpression(t, ce.Arguments[3], 3, "/", 5)
}

func TestCallArguments(t *testing.T) {
	tests := []struct {
		input          string
		expectedValues []interface{}
	}{
		{input: `add();`, expectedValues: []interface{}{}},
		{input: `min(true);`, expectedValues: []interface{}{true}},
		{input: `max(x, 23, 3);`, expectedValues: []interface{}{"x", 23, 3}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		call := stmt.Expression.(*ast.CallExpression)

		if len(call.Arguments) != len(tt.expectedValues) {
			t.Errorf("length of parameters is not %d. got %d",
				len(tt.expectedValues), len(call.Arguments))
		}

		for i, value := range tt.expectedValues {
			testLiteralExpression(t, call.Arguments[i], value)
		}
	}
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got %T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	i, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got %T", il)
		return false
	}

	if i.Value != value {
		t.Errorf("i.Value not %d. got %d", value, i.Value)
		return false
	}

	if i.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("i.TokenLiteral not %d. got %s", value, i.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, bl ast.Expression, value bool) bool {
	bo, ok := bl.(*ast.Boolean)
	if !ok {
		t.Errorf("bo not *ast.Boolean. got %T", bl)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got %t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got %s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	id, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got %T", exp)
		return false
	}

	if id.Value != value {
		t.Errorf("id.Value not %s. got %s", value, id.Value)
		return false
	}

	if id.TokenLiteral() != value {
		t.Errorf("id.TokenLiteral not %s. got %s", value, id.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
