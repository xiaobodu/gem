package player

import (
	"testing"

	"github.com/qur/gopy/lib"

	"github.com/sinusoids/gem/gem/game/interface/entity"
	"github.com/sinusoids/gem/gem/game/position"
)

func TestSnapshot(t *testing.T) {
	lock := py.NewLock()
	defer lock.Unlock()

	dummyPos, err := position.NewAbsolute(3500, 3500, 1)
	if err != nil {
		panic(err)
	}

	dummyPlayer := &PlayerSnapshot{
		flags:  entity.MobFlagWalkUpdate,
		region: dummyPos.RegionOf(),
		profile: &ProfileSnapshot{
			username: "a player",
			password: "hunter2",
			rights:   RightsPlayer,
			pos:      dummyPos,

			skills: &SkillsSnapshot{
				combatLevel: 8,
			},
			appearance: &AppearanceSnapshot{
				gender:   1,
				headIcon: 1,
				models:   map[BodyPart]int{},
				colors:   map[BodyPart]int{},
			},
			animations: &AnimationsSnapshot{
				anims: map[Anim]int{},
			},
		},
		session: &SessionSnapshot{},
		waypointQueue: &WaypointQueueSnapshot{
			currentDirection: 1,
			lastDirection:    2,
		},
	}

	snapshot := Snapshot(dummyPlayer)

	if !comparePlayers(dummyPlayer, snapshot) {
		t.Errorf("snapshotted player didn't match")
	}
}
