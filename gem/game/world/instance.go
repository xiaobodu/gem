package world

import (
	"github.com/qur/gopy/lib"

	"github.com/sinusoids/gem/gem/game/position"
)

type Instance struct {
	py.BaseObject

	sectors map[position.SectorHash]*Sector
}

func (i *Instance) Init() {
	i.sectors = make(map[position.SectorHash]*Sector)
}

// Sector returns a sector instance for a given position.Sector.
// If the sector is not currently active, the sector is instantiated first.
func (i *Instance) Sector(s *position.Sector) *Sector {
	hash := s.HashCode()
	if sector, ok := i.sectors[hash]; ok {
		return sector
	}

	// Sector not yet tracked; Create a new sector
	i.sectors[hash] = NewSector(s)
	return i.sectors[hash]
}

// GC removes any inactive sectors from memory
func (i *Instance) GC() {
	garbage := make([]position.SectorHash, 0)
	for hash, sector := range i.sectors {
		if !sector.Active() {
			garbage = append(garbage, hash)
		}
	}

	for _, hash := range garbage {
		delete(i.sectors, hash)
	}
}