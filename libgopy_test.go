package libgopy2

import (
	"fmt"
	"testing"
)

func TestLibgopy(t *testing.T) {
	p := &Python{}
	p.Init()

	err := p.Load("libtests.test_script1")
	if err != nil {
		t.Errorf("Load failed: %v", err)
		t.FailNow()
	}
	err = p.Load("libtests.test_script2")
	if err != nil {
		t.Errorf("Load failed: %v", err)
		t.FailNow()
	}
	err = p.Load("libtests.test_script4")
	if err != nil {
		t.Errorf("Load failed: %v", err)
		t.FailNow()
	}
	err = p.Load("libtests.test_script5")
	if err != nil {
		t.Errorf("Load failed: %v", err)
		t.FailNow()
	}

	res1, err := p.Call("func6")
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	if res1 != 3.14 {
		t.Errorf("Call failed: %v != %v", res1, 3.14)
	}

	res2, err := p.Call("func1")
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	if res2 != int64(4) {
		t.Errorf("Call failed: %v != %v", res2, 4)
	}

	res3, err := p.Call("func5")
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	res3_test := []byte("hello world")
	if string(res3.([]byte)) != string(res3_test) {
		t.Errorf("Call failed: %v != %v", res3, res3_test)
	}

	res4, err := p.Call("func7", 6.5, 10.0, 9.7, 8.2)
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	if res4 != 6.5 {
		t.Errorf("Call failed: %v != %v", res4, 6.5)
	}

	res5, err := p.Call("func7", 6, 10, 9, 8)
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	if res5 != int64(6) {
		t.Errorf("Call failed: %v != %v", res5, 6)
	}

	res6, err := p.Call("func7", "Hello", "World", "Go", "Python")
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	if res6.(string) != "Hello" {
		t.Errorf("Call failed: %v != %v", res6, "Hello")
	}

	res7, err := p.Call("func7", []byte("Hello"), []byte("World"), []byte("Go"), []byte("Python"))
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	res7_test := []byte("Hello")
	if string(res7.([]byte)) != string(res7_test) {
		t.Errorf("Call failed: %v != %v", res7, res7_test)
	}

	res8, err := p.Call("func7", 6.5, 10, "Hello", []byte("World"), int64(3))
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	if res8 != 6.5 {
		t.Errorf("Call failed: %v != %v", res8, 6.5)
	}

	res9, err := p.Call("func7",
		float64(71.5), float32(3.14),
		int64(3), int32(4), int16(5), int8(6), int(10),
		uint64(3), uint32(4), uint16(5), uint8(6), uint(10),
		rune(65),
		byte('A'),
		[]uint8("World"),
		[]byte("World"),
		"Hello",
		"Hello\x00World",
		[]byte("Hello\x00World"),
	)
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	if res9 != 71.5 {
		t.Errorf("Call failed: %v != %v", res9, 71.5)
	}

	res10, err := p.Call("func_list", 777, "Hello", 1.1, []byte("World"))
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	if res10 == nil {
		t.Errorf("Call failed: %v != %v", res10, nil)
	}
	tmp := res10.([]interface{})
	if len(tmp) != 4 {
		t.Errorf("Call failed: %v != %v", len(tmp), 4)
	}
	if tmp[0] != int64(777) {
		t.Errorf("Call failed: %v != %v", tmp[0], 777)
	}
	if tmp[1] != "Hello" {
		t.Errorf("Call failed: %v != %v", tmp[1], "Hello")
	}
	if tmp[2] != 1.1 {
		t.Errorf("Call failed: %v != %v", tmp[2], 1.1)
	}
	if string(tmp[3].([]byte)) != "World" {
		t.Errorf("Call failed: %v != %v", tmp[3], "World")
	}

	res11, err := p.Call("func_dict")
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	tmp2 := fmt.Sprintf("%+v", res11)
	out := "map[key1:[1 2 3] key2:[4 5.5 <nil>] key3:map[a:1 b:2] key4:[72 101 108 108 111 44 32 87 111 114 108 100 33]]"
	if tmp2 != out {
		t.Errorf("Call failed: %v != %v", tmp2, out)
	}

	res12, err := p.Call("func_tuple", nil)
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	tmp3 := fmt.Sprintf("%+v", res12)
	out2 := "[<nil>]"
	if tmp3 != out2 {
		t.Errorf("Call failed: %v != %v", tmp3, out2)
	}

	arg := map[any]any{
		"key1": 1,
		"key2": 2.2,
		"key3": "Hello",
		"key4": []byte("World"),
		"key5": []any{1, nil, 0},
		"key6": true,
	}
	res13, err := p.Call("func_arg", arg)
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	tmp4 := fmt.Sprintf("%+v", res13)
	out3 := fmt.Sprintf("%+v", arg)
	if tmp4 != out3 {
		t.Errorf("Call failed: %v != %v", tmp4, out3)
	}

	res14, err := p.Call("func_set")
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	tmp5 := fmt.Sprintf("%+v", res14)
	out4 := "map[1:true 2:true 3:true 4:true]"
	if tmp5 != out4 {
		t.Errorf("Call failed: %v != %v", tmp5, out4)
	}

	res15, err := p.Call("create_pyobj", "Hello World!")
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	res16, err := p.Call("test_pyobj", res15)
	if err != nil {
		t.Errorf("Call failed: %v", err)
	}
	tmp6 := fmt.Sprintf("%+v", res16)
	out5 := "Hello World!"
	if tmp6 != out5 {
		t.Errorf("Call failed: %v != %v", tmp6, out5)
	}

	p.Finalize()
}

func BenchmarkLibgopy(b *testing.B) {
	p := &Python{}
	p.Init()
	err := p.Load("libtests.test_script2")
	if err != nil {
		b.Errorf("Load failed: %v", err)
		return
	}
	var res any
	for n := 0; n < b.N; n++ {
		res, _ = p.Call("func8",
			[]byte("Hello"),
			[]byte("World"),
			[]byte("Go"),
			[]byte("Python"),
			[]byte("Hello"),
			[]byte("World"),
			[]byte("Go"),
			[]byte("Python"),
		)
		if string(res.([]uint8)) != "olleH" {
			b.Errorf("Call_byte failed: %v != %v", res, "olleH")
			return
		}
	}
	p.Finalize()
}
