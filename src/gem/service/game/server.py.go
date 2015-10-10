// Generated by gopygen; DO NOT EDIT
package game

import (
	"fmt"
	"gem/auth"
	"gem/runite"
	"sync"

	"github.com/qur/gopy/lib"
	"gopkg.in/tomb.v2"

	"github.com/tgascoigne/gopygen/gopygen"
)

// Sometimes we might generate code which doesn't use some of the above imports
// Use them here just in case
var _ = fmt.Sprintf("")
var _ = gopygen.Dummy

var ServerDef = py.Class{
	Name:    "Server",
	Pointer: (*Server)(nil),
}

// Registers this type with a python module
func RegisterServer(module *py.Module) error {
	var err error
	var class *py.Type
	if class, err = ServerDef.Create(); err != nil {
		return err
	}

	if err = module.AddObject("Server", class); err != nil {
		return err
	}

	return nil
}

// Alloc allocates an object for use in python land.
// Copies the member fields from this object to the newly allocated object
// Usage: obj := GoObject{X:1, Y: 2}.Alloc()
func (obj Server) Alloc() (*Server, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	// Allocate
	alloc_, err := ServerDef.Alloc(0)
	if err != nil {
		return nil, err
	}
	alloc := alloc_.(*Server)
	// Copy fields

	alloc.laddr = obj.laddr

	alloc.ln = obj.ln

	alloc.update = obj.update

	alloc.game = obj.game

	alloc.runite = obj.runite

	alloc.nextIndex = obj.nextIndex

	alloc.m = obj.m

	alloc.clients = obj.clients

	alloc.t = obj.t

	return alloc, nil
}

func (obj *Server) PyGet_m() (py.Object, error) {
	return gopygen.TypeConvOut(obj.m, "sync.Mutex")
}

func (obj *Server) PySet_m(arg py.Object) error {
	val, err := gopygen.TypeConvIn(arg, "sync.Mutex")
	if err != nil {
		return err
	}
	obj.m = val.(sync.Mutex)
	return nil
}

func (obj *Server) PyGet_t() (py.Object, error) {
	return gopygen.TypeConvOut(obj.t, "tomb.Tomb")
}

func (obj *Server) PySet_t(arg py.Object) error {
	val, err := gopygen.TypeConvIn(arg, "tomb.Tomb")
	if err != nil {
		return err
	}
	obj.t = val.(tomb.Tomb)
	return nil
}

func (s *Server) Py_Start(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 4 {
		return nil, fmt.Errorf("Py_Start: parameter length mismatch")
	}

	in_0, err := gopygen.TypeConvIn(args[0], "string")
	if err != nil {
		return nil, err
	}

	in_1, err := gopygen.TypeConvIn(args[1], "*runite.Context")
	if err != nil {
		return nil, err
	}

	in_2, err := gopygen.TypeConvIn(args[2], "string")
	if err != nil {
		return nil, err
	}

	in_3, err := gopygen.TypeConvIn(args[3], "auth.Provider")
	if err != nil {
		return nil, err
	}

	res0 := s.Start(in_0.(string), in_1.(*runite.Context), in_2.(string), in_3.(auth.Provider))

	out_0, err := gopygen.TypeConvOut(res0, "error")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}

func (s *Server) Py_Stop(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_Stop: parameter length mismatch")
	}

	res0 := s.Stop()

	out_0, err := gopygen.TypeConvOut(res0, "error")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}
