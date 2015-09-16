// Generated by framecc; DO NOT EDIT
package protocol

import (
	"bytes"
	"gem/encoding"
)

type TestFrame AnonStruct_49

func (frm *TestFrame) Encode(buf *bytes.Buffer, flags interface{}) (err error) {
	struc := (*AnonStruct_49)(frm)
	return struc.Encode(buf, flags)
}

func (frm *TestFrame) Decode(buf *bytes.Buffer, flags interface{}) (err error) {
	struc := (*AnonStruct_49)(frm)
	return struc.Decode(buf, flags)
}

type EmbeddedStruct struct {
	A encoding.Int32
	B encoding.Int32
	C encoding.Int32
}

func (struc *EmbeddedStruct) Encode(buf *bytes.Buffer, flags interface{}) (err error) {
	err = struc.A.Encode(buf, encoding.IntegerFlag(encoding.IntLittleEndian))
	if err != nil {
		return err
	}

	err = struc.B.Encode(buf, encoding.IntegerFlag(encoding.IntOffset128|encoding.IntPDPEndian))
	if err != nil {
		return err
	}

	err = struc.C.Encode(buf, encoding.IntegerFlag(encoding.IntInverse128|encoding.IntRPDPEndian))
	if err != nil {
		return err
	}

	return
}

func (struc *EmbeddedStruct) Decode(buf *bytes.Buffer, flags interface{}) (err error) {
	err = struc.A.Decode(buf, encoding.IntegerFlag(encoding.IntLittleEndian))
	if err != nil {
		return err
	}

	err = struc.B.Decode(buf, encoding.IntegerFlag(encoding.IntOffset128|encoding.IntPDPEndian))
	if err != nil {
		return err
	}

	err = struc.C.Decode(buf, encoding.IntegerFlag(encoding.IntInverse128|encoding.IntRPDPEndian))
	if err != nil {
		return err
	}

	return
}

type AnonStruct_49 struct {
	Message  encoding.JString
	Values8  [4]encoding.Int8
	Values16 [2]encoding.Int16
	Struc1   EmbeddedStruct
	Struc2   [2]EmbeddedStruct
}

func (struc *AnonStruct_49) Encode(buf *bytes.Buffer, flags interface{}) (err error) {
	err = struc.Message.Encode(buf, 16)
	if err != nil {
		return err
	}

	for i := 0; i < 4; i++ {
		err = struc.Values8[i].Encode(buf, encoding.IntegerFlag(encoding.IntNegate))
		if err != nil {
			return err
		}
	}

	for i := 0; i < 2; i++ {
		err = struc.Values16[i].Encode(buf, encoding.IntegerFlag(0))
		if err != nil {
			return err
		}
	}

	err = struc.Struc1.Encode(buf, encoding.NilFlags)
	if err != nil {
		return err
	}

	for i := 0; i < 2; i++ {
		err = struc.Struc2[i].Encode(buf, encoding.NilFlags)
		if err != nil {
			return err
		}
	}

	return
}

func (struc *AnonStruct_49) Decode(buf *bytes.Buffer, flags interface{}) (err error) {
	err = struc.Message.Decode(buf, 16)
	if err != nil {
		return err
	}

	for i := 0; i < 4; i++ {
		err = struc.Values8[i].Decode(buf, encoding.IntegerFlag(encoding.IntNegate))
		if err != nil {
			return err
		}
	}

	for i := 0; i < 2; i++ {
		err = struc.Values16[i].Decode(buf, encoding.IntegerFlag(0))
		if err != nil {
			return err
		}
	}

	err = struc.Struc1.Decode(buf, encoding.NilFlags)
	if err != nil {
		return err
	}

	for i := 0; i < 2; i++ {
		err = struc.Struc2[i].Decode(buf, encoding.NilFlags)
		if err != nil {
			return err
		}
	}

	return
}