package dbrepo

import (
	"errors"
	"time"

	"github.com/rohit-010/go-bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into a database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if room id = 2 then fail ; otherwise pass
	if res.RoomID == 2 {
		return 1, errors.New("Some error")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into a database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("some error")
	}

	return nil
}

// SearchAvailibilityByDatesByRoomID returns true if availibilty exists for roomID
//
//	and false if no availibity exist
func (m *testDBRepo) SearchAvailibilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return false, nil
}

// SearchAvailibilityForAllRooms returns a slice of available rooms if any for given date range
func (m *testDBRepo) SearchAvailibilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil

}

// GetRoomByID get room by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id > 2 {
		return room, errors.New("Some error")
	}

	return room, nil
}
