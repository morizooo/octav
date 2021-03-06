package model

// Automatically generated by genmodel utility. DO NOT EDIT!

import (
	"encoding/json"
	"time"

	"github.com/builderscon/octav/octav/db"
	"github.com/builderscon/octav/octav/tools"
	"github.com/lestrrat/go-pdebug"
)

var _ = time.Time{}

type ConferenceL10N struct {
	Conference
	L10N tools.LocalizedFields `json:"-"`
}
type ConferenceL10NList []ConferenceL10N

func (v ConferenceL10N) MarshalJSON() ([]byte, error) {
	buf, err := json.Marshal(v.Conference)
	if err != nil {
		return nil, err
	}
	return tools.MarshalJSONWithL10N(buf, v.L10N)
}

func (v *Conference) Load(tx *db.Tx, id string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("Conference.Load").BindError(&err)
		defer g.End()
	}
	vdb := db.Conference{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	if err := v.FromRow(vdb); err != nil {
		return err
	}
	return nil
}

func (v *Conference) FromRow(vdb db.Conference) error {
	v.ID = vdb.EID
	v.Title = vdb.Title
	if vdb.SubTitle.Valid {
		v.SubTitle = vdb.SubTitle.String
	}
	v.Slug = vdb.Slug
	return nil
}

func (v *Conference) ToRow(vdb *db.Conference) error {
	vdb.EID = v.ID
	vdb.Title = v.Title
	vdb.SubTitle.Valid = true
	vdb.SubTitle.String = v.SubTitle
	vdb.Slug = v.Slug
	return nil
}

func (v ConferenceL10N) GetPropNames() ([]string, error) {
	l, _ := v.L10N.GetPropNames()
	return append(l, "title", "sub_title"), nil
}

func (v ConferenceL10N) GetPropValue(s string) (interface{}, error) {
	switch s {
	case "id":
		return v.ID, nil
	case "title":
		return v.Title, nil
	case "sub_title":
		return v.SubTitle, nil
	case "slug":
		return v.Slug, nil
	case "dates":
		return v.Dates, nil
	case "administrators":
		return v.Administrators, nil
	case "venues":
		return v.Venues, nil
	default:
		return v.L10N.GetPropValue(s)
	}
}

func (v *ConferenceL10N) UnmarshalJSON(data []byte) error {
	var s Conference
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	v.Conference = s
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if err := tools.ExtractL10NFields(m, &v.L10N, []string{"title", "sub_title"}); err != nil {
		return err
	}

	return nil
}

func (v *ConferenceL10N) LoadLocalizedFields(tx *db.Tx) error {
	ls, err := db.LoadLocalizedStringsForParent(tx, v.Conference.ID, "Conference")
	if err != nil {
		return err
	}

	if len(ls) > 0 {
		v.L10N = tools.LocalizedFields{}
		for _, l := range ls {
			v.L10N.Set(l.Language, l.Name, l.Localized)
		}
	}
	return nil
}
