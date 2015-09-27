// Generated by bbc; DO NOT EDIT
package rt3

import (
	"gem/encoding"
	"io"
)

type FSIndex struct {
	Length     encoding.Int24
	StartBlock encoding.Int24
}

func (struc *FSIndex) Encode(buf io.Writer, flags interface{}) (err error) {
	err = struc.Length.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.StartBlock.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

func (struc *FSIndex) Decode(buf io.Reader, flags interface{}) (err error) {
	err = struc.Length.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.StartBlock.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

type FSBlock struct {
	FileID       encoding.Int16
	FilePosition encoding.Int16
	NextBlock    encoding.Int24
	Partition    encoding.Int8
	Data         encoding.Bytes
}

func (struc *FSBlock) Encode(buf io.Writer, flags interface{}) (err error) {
	err = struc.FileID.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.FilePosition.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.NextBlock.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Partition.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Data.Encode(buf, 512)
	if err != nil {
		return err
	}

	return err
}

func (struc *FSBlock) Decode(buf io.Reader, flags interface{}) (err error) {
	err = struc.FileID.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.FilePosition.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.NextBlock.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Partition.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Data.Decode(buf, 512)
	if err != nil {
		return err
	}

	return err
}

type CRCFile struct {
	Archives [9]encoding.Int32
	Sum      encoding.Int32
}

func (struc *CRCFile) Encode(buf io.Writer, flags interface{}) (err error) {
	for i := 0; i < 9; i++ {
		err = struc.Archives[i].Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
		if err != nil {
			return err
		}
	}

	err = struc.Sum.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

func (struc *CRCFile) Decode(buf io.Reader, flags interface{}) (err error) {
	for i := 0; i < 9; i++ {
		err = struc.Archives[i].Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
		if err != nil {
			return err
		}
	}

	err = struc.Sum.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}
