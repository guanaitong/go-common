package format

import "testing"

type T struct {
	A string `json:"a"`
	B int    `json:"b"`
	C []int  `json:"c"`
}

func TestAsString(t *testing.T) {
	v := map[string]interface{}{
		"a": "A",
		"b": 9,
		"c": []int{1, 2, 3},
	}
	if s := AsString(v); s == "" {
		t.Error("Format failure")
	} else {
		t.Log(s)
	}
}

func TestAsJson(t *testing.T) {
	d := `{"a":"A","b":9,"c":[1,2,3]}`
	v := &T{}
	if err := AsJson(d, v); err != nil {
		t.Error(err)
	} else {
		t.Log(v)
	}
}

func TestInt32Ptr(t *testing.T) {
	t.Log(Int32Ptr(1))
}
