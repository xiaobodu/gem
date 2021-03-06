package pybind

import (
	"reflect"

	"github.com/qur/gopy/lib"
)

type Constructor func(*py.Type, *py.Tuple, *py.Dict) (py.Object, error)
type Wrapper func(*py.Tuple, *py.Dict) (py.Object, error)

func Define(name string, ptr interface{}) *py.Class {
	class := &py.Class{
		Name:    name,
		Flags:   py.TPFLAGS_BASETYPE,
		Pointer: ptr,
	}

	if initFn := getInit(ptr); initFn != nil {
		class.New = WrapConstructor(initFn)
	}

	return class
}

func getInit(ptr interface{}) interface{} {
	typ := reflect.TypeOf(ptr)
	initMethod, ok := typ.MethodByName("Init")
	if ok {
		initFn := initMethod.Func.Interface()

		return initFn.(interface{})
	}

	return nil
}

func GenerateRegisterFunc(def *py.Class) func(*py.Module) error {
	return func(module *py.Module) error {
		return Register(def, module)
	}
}

func GenerateConstructor(def *py.Class) interface{} {
	init := getInit(def.Pointer)
	if init == nil {
		panic("can't generate constructor: no Init method")
	}

	initType := reflect.TypeOf(init)
	initVal := reflect.ValueOf(init)
	constructedType := reflect.TypeOf(def.Pointer)

	// [1:] to remove the first arg of Init, which should be the Alloc()ated object
	inTypes := InTypes(initType)[1:]
	// return types of the constructor are the object + whatever init returns
	outTypes := OutTypes(initType)
	outTypes = append([]reflect.Type{constructedType}, outTypes...)

	constructorType := reflect.FuncOf(inTypes, outTypes, false)

	genericNew := func(args []reflect.Value) (results []reflect.Value) {
		lock := py.NewLock()
		defer lock.Unlock()

		pyObj, err := def.Alloc(0)
		if err != nil {
			// We assume that a failure to allocate a python object is unrecoverable
			panic(err)
		}

		// prepend the constructed object
		args = append([]reflect.Value{reflect.ValueOf(pyObj)}, args...)
		results = initVal.Call(args)

		// prepend the constructed object
		results = append([]reflect.Value{reflect.ValueOf(pyObj)}, results...)
		return results
	}

	constructorFn := reflect.MakeFunc(constructorType, genericNew)
	return constructorFn.Interface()
}

func Register(def *py.Class, module *py.Module) error {
	var err error
	var class *py.Type
	if class, err = def.Create(); err != nil {
		return err
	}

	if err = module.AddObject(def.Name, class); err != nil {
		return err
	}
	addToAll(module, def.Name)

	return nil
}

func addToAll(module *py.Module, name string) {
	var all *py.List
	dict := module.Dict()
	if obj, _ := dict.GetItemString("__all__"); obj != nil {
		all = obj.(*py.List)
	} else {
		var err error
		all, err = py.NewList(0)
		if err != nil {
			panic(err)
		}
		dict.SetItemString("__all__", all)
	}
	all.Incref()
	defer all.Decref()

	nameObj, err := py.NewString(name)
	if err != nil {
		panic(err)
	}
	all.Append(nameObj)
}

func WrapConstructor(fn interface{}) Constructor {
	wrapper := Wrap(fn)
	return func(pyType *py.Type, pyArgs *py.Tuple, pyKwds *py.Dict) (py.Object, error) {
		lock := py.NewLock()
		defer lock.Unlock()

		pyObj, err := pyType.Alloc(0)
		if err != nil {
			return nil, err
		}

		// prepend the object to pyArgs
		argsSlice := pyArgs.Slice()
		argsSlice = append([]py.Object{pyObj}, argsSlice...)
		pyArgs, err = py.PackTuple(argsSlice...)
		if err != nil {
			return nil, err
		}

		_, err = wrapper(pyArgs, pyKwds)
		if err != nil {
			return nil, err
		}

		return pyObj, nil
	}
}

func Wrap(fn interface{}) Wrapper {
	val := reflect.ValueOf(fn)
	typ := reflect.TypeOf(fn)

	inTypes := InTypes(typ)

	return func(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
		lock := py.NewLock()
		defer lock.Unlock()

		// Convert args to []Value
		convertedArgs, err := ConvertIn(inTypes, args)
		if err != nil {
			return nil, err
		}

		// Call the go function
		outs := val.Call(convertedArgs)

		// Convert outs to []py.Object
		convertedOuts, err := ConvertOut(outs)
		if err != nil {
			return nil, err
		}

		// Return
		if len(convertedOuts) == 1 {
			return convertedOuts[0], nil
		} else {
			return py.PackTuple(convertedOuts...)
		}
	}
}
