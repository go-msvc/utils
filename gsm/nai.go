package gsm

import "fmt"

/*
  -- bits 765: nature of address indicator
  --  000  unknown
  --  001  international number
  --  010  national significant number
  --  011  network specific number
  --  100  subscriber number
  --  101  alphanumeric
  --  110  abbreviated number
  --  111  reserved for extension
*/

type NAI int

func (nai NAI) String() string {
	if s, ok := naiString[nai]; ok {
		return s
	}
	return fmt.Sprintf("0x%02x", int(nai))
}

const (
	NAI_UNKNOWN              NAI = 0x00
	NAI_INTERNATIONAL        NAI = 0x01
	NAI_NATIONAL             NAI = 0x02
	NAI_NETWORK              NAI = 0x03
	NAI_SUBSCRIBER           NAI = 0x04
	NAI_ALPHANUMERIC         NAI = 0x05
	NAI_ABBREVIATED          NAI = 0x06
	NAI_ALPHANUMERIC_FOR_EXT NAI = 0x07
	_NAI_NR_OF               NAI = 0x08
)

var (
	naiString = map[NAI]string{
		NAI_UNKNOWN:              "UNK",
		NAI_INTERNATIONAL:        "INT",
		NAI_NATIONAL:             "NAT",
		NAI_NETWORK:              "NET",
		NAI_SUBSCRIBER:           "SUB",
		NAI_ALPHANUMERIC:         "ALP",
		NAI_ABBREVIATED:          "ABR",
		NAI_ALPHANUMERIC_FOR_EXT: "AEX",
	}
)
