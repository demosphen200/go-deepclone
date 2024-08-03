package deepclone

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func commonTest[
	R any,
](t *testing.T, initSrc R, initDst R) {
	var src = initSrc
	var dst = initDst

	err := DeepCloneAny(&src, &dst)
	if assert.NoError(t, err) {
		assert.Equal(t, src, dst)
	} else {
		assert.Fail(t, fmt.Sprintf("DeepCloneAny failed on %s", reflect.TypeOf(src).String()))
	}
}

func TestDeepCloneAny_int(t *testing.T) {
	commonTest[int](t, 123, 0)
	commonTest[int8](t, 123, 0)
	commonTest[int16](t, 123, 0)
	commonTest[int32](t, 123, 0)
	commonTest[int64](t, 123, 0)
	commonTest[uint](t, 123, 0)
	commonTest[uint8](t, 123, 0)
	commonTest[uint16](t, 123, 0)
	commonTest[uint32](t, 123, 0)
	commonTest[uint64](t, 123, 0)
	commonTest[uintptr](t, 123, 0)
}

func TestDeepCloneAny_string(t *testing.T) {
	commonTest(t, "123", "")
}

func TestDeepCloneAny_slice(t *testing.T) {
	commonTest(t, []int{1, 2, 3}, nil)
	commonTest(t, nil, []int{1, 2, 3})

	src := []int{1, 2, 3}
	expected := []int{1, 2, 3}
	var dst []int

	err := DeepCloneAny(&src, &dst)
	src[1] = 222

	if assert.NoError(t, err) {
		assert.Equal(t, expected, dst)
	} else {
		assert.Fail(t, fmt.Sprintf("DeepCloneAny failed on %s", reflect.TypeOf(src).String()))
	}
}

func TestDeepCloneAny_map(t *testing.T) {
	commonTest(t, map[string]string{
		"a": "aaa",
		"b": "bbb",
	}, nil)

	commonTest(t, nil, map[string]string{
		"a": "aaa",
		"b": "bbb",
	})

	src := map[string]string{
		"a": "aaa",
		"b": "bbb",
	}
	expected := map[string]string{
		"a": "aaa",
		"b": "bbb",
	}
	var dst map[string]string

	err := DeepCloneAny(&src, &dst)
	clear(src)

	if assert.NoError(t, err) {
		assert.Equal(t, expected, dst)
	} else {
		assert.Fail(t, fmt.Sprintf("DeepCloneAny failed on %s", reflect.TypeOf(src).String()))
	}
}

type testStruct struct {
	A int
	B string
}

type testStruct2 struct {
	A int
	B string
	c int
}

func TestDeepCloneAny_struct(t *testing.T) {

	commonTest(t, testStruct{
		A: 1,
		B: "BBB",
	}, testStruct{})

	commonTest(t, &testStruct{
		A: 1,
		B: "BBB",
	}, &testStruct{})

	commonTest(t, nil, &testStruct{
		A: 1,
		B: "BBB",
	})
}

func TestDeepCloneAny_any_int(t *testing.T) {
	var initSrc any = 123
	var initDst any = ""
	commonTest(t, initSrc, initDst)
}

func TestDeepCloneAny_any_untyped_nil(t *testing.T) {
	var initSrc any = nil
	var initDst any = ""
	commonTest(t, initSrc, initDst)
}

func TestDeepCloneAny_any_typed_nil(t *testing.T) {
	var ptr *int
	var initSrc any = ptr
	var initDst any = ""
	commonTest(t, initSrc, initDst)
}

func TestDeepCloneAny_pointer(t *testing.T) {
	val1 := 123
	var initSrc = &val1
	var initDst *int
	commonTest(t, initSrc, initDst)

	commonTest(t, nil, initSrc)
}

func TestDeepCloneAny_pointer_shouldCreateNewInstance(t *testing.T) {
	val1 := 123
	val2 := 123

	src := &val1
	expected := &val2
	dst := &val1

	err := DeepCloneAny(&src, &dst)

	*src = 11111

	if assert.NoError(t, err) {
		assert.Equal(t, expected, dst)
	} else {
		assert.Fail(t, fmt.Sprintf("DeepCloneAny failed on %s", reflect.TypeOf(src).String()))
	}
}

func TestDeepClone_int(t *testing.T) {
	var src = 123
	var dst = 0

	err := DeepClone(&src, &dst)
	if assert.NoError(t, err) {
		assert.Equal(t, src, dst)
	} else {
		assert.Fail(t, fmt.Sprintf("DeepCloneAny failed on %s", reflect.TypeOf(src).String()))
	}
}

func TestDeepClone_time(t *testing.T) {
	commonTest(t, time.Now(), time.Time{})
}

func TestDeepClone_ignoresNotExportedInUnknownStructs(t *testing.T) {
	src := &testStruct2{
		A: 1,
		B: "BBB",
		c: 123,
	}
	expected := &testStruct2{
		A: 1,
		B: "BBB",
		c: 0,
	}
	dst := &testStruct2{}
	err := DeepClone(&src, &dst)
	if assert.NoError(t, err) {
		assert.Equal(t, expected, dst)
	} else {
		assert.Fail(t, fmt.Sprintf("DeepCloneAny failed on %s", reflect.TypeOf(src).String()))
	}

}

func TestDeepClone_UsesKnownStructCloners(t *testing.T) {

	RegisterStructCloner(
		func(src *testStruct2) (*testStruct2, error) {
			newValue := *src
			return &newValue, nil
		},
	)

	src := &testStruct2{
		A: 1,
		B: "BBB",
		c: 123,
	}
	dst := &testStruct2{}
	err := DeepClone(&src, &dst)
	if assert.NoError(t, err) {
		assert.Equal(t, src, dst)
	} else {
		assert.Fail(t, fmt.Sprintf("DeepCloneAny failed on %s", reflect.TypeOf(src).String()))
	}

}
