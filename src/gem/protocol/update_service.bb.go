// Generated by bbc; DO NOT EDIT
package protocol

import (
	"gem/encoding"
	"io"
)

type UpdateRequest struct {
	Index    encoding.Int8
	File     encoding.Int16
	Priority encoding.Int8
}

func (struc *UpdateRequest) Encode(buf io.Writer, flags interface{}) (err error) {
	err = struc.Index.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.File.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Priority.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

func (struc *UpdateRequest) Decode(buf io.Reader, flags interface{}) (err error) {
	err = struc.Index.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.File.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Priority.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

type UpdateResponse struct {
	Index encoding.Int8
	File  encoding.Int16
	Size  encoding.Int16
	Chunk encoding.Int8
	Data  encoding.Bytes
}

func (struc *UpdateResponse) Encode(buf io.Writer, flags interface{}) (err error) {
	err = struc.Index.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.File.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Size.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Chunk.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Data.Encode(buf, 500)
	if err != nil {
		return err
	}

	return err
}

func (struc *UpdateResponse) Decode(buf io.Reader, flags interface{}) (err error) {
	err = struc.Index.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.File.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Size.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Chunk.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.Data.Decode(buf, 500)
	if err != nil {
		return err
	}

	return err
}
