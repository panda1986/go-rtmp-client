package main

import (
	"bufio"
	"log"
	"fmt"
	"io"
)

// page 7,
// it consists of three static-sized chunks rather than consisting of variable-sized chunks with headers.
// c0 and s0, 1 byte, The version defined by this specification is 3.
// c1 and s1, 1536 bytes, 4 bytes timestamp, 4 bytes zero, 1528 bytes random
// c2 and s2, 1536 bytes, 4 bytes timestamp, 4 bytes timestamp1, 1528 bytes random

func Handshake(rw *bufio.ReadWriter) (err error) {
	var random [(1 + 1536 * 2) * 2]byte
	c0c1c2 := random[0: (1 + 1536 * 2)]
	s0s1s2 := random[(1 + 1536 * 2):]

	//conn.SetDeadline()

	c0c1 := c0c1c2[:1536+1]
	c0 := c0c1c2[:1]
	c0[0] = 3
	nn, err := rw.Write(c0c1) // Write writes the contents of p into the rw's buffer
	if err != nil {
		log.Println(fmt.Sprintf("write c0c1 to buffer failed, err is %v, nn=%v", err, nn))
		return err
	}
	if err = rw.Flush(); err != nil {// Flush writes any buffered data to the underlying io.Writer
		log.Println(fmt.Sprintf("flush data to net conn failed, err is %v", err))
		return err
	}
	log.Println(fmt.Sprintf("write %v c0c1 to conn", nn))

	nn, err = io.ReadFull(rw, s0s1s2)
	if err != nil {
		log.Println(fmt.Sprintf("read from conn failed, err is %v, nn=%v", err, nn))
		return err
	}
	log.Println(fmt.Sprintf("read %v bytes from server", len(s0s1s2)))

	s1 := s0s1s2[1: 1536+1]
	log.Println(fmt.Sprintf("read s1, timestamp:%v", s1[0:4]))
	c2 := c0c1c2[1536+1:]
	c2 = s1
	nn, err = rw.Write(c2)
	if err != nil {
		log.Println(fmt.Sprintf("write c2 to buffer failed, err is %v", err))
		return err
	}
	if err = rw.Flush(); err != nil {
		log.Println(fmt.Sprintf("flush c2 to net conn failed, err is %v", err))
		return err
	}
	log.Println(fmt.Sprintf("write %v c2 to conn", len(c2)))

	return nil
}

