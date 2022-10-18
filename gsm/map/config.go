package gsmmap

import (
	"context"

	"github.com/go-msvc/errors"
	"github.com/go-msvc/utils/gsm"
)

type Config struct {
}

func (c Config) Validate() error { return nil }

func (c Config) Create() (any /*Client*/, error) {
	return Client{}, nil
	//return nil, errors.Errorf("NYI")
}

type Client struct {
}

type SriSmRequest struct {
	Msisdn gsm.IsdnAddr
}

func (req SriSmRequest) Validate() error {
	if err := req.Msisdn.CheckIntlMsisdn(); err != nil {
		return errors.Wrapf(err, "invalid msisdn")
	}
	return nil
}

type SubscriberStatus int

const (
	SubscriberStatusUnknown SubscriberStatus = iota
	SubscriberStatusAbsent
	SubscriberStatusAvailable
)

type SriSmResponse struct {
	Imsi        gsm.Imsi
	Status      SubscriberStatus
	VlrGt       gsm.IsdnAddr
	VlrNetwork  *gsm.Network
	ImsiNetwork *gsm.Network
}

func (cli Client) DoSriSm(ctx context.Context, req SriSmRequest) (SriSmResponse, error) {
	return SriSmResponse{}, errors.Errorf("NYI")
}

type FwdSmRequest struct {
	VlrGt gsm.IsdnAddr `json:"vlr_gt"`
	Imsi  gsm.Imsi     `json:"imsi"`
	Text  string       `json:"text"` //todo: include encoding etc..
}

func (req FwdSmRequest) Validate() error {
	if err := req.VlrGt.CheckIntlMsisdn(); err != nil {
		return errors.Wrapf(err, "invalid vlr_gt")
	}
	//todo more checks
	return nil
}

type FwdSmResponse struct{}

func (cli Client) DoFwdSm(ctx context.Context, req FwdSmRequest) (FwdSmResponse, error) {
	return FwdSmResponse{}, errors.Errorf("NYI")
}
