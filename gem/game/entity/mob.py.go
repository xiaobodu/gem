package entity

import (
	"github.com/qur/gopy/lib"

	"github.com/gemrs/gem/gem/game/interface/entity"
	"github.com/gemrs/gem/pybind"
)

var GenericMobDef = pybind.Define("GenericMob", (*GenericMob)(nil))
var RegisterGenericMob = pybind.GenerateRegisterFunc(GenericMobDef)
var NewGenericMob = pybind.GenerateConstructor(GenericMobDef).(func(entity.WaypointQueue) *GenericMob)

func (mob *GenericMob) PyGet_position() (py.Object, error) {
	fn := pybind.Wrap(mob.Position)
	return fn(nil, nil)
}
