////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles implementation of the database backend

package storage

import "github.com/pkg/errors"

// Obtain User from backend by primary key
func (db *DatabaseImpl) GetUser(userId string) (*User, error) {
	u := &User{
		Id: userId,
	}
	err := db.pg.Select(u)
	if err != nil {
		return nil, errors.Errorf("Failed to retrieve user with ID %s: %+v", userId, err)
	}
	return u, nil
}

// Delete User from backend by primary key
func (db *DatabaseImpl) DeleteUser(userId string) error {
	err := db.pg.Delete(&User{
		Id: userId,
	})
	if err != nil {
		return errors.Errorf("Failed to delete user with ID %s: %+v", userId, err)
	}
	return nil
}

// Insert or Update User into backend
func (db *DatabaseImpl) UpsertUser(user *User) error {
	expectedToken := user.Token
	_, err := db.pg.Model(user).
		OnConflict("(Id) DO UPDATE").
		Set("Token = EXCLUDED.Token").Insert()
	if err != nil {
		return errors.Errorf("Failed to insert user %s: %+v", user.Id, err)
	}

	err = db.pg.Select(user)
	if err != nil || expectedToken != user.Token {
		return errors.Errorf("User was not inserted properly: %+v", err)
	}
	return nil
}
