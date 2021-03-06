package player

import (
	"github.com/gemrs/gem/gem/game/interface/entity"
)

// comparePlayers performs a deep comparison between two Players
// only used for testing
func comparePlayers(p1, p2 Player) bool {
	if !compareProfile(p1.Profile(), p2.Profile()) {
		return false
	}

	if !compareAnimations(p1.Animations(), p2.Animations()) {
		return false
	}

	if !compareWaypointQueue(p1.WaypointQueue(), p2.WaypointQueue()) {
		return false
	}

	if p1.Flags() != p2.Flags() {
		return false
	}

	if !p1.LoadedRegion().Compare(p2.LoadedRegion()) {
		return false
	}

	return true
}

func compareWaypointQueue(q1, q2 entity.WaypointQueue) bool {
	d1, d2 := q1.WalkDirection()
	d3, d4 := q1.WalkDirection()
	if d1 != d3 || d2 != d4 {
		return false
	}
	return true
}

func compareProfile(p1, p2 Profile) bool {
	if p1.Username() != p2.Username() {
		return false
	}

	if p1.Password() != p2.Password() {
		return false
	}

	if p1.Rights() != p2.Rights() {
		return false
	}

	if !p1.Position().Compare(p2.Position()) {
		return false
	}

	if !compareSkills(p1.Skills(), p2.Skills()) {
		return false
	}

	if !compareAppearance(p1.Appearance(), p2.Appearance()) {
		return false
	}

	return true
}

func compareSkills(s1, s2 Skills) bool {
	return s1.CombatLevel() == s2.CombatLevel()
}

func compareAppearance(a1, a2 Appearance) bool {
	if a1.Gender() != a2.Gender() {
		return false
	}

	if a1.HeadIcon() != a2.HeadIcon() {
		return false
	}

	for i := BodyPart(0); i < BodyPartMax; i++ {
		if a1.Model(i) != a2.Model(i) {
			return false
		}

		if a1.Color(i) != a2.Color(i) {
			return false
		}
	}

	return true
}

func compareAnimations(a1, a2 Animations) bool {
	for i := Anim(0); i < AnimMax; i++ {
		if a1.Animation(i) != a2.Animation(i) {
			return false
		}
	}

	return true
}
