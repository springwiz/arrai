package syntax

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"

	"github.com/arr-ai/arrai/rel"
	"github.com/arr-ai/wbnf/parser"
	"github.com/arr-ai/wbnf/wbnf"
)

var once sync.Once
var stdScopeVar rel.Scope

func stdScope() rel.Scope {
	once.Do(func() {
		stdScopeVar = rel.EmptyScope.
			With(".", rel.NewTuple(
				rel.NewAttr("math", rel.NewTuple(
					rel.NewAttr("pi", rel.NewNumber(math.Pi)),
					newFloatFuncAttr("sin", math.Sin),
					newFloatFuncAttr("cos", math.Cos),
				)),
				rel.NewAttr("grammar", rel.NewTuple(
					rel.NewNativeFunctionAttr("parse", parseGrammar),
					rel.NewAttr("lang", rel.NewTuple(
						rel.NewAttr("arrai", rel.ASTNodeToValue(wbnf.FromParserNode(
							wbnf.Core().Grammar(), *arraiParsers.Node()))),
						rel.NewAttr("wbnf", rel.ASTNodeToValue(wbnf.FromParserNode(wbnf.Core().Grammar(), *wbnf.Core().Node()))),
					)),
				)),
				rel.NewAttr("func", rel.NewTuple(
					rel.NewAttr("fix", parseLit(`(\f f(f))(\f \g \n g(f(f)(g))(n))`)),
					// rel.NewAttr("fixt", parseLit(`(\f f(f))(\f \g \n g(f(f)(g))(n))`)),
					rel.NewAttr("fixt", parseLit(`(\f f(f))(\f \t t :> \g \n g(f(f)(t))(n))`)),
				)),
				rel.NewAttr("log", rel.NewTuple(
					rel.NewNativeFunctionAttr("print", func(value rel.Value) rel.Value {
						log.Print(value)
						return value
					}),
				)),
			)).
			With("//./", rel.NewNativeFunction("//./", importLocalFile)).
			With("//", rel.NewNativeFunction("//", importURL))
	})
	return stdScopeVar
}

func parseLit(s string) rel.Value {
	e := MustParseString(s, "")
	v, err := e.Eval(rel.EmptyScope)
	if err != nil {
		panic(err)
	}
	return v
}

func newFloatFuncAttr(name string, f func(float64) float64) rel.Attr {
	return rel.NewNativeFunctionAttr(name, func(value rel.Value) rel.Value {
		return rel.NewNumber(f(value.(rel.Number).Float64()))
	})
}

func parseGrammar(v rel.Value) rel.Value {
	astNode := rel.ASTNodeFromValue(v).(wbnf.Branch)
	parserNode := wbnf.ToParserNode(wbnf.Core().Grammar(), astNode).(parser.Node)
	parsers := wbnf.NewFromNode(parserNode).Compile(&parserNode)
	return rel.NewNativeFunction("parse(<grammar>)", func(v rel.Value) rel.Value {
		rule := v.String()
		return rel.NewNativeFunction(fmt.Sprintf("parse(%s)", rule), func(v rel.Value) rel.Value {
			node, err := parsers.Parse(parser.Rule(rule), parser.NewScanner(v.String()))
			if err != nil {
				panic(err)
			}
			return rel.ASTNodeToValue(wbnf.FromParserNode(parsers.Grammar(), node))
		})
	})
}

func importLocalFile(v rel.Value) rel.Value {
	data, err := ioutil.ReadFile(v.String())
	if err != nil {
		panic(err)
	}
	return rel.NewString([]rune(string(data)))
}

func importURL(v rel.Value) rel.Value {
	url := v.String()
	if !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return rel.NewString([]rune(string(data)))
}