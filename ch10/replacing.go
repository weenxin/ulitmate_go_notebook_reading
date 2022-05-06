package ch10

import (
	"bytes"
)

var data = []struct {
	input  []byte
	output []byte
}{
	{[]byte("abc"), []byte("abc")},
	{[]byte("elvis"), []byte("Elvis")},
	{[]byte("aElvis"), []byte("aElvis")},
	{[]byte("abcelvis"), []byte("abcElvis")},
	{[]byte("eelvis"), []byte("eElvis")},
	{[]byte("aelvis"), []byte("aElvis")},
	{[]byte("aabeeeelvis"), []byte("aabeeeElvis")},
	{[]byte("e l v i s"), []byte("e l v i s")},
	{[]byte("aa bb e l v i saa"), []byte("aa bb e l v i saa")},
	{[]byte(" elvi s"), []byte(" elvi s")},
	{[]byte("elvielvis"), []byte("elviElvis")},
	{[]byte("elvielvielviselvi1"), []byte("elvielviElviselvi1")},
	{[]byte("elvielviselvis"), []byte("elviElvisElvis")},
}

func assembleInputStream() []byte {
	var in []byte
	for _, d := range data {
		in = append(in, d.input...)
	}
	return in
}
func assembleOutputStream() []byte {
	var out []byte
	for _, d := range data {
		out = append(out, d.output...)
	}
	return out
}

func algOne(data []byte, find []byte, repl []byte, output *bytes.Buffer) {
	input := bytes.NewBuffer(data)
	size := len(find)
	buf := make([]byte, 5)
	end := size - 1
	if n, err := input.Read(buf[:end]); err != nil { //直接调用对象的方法，没有隐式转换，所以不会逃逸
	//if n, err := io.ReadFull(input, buf[:end]); err != nil { //把input当作一个io.Reader接口来用，此时会导致input逃逸
		output.Write(buf[:n])
		return
	}
	for {
		if _, err := input.Read(buf[end:]); err != nil { //直接调用对象的方法，没有隐式转换，所以不会逃逸
		//if _, err := io.ReadFull(input, buf[end:]); err != nil {//把input当作一个io.Reader接口来用，此时会导致input逃逸
			output.Write(buf[:end])
			return
		}

		if equal(buf, find) { //3 占用时间比较多，因此可以优化下
		//if bytes.Equal(buf, find) { //3 转换成为一个string，会有一次cop，性能降低
			output.Write(repl)
			//if _, err := io.ReadFull(input, buf[end:]); err != nil { //当作io.Reader来用，input会逃逸
			//	output.Write(buf[:end])
			//	return
			//}

			data, err := input.ReadByte() //　调用方法，没有隐式转换
			buf[end:][0] = data
			if  err != nil {
				output.Write(buf[:end])
				return
			}

			continue
		}
		output.WriteByte(buf[0])
		copy(buf, buf[1:])
	}
}

func algTwo(data []byte, find []byte, repl []byte, output *bytes.Buffer) {
	input := bytes.NewReader(data)
	size := len(find)
	idx := 0
	for {
		b, err := input.ReadByte()
		if err != nil {
			break
		}
		if b == find[idx] {
			idx++
			if idx == size {
				output.Write(repl)
				idx = 0
			}
			continue
		}
		if idx != 0 {
			output.Write(find[:idx])
			input.UnreadByte()
			idx = 0
			continue
		}
		output.WriteByte(b)
		idx = 0
	}
}

func equal(a , b []byte) bool {
	//return  (*(*string)(unsafe.Pointer(&a))) ==  (*(*string)(unsafe.Pointer(&b)))
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index] != b[index] {
			return  false
		}
	}
	return true
}


