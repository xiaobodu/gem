package player

import (
	"fmt"
	"github.com/qur/gopy/lib"
)

type Rights int

const (
	RightsPlayer Rights = iota
	RightsModerator
	RightsAdmin
)

//go:generate gopygen -type Profile $GOFILE
// Profile represents the saved state of a user
type Profile struct {
	py.BaseObject

	Username string
	Password string /* todo: hash */
	Rights   Rights
}

func (p *Profile) String() string {
	return fmt.Sprintf("Username: %v", p.Username)
}
