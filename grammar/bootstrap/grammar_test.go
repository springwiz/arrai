package bootstrap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/arr-ai/arrai/grammar/parse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertUnparse(t *testing.T, expected string, parsers Parsers, v interface{}) bool { //nolint:unparam
	var sb strings.Builder
	_, err := parsers.Unparse(v, &sb)
	return assert.NoError(t, err) && assert.Equal(t, expected, sb.String())
}

var expr = Rule("expr")

var exprGrammarSrc = `
// Simple expression grammar
expr -> expr:/{([-+])}
      ^ expr:/{([*\/])}
      ^ "-"? expr
	  ^ /{(\d+)} | expr
	  ^ expr<:"**"
      ^ "(" expr ")";
`

var exprGrammar = Grammar{
	expr: Stack{
		Delim{Term: expr, Sep: RE(`([-+])`)},
		Delim{Term: expr, Sep: RE(`([*/])`)},
		Seq{Opt(S("-")), expr},
		Oneof{RE(`(\d+)`), expr},
		R2L(expr, S("**")),
		Seq{S("("), expr, S(")")},
	},
}

func assertEqualObjects(t *testing.T, expected, actual interface{}) bool { //nolint:unparam
	if assert.True(t, reflect.DeepEqual(expected, actual)) {
		return true
	}
	t.Logf("raw expected: %#v", expected)
	t.Logf("raw actual:   %#v", actual)

	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)
	actualJSON, err := json.Marshal(actual)
	require.NoError(t, err)
	t.Log("JSON(expected): ", string(expectedJSON))
	t.Log("JSON(actual):   ", string(actualJSON))

	assert.JSONEq(t, string(expectedJSON), string(actualJSON))

	return false
}

func assertEqualNodes(t *testing.T, expected, actual parse.Node) bool {
	if diff := parse.NewNodeDiff(&expected, &actual); !assert.True(t, diff.Equal()) {
		t.Logf("\nexpected: %v\nactual  : %v\ndiff: %v", expected, actual, diff)
		return false
	}
	return true
}

func assertParseToNode(t *testing.T, expected parse.Node, rule Rule, input *parse.Scanner) bool { //nolint:unparam
	parsers := Core()
	v, err := parsers.Parse(rule, input)
	if assert.NoError(t, err) {
		if assert.NoError(t, parsers.ValidateParse(v)) {
			return assertEqualNodes(t, expected, v.(parse.Node))
		}
	} else {
		t.Logf("input: %s", input.Context())
	}
	return false
}

type testNodes struct {
	r      *parse.Scanner
	assert *require.Assertions
}

func (n testNodes) slice(s string, a, b int) *parse.Scanner {
	slice := n.r.Slice(a, b)
	n.assert.Equal(s, slice.String())
	return slice
}

func (n testNodes) term(children ...interface{}) parse.Node {
	return parse.Node{Tag: "term", Extra: NonAssociative, Children: children}
}

func (n testNodes) atom(i int, a interface{}) parse.Node {
	return parse.Node{Tag: "atom", Extra: i, Children: []interface{}{a}}
}

func (n testNodes) ident(s *parse.Scanner) parse.Node {
	return n.atom(0, s)
}

func (n testNodes) diff(a, b parse.Node) parse.Node {
	return stack(`term\:`, NonAssociative).a(`\:`, NonAssociative).a(`\_`).z(
		stack(`term#3\?`).a(`\_`).z(
			stack(`named\_`).z(
				stack(`?`).z(),
				a,
			),
			stack(`?`).z(),
		),
		stack(`?`).z(),
	)
}

func (n testNodes) quant(a, b parse.Node) parse.Node {
}

func (n testNodes) named(name *parse.Scanner, atom parse.Node) parse.Node {
	return stack(`term\:`, NonAssociative).a(`\:`, NonAssociative).a(`\_`).z(
		stack(`term#3\?`).a(`\_`).z(
			stack(`named\_`).z(
				stack(`?`).a(`_`).z(name, atom),
				atom,
			),
			stack(`?`).z(),
		),
		stack(`?`).z(),
	)
}

func TestParseNamedTerm(t *testing.T) {
	r := parse.NewScanner(`opt=""`)
	x := stack(`term\:`, NonAssociative).a(`\:`, NonAssociative).a(`\_`).z(
		stack(`term#3\?`).a(`\_`).z(
			stack(`named\_`).z(
				stack(`?`).a(`_`).z(r.Slice(0, 3), r.Slice(3, 4)),
				stack(`atom\|`, 1).z(r.Slice(4, 6)),
			),
			stack(`?`).z(),
		),
		stack(`?`).z(),
	)
	assertParseToNode(t, x, term, r)
}

func TestParseNamedTermInDelim(t *testing.T) {
	r := parse.NewScanner(`"1":op=","`)
	x := stack(`term\:`, NonAssociative).a(`\:`, NonAssociative).a(`\_`).z(
		stack(`term#3\?`).a(`\_`).z(
			stack(`named\_`).z(
				stack(`?`).z(),
				stack(`atom\|`, 1).z(r.Slice(0, 3)),
			),
			stack(`?`).a(`quant\|`, 2).a(`_`).z(
				r.Slice(3, 4),
				stack(`?`).z(),
				stack(`named\_`).z(
					stack(`?`).a(`_`).z(r.Slice(4, 6), r.Slice(6, 7)),
					stack(`atom\|`, 1).z(r.Slice(7, 10)),
				),
				stack(`?`).z(),
			),
		),
		stack(`?`).z(),
	)
	assertParseToNode(t, x, term, r)
}

func TestGrammarParser(t *testing.T) {
	t.Parallel()

	parsers := exprGrammar.Compile()

	r := parse.NewScanner("1+2*3")
	v, err := parsers.Parse(expr, r)
	require.NoError(t, err)
	assert.NoError(t, parsers.ValidateParse(v))
	assertUnparse(t, "1+2*3", parsers, v)
	assert.Equal(t,
		`expr\:║:(expr#1\:║:(expr#2\_(?(), expr#3\|║0(1))), `+
			`+, `+
			`expr#1\:║:(expr#2\_(?(), expr#3\|║0(2)), *, expr#2\_(?(), expr#3\|║0(3))))`,
		fmt.Sprintf("%v", v),
	)

	r = parse.NewScanner("1+(2-3/4)")
	v, err = parsers.Parse(expr, r)
	assert.NoError(t, err)
	assert.NoError(t, parsers.ValidateParse(v))
	assertUnparse(t, "1+(2-3/4)", parsers, v)
	assert.Equal(t,
		`expr\:║:(`+
			`expr#1\:║:(expr#2\_(?(), expr#3\|║0(1))), `+
			`+, `+
			`expr#1\:║:(expr#2\_(?(), expr#3\|║1(expr#4\:║:(expr#5\_((, `+
			`expr\:║:(expr#1\:║:(expr#2\_(?(), expr#3\|║0(2))), `+
			`-, `+
			`expr#1\:║:(expr#2\_(?(), expr#3\|║0(3)), `+
			`/, `+
			`expr#2\_(?(), expr#3\|║0(4)))), `+
			`)))))))`,
		fmt.Sprintf("%v", v),
	)
}

func TestExprGrammarGrammar(t *testing.T) {
	t.Parallel()

	parsers := Core()
	r := parse.NewScanner(exprGrammarSrc)
	v, err := parsers.Parse(grammarR, r)
	require.NoError(t, err, "r=%v\nv=%v", r.Context(), v)
	require.Equal(t, len(exprGrammarSrc), r.Offset(), "r=%v\nv=%v", r.Context(), v)
	assert.NoError(t, parsers.ValidateParse(v))
	assertUnparse(t,
		`// Simple expression grammar`+
			`expr->expr:([-+])`+
			`^expr:([*\/])`+
			`^"-"?expr`+
			`^(\d+)|expr`+
			`^expr<:"**"`+
			`^"("expr")";`,
		parsers,
		v,
	)
}

func TestGrammarSnippet(t *testing.T) {
	t.Parallel()

	parsers := Core()
	r := parse.NewScanner(`prod+`)
	v, err := parsers.Parse(term, r)
	require.NoError(t, err)
	assert.Equal(t,
		`term\:║:(term#1\:║:(term#2\_(term#3\?(term#4\_(named\_(?(), atom\|║0(prod)), ?(quant\|║0(+)))), ?())))`,
		fmt.Sprintf("%v", v),
	)
	assert.NoError(t, parsers.ValidateParse(v))
	assertUnparse(t, "prod+", parsers, v)
}

func TestTinyGrammarGrammarGrammar(t *testing.T) {
	t.Parallel()

	tiny := Rule("tiny")
	tinyGrammar := Grammar{tiny: S("x")}
	tinyGrammarSrc := `tiny -> "x";`

	parsers := Core()
	r := parse.NewScanner(tinyGrammarSrc)
	v, err := parsers.Parse(grammarR, r)
	require.NoError(t, err)
	e := v.(parse.Node)
	assert.NoError(t, parsers.ValidateParse(v))

	grammar2 := NewFromNode(e)
	assertEqualObjects(t, tinyGrammar, grammar2)
}

func TestExprGrammarGrammarGrammar(t *testing.T) {
	t.Parallel()

	parsers := Core()
	r := parse.NewScanner(exprGrammarSrc)
	v, err := parsers.Parse(grammarR, r)
	require.NoError(t, err)
	e := v.(parse.Node)
	assert.NoError(t, parsers.ValidateParse(v))

	grammar2 := NewFromNode(e)
	assertEqualObjects(t, exprGrammar, grammar2)
}

func TestDiff(t *testing.T) {
	n := testNodes{r: parse.NewScanner(`a~b`), assert: require.New(t)}
	// assertParseToNode(t, r.diff(r.ident(0, 1))), term, n.r)
}
