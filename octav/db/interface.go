package db

import (
	"database/sql"
	"time"
)

type Config struct {
	DSN string // DSN, can be a template string
}

type InsertOption interface{}
type insertIgnoreOption bool

// +DB tablename:"conferences"
type Conference struct {
	OID        int64  // intenral id, used for sorting and what not
	EID        string // ID that is visible to the outside
	Slug       string
	Title      string
	SubTitle   sql.NullString
	CreatedBy  string // User ID that creates this conference
	CreatedOn  time.Time
	ModifiedOn NullTime
}

// +DB tablename:"conference_dates"
type ConferenceDate struct {
	OID          int64
	ConferenceID string
	Date         string
	Open         sql.NullString
	Close        sql.NullString
}

// +DB tablename:"conference_administrators"
type ConferenceAdministrator struct {
	OID          int64 // OID is the internal id, used for sorting and what not
	ConferenceID string
	UserID       string
	CreatedOn    time.Time
	ModifiedOn   NullTime
}

// +DB tablename:"conference_venues"
type ConferenceVenue struct {
	OID          int64 // OID is the internal id, used for sorting and what not
	ConferenceID string
	VenueID      string
	CreatedOn    time.Time
	ModifiedOn   NullTime
}

// +DB tablename:"rooms"
type Room struct {
	OID        int64  // intenral id, used for sorting and what not
	EID        string // ID that is visible to the outside
	VenueID    string // ID of the venue that this room belongs to
	Name       string // Name of the room (English)
	Capacity   uint   // How many people fit in this room? Approximation.
	CreatedOn  time.Time
	ModifiedOn NullTime
}

// +DB tablename:"sessions"
type Session struct {
	OID               int64          // OID is the internal id, used for sorting and what not
	EID               string         // EID is the ID that is visible to the outside
	ConferenceID      string         // ConferenceID is the ID of the conference that this session belongs to
	RoomID            sql.NullString // ID of the room where this session will be held at.
	SpeakerID         string         // ID of the speaker that this session belongs to
	Title             sql.NullString // Title of the session (English)
	Abstract          sql.NullString // Abstract of the session (English)
	Memo              sql.NullString // Correspondence between the speaker and the organizer. Should not be publicly available
	StartsOn          NullTime       // Time that this session is scheduled to start on
	Duration          int            // Length of this session in minutes.
	MaterialLevel     sql.NullString
	Tags              sql.NullString // Comma separated tags
	Category          sql.NullString
	SpokenLanguage    sql.NullString
	SlideLanguage     sql.NullString
	SlideSubtitles    sql.NullString
	SlideURL          sql.NullString
	VideoURL          sql.NullString
	PhotoPermission   sql.NullString
	VideoPermission   sql.NullString
	HasInterpretation bool
	Status            string
	SortOrder         int
	Confirmed         bool
	CreatedOn         time.Time
	ModifiedOn        NullTime
}

// +DB tablename:"users"
type User struct {
	OID        int64
	EID        string
	AuthVia    sql.NullString
	AuthUserID sql.NullString
	AvatarURL  sql.NullString
	FirstName  sql.NullString
	LastName   sql.NullString
	Nickname   string
	Email      sql.NullString
	TshirtSize sql.NullString
	IsAdmin    bool
	CreatedOn  time.Time
	ModifiedOn NullTime
}

// +DB tablename:"venues"
type Venue struct {
	OID        int64  // intenral id, used for sorting and what not
	EID        string // ID that is visible to the outside
	Name       string // Name of the venue (English)
	Address    string
	Latitude   float64
	Longitude  float64
	CreatedOn  time.Time
	ModifiedOn NullTime
}

// +DB tablename:"localized_strings"
type LocalizedString struct {
	OID        int64
	ParentID   string // EID of the parent object
	ParentType string // Type of the parent object
	Name       string
	Language   string
	Localized  string
}
