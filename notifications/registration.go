////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package notifications

import (
	"encoding/base64"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/notifications"
	"gitlab.com/elixxir/crypto/registration"
	"gitlab.com/elixxir/crypto/rsa"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/id/ephemeral"
	"time"
)

var timestampError = "Timestamp of request must be within last 5 seconds.  Request timestamp: %s, current time: %s"

// RegisterToken registers the given token. It evaluates that the TransmissionRsaRegistarSig is
// correct. The RSA->PEM relationship is one to many. It will succeed if the token is already
// registered.
func (nb *Impl) RegisterToken(msg *pb.RegisterTokenRequest) error {
	jww.INFO.Println("RegisterToken")
	requestTimestamp := time.Unix(0, msg.RequestTimestamp)
	if time.Now().Sub(requestTimestamp) > time.Second*5 {
		return errors.Errorf(timestampError, requestTimestamp.String(), time.Now().String())
	}
	// Verify permissioning RSA signature
	permHost, ok := nb.Comms.GetHost(&id.Permissioning)
	if !ok {
		return errors.New("Could not find permissioning host to verify client signature")
	}
	jww.INFO.Printf("Verifying perm sig with params:\n\tPubKey: %s\n\tTimestamp: %d\n\tTRSA: %s\n\tSIG: %s\n", base64.StdEncoding.EncodeToString(permHost.GetPubKey().Bytes()), msg.RegistrationTimestamp, base64.StdEncoding.EncodeToString(msg.TransmissionRsaPem), base64.StdEncoding.EncodeToString(msg.TransmissionRsaRegistrarSig))
	err := registration.VerifyWithTimestamp(permHost.GetPubKey(), msg.RegistrationTimestamp,
		string(msg.TransmissionRsaPem), msg.TransmissionRsaRegistrarSig)
	if err != nil {
		return errors.WithMessage(err, "Failed to verify permissioning signature")
	}

	// Verify token signature
	pub, err := rsa.GetScheme().UnmarshalPublicKeyPEM(msg.TransmissionRsaPem)
	if err != nil {
		return errors.WithMessage(err, "Failed to unmarshal public key")
	}
	err = notifications.VerifyToken(pub, msg.Token, msg.App, requestTimestamp, notifications.RegisterTokenTag, msg.TokenSignature)
	if err != nil {
		return errors.WithMessage(err, "Failed to verify token signature")
	}

	return nb.Storage.RegisterToken(msg.Token, msg.App, msg.TransmissionRsaPem)
}

// RegisterTrackedID registers the given ID to be tracked. The request is signed
// Returns an error if TransmissionRSA is not registered with a valid token.
// The actual ID is not revealed, instead an intermediary value is sent which cannot
// be revered to get the ID, but is repeatable. So it can be rainbow-tabled.
func (nb *Impl) RegisterTrackedID(msg *pb.RegisterTrackedIdRequest) error {
	jww.INFO.Println("RegisterTrackedID")
	requestTimestamp := time.Unix(0, msg.Request.RequestTimestamp)
	if time.Now().Sub(requestTimestamp) > time.Second*5 {
		return errors.Errorf(timestampError, requestTimestamp.String(), time.Now().String())
	}

	// Verify permissioning RSA signature
	permHost, ok := nb.Comms.GetHost(&id.Permissioning)
	if !ok {
		return errors.New("Could not find permissioning host to verify client signature")
	}
	jww.INFO.Printf("Verifying perm sig with params:\n\tPubKey: %s\n\tTimestamp: %d\n\tTRSA: %s\n\tSIG: %s\n", base64.StdEncoding.EncodeToString(permHost.GetPubKey().Bytes()), msg.RegistrationTimestamp, base64.StdEncoding.EncodeToString(msg.Request.TransmissionRsaPem), base64.StdEncoding.EncodeToString(msg.TransmissionRsaRegistrarSig))
	err := registration.VerifyWithTimestamp(permHost.GetPubKey(), msg.RegistrationTimestamp,
		string(msg.Request.TransmissionRsaPem), msg.TransmissionRsaRegistrarSig)
	if err != nil {
		return errors.WithMessage(err, "Failed to verify permissioning signature")
	}

	pub, err := rsa.GetScheme().UnmarshalPublicKeyPEM(msg.Request.TransmissionRsaPem)
	if err != nil {
		return errors.WithMessage(err, "Failed to unmarshal public key")
	}

	err = notifications.VerifyIdentity(pub, msg.Request.TrackedIntermediaryID, requestTimestamp, notifications.RegisterTrackedIDTag, msg.Request.Signature)
	if err != nil {
		return errors.WithMessage(err, "Failed to verify identity signature")
	}
	_, epoch := ephemeral.HandleQuantization(time.Now())

	return nb.Storage.RegisterTrackedID(msg.Request.TrackedIntermediaryID, msg.Request.TransmissionRsaPem, epoch, nb.inst.GetPartialNdf().Get().AddressSpace[0].Size)
}

// UnregisterToken unregisters the given device token. The request is signed.
// Does not return an error if the token cannot be found
func (nb *Impl) UnregisterToken(msg *pb.UnregisterTokenRequest) error {
	jww.INFO.Println("UnregisterToken")
	requestTimestamp := time.Unix(0, msg.RequestTimestamp)
	if time.Now().Sub(requestTimestamp) > time.Second*5 {
		return errors.Errorf(timestampError, requestTimestamp.String(), time.Now().String())
	}

	pub, err := rsa.GetScheme().UnmarshalPublicKeyPEM(msg.TransmissionRsaPem)
	if err != nil {
		return errors.WithMessage(err, "Failed to unmarshal public key")
	}

	err = notifications.VerifyToken(pub, msg.Token, msg.App, requestTimestamp, notifications.UnregisterTokenTag, msg.TokenSignature)
	if err != nil {
		return errors.WithMessage(err, "Failed to verify token signature")
	}

	return nb.Storage.UnregisterToken(msg.Token, msg.TransmissionRsaPem)
}

// UnregisterTrackedID unregisters the given tracked ID. The request is signed.
// Does not return an error if the ID cannot be found
func (nb *Impl) UnregisterTrackedID(msg *pb.TrackedIntermediaryIdRequest) error {
	jww.INFO.Println("UnregisterTrackedID")
	requestTimestamp := time.Unix(0, msg.RequestTimestamp)
	if time.Now().Sub(requestTimestamp) > time.Second*5 {
		return errors.Errorf(timestampError, requestTimestamp.String(), time.Now().String())
	}

	pub, err := rsa.GetScheme().UnmarshalPublicKeyPEM(msg.TransmissionRsaPem)
	if err != nil {
		return errors.WithMessage(err, "Failed to unmarshal public key")
	}

	err = notifications.VerifyIdentity(pub, msg.TrackedIntermediaryID, requestTimestamp, notifications.UnregisterTrackedIDTag, msg.Signature)
	if err != nil {
		return errors.WithMessage(err, "Failed to verify identity signature")
	}

	return nb.Storage.UnregisterTrackedIDs(msg.TrackedIntermediaryID, msg.TransmissionRsaPem)
}
