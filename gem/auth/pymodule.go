package auth

import (
	"github.com/qur/gopy/lib"

	"github.com/gemrs/gem/gem/python/modules"
)

type registerFunc func(*py.Module) error

var moduleRegisterFuncs = []registerFunc{
	RegisterProviderImpl,
	createAuthConstants,
}

func init() {
	lock := py.NewLock()
	defer lock.Unlock()

	/* Create package */
	var err error
	var module *py.Module
	if module, err = modules.Init("gem.auth", []py.Method{}); err != nil {
		panic(err)
	}

	/* Register modules */
	for _, registerFunc := range moduleRegisterFuncs {
		if err = registerFunc(module); err != nil {
			panic(err)
		}
	}
}
