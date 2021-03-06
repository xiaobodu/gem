package auth

import (
	"fmt"

	"github.com/qur/gopy/lib"

	"github.com/gemrs/gem/gem/game/interface/player"
)

//go:generate stringer -type=AuthResponse
type AuthResponse int

const (
	AuthPending AuthResponse = iota
	AuthDelay
	AuthOkay
	AuthInvalidCredentials
	AuthDisabled
	AuthDuplicateSession
	AuthUpdates
	AuthServerFull
	AuthNoLoginServer
	AuthTooManyConnections
	AuthBadSessionId
	AuthRejected
	AuthMembersWorld
	AuthIncomplete
	AuthUpdating
	AuthUnknown
	AuthAttemptsExceeded
	AuthMembersArea
	_
	_
	AuthInvalidLoginServer
	AuthInvalidTransferring
	AuthEnd
)

// Provider is a provider of authorization.
type Provider interface {
	LookupProfile(name, password string) (player.Profile, AuthResponse)
}

// ProviderImpl is a base class to be extended in Python
type ProviderImpl struct {
	py.BaseObject
}

func (p *ProviderImpl) Init() {}

func (p *ProviderImpl) LookupProfile(name, password string) (player.Profile, AuthResponse) {
	lock := py.NewLock()
	defer lock.Unlock()

	obj, err := p.CallMethod("lookup_profile", "(ss)", name, password)
	if err != nil {
		panic(fmt.Sprintf("cannot call LookupProfile: %v", err))
	}

	var profile py.Object
	var response int

	err = py.ParseTuple(obj.(*py.Tuple), "Oi", &profile, &response)
	if err != nil {
		return nil, AuthIncomplete
	}

	if AuthResponse(response) != AuthOkay {
		return nil, AuthResponse(response)
	}

	var profileIface interface{}
	profileIface = profile

	if p, ok := profileIface.(player.Profile); ok {
		return p, AuthResponse(response)
	}
	panic("invalid profile returned: not a player.Profile")
}
