package libgopy

import "testing"

// WARNING: This tests can panic if called mutiple Init and/or Finalize

// func TestLibgopy(t *testing.T) {
// 	Init()

// 	err := Load("tests.test_script1")
// 	if err != nil {
// 		t.Errorf("Load failed: %v", err)
// 		t.FailNow()
// 	}
// 	err = Load("tests.test_script2")
// 	if err != nil {
// 		t.Errorf("Load failed: %v", err)
// 		t.FailNow()
// 	}

// 	res1 := Call_f64("func6")
// 	if res1 != 3.14 {
// 		t.Errorf("Call_f64 failed: %v != %v", res1, 3.14)
// 	}
// 	res2 := Call_i64("func1")
// 	if res2 != 4 {
// 		t.Errorf("Call_i64 failed: %v != %v", res2, 4)
// 	}
// 	res3 := Call_byte("func5")
// 	res3_test := []byte("hello world")
// 	if string(res3) != string(res3_test) {
// 		t.Errorf("Call_byte failed: %v != %v", res3, res3_test)
// 	}

// 	res4 := Call_f64("func7", 6.5, 10.0, 9.7, 8.2)
// 	if res4 != 6.5 {
// 		t.Errorf("Call_f64 failed: %v != %v", res4, 6.5)
// 	}
// 	res5 := Call_i64("func7", 6, 10, 9, 8)
// 	if res5 != 6 {
// 		t.Errorf("Call_i64 failed: %v != %v", res5, 6)
// 	}
// 	res6 := Call_byte("func7", "Hello", "World", "Go", "Python")
// 	res6_test := []byte("Hello")
// 	if string(res6) != string(res6_test) {
// 		t.Errorf("Call_byte failed: %v != %v", res6, res6_test)
// 	}
// 	res7 := Call_byte("func7", []byte("Hello"), []byte("World"), []byte("Go"), []byte("Python"))
// 	res7_test := []byte("Hello")
// 	if string(res7) != string(res7_test) {
// 		t.Errorf("Call_byte failed: %v != %v", res7, res7_test)
// 	}
// 	res8 := Call_f64("func7", 6.5, 10, "Hello", []byte("World"), int64(3))
// 	if res8 != 6.5 {
// 		t.Errorf("Call_f64 failed: %v != %v", res8, 6.5)
// 	}

// 	Call("func7",
// 		float64(71.5), float32(3.14),
// 		int64(3), int32(4), int16(5), int8(6), int(10),
// 		uint64(3), uint32(4), uint16(5), uint8(6), uint(10),
// 		uintptr(65), // Unknown type: uintptr
// 		rune(65),
// 		byte('A'),
// 		[]uint8("World"),
// 		[]byte("World"),
// 		"Hello",
// 		"Hello\x00World",
// 		[]byte("Hello\x00World"),
// 	)
// 	Finalize()
// }

func BenchmarkLibgopy(b *testing.B) {
	var res []byte
	for n := 0; n < b.N; n++ {
		res, _ = Call("func8", []byte("Hello"), []byte("World"), []byte("Go"), []byte("Python"), []byte("Hello"), []byte("World"), []byte("Go"), []byte("Python"))
		if string(res) != "olleH" {
			b.Errorf("Call_byte failed: %v != %v", res, "olleH")
		}
	}
}
