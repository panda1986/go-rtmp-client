package main

import (
	"bytes"
	"bufio"
)

/*
+---------+-----------------+-----------------+-----------------+
|         |Message Stream ID| Message TYpe ID | Time  | Length  |
+---------+-----------------+-----------------+-------+---------+
| Msg # 1 |    12345        |         8       | 1000  |   32    |
+---------+-----------------+-----------------+-------+---------+
| Msg # 2 |    12345        |         8       | 1020  |   32    |
+---------+-----------------+-----------------+-------+---------+
| Msg # 3 |    12345        |         8       | 1040  |   32    |
+---------+-----------------+-----------------+-------+---------+
| Msg # 4 |    12345        |         8       | 1060  |   32    |
+---------+-----------------+-----------------+-------+---------+
*/

func WriteCommandMessage(rw *bufio.ReadWriter, args ...interface{}) (err error) {
	buf := bytes.NewBuffer(nil)
	e := &Encoder{}
	for _, arg := range args {
		e.EncodeAmf0(buf, arg)
	}

	data := buf.Bytes()
	cs := &ChunkStream{
		Format: 0,
		CSID: 3,
		Timestamp: 0,
		TypeID: 20,
		StreamID: 0,
		Length: uint32(len(data)),
		Data: data,
	}

	cs.writeChunk(rw)
	rw.Flush()
	return
}

func WriteConnect(rw *bufio.ReadWriter) (err error) {
	event := make(Object)
	event["app"] = "live"
	event["type"] = "nonprivate"
	event["flashVer"] = "FMS.3.1"
	event["tcUrl"] = "rtmp://127.0.0.1/live"
	return WriteCommandMessage(rw, "connect", 1, event)
}
