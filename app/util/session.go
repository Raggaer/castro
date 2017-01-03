package util

import (
	"github.com/raggaer/castro/app/models"
	"github.com/raggaer/castro/app/database"
	"github.com/jinzhu/gorm"
	"time"
	"bytes"
	"encoding/gob"
)

type Session struct {
	Token string
	Logged bool
	Flash map[string]string
	Data map[string]interface{}
}

func GetSession(token string) (*Session, error) {
	// Hold session information
	s := models.Session{}

	// Try to fetch data from database
	if err := database.DB.Table("sessions").Where("token = ?", token).Scan(&s).Error; err != nil {

		// If there is another error return
		if err != gorm.ErrRecordNotFound {

			return nil, err
		}

		// Create record
		s.Token = token
		s.CreatedAt = time.Now()

		// Create session data
		data, err := Encode(
			&Session{
				Token: s.Token,
				Data: make(map[string]interface{}),
				Flash: make(map[string]string),
			},
		)

		if err != nil {
			return nil, err
		}

		s.Data = data

		// Create record
		if err := database.DB.Create(&s).Error; err != nil {

			return nil, err
		}
	}

	// Decode data and return
	sess, err := Decode(s.Data)

	return sess, err
}

// Encode uses gob to encode the given session struct
// return a byte array
func Encode(s *Session) ([]byte, error) {
	// Byte buffer to hold data
	buffer := &bytes.Buffer{}

	// Create gob encoder
	enc := gob.NewEncoder(buffer)

	// Encode session struct
	err := enc.Encode(s)

	return buffer.Bytes(), err
}

func Decode(data []byte) (*Session, error) {
	// Struct to hold session data
	s := &Session{}

	// Byte buffer to hold data
	buffer := bytes.NewBuffer(data)

	// Create gob decoder
	dec := gob.NewDecoder(buffer)

	// Decode buffer
	err := dec.Decode(&s)

	return s, err
}

// Save stores the session data into database
func (s *Session) Save() error {
	// Encode session struct
	data, err := Encode(s)

	if err != nil {
		return err
	}

	// Populate model with session data
	sess := models.Session{
		Token: s.Token,
		Data: data,
		UpdatedAt: time.Now(),
	}

	// Try to update session data
	return database.DB.Table("sessions").Where("token = ?", s.Token).Update(&sess).Error
}

// Destroy removes the session data from the database
func (s *Session) Destroy() error {
	// Try to remove data from database
	return database.DB.Table("sessions").Delete("token = ?", s.Token).Error
}

// DeleteSession removes the given session data from the database
func DeleteSession(token string) error {
	// Try to remove data from database
	return database.DB.Table("sessions").Delete("token = ?", token).Error
}