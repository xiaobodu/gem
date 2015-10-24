package task

import (
	"github.com/qur/gopy/lib"

	"github.com/sinusoids/gem/gem/python"
)

type registerFunc func(*py.Module) error

var moduleRegisterFuncs = []registerFunc{}

func init() {
	lock := py.NewLock()
	defer lock.Unlock()

	/* Create package */
	var err error
	var module *py.Module
	methods := []py.Method{
		{"submit", Py_Submit, "submit a task to the scheduler"},
	}
	if module, err = python.InitModule("gem.task", methods); err != nil {
		panic(err)
	}

	createTaskHookConstants(module)

	/* Register modules */
	for _, registerFunc := range moduleRegisterFuncs {
		if err = registerFunc(module); err != nil {
			panic(err)
		}
	}
}
