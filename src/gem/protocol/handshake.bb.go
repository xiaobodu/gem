// Generated by bbc; DO NOT EDIT
package protocol

import (
	"gem/encoding"
	"io"
)

type ServiceSelect struct {
	Service encoding.Int8
}

func (struc *ServiceSelect) Encode(buf io.Writer, flags interface{}) (err error) {
	err = struc.Service.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

func (struc *ServiceSelect) Decode(buf io.Reader, flags interface{}) (err error) {
	err = struc.Service.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

type GameHandshake struct {
	NameHash encoding.Int8
}

func (struc *GameHandshake) Encode(buf io.Writer, flags interface{}) (err error) {
	err = struc.NameHash.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

func (struc *GameHandshake) Decode(buf io.Reader, flags interface{}) (err error) {
	err = struc.NameHash.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

type UpdateHandshakeResponse struct {
	ignored [8]encoding.Int8
}

func (struc *UpdateHandshakeResponse) Encode(buf io.Writer, flags interface{}) (err error) {
	for i := 0; i < 8; i++ {
		err = struc.ignored[i].Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
		if err != nil {
			return err
		}
	}

	return err
}

func (struc *UpdateHandshakeResponse) Decode(buf io.Reader, flags interface{}) (err error) {
	for i := 0; i < 8; i++ {
		err = struc.ignored[i].Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
		if err != nil {
			return err
		}
	}

	return err
}

type GameHandshakeResponse struct {
	ignored         [8]encoding.Int8
	loginRequest    encoding.Int8
	ServerISAACSeed encoding.Int64
}

func (struc *GameHandshakeResponse) Encode(buf io.Writer, flags interface{}) (err error) {
	for i := 0; i < 8; i++ {
		err = struc.ignored[i].Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
		if err != nil {
			return err
		}
	}

	err = struc.loginRequest.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.ServerISAACSeed.Encode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}

func (struc *GameHandshakeResponse) Decode(buf io.Reader, flags interface{}) (err error) {
	for i := 0; i < 8; i++ {
		err = struc.ignored[i].Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
		if err != nil {
			return err
		}
	}

	err = struc.loginRequest.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	err = struc.ServerISAACSeed.Decode(buf, encoding.IntegerFlag(encoding.IntNilFlag))
	if err != nil {
		return err
	}

	return err
}
