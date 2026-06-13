package jsonx

import (
	"encoding/json"
	"testing"

	"github.com/oylshe1314/framework/util"
)

type AT struct {
	Val1 string `json:"val1"`
}

type PAT struct {
	Val2 string `json:"val2"`
}

type IAT struct {
	Val3 string `json:"val3"`
}

type ST struct {
	Val4 string `json:"val4"`
}

type TT struct {
	Bool1 bool `json:"bool1"`
	Bool2 bool `json:"-"`

	Int   int   `json:"int"`
	Int8  int8  `json:"int8"`
	Int16 int16 `json:"int16"`
	Int32 int32 `json:"int32"`
	Int64 int64 `json:"int64"`

	Uint   int   `json:"uint"`
	Uint8  int8  `json:"uint8"`
	Uint16 int16 `json:"uint16"`
	Uint32 int32 `json:"uint32"`
	Uint64 int64 `json:"uint64"`

	Uintptr uintptr `json:"uintptr"`

	Float32 float32 `json:"float32"`
	Float64 float64 `json:"float64"`

	IntAry [2]int    `json:"intAry"`
	StrAry [3]string `json:"strAry"`

	Any1 any         `json:"any1"`
	Any2 interface{} `json:"any2"`

	Map1 map[string]any `json:"map1"`
	Map2 map[string]int `json:"map2"`
	Map3 map[int]string `json:"map3"`
	Map4 map[string]ST  `json:"map4"`

	Slice1 []string `json:"slice1"`
	Slice2 []int    `json:"slice2"`
	Slice3 []ST     `json:"slice3"`
	Slice4 []*ST    `json:"slice4"`

	Str1 string `json:"str1"`
	Str2 string `json:"-"`

	AT   `json:"at"`
	*PAT `json:"pat"`
	*IAT `json:",inline"`

	St1 ST  `json:"st1"`
	St2 *ST `json:"st2"`
	St3 *ST `json:",inline"`
}

func TestObjectRead(t *testing.T) {
	var str = `
{
  "bool1": true,
  "bool2": true,
  "int": 1,
  "int8": 2,
  "int16": 3,
  "int32": 4,
  "int64": 5,
  "uint": 6,
  "uint8": 7,
  "uint16": 8,
  "uint32": 9,
  "uint64": 10,
  "uintptr": 11,
  "float32": 12.5,
  "float64": 13.5,
  "intAry": [
    14,
    15,
    16,
    17,
    18
  ],
  "strAry": [
    "str19",
    "str20",
    "str21",
    "str22",
    "str23"
  ],
  "any1": 24,
  "any2": "str25",
  "map1": {
    "k26": "str26",
    "k27": 27
  },
  "map2": {
    "k28": 28,
    "k29": 29
  },
  "map3": {
    "30": "str30",
    "31": "str31"
  },
  "map4": {
    "k32": {
      "val4": "str32"
    },
    "k33": {
      "val4": "str33"
    }
  },
  "slice1": [
    "str34",
    "str35",
    "str36"
  ],
  "slice2": [
    37,
    38
  ],
  "slice3": [
    {
      "val4": "str39"
    },
    {
      "val4": "str40"
    }
  ],
  "slice4": [
    {
      "val4": "str31"
    },
    {
      "val4": "str42"
    }
  ],
  "str1": "str43",
  "str2": "str44",
  "at": {
    "val1": "str45"
  },
  "pat": {
    "val2": "str46"
  },
  "iat": {
    "val3": "str47"
  },
  "val3": "str48",
  "st1": {
    "val4": "str49"
  },
  "st2": {
    "val4": "str50"
  },
  "st3": {
    "val4": "str51"
  },
  "val4": "str52"
}
`

	var obj map[string]any

	var err = json.Unmarshal([]byte(str), &obj)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(obj)

	var st = TT{Slice1: make([]string, 2), Slice2: make([]int, 5)}

	t.Logf("%+v", st)

	var b = util.NowMicro()
	err = Object(obj).ToStruct(&st)
	if err != nil {
		t.Error(err)
		return
	}
	var e = util.NowMicro()

	t.Logf("%+v", st)
	t.Log("Time: ", e-b, "us")
}
