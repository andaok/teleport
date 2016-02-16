/*
Copyright 2015 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"time"

	"github.com/gravitational/teleport/lib/backend/encryptedbk/encryptor"
	"github.com/gravitational/teleport/lib/events"
	"github.com/gravitational/teleport/lib/recorder"
	"github.com/gravitational/teleport/lib/services"
	"github.com/gravitational/teleport/lib/session"

	"github.com/codahale/lunk"
	websession "github.com/gravitational/session"
	"golang.org/x/crypto/ssh"
)

type AuthWithRoles struct {
	authServer  *AuthServer
	permChecker PermissionChecker
	elog        events.Log
	sessions    session.SessionServer
	role        string
	recorder    recorder.Recorder
}

func NewAuthWithRoles(authServer *AuthServer, permChecker PermissionChecker,
	elog events.Log, sessions session.SessionServer,
	role string, recorder recorder.Recorder) *AuthWithRoles {

	return &AuthWithRoles{
		authServer:  authServer,
		permChecker: permChecker,
		sessions:    sessions,
		role:        role,
		recorder:    recorder,
	}
}

func (a *AuthWithRoles) GetSessions() ([]session.Session, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetSessions); err != nil {
		return nil, err
	} else {
		return a.sessions.GetSessions()
	}
}

func (a *AuthWithRoles) GetSession(id string) (*session.Session, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetSession); err != nil {
		return nil, err
	} else {
		return a.sessions.GetSession(id)
	}
}
func (a *AuthWithRoles) DeleteSession(id string) error {
	if err := a.permChecker.HasPermission(a.role, ActionDeleteSession); err != nil {
		return err
	} else {
		return a.sessions.DeleteSession(id)
	}
}
func (a *AuthWithRoles) UpsertSession(id string, ttl time.Duration) error {
	if err := a.permChecker.HasPermission(a.role, ActionUpsertSession); err != nil {
		return err
	} else {
		return a.sessions.UpsertSession(id, ttl)
	}
}
func (a *AuthWithRoles) UpsertParty(id string, p session.Party, ttl time.Duration) error {
	if err := a.permChecker.HasPermission(a.role, ActionUpsertParty); err != nil {
		return err
	} else {
		return a.sessions.UpsertParty(id, p, ttl)
	}
}
func (a *AuthWithRoles) UpsertRemoteCertificate(cert services.CertificateAuthority, ttl time.Duration) error {
	if err := a.permChecker.HasPermission(a.role, ActionUpsertRemoteCertificate); err != nil {
		return err
	} else {
		return a.authServer.UpsertRemoteCertificate(cert, ttl)
	}
}
func (a *AuthWithRoles) GetRemoteCertificates(ctype string, domainName string) ([]services.CertificateAuthority, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetRemoteCertificates); err != nil {
		return nil, err
	} else {
		return a.authServer.GetRemoteCertificates(ctype, domainName)
	}
}
func (a *AuthWithRoles) GetTrustedCertificates(ctype string) ([]services.CertificateAuthority, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetRemoteCertificates); err != nil {
		return nil, err
	} else {
		return a.authServer.GetTrustedCertificates(ctype)
	}
}
func (a *AuthWithRoles) DeleteRemoteCertificate(ctype string, domainName, id string) error {
	if err := a.permChecker.HasPermission(a.role, ActionDeleteRemoteCertificate); err != nil {
		return err
	} else {
		return a.authServer.DeleteRemoteCertificate(ctype, domainName, id)
	}
}
func (a *AuthWithRoles) GenerateToken(domainName, role string, ttl time.Duration) (string, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGenerateToken); err != nil {
		return "", err
	} else {
		return a.authServer.GenerateToken(domainName, role, ttl)
	}
}
func (a *AuthWithRoles) RegisterUsingToken(token, domainName, role string) (keys PackedKeys, e error) {
	if err := a.permChecker.HasPermission(a.role, ActionRegisterUsingToken); err != nil {
		return PackedKeys{}, err
	} else {
		return a.authServer.RegisterUsingToken(token, domainName, role)
	}
}
func (a *AuthWithRoles) RegisterNewAuthServer(domainName, token string,
	publicSealKey encryptor.Key) (masterKey encryptor.Key, e error) {

	if err := a.permChecker.HasPermission(a.role, ActionRegisterNewAuthServer); err != nil {
		return encryptor.Key{}, err
	} else {
		return a.authServer.RegisterNewAuthServer(domainName, token, publicSealKey)
	}
}
func (a *AuthWithRoles) Log(id lunk.EventID, e lunk.Event) {
	if err := a.permChecker.HasPermission(a.role, ActionLog); err != nil {
		return
	} else {
		a.elog.Log(id, e)
	}
}
func (a *AuthWithRoles) LogEntry(en lunk.Entry) error {
	if err := a.permChecker.HasPermission(a.role, ActionLogEntry); err != nil {
		return err
	} else {
		return a.elog.LogEntry(en)
	}
}
func (a *AuthWithRoles) GetEvents(filter events.Filter) ([]lunk.Entry, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetEvents); err != nil {
		return nil, err
	} else {
		return a.elog.GetEvents(filter)
	}
}
func (a *AuthWithRoles) GetChunkWriter(id string) (recorder.ChunkWriteCloser, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetChunkWriter); err != nil {
		return nil, err
	} else {
		return a.recorder.GetChunkWriter(id)
	}
}
func (a *AuthWithRoles) GetChunkReader(id string) (recorder.ChunkReadCloser, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetChunkReader); err != nil {
		return nil, err
	} else {
		return a.recorder.GetChunkReader(id)
	}
}
func (a *AuthWithRoles) UpsertServer(s services.Server, ttl time.Duration) error {
	if err := a.permChecker.HasPermission(a.role, ActionUpsertServer); err != nil {
		return err
	} else {
		return a.authServer.UpsertServer(s, ttl)
	}
}
func (a *AuthWithRoles) GetServers() ([]services.Server, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetServers); err != nil {
		return nil, err
	} else {
		return a.authServer.GetServers()
	}
}
func (a *AuthWithRoles) UpsertWebTun(wt services.WebTun, ttl time.Duration) error {
	if err := a.permChecker.HasPermission(a.role, ActionUpsertWebTun); err != nil {
		return err
	} else {
		return a.authServer.UpsertWebTun(wt, ttl)
	}
}
func (a *AuthWithRoles) GetWebTuns() ([]services.WebTun, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetWebTuns); err != nil {
		return nil, err
	} else {
		return a.authServer.GetWebTuns()
	}
}
func (a *AuthWithRoles) GetWebTun(prefix string) (*services.WebTun, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetWebTun); err != nil {
		return nil, err
	} else {
		return a.authServer.GetWebTun(prefix)
	}
}
func (a *AuthWithRoles) DeleteWebTun(prefix string) error {
	if err := a.permChecker.HasPermission(a.role, ActionDeleteWebTun); err != nil {
		return err
	} else {
		return a.authServer.DeleteWebTun(prefix)
	}
}
func (a *AuthWithRoles) UpsertPassword(user string, password []byte) (hotpURL string, hotpQR []byte, err error) {
	if err := a.permChecker.HasPermission(a.role, ActionUpsertPassword); err != nil {
		return "", nil, err
	} else {
		return a.authServer.UpsertPassword(user, password)
	}
}
func (a *AuthWithRoles) CheckPassword(user string, password []byte, hotpToken string) error {
	if err := a.permChecker.HasPermission(a.role, ActionCheckPassword); err != nil {
		return err
	} else {
		return a.authServer.CheckPassword(user, password, hotpToken)
	}
}
func (a *AuthWithRoles) SignIn(user string, password []byte) (*Session, error) {
	if err := a.permChecker.HasPermission(a.role, ActionSignIn); err != nil {
		return nil, err
	} else {
		return a.authServer.SignIn(user, password)
	}
}
func (a *AuthWithRoles) GetWebSession(user string, sid websession.SecureID) (*Session, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetWebSession); err != nil {
		return nil, err
	} else {
		return a.authServer.GetWebSession(user, sid)
	}
}
func (a *AuthWithRoles) GetWebSessionsKeys(user string) ([]services.AuthorizedKey, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetWebSessionsKeys); err != nil {
		return nil, err
	} else {
		return a.authServer.GetWebSessionsKeys(user)
	}
}
func (a *AuthWithRoles) DeleteWebSession(user string, sid websession.SecureID) error {
	if err := a.permChecker.HasPermission(a.role, ActionDeleteWebSession); err != nil {
		return err
	} else {
		return a.authServer.DeleteWebSession(user, sid)
	}
}
func (a *AuthWithRoles) GetUsers() ([]services.User, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetUsers); err != nil {
		return nil, err
	} else {
		return a.authServer.GetUsers()
	}
}
func (a *AuthWithRoles) DeleteUser(user string) error {
	if err := a.permChecker.HasPermission(a.role, ActionDeleteUser); err != nil {
		return err
	} else {
		return a.authServer.DeleteUser(user)
	}
}

func (a *AuthWithRoles) GetHostCertificateAuthority() (*services.CertificateAuthority, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetHostCertificateAuthority); err != nil {
		return nil, err
	} else {
		return a.authServer.GetHostCertificateAuthority()
	}
}
func (a *AuthWithRoles) GetUserCertificateAuthority() (*services.CertificateAuthority, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetUserCertificateAuthority); err != nil {
		return nil, err
	} else {
		return a.authServer.GetUserCertificateAuthority()
	}
}
func (a *AuthWithRoles) GenerateKeyPair(pass string) ([]byte, []byte, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGenerateKeyPair); err != nil {
		return nil, nil, err
	} else {
		return a.authServer.GenerateKeyPair(pass)
	}
}
func (a *AuthWithRoles) GenerateHostCert(
	key []byte, id, hostname, role string,
	ttl time.Duration) ([]byte, error) {

	if err := a.permChecker.HasPermission(a.role, ActionGenerateHostCert); err != nil {
		return nil, err
	} else {
		return a.authServer.GenerateHostCert(key, id, hostname, role, ttl)
	}
}
func (a *AuthWithRoles) GenerateUserCert(key []byte, id, user string, ttl time.Duration) ([]byte, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGenerateUserCert); err != nil {
		return nil, err
	} else {
		return a.authServer.GenerateUserCert(key, id, user, ttl)
	}
}
func (a *AuthWithRoles) ResetHostCertificateAuthority(pass string) error {
	if err := a.permChecker.HasPermission(a.role, ActionResetHostCertificateAuthority); err != nil {
		return err
	} else {
		return a.authServer.ResetHostCertificateAuthority(pass)
	}
}
func (a *AuthWithRoles) ResetUserCertificateAuthority(pass string) error {
	if err := a.permChecker.HasPermission(a.role, ActionResetUserCertificateAuthority); err != nil {
		return err
	} else {
		return a.authServer.ResetUserCertificateAuthority(pass)
	}
}
func (a *AuthWithRoles) GetSealKeys() ([]encryptor.Key, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetSealKeys); err != nil {
		return nil, err
	} else {
		return a.authServer.GetSealKeys()
	}
}

func (a *AuthWithRoles) GenerateSealKey(keyName string) (encryptor.Key, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGenerateSealKey); err != nil {
		return encryptor.Key{}, err
	} else {
		return a.authServer.GenerateSealKey(keyName)
	}
}

func (a *AuthWithRoles) DeleteSealKey(keyID string) error {
	if err := a.permChecker.HasPermission(a.role, ActionDeleteSealKey); err != nil {
		return err
	} else {
		return a.authServer.DeleteSealKey(keyID)
	}
}

func (a *AuthWithRoles) AddSealKey(key encryptor.Key) error {
	if err := a.permChecker.HasPermission(a.role, ActionAddSealKey); err != nil {
		return err
	} else {
		return a.authServer.AddSealKey(key)
	}
}

func (a *AuthWithRoles) GetSealKey(keyID string) (encryptor.Key, error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetSealKey); err != nil {
		return encryptor.Key{}, err
	} else {
		return a.authServer.GetSealKey(keyID)
	}
}

func (a *AuthWithRoles) CreateSignupToken(user string, mappings []string) (token string, e error) {
	if err := a.permChecker.HasPermission(a.role, ActionCreateSignupToken); err != nil {
		return "", err
	} else {
		return a.authServer.CreateSignupToken(user, mappings)
	}
}

func (a *AuthWithRoles) GetSignupTokenData(token string) (user string,
	QRImg []byte, hotpFirstValues []string, e error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetSignupTokenData); err != nil {
		return "", nil, nil, err
	} else {
		return a.authServer.GetSignupTokenData(token)
	}
}

func (a *AuthWithRoles) CreateUserWithToken(token, password, hotpToken string) error {
	if err := a.permChecker.HasPermission(a.role, ActionCreateUserWithToken); err != nil {
		return err
	} else {
		return a.authServer.CreateUserWithToken(token, password, hotpToken)
	}
}

func (a *AuthWithRoles) GetCertificateID(certType string, key ssh.PublicKey) (ID string, found bool, e error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetCertificateID); err != nil {
		return "", false, err
	} else {
		return a.authServer.GetCertificateID(certType, key)
	}
}

func (a *AuthWithRoles) UpsertUserMapping(certificateID, teleportUser, osUser string, ttl time.Duration) error {
	if err := a.permChecker.HasPermission(a.role, ActionUpsertUserMapping); err != nil {
		return err
	} else {
		return a.authServer.UpsertUserMapping(certificateID, teleportUser, osUser, ttl)
	}
}

func (a *AuthWithRoles) DeleteUserMapping(certificateID, teleportUser, osUser string) error {
	if err := a.permChecker.HasPermission(a.role, ActionDeleteUserMapping); err != nil {
		return err
	} else {
		return a.authServer.DeleteUserMapping(certificateID, teleportUser, osUser)
	}
}

func (a *AuthWithRoles) UserMappingExists(certificateID, teleportUser, osUser string) (bool, error) {
	if err := a.permChecker.HasPermission(a.role, ActionUserMappingExists); err != nil {
		return false, err
	} else {
		return a.authServer.UserMappingExists(certificateID, teleportUser, osUser)
	}
}

func (a *AuthWithRoles) GetAllUserMappings() (hashes []string, e error) {
	if err := a.permChecker.HasPermission(a.role, ActionGetAllUserMappings); err != nil {
		return nil, err
	} else {
		return a.authServer.GetAllUserMappings()
	}
}