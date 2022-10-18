package gsm

import (
	"encoding/json"
	"fmt"

	"github.com/go-msvc/errors"
)

type IsdnAddr struct {
	country *Country
	nai     NAI
	npi     NPI
	addr    string //e.g. "821234567" (exclude 0 or 27)
}

func (ia IsdnAddr) NPI() NPI          { return ia.npi }
func (ia IsdnAddr) NAI() NAI          { return ia.nai }
func (ia IsdnAddr) Addr() string      { return ia.addr }
func (ia IsdnAddr) Country() *Country { return ia.country }

func (ia IsdnAddr) CheckIntlMsisdn() error {
	if ia.nai != NAI_INTERNATIONAL {
		return errors.Errorf("nai=%s (expecting %s)", ia.nai, NAI_INTERNATIONAL)
	}

	if ia.npi != NPI_E164 {
		return errors.Errorf("npi=%s (expecting %s)", ia.npi, NPI_E164)
	}
	if len(ia.addr) < 6 {
		return errors.Errorf("addr=\"%s\" (expect 6+ digits)", ia.addr)
	}
	return nil
}

func (ia IsdnAddr) String() string {
	//use first CC...
	if ia.country == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s%s", ia.country.ccList[0], ia.addr) //international format, e.g. "27821234567"
}

func (ia IsdnAddr) StringNat() string {
	return fmt.Sprintf("0%s", ia.addr) //international format, e.g. "27821234567"
}

func (m *IsdnAddr) Scan(v interface{}) error {
	//must be international format because scan does not take country indication
	//and the code must be able to run in the cloud across multiple countries

	//get string value (must exclude quotes)
	s, ok := v.(string)
	if !ok {
		b, ok := v.([]byte)
		if !ok {
			return errors.Errorf("cannot scan IsdnAddr from %T", v)
		}
		s = string(b)
	}

	//todo: add support for detailed "INT,E164,27821234567"

	//match the longest possible country code from the start
	s, m.country = countryFromCCPrefix(s)
	if m.country == nil {
		return errors.Errorf("not starting with country code \"%s\"", s)
	}

	//rest of address are digits
	m.addr = s
	m.npi = NPI_E164
	m.nai = NAI_INTERNATIONAL
	return nil
}

type IsdnAddrJSONStruct struct {
	NAI  string
	NPI  string
	Addr string
}

func (m IsdnAddr) MarshalJSON() ([]byte, error) {
	j := IsdnAddrJSONStruct{
		NAI:  m.nai.String(),
		NPI:  m.npi.String(),
		Addr: m.addr,
	}
	return json.Marshal(j)
}

func (m *IsdnAddr) UnmarshalJSON(jsonIsdnAddr []byte) error {
	var j IsdnAddrJSONStruct
	if err := json.Unmarshal(jsonIsdnAddr, &j); err != nil {
		return errors.Wrapf(err, "failed to decode JSON into IsdnAddr")
	}
	return nil
}
