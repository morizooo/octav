package service

import (
	"database/sql"
	"strings"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
)

func (v *Conference) populateRowForCreate(vdb *db.Conference, payload model.CreateConferenceRequest) error {
	vdb.EID = tools.UUID()
	vdb.Slug = payload.Slug
	vdb.Title = payload.Title

	if payload.SubTitle.Valid() {
		vdb.SubTitle.Valid = true
		vdb.SubTitle.String = payload.SubTitle.String
	}
	return nil
}

func (v *Conference) populateRowForUpdate(vdb *db.Conference, payload model.UpdateConferenceRequest) error {
	if payload.Slug.Valid() {
		vdb.Slug = payload.Slug.String
	}

	if payload.Title.Valid() {
		vdb.Title = payload.Title.String
	}

	if payload.SubTitle.Valid() {
		vdb.SubTitle.Valid = true
		vdb.SubTitle.String = payload.SubTitle.String
	}
	return nil
}

func (v *Conference) AddAdministrator(tx *db.Tx, cid, uid string) error {
	c := db.ConferenceAdministrator{
		ConferenceID: cid,
		UserID:       uid,
	}
	return c.Create(tx)
}

const datefmt = `2006-01-02`

func (v *Conference) LoadByRange(tx *db.Tx, vdbl *db.ConferenceList, since, lang, rangeStart, rangeEnd string, limit int) error {
	var rs time.Time
	var re time.Time
	var err error

	if rangeStart != "" {
		rs, err = time.Parse(datefmt, rangeStart)
		if err != nil {
			return err
		}
	}

	if rangeEnd != "" {
		re, err = time.Parse(datefmt, rangeEnd)
		if err != nil {
			return err
		}
	}

	if err := vdbl.LoadByRange(tx, since, rs, re, limit); err != nil {
		return err
	}

	return nil
}

func (v *Conference) AddDates(tx *db.Tx, cid string, dates ...model.ConferenceDate) error {
	for _, date := range dates {
		cd := db.ConferenceDate{
			ConferenceID: cid,
			Date:         date.Date.String(),
			Open:         sql.NullString{String: date.Open.String(), Valid: true},
			Close:        sql.NullString{String: date.Close.String(), Valid: true},
		}
		if err := cd.Create(tx, db.WithInsertIgnore(true)); err != nil {
			return err
		}
	}

	return nil
}

func (v *Conference) DeleteDates(tx *db.Tx, cid string, dates ...model.Date) error {
	vdb := db.ConferenceDate{}
	sdatelist := make([]string, len(dates))
	for i, dt := range dates {
		sdatelist[i] = dt.String()
	}
	return vdb.DeleteDates(tx, cid, sdatelist...)
}

func (v *Conference) LoadDates(tx *db.Tx, cdl *model.ConferenceDateList, cid string) error {
	vdbl := db.ConferenceDateList{}
	if err := vdbl.LoadByConferenceID(tx, cid); err != nil {
		return err
	}

	res := make(model.ConferenceDateList, len(vdbl))
	for i, vdb := range vdbl {
		dt := vdb.Date
		if i := strings.IndexByte(dt, 'T'); i > -1 { // Cheat. Loading from DB contains time....!!!!
			dt = dt[:i]
		}
		if err := res[i].Date.Parse(dt); err != nil {
			return err
		}

		if vdb.Open.Valid {
			t := vdb.Open.String
			if len(t) > 5 {
				t = t[:5]
			}
			if err := res[i].Open.Parse(t); err != nil {
				return err
			}
		}

		if vdb.Close.Valid {
			t := vdb.Close.String
			if len(t) > 5 {
				t = t[:5]
			}
			if err := res[i].Close.Parse(t); err != nil {
				return err
			}
		}
	}
	*cdl = res
	return nil
}

func (v *Conference) AddAdmin(tx *db.Tx, cid, uid string) error {
	cd := db.ConferenceAdministrator{
		ConferenceID: cid,
		UserID:       uid,
	}
	if err := cd.Create(tx, db.WithInsertIgnore(true)); err != nil {
		return err
	}

	return nil
}

func (v *Conference) DeleteAdmin(tx *db.Tx, cid, uid string) error {
	return db.DeleteConferenceAdministrator(tx, cid, uid)
}

func (v *Conference) LoadAdmins(tx *db.Tx, cdl *model.UserList, cid string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.LoadAdmins").BindError(&err)
		defer g.End()
	}

	var vdbl db.UserList
	if err := db.LoadConferenceAdministrators(tx, &vdbl, cid); err != nil {
		return err
	}

	if pdebug.Enabled {
		pdebug.Printf("Loaded %d admins", len(vdbl))
	}

	res := make(model.UserList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.User
		if err := u.FromRow(vdb); err != nil {
			return err
		}
		res[i] = u
	}
	*cdl = res
	return nil
}

func (v *Conference) AddVenue(tx *db.Tx, cid, vid string) error {
	cd := db.ConferenceVenue{
		ConferenceID: cid,
		VenueID:      vid,
	}
	if err := cd.Create(tx, db.WithInsertIgnore(true)); err != nil {
		return err
	}

	return nil
}

func (v *Conference) DeleteVenue(tx *db.Tx, cid, uid string) error {
	return db.DeleteConferenceVenue(tx, cid, uid)
}

func (v *Conference) LoadVenues(tx *db.Tx, cdl *model.VenueList, cid string) error {
	var vdbl db.VenueList
	if err := db.LoadConferenceVenues(tx, &vdbl, cid); err != nil {
		return err
	}

	res := make(model.VenueList, len(vdbl))
	for i, vdb := range vdbl {
		var u model.Venue
		if err := u.FromRow(vdb); err != nil {
			return err
		}
		res[i] = u
	}
	*cdl = res
	return nil
}

func (v *Conference) Decorate(tx *db.Tx, c *model.Conference) error {
	if err := v.LoadDates(tx, &c.Dates, c.ID); err != nil {
		return err
	}

	if err := v.LoadAdmins(tx, &c.Administrators, c.ID); err != nil {
		return err
	}

	if err := v.LoadVenues(tx, &c.Venues, c.ID); err != nil {
		return err
	}
	return nil
}
