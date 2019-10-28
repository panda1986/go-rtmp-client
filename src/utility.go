package main

import (
	"log"
	"fmt"
)

func GetUnitBE(v uint32, n int) (data []byte) {
	data = make([]byte, n)
	log.Println(fmt.Sprintf("get bigendian, v=%v, n=%v", v, n))
	// 0x1234
	for  i := 0; i < n; i++ {
		b := byte(v >> (uint32(n-i-1) << 3))& 0xff
		data[i] = b
		log.Println(fmt.Sprintf("i=%v, %v", i, b))
	}
	return
}
