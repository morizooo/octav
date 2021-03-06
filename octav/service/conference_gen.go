package service

// Automatically generated by genmodel utility. DO NOT EDIT!

import (
	"errors"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/model"
	"github.com/lestrrat/go-pdebug"
)

var _ = time.Time{}

// Create takes in the transaction, the incoming payload, and a reference to
// a database row. The database row is initialized/populated so that the
// caller can use it afterwards
func (v *Conference) Create(tx *db.Tx, vdb *db.Conference, payload model.CreateConferenceRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.Create").BindError(&err)
		defer g.End()
	}

	if err := v.populateRowForCreate(vdb, payload); err != nil {
		return err
	}

	if err := vdb.Create(tx); err != nil {
		return err
	}

	if err := payload.L10N.CreateLocalizedStrings(tx, "Conference", vdb.EID); err != nil {
		return err
	}
	return nil
}

func (v *Conference) Update(tx *db.Tx, vdb *db.Conference, payload model.UpdateConferenceRequest) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.Update (%s)", vdb.EID).BindError(&err)
		defer g.End()
	}

	if vdb.EID == "" {
		return errors.New("vdb.EID is required (did you forget to call vdb.Load(tx) before hand?)")
	}

	if err := v.populateRowForUpdate(vdb, payload); err != nil {
		return err
	}

	if err := vdb.Update(tx); err != nil {
		return err
	}

	return payload.L10N.Foreach(func(l, k, x string) error {
		if pdebug.Enabled {
			pdebug.Printf("Updating l10n string for '%s' (%s)", k, l)
		}
		ls := db.LocalizedString{
			ParentType: "Conference",
			ParentID:   vdb.EID,
			Language:   l,
			Name:       k,
			Localized:  x,
		}
		return ls.Upsert(tx)
	})
}

func (v *Conference) ReplaceL10NStrings(tx *db.Tx, m *model.Conference, lang string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("service.Conference.ReplaceL10NStrings")
		defer g.End()
	}
	rows, err := tx.Query(`SELECT oid, parent_id, parent_type, name, language, localized FROM localized_strings WHERE parent_type = ? AND parent_id = ? AND language = ?`, "Conference", m.ID, lang)
	if err != nil {
		return err
	}

	var l db.LocalizedString
	for rows.Next() {
		if err := l.Scan(rows); err != nil {
			return err
		}

		switch l.Name {
		case "title":
			if pdebug.Enabled {
				pdebug.Printf("Replacing for key 'title'")
			}
			m.Title = l.Localized
		case "sub_title":
			if pdebug.Enabled {
				pdebug.Printf("Replacing for key 'sub_title'")
			}
			m.SubTitle = l.Localized
		}
	}
	return nil
}

func (v *Conference) Delete(tx *db.Tx, id string) error {
	if pdebug.Enabled {
		g := pdebug.Marker("Conference.Delete (%s)", id)
		defer g.End()
	}

	vdb := db.Conference{EID: id}
	if err := vdb.Delete(tx); err != nil {
		return err
	}
	if err := db.DeleteLocalizedStringsForParent(tx, id, "Conference"); err != nil {
		return err
	}
	return nil
}

func (v *Conference) LoadList(tx *db.Tx, vdbl *db.ConferenceList, since string, limit int) error {
	return vdbl.LoadSinceEID(tx, since, limit)
}
