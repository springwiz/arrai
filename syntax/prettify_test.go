package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFmtPrettyForDict(t *testing.T) { //nolint:dupl
	t.Parallel()
	simpleDict, err := EvaluateExpr(".", `{'c':3, 'a':1, 'b':'b2'}`)
	assert.Nil(t, err)
	str, err := PrettifyString(simpleDict, 0)
	assert.Nil(t, err)
	assert.Equal(t, `{
  a: 1,
  b: 'b2',
  c: 3
}`, str)

	complexDict, err := EvaluateExpr(".", `{'a':1, 'c':(d:11, e:'e22', f:{111, 222}), 'b':2}`)
	assert.Nil(t, err)
	str, err = PrettifyString(complexDict, 0)
	assert.Nil(t, err)
	assert.Equal(t,
		`{
  a: 1,
  b: 2,
  c: (
    d: 11,
    e: 'e22',
    f: {
      111,
      222
    }
  )
}`,
		str)
}

func TestFmtPrettyForSet(t *testing.T) { //nolint:dupl
	t.Parallel()
	simpleSet, err := EvaluateExpr(".", `{26, '24', 25}`)
	assert.Nil(t, err)
	str, err := PrettifyString(simpleSet, 0)
	assert.Nil(t, err)
	assert.Equal(t, `{
  25,
  26,
  '24'
}`, str)

	complexSet, err := EvaluateExpr(".", `{(a: 'a1'),(b: 2),(c:{11,22,'33'})}`)
	assert.Nil(t, err)
	str, err = PrettifyString(complexSet, 0)
	assert.Nil(t, err)
	assert.Equal(t,
		`{
  (
    a: 'a1'
  ),
  (
    b: 2
  ),
  (
    c: {
      11,
      22,
      '33'
    }
  )
}`,
		str)
}

func TestFmtPrettyForTuple(t *testing.T) { //nolint:dupl
	t.Parallel()
	simpleSet, err := EvaluateExpr(".", `(a:1, c:'c3', b:2)`)
	assert.Nil(t, err)
	str, err := PrettifyString(simpleSet, 0)
	assert.Nil(t, err)
	assert.Equal(t, `(
  a: 1,
  b: 2,
  c: 'c3'
)`, str)

	complexTuple, err := EvaluateExpr(".", `(a:1, b:(d:11, e:12, f:{1, '2'}), c:'3')`)
	assert.Nil(t, err)
	str, err = PrettifyString(complexTuple, 0)
	assert.Nil(t, err)
	assert.Equal(t,
		`(
  a: 1,
  b: (
    d: 11,
    e: 12,
    f: {
      1,
      '2'
    }
  ),
  c: '3'
)`,
		str)
}

func TestFmtPrettyForArray(t *testing.T) { //nolint:dupl
	t.Parallel()
	array, err := EvaluateExpr(".", `[1, 2, 3, 5, 6, 4, 10]`)
	assert.Nil(t, err)
	str, err := PrettifyString(array, 0)
	assert.Nil(t, err)
	assert.Equal(t, "[1, 2, 3, 5, 6, 4, 10]", str)
}
