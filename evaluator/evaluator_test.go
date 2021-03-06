package evaluator

import (
	"github.com/atrn0/go-monkey/lexer"
	"github.com/atrn0/go-monkey/object"
	"github.com/atrn0/go-monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5", 10},
		{"2 + 4 * 8", 34},
		{"4 / (2 + 2)", 1},
		{"3 - 5 + 87", 85},
		{"90 - -9 * 3", 117},
	}

	for _, tt := range tests {
		obj := testEval(tt.input)
		testIntegerObject(t, obj, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %d, expected %d.", result.Value, expected)
		return false
	}

	return false
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 == 2", false},
		{"1 != 1", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != true", false},
		{"false != false", false},
		{"true != false", true},
		{"(1 < 2) == true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got %t, want %t", result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (0 == 0) { 10 }", 10},
		{"if (0 != 0) { 10 }", nil},
		{"if (0 == 0) { 10 } else { 20 }", 10},
		{"if (0 != 0) { 10 } else { 20 }", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got %T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 0;", 10},
		{"9; return 2 * 5; 0;", 10},
		{`
if (10 > 1) {
	if (10 > 1) {
		return 10;
	}
	return 1;
}
`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input       string
		expectedMsg string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + true; 5;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`
if (10 > 1) {
	if (10 > 1) {
		return true + false;
	}
	return true + 5;
}
`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("error object expected. got %T (%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMsg {
			t.Errorf("wrong error message. got '%s', expected '%s'.", errObj.Message, tt.expectedMsg)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; }"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not function. got %T, (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAddr = fn(x) {
	fn(y) {x + y}
}
let addTwo = newAddr(2)
addTwo(3)
`

	testIntegerObject(t, testEval(input), 5)
}

func testEval(input string) object.Object {
	env := object.NewEnvironment()
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program, env)
}
