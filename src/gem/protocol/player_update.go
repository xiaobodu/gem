// NOT Generated by bbc (but should be in future)
package protocol

import (
	"io"

	"gem/encoding"
	"gem/game/entity"
)

type PlayerUpdateBlock struct {
	OurPlayer entity.Player
}

func (struc *PlayerUpdateBlock) Encode(w io.Writer, flags interface{}) error {
	buf := encoding.NewBitBuffer(w)

	updateBlock := encoding.NewBuffer()
	struc.buildUpdateBlock(updateBlock, struc.OurPlayer)

	err := struc.buildMovementBlock(buf)
	if err != nil {
		return err
	}

	buf.Write(8, 0) // count of other players to update

	updateBlockBytes := updateBlock.Bytes()
	if len(updateBlockBytes) > 0 {
		buf.Write(11, 0x7FF)
		buf.Close()
		w.Write(updateBlockBytes)
	} else {
		buf.Close()
	}

	return nil
}

func (struc *PlayerUpdateBlock) Decode(buf io.Reader, flags interface{}) (err error) {
	panic("not implemented")
}

func (struc *PlayerUpdateBlock) buildMovementBlock(buf *encoding.BitBuffer) error {
	player := struc.OurPlayer
	flags := player.Flags()

	// Anything to do?
	if flags == 0 {
		buf.Write(1, 0) // No updates
		return nil
	}
	buf.Write(1, 1) // This player has updates

	// Do we have any non-movement updates to perform?
	otherUpdateFlags := (flags & ^entity.MobFlagMovementUpdate)

	switch {
	case (flags & entity.MobFlagRegionUpdate) != 0:
		localPos := player.Position().LocalTo(player.Region())

		buf.Write(2, 3) // update type 3 = warp to location
		buf.Write(2, uint32(localPos.Z))
		buf.WriteBit(true) // discard walk queue? not sure when/if we need this
		buf.WriteBit(otherUpdateFlags != 0)
		buf.Write(7, uint32(localPos.Y))
		buf.Write(7, uint32(localPos.X))

	case (flags & entity.MobFlagRunUpdate) != 0:
		current, last := player.WalkDirection()

		buf.Write(2, 2) // update type 2 = running
		buf.Write(3, uint32(last))
		buf.Write(3, uint32(current))
		buf.WriteBit(otherUpdateFlags != 0)

	case (flags & entity.MobFlagWalkUpdate) != 0:
		current, _ := player.WalkDirection()

		buf.Write(2, 1) // update type 1 = walking
		buf.Write(3, uint32(current))
		buf.WriteBit(otherUpdateFlags != 0)

	default:
		buf.Write(2, 0) // update type 0 = no movement updates
	}
	return nil
}

func (struc *PlayerUpdateBlock) buildUpdateBlock(w io.Writer, player entity.Player) error {
	flags := player.Flags() & ^entity.MobFlagMovementUpdate
	if flags == 0 {
		return nil
	}

	if flags >= 256 {
		flags |= 64
		flagsEnc := encoding.Int16(flags)
		err := flagsEnc.Encode(w, encoding.IntLittleEndian)
		if err != nil {
			return err
		}
	} else {
		flagsEnc := encoding.Int8(flags)
		err := flagsEnc.Encode(w, encoding.IntNilFlag)
		if err != nil {
			return err
		}
	}

	/* Update appearance */
	if (flags & entity.MobFlagIdentityUpdate) != 0 {
		buf := encoding.NewBuffer()
		appearance := player.Profile().Appearance
		anims := player.Profile().Animations
		appearanceBlock := OutboundPlayerAppearance{
			Gender:   encoding.Int8(appearance.Gender),
			HeadIcon: encoding.Int8(appearance.HeadIcon),

			HelmModel:       encoding.Int8(0),
			CapeModel:       encoding.Int8(0),
			AmuletModel:     encoding.Int8(0),
			RightWieldModel: encoding.Int8(0),
			TorsoModel:      encoding.Int16(256 + appearance.TorsoModel),
			LeftWieldModel:  encoding.Int8(0),
			ArmsModel:       encoding.Int16(256 + appearance.ArmsModel),
			LegsModel:       encoding.Int16(256 + appearance.LegsModel),
			HeadModel:       encoding.Int16(256 + appearance.HeadModel),
			HandsModel:      encoding.Int16(256 + appearance.HandsModel),
			FeetModel:       encoding.Int16(256 + appearance.FeetModel),
			BeardModel:      encoding.Int16(256 + appearance.BeardModel),

			HairColor:  encoding.Int8(appearance.HairColor),
			TorsoColor: encoding.Int8(appearance.TorsoColor),
			LegColor:   encoding.Int8(appearance.LegColor),
			FeetColor:  encoding.Int8(appearance.FeetColor),
			SkinColor:  encoding.Int8(appearance.SkinColor),

			AnimIdle:       encoding.Int16(anims.AnimIdle),
			AnimSpotRotate: encoding.Int16(anims.AnimSpotRotate),
			AnimWalk:       encoding.Int16(anims.AnimWalk),
			AnimRotate180:  encoding.Int16(anims.AnimRotate180),
			AnimRotateCCW:  encoding.Int16(anims.AnimRotateCCW),
			AnimRotateCW:   encoding.Int16(anims.AnimRotateCW),
			AnimRun:        encoding.Int16(anims.AnimRun),
		}

		err := appearanceBlock.Encode(buf, nil)
		if err != nil {
			return err
		}

		block := buf.Bytes()
		blockSize := encoding.Int8(len(block))
		err = blockSize.Encode(w, encoding.IntNegate)
		if err != nil {
			return err
		}

		_, err = w.Write(block)
		if err != nil {
			return err
		}
	}
	return nil
}

type PlayerUpdate PlayerUpdateBlock

var PlayerUpdateDefinition = encoding.PacketHeader{
	Type:   (*PlayerUpdate)(nil),
	Number: 81,
	Size:   encoding.SzVar16,
}

func (frm *PlayerUpdate) Encode(buf io.Writer, flags interface{}) (err error) {
	struc := (*PlayerUpdateBlock)(frm)
	hdr := encoding.PacketHeader{
		Number: PlayerUpdateDefinition.Number,
		Size:   PlayerUpdateDefinition.Size,
		Object: struc,
	}
	return hdr.Encode(buf, flags)
}

func (frm *PlayerUpdate) Decode(buf io.Reader, flags interface{}) (err error) {
	struc := (*PlayerUpdateBlock)(frm)
	hdr := encoding.PacketHeader{
		Number: PlayerUpdateDefinition.Number,
		Size:   PlayerUpdateDefinition.Size,
		Object: struc,
	}
	return hdr.Decode(buf, flags)
}
