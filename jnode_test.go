package jpath

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"

	"github.com/bmizerany/assert"
)

func TestSimplejson(t *testing.T) {
	var ok bool
	var err error

	js, err := New([]byte(`{
		"test": {
			"string_list": ["asdf", "ghjk", "zxcv"],
			"string_list_null": ["abc", null, "efg"],
			"list": [1, "2", 3],
			"listwithsubs": [{"subkeyone": 1},
			{"subkeytwo": 2, "subkeythree": 3}],
			"int": 10,
			"float": 5.150,
			"string": "simplejson",
			"bool": true,
			"sub_obj": {"a": 1}
		}
	}`))

	assert.NotEqual(t, nil, js)
	assert.Equal(t, nil, err)

	_, ok = js.CheckGet("test")
	assert.Equal(t, true, ok)

	_, ok = js.CheckGet("missing_key")
	assert.Equal(t, false, ok)

	aws := js.Get("test").Get("listwithsubs")
	assert.NotEqual(t, nil, aws)
	var awsval int
	awsval, _ = aws.Get(0).Get("subkeyone").CheckInt()
	assert.Equal(t, 1, awsval)
	awsval, _ = aws.Get(1).Get("subkeytwo").CheckInt()
	assert.Equal(t, 2, awsval)
	awsval, _ = aws.Get(1).Get("subkeythree").CheckInt()
	assert.Equal(t, 3, awsval)

	i, _ := js.Get("test").Get("int").CheckInt()
	assert.Equal(t, 10, i)

	f, _ := js.Get("test").Get("float").CheckFloat64()
	assert.Equal(t, 5.150, f)

	s, _ := js.Get("test").Get("string").CheckString()
	assert.Equal(t, "simplejson", s)

	b, _ := js.Get("test").Get("bool").CheckBool()
	assert.Equal(t, true, b)

	mi := js.Get("test").Get("int").Int()
	assert.Equal(t, 10, mi)

	mi2 := js.Get("test").Get("missing_int").Int(5150)
	assert.Equal(t, 5150, mi2)

	ms := js.Get("test").Get("string").String()
	assert.Equal(t, "simplejson", ms)

	ms2 := js.Get("test").Get("missing_string").String("fyea")
	assert.Equal(t, "fyea", ms2)

	ma2 := js.Get("test").Get("missing_list").List([]any{"1", 2, "3"})
	assert.Equal(t, ma2, []any{"1", 2, "3"})

	mm2 := js.Get("test").Get("missing_map").Map(map[string]any{"found": false})
	assert.Equal(t, mm2, map[string]any{"found": false})

	gp, _ := js.Get("test", "string").CheckString()
	assert.Equal(t, "simplejson", gp)

	gp2, _ := js.Get("test", "int").CheckInt()
	assert.Equal(t, 10, gp2)

	gpa, _ := js.Get("test", "string_list", 0).CheckString()
	assert.Equal(t, "asdf", gpa)

	gpa2, _ := js.Get("test", "listwithsubs", 1, "subkeythree").CheckInt()
	assert.Equal(t, 3, gpa2)

	jm, ok := js.Get("test").CheckNodeMap()
	assert.Equal(t, ok, true)
	jmbool, _ := jm["bool"].CheckBool()
	assert.Equal(t, true, jmbool)

	ja, ok := js.Get("test", "string_list").CheckNodeList()
	assert.Equal(t, ok, true)
	jastr, _ := ja[0].CheckString()
	assert.Equal(t, "asdf", jastr)

	assert.Equal(t, js.Get("test").Get("bool").Bool(), true)

	js.Set("float2", 300.0)
	assert.Equal(t, js.Get("float2").Float64(), 300.0)

	js.Set("test2", "setTest")
	assert.Equal(t, "setTest", js.Get("test2").String())

	js.DelKey("test2")
	assert.NotEqual(t, "setTest", js.Get("test2").String())

	js.Get("test").Get("sub_obj").Set("a", 2)
	assert.Equal(t, 2, js.Get("test").Get("sub_obj").Get("a").Int())

	js.Get("test", "sub_obj").Set("a", 3)
	assert.Equal(t, 3, js.Get("test", "sub_obj", "a").Int())

	jmm := js.Get("missing_map").NodeMap(NodeMap{"js1": js})
	assert.Equal(t, js, jmm["js1"])

	jma := js.Get("missing_list").NodeList(NodeList{js})
	assert.Equal(t, js, jma[0])
}

func TestStdlibInterfaces(t *testing.T) {
	val := new(struct {
		Name   string `json:"name"`
		Params *Node  `json:"params"`
	})
	val2 := new(struct {
		Name   string `json:"name"`
		Params *Node  `json:"params"`
	})

	raw := `{"name":"myobject","params":{"string":"simplejson"}}`

	assert.Equal(t, nil, json.Unmarshal([]byte(raw), val))

	assert.Equal(t, "myobject", val.Name)
	assert.NotEqual(t, nil, val.Params.data)
	s, _ := val.Params.Get("string").CheckString()
	assert.Equal(t, "simplejson", s)

	p, err := json.Marshal(val)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, json.Unmarshal(p, val2))
	assert.Equal(t, val, val2) // stable
}

func TestSet(t *testing.T) {
	js, err := New([]byte(`{}`))
	assert.Equal(t, nil, err)

	js.Set("baz", "bing")

	s, ok := js.Get("baz").CheckString()
	assert.Equal(t, true, ok)
	assert.Equal(t, "bing", s)
}

func TestReplace(t *testing.T) {
	js, err := New([]byte(`{}`))
	assert.Equal(t, nil, err)

	err = js.UnmarshalJSON([]byte(`{"baz":"bing"}`))
	assert.Equal(t, nil, err)

	s, ok := js.Get("baz").CheckString()
	assert.Equal(t, true, ok)
	assert.Equal(t, "bing", s)
}

func TestSetBranch(t *testing.T) {
	js, err := New([]byte(`{}`))
	assert.Equal(t, nil, err)

	js.SetBranch([]string{"foo", "bar"}, "baz")

	s, ok := js.Get("foo", "bar").CheckString()
	assert.Equal(t, true, ok)
	assert.Equal(t, "baz", s)
}

func TestSetBranchNoPath(t *testing.T) {
	js, err := New([]byte(`{"some":"data","some_number":1.0,"some_bool":false}`))
	assert.Equal(t, nil, err)

	f := js.Get("some_number").Float64(99.0)
	assert.Equal(t, f, 1.0)

	js.SetBranch([]string{}, map[string]any{"foo": "bar"})

	s, ok := js.Get("foo").CheckString()
	assert.Equal(t, true, ok)
	assert.Equal(t, "bar", s)

	f = js.Get("some_number").Float64(99.0)
	assert.Equal(t, f, 99.0)
}

func TestPathWillAugmentExisting(t *testing.T) {
	js, err := New([]byte(`{"this":{"a":"aa","b":"bb","c":"cc"}}`))
	assert.Equal(t, nil, err)

	js.SetBranch([]string{"this", "d"}, "dd")

	cases := []struct {
		path    []any
		outcome string
	}{
		{
			path:    []any{"this", "a"},
			outcome: "aa",
		},
		{
			path:    []any{"this", "b"},
			outcome: "bb",
		},
		{
			path:    []any{"this", "c"},
			outcome: "cc",
		},
		{
			path:    []any{"this", "d"},
			outcome: "dd",
		},
	}

	for _, tc := range cases {
		s, ok := js.Get(tc.path...).CheckString()
		assert.Equal(t, true, ok)
		assert.Equal(t, tc.outcome, s)
	}
}

func TestPathWillOverwriteExisting(t *testing.T) {
	// notice how "a" is 0.1 - but then we'll try to set at path a, foo
	js, err := New([]byte(`{"this":{"a":0.1,"b":"bb","c":"cc"}}`))
	assert.Equal(t, nil, err)

	js.SetBranch([]string{"this", "a", "foo"}, "bar")

	s, ok := js.Get("this", "a", "foo").CheckString()
	assert.Equal(t, true, ok)
	assert.Equal(t, "bar", s)
}

func TestNewFromReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`{
		"test": {
			"list": [1, "2", 3],
			"listwithsubs": [
				{"subkeyone": 1},
				{"subkeytwo": 2, "subkeythree": 3}
			],
			"bignum": 8000000000
		}
	}`))
	js, err := NewFromReader(buf)

	//Standard Test Case
	assert.NotEqual(t, nil, js)
	assert.Equal(t, nil, err)

	arr, _ := js.Get("test").Get("list").CheckList()
	assert.NotEqual(t, nil, arr)
	for i, v := range arr {
		var iv int
		switch v := v.(type) {
		case float64:
			iv = int(v)
		case string:
			iv, _ = strconv.Atoi(v)
		}
		assert.Equal(t, i+1, iv)
	}

	ma := js.Get("test").Get("list").List()
	assert.Equal(t, ma, []any{float64(1), "2", float64(3)})

	mm := js.Get("test").Get("listwithsubs").Get(0).Map()
	assert.Equal(t, mm, map[string]any{"subkeyone": float64(1)})

	assert.Equal(t, js.Get("test").Get("bignum").Int64(), int64(8000000000))
}

func TestSimplejson2(t *testing.T) {
	js, err := New([]byte(`{
		"test": {
			"list": [1, "2", 3],
			"listwithsubs": [
				{"subkeyone": 1},
				{"subkeytwo": 2, "subkeythree": 3}
			],
			"bignum": 8000000000
		}
	}`))

	assert.NotEqual(t, nil, js)
	assert.Equal(t, nil, err)

	arr, _ := js.Get("test").Get("list").CheckList()
	assert.NotEqual(t, nil, arr)
	for i, v := range arr {
		var iv int
		switch v := v.(type) {
		case float64:
			iv = int(v)
		case string:
			iv, _ = strconv.Atoi(v)
		}
		assert.Equal(t, i+1, iv)
	}

	ma := js.Get("test").Get("list").List()
	assert.Equal(t, ma, []any{float64(1), "2", float64(3)})

	mm := js.Get("test").Get("listwithsubs").Get(0).Map()
	assert.Equal(t, mm, map[string]any{"subkeyone": float64(1)})

	assert.Equal(t, js.Get("test").Get("bignum").Int64(), int64(8000000000))
}

// Test simple JSON path expressions
func TestPath(t *testing.T) {
	js, _ := New([]byte(`{
		"test": {
			"string_list": ["asdf", "ghjk", "zxcv"],
			"string_list_null": ["abc", null, "efg"],
			"list": [1, "2", 3],
			"listwithsubs": [{"subkeyone": 1},
			{"subkeytwo": 2, "subkeythree": 3}],
			"int": 10,
			"float": 5.150,
			"string": "simplejson",
			"bool": true,
			"sub_obj": {"a": 1}
		}
	}`))
	n := js.GetNode("x")
	assert.NotEqual(t, n, NilNode)
	n = js.GetNode("x.test")
	assert.NotEqual(t, n, NilNode)
	n = js.GetNode("x.test.string_list")
	assert.NotEqual(t, n, NilNode)
	n = js.GetNode("x.test.string_list[1]")
	assert.NotEqual(t, n, NilNode)
	assert.Equal(t, "ghjk", n.String())
	n = js.GetNode(".test.listwithsubs.[1].subkeytwo")
	assert.Equal(t, 2, n.Int())
	v := js.GetNode("x.test.sub_obj.a").Int()
	assert.Equal(t, 1, v)
}

func TestAddEmpty(t *testing.T) {
	// Set up a JSON document, data to be added and the correct output to compare with
	documentJSON := []byte("")
	someJSON := []byte(`{"x":"7", "y":"15"}`)
	// Expect brackets to be created automatically
	correctNode, err := New(badd([]byte("["), someJSON, []byte("]")))
	assert.Equal(t, nil, err)
	correctJSON, err := correctNode.JSON()
	assert.Equal(t, nil, err)

	// Create a JSON document, add the data and compare
	document, err := New(documentJSON)
	assert.NotEqual(t, nil, document)
	assert.Equal(t, nil, err)
	err = document.AddJSON("x", someJSON)
	assert.Equal(t, nil, err)
	newJSON, err := document.JSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, bytes.Equal(newJSON, correctJSON))
}

func TestAddEmpty2(t *testing.T) {
	// Set up an empty JSON document, data to be added and the correct output to compare with
	documentJSON := []byte("[]")
	someJSON := []byte(`{"x":"7", "y":"15"}`)
	correctNode, err := New(badd([]byte("["), someJSON, []byte("]")))
	assert.Equal(t, nil, err)
	correctJSON, err := correctNode.JSON()
	assert.Equal(t, nil, err)

	// Create a JSON document, add the data and compare
	document, err := New(documentJSON)
	assert.NotEqual(t, nil, document)
	assert.Equal(t, nil, err)
	err = document.AddJSON("x", someJSON)
	assert.Equal(t, nil, err)
	newJSON, err := document.JSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, bytes.Equal(newJSON, correctJSON))
}

func TestAdd(t *testing.T) {
	// Set up a JSON document, data to be added and the correct output to compare with
	contentJSON := []byte(`{"x":"2", "y":"3"}`)
	documentJSON := badd([]byte("["), contentJSON, []byte("]"))
	someJSON := []byte(`{"x":"7", "y":"15"}`)
	correctNode, err := New(badd([]byte("["), contentJSON, []byte(","), someJSON, []byte("]")))
	assert.Equal(t, nil, err)
	correctJSON, err := correctNode.JSON()
	assert.Equal(t, nil, err)

	// Create a JSON document, add the data and compare
	document, err := New(documentJSON)
	assert.NotEqual(t, nil, document)
	assert.Equal(t, nil, err)
	err = document.AddJSON("x", someJSON)
	assert.Equal(t, nil, err)
	newJSON, err := document.JSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, bytes.Equal(newJSON, correctJSON))
}

func TestAdd2(t *testing.T) {
	// Set up a JSON document, data to be added and the correct output to compare with
	contentJSON := []byte(`{"x":"2", "y":"3"}`)
	documentJSON := badd([]byte("["), contentJSON, []byte(","), contentJSON, []byte("]"))
	someJSON := []byte(`{"x":"7", "y":"15"}`)
	correctNode, err := New(badd([]byte("["), contentJSON, []byte(","), contentJSON, []byte(","), someJSON, []byte("]")))
	assert.Equal(t, nil, err)
	correctJSON, err := correctNode.JSON()
	assert.Equal(t, nil, err)

	// Create a JSON document, add the data and compare
	document, err := New(documentJSON)
	assert.NotEqual(t, nil, document)
	assert.Equal(t, nil, err)
	err = document.AddJSON("x", someJSON)
	assert.Equal(t, nil, err)
	newJSON, err := document.JSON()
	assert.Equal(t, nil, err)
	assert.Equal(t, true, bytes.Equal(newJSON, correctJSON))
}

// Test removal of keys
func TestDel(t *testing.T) {
	js, err := New([]byte(`{
		"test": {
			"string_list": ["asdf", "ghjk", "zxcv"],
			"list": [1, "2", 3]
		}
	}`))
	assert.NotEqual(t, nil, js, err)
	assert.Equal(t, nil, err)

	correctNode, err := New([]byte(`{
		"test": {
			"list": [1, "2", 3]
		}
	}`))
	assert.NotEqual(t, correctNode, nil)
	assert.Equal(t, err, nil)

	correctJSON, err := correctNode.JSON()
	assert.Equal(t, err, nil)

	err = js.DelKey("blublublu")
	assert.NotEqual(t, err, nil)

	err = js.DelKey("x.test.string_list")
	assert.Equal(t, err, nil)

	newJSON, err := js.JSON()
	assert.Equal(t, err, nil)

	assert.Equal(t, true, bytes.Equal(newJSON, correctJSON))
}

func TestSetPath(t *testing.T) {
	js, err := New([]byte(`{"user":{"name":"Alice"}}`))
	assert.Equal(t, nil, err)

	js.SetPath("user.name", "Bob")
	assert.Equal(t, "Bob", js.Get("user", "name").String())

	js.SetPath("x.user.name", "Charlie")
	assert.Equal(t, "Charlie", js.Get("user", "name").String())

	js.SetPath(".user.name", "Dave")
	assert.Equal(t, "Dave", js.Get("user", "name").String())

	js.SetPath("settings.theme.color", "blue")
	assert.Equal(t, "blue", js.Get("settings", "theme", "color").String())
}

func TestStringValue(t *testing.T) {
	js, err := New([]byte(`{"name":"Algernon","count":42,"active":true,"items":[1,2,3],"meta":{"k":"v"}}`))
	assert.Equal(t, nil, err)

	assert.Equal(t, "Algernon", js.Get("name").StringValue())
	assert.Equal(t, "42", js.Get("count").StringValue())
	assert.Equal(t, "true", js.Get("active").StringValue())
	assert.Equal(t, "[1,2,3]", js.Get("items").StringValue())
	assert.Equal(t, `{"k":"v"}`, js.Get("meta").StringValue())
	assert.Equal(t, "", js.Get("missing").StringValue())
	assert.Equal(t, "", NilNode.StringValue())
}

func TestGetNodePathEquivalence(t *testing.T) {
	js, err := New([]byte(`{"user":{"name":"Alice","scores":[10,20]}}`))
	assert.Equal(t, nil, err)

	assert.Equal(t, "Alice", js.GetNode("x.user.name").String())
	assert.Equal(t, "Alice", js.GetNode(".user.name").String())
	assert.Equal(t, "Alice", js.GetNode("user.name").String())

	n := js.GetNode("user")
	_, ok := n.CheckMap()
	assert.Equal(t, true, ok)
}

func TestGetNodeBracketPaths(t *testing.T) {
	js, err := New([]byte(`[{"id":"first"},{"id":"second"}]`))
	assert.Equal(t, nil, err)

	// Bare bracket path without "x" prefix
	assert.Equal(t, "first", js.GetNode("[0].id").String())
	assert.Equal(t, "second", js.GetNode("[1].id").String())

	// With "x" prefix
	assert.Equal(t, "first", js.GetNode("x[0].id").String())

	// Bracket path on a nested list
	js2, err := New([]byte(`{"items":[{"name":"a"},{"name":"b"}]}`))
	assert.Equal(t, nil, err)

	assert.Equal(t, "a", js2.GetNode("items[0].name").String())
	assert.Equal(t, "b", js2.GetNode("x.items[1].name").String())
	assert.Equal(t, "b", js2.GetNode(".items[1].name").String())
}

func TestGetIndexNegative(t *testing.T) {
	js, err := New([]byte(`{"list":[1,2,3]}`))
	assert.Equal(t, nil, err)

	// Negative index should not panic, should return NilNode
	n := js.Get("list").Get(-1)
	assert.Equal(t, NilNode, n)
}
