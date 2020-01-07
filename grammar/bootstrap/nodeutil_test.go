package bootstrap

import (
	"github.com/arr-ai/arrai/grammar/parse"
	"github.com/stretchr/testify/require"
)

type testNodes struct {
	r      *parse.Scanner
	assert *require.Assertions
}

func (n testNodes) slice(s string, a, b int) *parse.Scanner {
	slice := n.r.Slice(a, b)
	n.assert.Equal(s, slice.String())
	return slice
}

func (n testNodes) node(tag string, extra interface{}, children ...interface{}) parse.Node {
	return parse.Node{Tag: tag, Extra: extra /* Read all about it! */, Children: children}
}

func (n testNodes) grammar(stmt ...interface{}) parse.Node {
	return n.node("grammar", nil, stmt...)
}

func (n testNodes) stmt() parse.Node {
	panic(Unfinished)
}

func (n testNodes) comment() parse.Node {
	panic(Unfinished)
}

func (n testNodes) prod() parse.Node {
	panic(Unfinished)
}

func (n testNodes) term(children ...interface{}) parse.Node {
	return n.node("term", NonAssociative, children...)
}

func (n testNodes) choice() parse.Node {
	panic(Unfinished)
}

func (n testNodes) diff(a, parse.Node, tilde *parse.Scanner, b parse.Node) parse.Node {
	n.term(n.choice(n.node("term#2", nil, a, n.opt())))
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

func (n testNodes) seq() parse.Node {
	panic(Unfinished)
}

func (n testNodes) nameds() parse.Node {
	panic(Unfinished)
}

func (n testNodes) quant(s *parse.Scanner) parse.Node {
}

func (n testNodes) repeat() parse.Node {
	panic(Unfinished)
}

func (n testNodes) opt() parse.Node {
	panic(Unfinished)
}

func (n testNodes) any() parse.Node {
	panic(Unfinished)
}

func (n testNodes) some() parse.Node {
	panic(Unfinished)
}

func (n testNodes) minmax() parse.Node {
	panic(Unfinished)
}

func (n testNodes) delim() parse.Node {
	panic(Unfinished)
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

func (n testNodes) atom(i int, a interface{}) parse.Node {
	return n.node("atom", i, a)
}

func (n testNodes) ident(s *parse.Scanner) parse.Node {
	return n.atom(0, s)
}

func (n testNodes) str(s *parse.Scanner) parse.Node {
	return n.atom(1, s)
}

func (n testNodes) re(s *parse.Scanner) parse.Node {
	return n.atom(1, s)
}

func (n testNodes) paren() parse.Node {
	panic(Unfinished)
}

func (n testNodes) empty() parse.Node {
	panic(Unfinished)
}
