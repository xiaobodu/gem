package gem

import (
	"gem/archive"
	"gem/auth"
	"gem/event"
	"gem/game"
	"gem/log"
	"gem/runite"
	"gem/task"

	"github.com/qur/gopy/lib"
)

type registerFunc func(*py.Module) error

var moduleRegisterFuncs = []registerFunc{
	RegisterEngine,
	log.RegisterSysLog,
	log.RegisterModule,

	auth.InitPyModule,
	event.InitPyModule,
	task.InitPyModule,
	runite.InitPyModule,

	archive.InitPyModule,
	game.InitPyModule,
}

func InitPyModule() error {
	/* Create package */
	var err error
	var module *py.Module
	if module, err = py.InitModule("gem", []py.Method{}); err != nil {
		return err
	}

	/* Register modules */
	for _, registerFunc := range moduleRegisterFuncs {
		if err = registerFunc(module); err != nil {
			return err
		}
	}

	/* Create our logger object */
	log.InitSysLog()
	if err := module.AddObject("syslog", log.Sys); err != nil {
		return err
	}

	return nil
}
