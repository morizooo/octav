package db

// Automatically generated by gendb utility. DO NOT EDIT!

import (
	"bytes"
	"database/sql"
	"errors"
	"strconv"
	"time"
)

const UserStdSelectColumns = "users.oid, users.eid, users.first_name, users.last_name, users.nickname, users.email, users.tshirt_size, users.created_on, users.modified_on"
const UserTable = "users"

type UserList []User

func (u *User) Scan(scanner interface {
	Scan(...interface{}) error
}) error {
	return scanner.Scan(&u.OID, &u.EID, &u.FirstName, &u.LastName, &u.Nickname, &u.Email, &u.TshirtSize, &u.CreatedOn, &u.ModifiedOn)
}

func (u *User) LoadByEID(tx *Tx, eid string) error {
	row := tx.QueryRow(`SELECT `+UserStdSelectColumns+` FROM `+UserTable+` WHERE users.eid = ?`, eid)
	if err := u.Scan(row); err != nil {
		return err
	}
	return nil
}

func (u *User) Create(tx *Tx, opts ...InsertOption) error {
	if u.EID == "" {
		return errors.New("create: non-empty EID required")
	}

	u.CreatedOn = time.Now()
	doIgnore := false
	for _, opt := range opts {
		switch opt.(type) {
		case insertIgnoreOption:
			doIgnore = true
		}
	}

	stmt := bytes.Buffer{}
	stmt.WriteString("INSERT ")
	if doIgnore {
		stmt.WriteString("IGNORE ")
	}
	stmt.WriteString("INTO ")
	stmt.WriteString(UserTable)
	stmt.WriteString(` (eid, first_name, last_name, nickname, email, tshirt_size, created_on, modified_on) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	result, err := tx.Exec(stmt.String(), u.EID, u.FirstName, u.LastName, u.Nickname, u.Email, u.TshirtSize, u.CreatedOn, u.ModifiedOn)
	if err != nil {
		return err
	}

	lii, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.OID = lii
	return nil
}

func (u User) Update(tx *Tx) error {
	if u.OID != 0 {
		_, err := tx.Exec(`UPDATE `+UserTable+` SET eid = ?, first_name = ?, last_name = ?, nickname = ?, email = ?, tshirt_size = ? WHERE oid = ?`, u.EID, u.FirstName, u.LastName, u.Nickname, u.Email, u.TshirtSize, u.OID)
		return err
	}
	if u.EID != "" {
		_, err := tx.Exec(`UPDATE `+UserTable+` SET first_name = ?, last_name = ?, nickname = ?, email = ?, tshirt_size = ? WHERE eid = ?`, u.FirstName, u.LastName, u.Nickname, u.Email, u.TshirtSize, u.EID)
		return err
	}
	return errors.New("either OID/EID must be filled")
}

func (u User) Delete(tx *Tx) error {
	if u.OID != 0 {
		_, err := tx.Exec(`DELETE FROM `+UserTable+` WHERE oid = ?`, u.OID)
		return err
	}

	if u.EID != "" {
		_, err := tx.Exec(`DELETE FROM `+UserTable+` WHERE eid = ?`, u.EID)
		return err
	}

	return errors.New("either OID/EID must be filled")
}

func (v *UserList) FromRows(rows *sql.Rows, capacity int) error {
	var res []User
	if capacity > 0 {
		res = make([]User, 0, capacity)
	} else {
		res = []User{}
	}

	for rows.Next() {
		vdb := User{}
		if err := vdb.Scan(rows); err != nil {
			return err
		}
		res = append(res, vdb)
	}
	*v = res
	return nil
}

func (v *UserList) LoadSinceEID(tx *Tx, since string, limit int) error {
	var s int64
	if id := since; id != "" {
		vdb := User{}
		if err := vdb.LoadByEID(tx, id); err != nil {
			return err
		}

		s = vdb.OID
	}
	return v.LoadSince(tx, s, limit)
}

func (v *UserList) LoadSince(tx *Tx, since int64, limit int) error {
	rows, err := tx.Query(`SELECT `+UserStdSelectColumns+` FROM `+UserTable+` WHERE users.oid > ? ORDER BY oid ASC LIMIT `+strconv.Itoa(limit), since)
	if err != nil {
		return err
	}

	if err := v.FromRows(rows, limit); err != nil {
		return err
	}
	return nil
}
