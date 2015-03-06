package main

import (
	"github.com/keybase/client/go/engine"
	"github.com/keybase/client/go/libkb"
	keybase_1 "github.com/keybase/client/protocol/go"
	"github.com/maxtaco/go-framed-msgpack-rpc/rpc2"
)

type MykeyHandler struct {
	BaseHandler
}

func NewMykeyHandler(xp *rpc2.Transport) *MykeyHandler {
	return &MykeyHandler{BaseHandler{xp: xp}}
}

func (h *MykeyHandler) PgpKeyGen(arg keybase_1.PgpKeyGenArg) (err error) {
	earg := engine.ImportPGPEngineArg(arg)
	return h.keygen(earg, true)
}

func (h *MykeyHandler) keygen(earg engine.PGPEngineArg, doInteractive bool) (err error) {
	sessionId := nextSessionId()
	ctx := &engine.Context{LogUI: h.getLogUI(sessionId), SecretUI: h.getSecretUI(sessionId)}
	earg.Gen.AddDefaultUid()
	eng := engine.NewPGPEngine(earg)
	err = engine.RunEngine(eng, ctx, nil, nil)
	return err
}

func (h *MykeyHandler) PgpKeyGenDefault(arg keybase_1.PgpCreateUids) (err error) {
	earg := engine.PGPEngineArg{
		Gen: &libkb.PGPGenArg{
			Ids:         libkb.ImportPgpIdentities(arg.Ids),
			NoDefPGPUid: !arg.UseDefault,
		},
	}
	return h.keygen(earg, false)
}

func (h *MykeyHandler) DeletePrimary() (err error) {
	return libkb.DeletePrimary()
}

func (h *MykeyHandler) Show() (err error) {
	sessionId := nextSessionId()
	return libkb.ShowKeys(h.getLogUI(sessionId))
}

func (h *MykeyHandler) Select(sarg keybase_1.SelectArg) error {
	sessionID := nextSessionId()
	gpgui := NewRemoteGPGUI(sessionID, h.getRpcClient())
	secretui := h.getSecretUI(sessionID)
	arg := engine.GPGArg{Query: sarg.Query, AllowMulti: sarg.AllowMulti, SkipImport: sarg.SkipImport}
	gpg := engine.NewGPG(&arg)
	ctx := &engine.Context{
		GPGUI:    gpgui,
		SecretUI: secretui,
		LogUI:    h.getLogUI(sessionID),
		LoginUI:  h.getLoginUI(sessionID),
	}
	return engine.RunEngine(gpg, ctx, nil, nil)
}

func (h *MykeyHandler) SaveArmoredPGPKey(arg keybase_1.SaveArmoredPGPKeyArg) error {
	sessionID := nextSessionId()
	ctx := &engine.Context{
		SecretUI: h.getSecretUI(sessionID),
		LogUI:    h.getLogUI(sessionID),
	}
	eng := engine.NewPGPSaveArmored(arg.Key, arg.PushPublic, arg.PushPrivate)
	return engine.RunEngine(eng, ctx, nil, nil)
}

func (h *MykeyHandler) SavePGPKey(arg keybase_1.SavePGPKeyArg) error {
	sessionID := nextSessionId()
	ctx := &engine.Context{
		SecretUI: h.getSecretUI(sessionID),
		LogUI:    h.getLogUI(sessionID),
	}
	eng := engine.NewPGPSaveRaw(arg.Key, arg.PushPublic, arg.PushPrivate)
	return engine.RunEngine(eng, ctx, nil, nil)
}
