// Generated by gopygen; DO NOT EDIT
package gem

import (
	"fmt"

	"gopkg.in/tomb.v2"

	"github.com/qur/gopy/lib"
	"github.com/tgascoigne/gopygen/gopygen"
)

// Sometimes we might generate code which doesn't use some of the above imports
// Use them here just in case
var _ = fmt.Sprintf("")
var _ = gopygen.Dummy

var EngineDef = py.Class{
	Name:    "Engine",
	Flags:   py.TPFLAGS_BASETYPE,
	Pointer: (*Engine)(nil),
}

// Registers this type with a python module
func RegisterEngine(module *py.Module) error {
	var err error
	var class *py.Type
	if class, err = EngineDef.Create(); err != nil {
		return err
	}

	if err = module.AddObject("Engine", class); err != nil {
		return err
	}

	return nil
}

// Alloc allocates an object for use in python land.
// Copies the member fields from this object to the newly allocated object
// Usage: obj := GoObject{X:1, Y: 2}.Alloc()
func NewEngine() (*Engine, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	// Allocate
	alloc_, err := EngineDef.Alloc(0)
	if err != nil {
		return nil, err
	}
	alloc := alloc_.(*Engine)
	err = alloc.Init()
	return alloc, err
}

func (obj *Engine) PyInit(_args *py.Tuple, kwds *py.Dict) error {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return fmt.Errorf("(Engine) PyInit: parameter length mismatch")
	}

	err = obj.Init()

	return err
}

func (obj *Engine) PyGet_t() (py.Object, error) {
	return gopygen.TypeConvOut(obj.t, "tomb.Tomb")
}

func (obj *Engine) PySet_t(arg py.Object) error {
	arg.Incref()
	val, err := gopygen.TypeConvIn(arg, "tomb.Tomb")
	if err != nil {
		return err
	}

	if _, ok := val.(py.Object); ok {
		// If we're not converting it from a python object, we should refcount it properly
		val.(py.Object).Incref()
	}
	arg.Decref()

	var tmp interface{}
	tmp = &obj.t
	obj.t = val.(tomb.Tomb)

	if oldObj, ok := tmp.(py.Object); ok {
		// If we're not converting it from a python object, we should refcount it properly
		oldObj.Decref()
	}
	return nil
}

func (e *Engine) Py_Start(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_Start: parameter length mismatch")
	}
	// Convert parameters

	// Make the function call

	e.Start()

	// Remove local references

	py.None.Incref()
	return py.None, nil

}

func (e *Engine) Py_Join(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_Join: parameter length mismatch")
	}
	// Convert parameters

	// Make the function call

	res0 := e.Join()

	// Remove local references

	out_0, err := gopygen.TypeConvOut(res0, "bool")
	if err != nil {
		return nil, err
	}
	out_0.Incref()

	return out_0, nil

}

func (e *Engine) Py_Stop(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 0 {
		return nil, fmt.Errorf("Py_Stop: parameter length mismatch")
	}
	// Convert parameters

	// Make the function call

	e.Stop()

	// Remove local references

	py.None.Incref()
	return py.None, nil

}
