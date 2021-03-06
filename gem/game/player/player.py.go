package player

import (
	"github.com/qur/gopy/lib"

	"github.com/gemrs/gem/gem/game/server"
	"github.com/gemrs/gem/gem/game/world"
	"github.com/gemrs/gem/pybind"
)

var PlayerDef = pybind.Define("Player", (*Player)(nil))
var RegisterPlayer = pybind.GenerateRegisterFunc(PlayerDef)
var NewPlayer = pybind.GenerateConstructor(PlayerDef).(func(*server.Connection, *world.Instance) *Player)

func (client *Player) PyGet_username() (py.Object, error) {
	fn := pybind.Wrap(client.Profile().Username)
	return fn(nil, nil)
}

func (client *Player) Py_send_message(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	fn := pybind.Wrap(client.SendMessage)
	return fn(args, kwds)
}

func (client *Player) Py_warp(args *py.Tuple, kwds *py.Dict) (py.Object, error) {
	fn := pybind.Wrap(client.SetPosition)
	return fn(args, kwds)
}

func (client *Player) PyGet_skills() (py.Object, error) {
	fn := pybind.Wrap(client.Profile().Skills)
	return fn(nil, nil)
}

func (client *Player) PyGet_loaded_region() (py.Object, error) {
	fn := pybind.Wrap(client.LoadedRegion)
	return fn(nil, nil)
}

func (client *Player) PyGet_appearance() (py.Object, error) {
	fn := pybind.Wrap(client.Appearance)
	return fn(nil, nil)
}

func (client *Player) PySet_appearance(value py.Object) error {
	fn := pybind.Wrap(client.SetAppearance)
	args, err := py.PackTuple(value)
	if err != nil {
		return err
	}
	_, err = fn(args, nil)
	return err
}

func (client *Player) PyGet_entity_type() (py.Object, error) {
	fn := pybind.Wrap(client.EntityType)
	return fn(nil, nil)
}

func (client *Player) PyGet_index() (py.Object, error) {
	fn := pybind.Wrap(client.Index)
	return fn(nil, nil)
}
