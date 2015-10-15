// Generated by gopygen; DO NOT EDIT
package server

import (
	"fmt"
	"gem/encoding"

	"github.com/qur/gopy/lib"
	"github.com/tgascoigne/gopygen/gopygen"
)

// Sometimes we might generate code which doesn't use some of the above imports
// Use them here just in case
var _ = fmt.Sprintf("")
var _ = gopygen.Dummy

var ConnectionDef = py.Class{
	Name:    "Connection",
	Pointer: (*Connection)(nil),
}

// Registers this type with a python module
func RegisterConnection(module *py.Module) error {
	var err error
	var class *py.Type
	if class, err = ConnectionDef.Create(); err != nil {
		return err
	}

	if err = module.AddObject("Connection", class); err != nil {
		return err
	}

	return nil
}

// Alloc allocates an object for use in python land.
// Copies the member fields from this object to the newly allocated object
// Usage: obj := GoObject{X:1, Y: 2}.Alloc()
func (obj Connection) Alloc() (*Connection, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	// Allocate
	alloc_, err := ConnectionDef.Alloc(0)
	if err != nil {
		return nil, err
	}
	alloc := alloc_.(*Connection)
	// Copy fields

	alloc.ReadBuffer = obj.ReadBuffer

	alloc.WriteBuffer = obj.WriteBuffer

	alloc.Read = obj.Read

	alloc.Write = obj.Write

	alloc.DisconnectChan = obj.DisconnectChan

	alloc.log = obj.log

	alloc.index = obj.index

	alloc.conn = obj.conn

	return alloc, nil
}

func (c *Connection) Py_Log(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_Log: parameter length mismatch")
	}

	res0 := c.Log()

	out_0, err := gopygen.TypeConvOut(res0, "*log.Module")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}

func (conn *Connection) Py_WaitForDisconnect(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_WaitForDisconnect: parameter length mismatch")
	}

	conn.WaitForDisconnect()

	py.None.Incref()
	return py.None, nil

}

func (conn *Connection) Py_IsDisconnecting(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_IsDisconnecting: parameter length mismatch")
	}

	res0 := conn.IsDisconnecting()

	out_0, err := gopygen.TypeConvOut(res0, "bool")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}

func (conn *Connection) Py_Disconnect(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_Disconnect: parameter length mismatch")
	}

	conn.Disconnect()

	py.None.Incref()
	return py.None, nil

}

func (conn *Connection) Py_Index(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_Index: parameter length mismatch")
	}

	res0 := conn.Index()

	out_0, err := gopygen.TypeConvOut(res0, "int")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}

func (conn *Connection) Py_SetIndex(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 1 {
		return nil, fmt.Errorf("Py_SetIndex: parameter length mismatch")
	}

	in_0, err := gopygen.TypeConvIn(args[0], "int")
	if err != nil {
		return nil, err
	}

	conn.SetIndex(in_0.(int))

	py.None.Incref()
	return py.None, nil

}

func (conn *Connection) Py_WriteEncodable(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 1 {
		return nil, fmt.Errorf("Py_WriteEncodable: parameter length mismatch")
	}

	in_0, err := gopygen.TypeConvIn(args[0], "encoding.Encodable")
	if err != nil {
		return nil, err
	}

	conn.WriteEncodable(in_0.(encoding.Encodable))

	py.None.Incref()
	return py.None, nil

}

func (conn *Connection) Py_flushWriteBuffer(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_flushWriteBuffer: parameter length mismatch")
	}

	res0 := conn.flushWriteBuffer()

	out_0, err := gopygen.TypeConvOut(res0, "error")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}

func (conn *Connection) Py_fillReadBuffer(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_fillReadBuffer: parameter length mismatch")
	}

	res0 := conn.fillReadBuffer()

	out_0, err := gopygen.TypeConvOut(res0, "error")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}
