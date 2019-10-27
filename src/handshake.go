package main

import (
	"net"
	"bufio"
)

// page 7,
// it consists of three static-sized chunks rather than consisting of variable-sized chunks with headers.
// c0 and s0, 1 byte, The version defined by this specification is 3.
// c1 and s1, 1536 bytes, 4 bytes timestamp, 4 bytes zero, 1528 bytes random
// c2 and s2, 1536 bytes, 4 bytes timestamp, 4 bytes timestamp1, 1528 bytes random

func Handshake(conn net.Conn)  {
	var random [(1 + 1536 * 2) * 2]byte
	c0c1c2 := random[0: (1 + 1536 * 2)]
	s0s1s2 := random[(1 + 1536 * 2):]

	//conn.SetDeadline()
	size := 4 * 1024
	rw := bufio.NewReadWriter(bufio.NewReaderSize(conn, size), bufio.NewWriterSize(conn, size))

	c0c1 := c0c1c2[:1536+1]
	c0 := c0c1c2[:1]
	rw.Write()
	rw.Flush()
}

