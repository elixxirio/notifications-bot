package storage

import (
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/id/ephemeral"
	"testing"
	"time"
)

func TestStorage_AddUser(t *testing.T) {
	s, err := NewStorage("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create new storage object: %+v", err)
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	iid, err := ephemeral.GetIntermediaryId(uid)
	if err != nil {
		t.Errorf("Failed to create iid: %+v", err)
	}
	testTime, err := time.Parse(time.RFC3339,
		"2012-12-21T22:08:41+00:00")
	if err != nil {
		t.Errorf("Could not parse precanned time: %v", err.Error())
	}
	_, err = s.AddUser(iid, []byte("transmissionrsa"), []byte("signature"),testTime, "token")
	if err != nil {
		t.Errorf("Failed to add user: %+v", err)
	}
}

func TestStorage_DeleteUser(t *testing.T) {
	s, err := NewStorage("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create new storage object: %+v", err)
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	iid, err := ephemeral.GetIntermediaryId(uid)
	if err != nil {
		t.Errorf("Failed to create iid: %+v", err)
	}
	testTime, err := time.Parse(time.RFC3339,
		"2012-12-21T22:08:41+00:00")
	if err != nil {
		t.Errorf("Could not parse precanned time: %v", err.Error())
	}
	u, err := s.AddUser(iid, []byte("transmissionrsa"), []byte("signature"),testTime, "token")
	if err != nil {
		t.Errorf("Failed to add user: %+v", err)
	}
	err = s.DeleteUser(u.TransmissionRSA)
	if err != nil {
		t.Errorf("Failed to delete user: %+v", err)
	}
}

func TestStorage_AddLatestEphemeral(t *testing.T) {
	s, err := NewStorage("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create new storage object: %+v", err)
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	iid, err := ephemeral.GetIntermediaryId(uid)
	if err != nil {
		t.Errorf("Failed to create iid: %+v", err)
	}
	testTime, err := time.Parse(time.RFC3339,
		"2012-12-21T22:08:41+00:00")
	if err != nil {
		t.Errorf("Could not parse precanned time: %v", err.Error())
	}
	u, err := s.AddUser(iid, []byte("transmissionrsa"), []byte("signature"),testTime, "token")
	if err != nil {
		t.Errorf("Failed to add user: %+v", err)
	}
	_, err = s.AddLatestEphemeral(u, 5, 16)
	if err != nil {
		t.Errorf("Failed to add latest ephemeral: %+v", err)
	}
}

func TestStorage_AddEphemeralsForOffset(t *testing.T) {
	_, err := NewStorage("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to create new storage object: %+v", err)
	}
}
