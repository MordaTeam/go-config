package config

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestMergeLeft(t *testing.T) {
	type testBar struct {
		F1 int
		F2 string
	}

	type testFoo struct {
		F1  bool
		F2  int
		F3  int8
		F4  int16
		F5  int32
		F6  int64
		F7  uint
		F8  uint16
		F9  uint32
		F10 uint64
		F11 float32
		F12 float64
		F13 complex64
		F14 complex128
		F15 string
		F16 [4]byte
		F17 map[string]string
		F18 []byte
		F19 *int
		F20 unsafe.Pointer
		F21 any

		Nested testBar

		ignored int
	}

	type testBuz struct {
		NestedPtr *testBar
	}

	num := 666
	bar := testBar{
		F1: 11,
		F2: "12",
	}

	testCases := []struct {
		Name string
		Test func(t *testing.T)
	}{
		{
			Name: "int",
			Test: mergeWith(0, 1, 1),
		},
		{
			Name: "string",
			Test: mergeWith("", "def", "def"),
		},
		{
			Name: "bool",
			Test: mergeWith(false, true, true),
		},
		{
			Name: "[]string",
			Test: mergeWith(nil, []string{"1", "2", "3"}, []string{"1", "2", "3"}),
		},
		{
			Name: "struct",
			Test: mergeWith(
				testFoo{
					ignored: 222,
					F1:      false,
					F2:      0,
					F3:      0,
					F4:      0,
					F5:      0,
					F6:      0,
					F7:      0,
					F8:      0,
					F9:      0,
					F10:     0,
					F11:     0.0,
					F12:     0.0,
					F13:     complex64(0.0),
					F14:     complex128(0.0),
					F15:     "",
					F16:     [4]byte{},
					F17:     nil,
					F18:     nil,
					F19:     nil,
					F20:     unsafe.Pointer(nil),
					F21:     nil,
					Nested: testBar{
						F1: 0,
						F2: "hello",
					},
				},
				testFoo{
					ignored: 111,
					F1:      true,
					F2:      1,
					F3:      2,
					F4:      3,
					F5:      4,
					F6:      5,
					F7:      6,
					F8:      7,
					F9:      8,
					F10:     9,
					F11:     10.0,
					F12:     11.0,
					F13:     complex64(12.0),
					F14:     complex128(13.0),
					F15:     "14",
					F16:     [4]byte{15, 16, 17, 18},
					F17:     map[string]string{},
					F18:     []byte{},
					F19:     &num,
					F20:     unsafe.Pointer(&num),
					F21:     21,
					Nested: testBar{
						F1: 22,
						F2: "23",
					},
				},
				testFoo{
					ignored: 222,
					F1:      true,
					F2:      1,
					F3:      2,
					F4:      3,
					F5:      4,
					F6:      5,
					F7:      6,
					F8:      7,
					F9:      8,
					F10:     9,
					F11:     10.0,
					F12:     11.0,
					F13:     complex64(12.0),
					F14:     complex128(13.0),
					F15:     "14",
					F16:     [4]byte{15, 16, 17, 18},
					F17:     map[string]string{},
					F18:     []byte{},
					F19:     &num,
					F20:     unsafe.Pointer(&num),
					F21:     21,
					Nested: testBar{
						F1: 22,
						F2: "hello",
					},
				},
			),
		},
		{
			Name: "struct{ptr *struct} left=nil",
			Test: mergeWith(
				testBuz{
					NestedPtr: nil,
				},
				testBuz{
					NestedPtr: &bar,
				},
				testBuz{
					NestedPtr: &bar,
				},
			),
		},
		{
			Name: "struct{ptr *struct} right=nil",
			Test: mergeWith(
				testBuz{
					NestedPtr: &bar,
				},
				testBuz{
					NestedPtr: nil,
				},
				testBuz{
					NestedPtr: &bar,
				},
			),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, testCase.Test)
	}
}

func mergeWith[T any](left, right, expected T) (f func(t *testing.T)) {
	return func(t *testing.T) {
		assert.Equal(t, expected, mergeLeft(left, right))
	}
}
