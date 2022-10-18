package gsm

import "fmt"

/*
	-- bits 4321: numbering plan indicator
	--  0000  unknown
	--  0001  ISDN/Telephony Numbering Plan (Rec CCITT E.164)
	--  0010  spare
	--  0011  data numbering plan (CCITT Rec X.121)
	--  0100  telex numbering plan (CCITT Rec F.69)
	--  0101  spare
	--  0110  land mobile numbering plan (CCITT Rec E.212)
	--  0111  spare
	--  1000  national numbering plan
	--  1001  private numbering plan
	--  1111  reserved for extension
	--  all other values are reserved.
*/
type NPI int

func (npi NPI) String() string {
	if s, ok := npiString[npi]; ok {
		return s
	}
	return fmt.Sprintf("0x%02x", int(npi))
}

const (
	NPI_UNKNOWN  NPI = 0x00
	NPI_E164     NPI = 0x01
	NPI_SPARE_2  NPI = 0x02
	NPI_X121     NPI = 0x03
	NPI_F69      NPI = 0x04
	NPI_SPARE_5  NPI = 0x05
	NPI_E212     NPI = 0x06
	NPI_SPARE_7  NPI = 0x07
	NPI_NATIONAL NPI = 0x08
	NPI_PRIVATE  NPI = 0x09
	NPI_RESERVED NPI = 0x0f
	_NPI_NR_OF
)

var (
	npiString = map[NPI]string{
		NPI_UNKNOWN:  "UNKNOWN",
		NPI_E164:     "E164",
		NPI_SPARE_2:  "SPARE_2",
		NPI_X121:     "X121",
		NPI_F69:      "F69",
		NPI_SPARE_5:  "SPARE_5",
		NPI_E212:     "E212",
		NPI_SPARE_7:  "SPARE_7",
		NPI_NATIONAL: "NAT",
		NPI_PRIVATE:  "PRIV",
		NPI_RESERVED: "RESERVED",
	}
)
