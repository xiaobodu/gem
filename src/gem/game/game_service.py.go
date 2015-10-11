// Generated by gopygen; DO NOT EDIT
package game

import (
	"fmt"
	"gem/auth"
	"gem/game/server"
	"gem/runite"

	"github.com/qur/gopy/lib"
	"github.com/tgascoigne/gopygen/gopygen"
)

// Sometimes we might generate code which doesn't use some of the above imports
// Use them here just in case
var _ = fmt.Sprintf("")
var _ = gopygen.Dummy

var GameServiceDef = py.Class{
	Name:    "GameService",
	Pointer: (*GameService)(nil),
}

// Registers this type with a python module
func RegisterGameService(module *py.Module) error {
	var err error
	var class *py.Type
	if class, err = GameServiceDef.Create(); err != nil {
		return err
	}

	if err = module.AddObject("GameService", class); err != nil {
		return err
	}

	return nil
}

// Alloc allocates an object for use in python land.
// Copies the member fields from this object to the newly allocated object
// Usage: obj := GoObject{X:1, Y: 2}.Alloc()
func (obj GameService) Alloc() (*GameService, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	// Allocate
	alloc_, err := GameServiceDef.Alloc(0)
	if err != nil {
		return nil, err
	}
	alloc := alloc_.(*GameService)
	// Copy fields

	alloc.runite = obj.runite

	alloc.key = obj.key

	alloc.auth = obj.auth

	return alloc, nil
}

func (svc *GameService) Py_Init(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 3 {
		return nil, fmt.Errorf("Py_Init: parameter length mismatch")
	}

	in_0, err := gopygen.TypeConvIn(args[0], "*runite.Context")
	if err != nil {
		return nil, err
	}

	in_1, err := gopygen.TypeConvIn(args[1], "string")
	if err != nil {
		return nil, err
	}

	in_2, err := gopygen.TypeConvIn(args[2], "auth.Provider")
	if err != nil {
		return nil, err
	}

	res0 := svc.Init(in_0.(*runite.Context), in_1.(string), in_2.(auth.Provider))

	out_0, err := gopygen.TypeConvOut(res0, "error")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}

func (svc *GameService) Py_NewClient(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 2 {
		return nil, fmt.Errorf("Py_NewClient: parameter length mismatch")
	}

	in_0, err := gopygen.TypeConvIn(args[0], "*server.Connection")
	if err != nil {
		return nil, err
	}

	in_1, err := gopygen.TypeConvIn(args[1], "int")
	if err != nil {
		return nil, err
	}

	res0 := svc.NewClient(in_0.(*server.Connection), in_1.(int))

	out_0, err := gopygen.TypeConvOut(res0, "server.Client")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}

func (svc *GameService) Py_decodePacket(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 1 {
		return nil, fmt.Errorf("Py_decodePacket: parameter length mismatch")
	}

	in_0, err := gopygen.TypeConvIn(args[0], "*GameClient")
	if err != nil {
		return nil, err
	}

	res0 := svc.decodePacket(in_0.(*GameClient))

	out_0, err := gopygen.TypeConvOut(res0, "error")
	if err != nil {
		return nil, err
	}

	return out_0, nil

}

func (svc *GameService) Py_packetConsumer(_args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	lock := py.NewLock()
	defer lock.Unlock()

	var err error
	_ = err
	args := _args.Slice()
	if len(args) != 1 {
		return nil, fmt.Errorf("Py_packetConsumer: parameter length mismatch")
	}

	in_0, err := gopygen.TypeConvIn(args[0], "*GameClient")
	if err != nil {
		return nil, err
	}

	svc.packetConsumer(in_0.(*GameClient))

	py.None.Incref()
	return py.None, nil

}
