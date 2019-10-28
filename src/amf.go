package main

import (
	"io"
	"fmt"
	"reflect"
	"encoding/binary"
)

const (
	AMF0_NUMBER_MARKER         = 0x00
	AMF0_BOOLEAN_MARKER        = 0x01
	AMF0_STRING_MARKER         = 0x02
	AMF0_OBJECT_MARKER         = 0x03
	AMF0_MOVIECLIP_MARKER      = 0x04 // reserved, not supported
	AMF0_NULL_MARKER           = 0x05
	AMF0_UNDEFINED_MARKER      = 0x06
	AMF0_REFERENCE_MARKER      = 0x07
	AMF0_ECMA_ARRAY_MARKER     = 0x08
	AMF0_OBJECT_END_MARKER     = 0x09
	AMF0_STRICT_ARRAY_MARKER   = 0x0a
	AMF0_DATE_MARKER           = 0x0b
	AMF0_LONG_STRING_MARKER    = 0x0c
	AMF0_UNSUPPORTED_MARKER    = 0x0d
	AMF0_RECORDSET_MARKER      = 0x0e //reserved, not supported
	AMF0_XML_DOCUMENT_MARKER   = 0x0f
	AMF0_TYPED_OBJECT_MARKER   = 0x10
	AMF0_ACMPLUS_OBJECT_MARKER = 0x11
)

func WriteMarker(w io.Writer, m byte) (err error) {
	data := []byte{}
	data = append(data, m)
	_, err = w.Write(data)
	return
}

type Encoder struct {

}

func (v *Encoder) EncodeAmf0(w io.Writer, val interface{}) (int, error) {
	if val == nil {
		return v.EncodeAmf0Null(w, true)
	}
	vt := reflect.ValueOf(val)
	if !vt.IsValid() {
		return v.EncodeAmf0Null(w, true)
	}

	switch vt.Kind() {
	case reflect.String:
		return v.EncodeAmf0String(w, vt.String(), true)
	case reflect.Map:
		obj, ok := val.(Object)
		if ok != true {
			return 0, fmt.Errorf("encode amf0: unable to create object from map")
		}
		return v.EncodeAmf0Object(w, obj, true)
	}

	return 0, fmt.Errorf("encode amf0 failed, unsupported type:%v", vt.Type())
}

// marker: 1 byte 0x05
// no additional data
func (e *Encoder) EncodeAmf0Null(w io.Writer, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_NULL_MARKER); err != nil {
			return
		}
		n += 1
	}
	return
}

// marker: 1byte, 0x02
// format:
// - 2 byte big endian uint16  header to determine size
// -n byte utf8 string
func (e *Encoder) EncodeAmf0String(w io.Writer, val string, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_STRING_MARKER); err != nil {
			return
		}
		n += 1
	}

	length := uint16(len(val))
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, length)
	if _, err = w.Write(data); err != nil {
		return n, fmt.Errorf("write string length failed, err is %v", err)
	}
	n += 2

	if _, err = w.Write([]byte(val)); err != nil {
		return n, fmt.Errorf("write string failed, err is %v", err)
	}
	n += len(val)
	return
}

// marker: 1 byte 0x03
// format:
// - loop encoded string followed by encoded value
// - terminated with empty string followed by 1 byte 0x09
type Object map[string]interface{}

func (e *Encoder) EncodeAmf0Object(w io.Writer, val Object, encodeMarker bool) (n int, err error) {
	if encodeMarker {
		if err = WriteMarker(w, AMF0_OBJECT_MARKER); err != nil {
			return
		}
		n += 1
	}

	var m int
	for k, v := range val {
		m, err = e.EncodeAmf0String(w, k, false)
		if err != nil {
			return n, fmt.Errorf("encode amf0: unable to encode object key: %s", err)
		}
		n += m

		m, err = e.EncodeAmf0(w, v)
		if err != nil {
			return n, fmt.Errorf("encode amf0: unable to encode object value: %s", err)
		}
		n += m
	}

	m, err = e.EncodeAmf0String(w, "", false)
	if err != nil {
		return n, fmt.Errorf("encode amf0: unable to encode object empty string: %s", err)
	}
	n += m

	err = WriteMarker(w, AMF0_OBJECT_END_MARKER)
	if err != nil {
		return n, fmt.Errorf("encode amf0: unable to object end marker: %s", err)
	}
	n += 1

	return
}


