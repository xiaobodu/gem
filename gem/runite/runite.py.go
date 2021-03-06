package runite

import (
	"github.com/qur/gopy/lib"

	"github.com/gemrs/gem/pybind"
)

var ContextDef = pybind.Define("Context", (*Context)(nil))
var RegisterContext = pybind.GenerateRegisterFunc(ContextDef)
var NewContext = pybind.GenerateConstructor(ContextDef).(func() *Context)

func (c *Context) Py_unpack(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	fn := pybind.Wrap(c.Unpack)
	return fn(args, kwds)
}
