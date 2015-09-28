package service

import (
	"github.com/qur/gopy/lib"

	"gem/service/archive"
	"gem/service/game"
)

type registerFunc func(*py.Module) error

var moduleRegisterFuncs = []registerFunc{
	archive.InitPyModule,
	game.InitPyModule,
}

func InitPyModule(parent *py.Module) error {
	/* Create package */
	var err error
	var module *py.Module
	if module, err = py.InitModule("gem.service", []py.Method{}); err != nil {
		return err
	}

	/* Register modules */
	for _, registerFunc := range moduleRegisterFuncs {
		if err = registerFunc(module); err != nil {
			return err
		}
	}

	if err = parent.AddObject("service", module); err != nil {
		return err
	}

	return nil
}