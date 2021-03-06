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

type UserL10N struct {
	User
	L10N tools.LocalizedFields `json:"-"`
}
type UserL10NList []UserL10N

func (v UserL10N) MarshalJSON() ([]byte, error) {
	buf, err := json.Marshal(v.User)
	if err != nil {
		return nil, err
	}
	return tools.MarshalJSONWithL10N(buf, v.L10N)
}

func (v *User) Load(tx *db.Tx, id string) (err error) {
	if pdebug.Enabled {
		g := pdebug.Marker("User.Load").BindError(&err)
		defer g.End()
	}
	vdb := db.User{}
	if err := vdb.LoadByEID(tx, id); err != nil {
		return err
	}

	if err := v.FromRow(vdb); err != nil {
		return err
	}
	return nil
}

func (v *User) FromRow(vdb db.User) error {
	v.ID = vdb.EID
	if vdb.AuthVia.Valid {
		v.AuthVia = vdb.AuthVia.String
	}
	if vdb.AuthUserID.Valid {
		v.AuthUserID = vdb.AuthUserID.String
	}
	if vdb.AvatarURL.Valid {
		v.AvatarURL = vdb.AvatarURL.String
	}
	if vdb.FirstName.Valid {
		v.FirstName = vdb.FirstName.String
	}
	if vdb.LastName.Valid {
		v.LastName = vdb.LastName.String
	}
	v.Nickname = vdb.Nickname
	if vdb.Email.Valid {
		v.Email = vdb.Email.String
	}
	if vdb.TshirtSize.Valid {
		v.TshirtSize = vdb.TshirtSize.String
	}
	v.IsAdmin = vdb.IsAdmin
	return nil
}

func (v *User) ToRow(vdb *db.User) error {
	vdb.EID = v.ID
	vdb.AuthVia.Valid = true
	vdb.AuthVia.String = v.AuthVia
	vdb.AuthUserID.Valid = true
	vdb.AuthUserID.String = v.AuthUserID
	vdb.AvatarURL.Valid = true
	vdb.AvatarURL.String = v.AvatarURL
	vdb.FirstName.Valid = true
	vdb.FirstName.String = v.FirstName
	vdb.LastName.Valid = true
	vdb.LastName.String = v.LastName
	vdb.Nickname = v.Nickname
	vdb.Email.Valid = true
	vdb.Email.String = v.Email
	vdb.TshirtSize.Valid = true
	vdb.TshirtSize.String = v.TshirtSize
	vdb.IsAdmin = v.IsAdmin
	return nil
}

func (v UserL10N) GetPropNames() ([]string, error) {
	l, _ := v.L10N.GetPropNames()
	return append(l, "first_name", "last_name"), nil
}

func (v UserL10N) GetPropValue(s string) (interface{}, error) {
	switch s {
	case "id":
		return v.ID, nil
	case "auth_via":
		return v.AuthVia, nil
	case "auth_user_id":
		return v.AuthUserID, nil
	case "avatar_url":
		return v.AvatarURL, nil
	case "first_name":
		return v.FirstName, nil
	case "last_name":
		return v.LastName, nil
	case "nickname":
		return v.Nickname, nil
	case "email":
		return v.Email, nil
	case "tshirt_size":
		return v.TshirtSize, nil
	case "is_admin":
		return v.IsAdmin, nil
	default:
		return v.L10N.GetPropValue(s)
	}
}

func (v *UserL10N) UnmarshalJSON(data []byte) error {
	var s User
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	v.User = s
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if err := tools.ExtractL10NFields(m, &v.L10N, []string{"first_name", "last_name"}); err != nil {
		return err
	}

	return nil
}

func (v *UserL10N) LoadLocalizedFields(tx *db.Tx) error {
	ls, err := db.LoadLocalizedStringsForParent(tx, v.User.ID, "User")
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
