package notifications

import (
	"fmt"
	"gitlab.com/elixxir/notifications-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/id/ephemeral"
	"testing"
	"time"
)

func TestImpl_InitDeleter(t *testing.T) {
	s, err := storage.NewStorage("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to init storage: %+v", err)
	}
	impl := &Impl{
		Storage: s,
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	iid, err := ephemeral.GetIntermediaryId(uid)
	if err != nil {
		t.Errorf("Failed to get intermediary ephemeral id: %+v", err)
	}
	u, err := s.AddUser(iid, []byte("trsa"), []byte("Sig"), "token")
	if err != nil {
		t.Errorf("Failed to add user to storage: %+v", err)
	}
	_, epoch := ephemeral.HandleQuantization(time.Now().Add(-30 * time.Hour))
	e, err := s.AddLatestEphemeral(u, epoch, 16)
	if err != nil {
		t.Errorf("Failed to add latest ephemeral for user: %+v", err)
	}
	elist, err := s.GetEphemeral(e.EphemeralId)
	if err != nil {
		t.Errorf("Failed to get latest ephemeral for user: %+v", err)
	}
	if elist == nil {
		t.Error("Did not receive ephemeral for user")
	}
	impl.initDeleter()
	time.Sleep(time.Second * 5)
	elist, err = s.GetEphemeral(e.EphemeralId)
	if err == nil {
		t.Errorf("Ephemeral should have been deleted, did not receive error: %+v", e)
	}
}

func TestImpl_InitCreator(t *testing.T) {
	s, err := storage.NewStorage("", "", "", "", "")
	if err != nil {
		t.Errorf("Failed to init storage: %+v", err)
	}
	impl, err := StartNotifications(Params{
		Address:  "",
		CertPath: "",
		KeyPath:  "",
		FBCreds:  "",
	}, true, true)
	if err != nil {
		t.Errorf("Failed to create impl: %+v", err)
	}
	impl.Storage = s
	uid := id.NewIdFromString("zezima", id.User, t)
	iid, err := ephemeral.GetIntermediaryId(uid)
	if err != nil {
		t.Errorf("Failed to get intermediary ephemeral id: %+v", err)
	}
	u, err := s.AddUser(iid, []byte("trsa"), []byte("Sig"), "token")
	if err != nil {
		t.Errorf("Failed to add user to storage: %+v", err)
	}
	fmt.Println(u.OffsetNum)
	impl.initCreator()
	e, err := s.GetLatestEphemeral()
	if err != nil {
		t.Errorf("Failed to get latest ephemeral: %+v", err)
	}
	if e == nil {
		t.Error("Did not receive ephemeral for user")
	}
}
