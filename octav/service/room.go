package service

import (
	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
)

func (v *Room) populateRowForCreate(vdb *db.Room, payload model.CreateRoomRequest) error {
	vdb.EID = tools.UUID()

	if payload.VenueID.Valid() {
		vdb.VenueID = payload.VenueID.String
	}

	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.Capacity.Valid() {
		vdb.Capacity = uint(payload.Capacity.Uint)
	}

	return nil
}

func (v *Room) populateRowForUpdate(vdb *db.Room, payload model.UpdateRoomRequest) error {
	if payload.VenueID.Valid() {
		vdb.VenueID = payload.VenueID.String
	}

	if payload.Name.Valid() {
		vdb.Name = payload.Name.String
	}

	if payload.Capacity.Valid() {
		vdb.Capacity = uint(payload.Capacity.Uint)
	}

	return nil
}