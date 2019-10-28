package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
)

/*
+--------+---------+-----+------------+------- ---+------------+
|        | Chunk   |Chunk|Header Data |No.of Bytes|Total No.of |
|        |Stream ID|Type |            |  After    |Bytes in the|
|        |         |     |            |Header     |Chunk       |
+--------+---------+-----+------------+-----------+------------+
|Chunk#1 |    3    |  0  | delta: 1000|   32      |    44      |
|        |         |     | length: 32,|           |
|        |         |     | type: 8,   |           |            |
|        |         |     | stream ID: |           |            |
|        |         |     | 12345 (11  |           |            |
|        |         |     | bytes)     |           |            |
+--------+---------+-----+------------+-----------+------------+
|Chunk#2 |    3    |  2  | 20 (3      |   32      |    36      |
|        |         |     | bytes)     |           |            |
+--------+---------+-----+----+-------+-----------+------------+
|Chunk#3 |    3    |  3  | none (0    |   32      |    33      |
|        |         |     | bytes)     |           |            |
+--------+---------+-----+------------+-----------+------------+
|Chunk#4 |    3    |  3  | none (0    |   32      |    33      |
|        |         |     | bytes)     |           |            |
+--------+---------+-----+------------+-----------+------------+
*/

type ChunkStream struct {
	Format uint32
	CSID uint32
	Timestamp uint32
	Length uint32
	TypeID uint32
	StreamID uint32
	Data []byte
}

/*
Command messages carry the AMF-encoded commands between the client and the server.
*/

func (v *ChunkStream) writeChunk(w *bufio.ReadWriter) (error) {
	chunkSize := 128
	totalLen := uint32(0)
	numChunks := v.Length / uint32(chunkSize)
	for i := uint32(0); i <= numChunks; i ++ {
		if totalLen == v.Length {
			break
		}
		if i == 0 {
			v.Format = 0
		} else {
			v.Format = 3
		}
		// write chunk header
		v.writeHeader(w)

		// write chunk payload
		inc := uint32(chunkSize)
		start := uint32(i) * uint32(chunkSize)
		if uint32(len(v.Data))-start <= inc {
			inc = uint32(len(v.Data)) - start
		}
		totalLen += inc
		end := start + inc
		buf := v.Data[start:end]
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}
	return nil
}


func (v *ChunkStream) writeHeader(w *bufio.ReadWriter) error {
	// write basic header
	// fmt - 2 bits （big endian）
	// csid (little endian)
	bh := v.Format << 6
	if v.CSID < 64 {
		bh |= v.CSID
		binary.Write(w, binary.BigEndian, uint(bh))
	} else if (v.CSID - 64) < 256 {
		bh |= 0
		binary.Write(w, binary.BigEndian, uint(bh))
		csid := v.CSID - 64
		binary.Write(w, binary.LittleEndian, uint(csid))
	} else if (v.CSID - 64) < 65536 {
		bh |= 1
		binary.Write(w, binary.BigEndian, uint(bh))
		csid := v.CSID - 64
		binary.Write(w, binary.LittleEndian, uint16(csid))
	}

	// write message header
	ts := v.Timestamp
	if v.Format == 3 {
		goto END
	}

	if ts > 0xffffff {
		ts = 0xffffff
	}
	w.Write(GetUnitBE(ts, 3))
	if v.Format == 2 {
		goto END
	}
	if v.Length > 0xffffff {
		return fmt.Errorf("length=%v, exceed 0xffffff", v.Length)
	}

	w.Write(GetUnitBE(v.Length, 3))
	binary.Write(w, binary.BigEndian, uint(v.TypeID))
	if v.Format == 1 {
		goto END
	}

	binary.Write(w, binary.LittleEndian, v.StreamID)

END:
// Extended Timestamp
	if v.Timestamp >= 0xffffff {
		binary.Write(w, binary.BigEndian, v.Timestamp)
	}
	return nil
}
