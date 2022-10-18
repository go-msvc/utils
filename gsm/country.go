package gsm

import "regexp"

//CountryByMCC() looks for a country description using the 3-digit MCC
//that any IMSI starts with, e.g. "655011234567890" starts with 655
//which is MCC for South Africa
func CountryByMCC(digits string) (*Country, bool) {
	if len(digits) < 3 {
		return nil, false
	}
	if c, ok := countryByMCC[digits[0:3]]; ok {
		return c, true
	}
	return nil, false
}

//CountryByCC() looks for the longest country code match from the start
//of an international phone number, e.g. 27 for South Africa
//if not found: returns nil, all digits, false
//if found:     returns country, remaining digits, true
//input digits must only be digits, no leading '+' or '00'
func CountryByCC(digits string) (*Country, string, bool) {
	cc := digits
	for len(cc) > 0 {
		if c, ok := countryByCC[cc]; ok {
			return c, digits[len(cc):], true
		}
		cc = cc[0 : len(cc)-1]
	}
	return nil, digits, false
}

//Mobile Country Code is 3 digits e.g. 655 = South Africa
const mccPattern = `[0-9][0-9][0-9]`

var mccRegex = regexp.MustCompile(`^` + mccPattern + `$`)

//Mobile Network Code is 2 digits e.g. 01 for Vodacom in South Africa
const mncPattern = `[0-9][0-9]`

var mncRegex = regexp.MustCompile(`^` + mncPattern + `$`)

var (
	countryByName = map[string]*Country{}
	countryByCC   = map[string]*Country{}
	countryByMCC  = map[string]*Country{}
	countryByISO  = map[string]*Country{}
)

type Country struct {
	mcc           string   //e.g. "655"
	nameList      []string //e.g. "South Africa"
	isoList       []string //e.g. "za"
	ccList        []string //e.g. "27"
	networkByName map[string]*Network
	networkByMNC  map[string]*Network
}

//call e.g. for phone nr
//returns remaining string, because CC has variable length
//if not found, return whole string and Country==nil
func countryFromCCPrefix(digits string) (string, *Country) {
	return "", nil //todo
}

func (c Country) HasName(n string) bool {
	for _, cn := range c.nameList {
		if n == cn {
			return true
		}
	}
	return false
}

func (c Country) NetworkByMNC(mnc string) (*Network, bool) {
	if len(mnc) < 2 {
		return nil, false
	}
	n, ok := c.networkByMNC[mnc[0:2]]
	if ok {
		return n, true
	}
	return nil, false
}

func NetworkByMCCMNC(digits string) (*Network, bool) {
	if len(digits) >= 5 {
		if c, ok := CountryByMCC(digits); ok {
			return c.NetworkByMNC(digits[3:])
		}
	}
	return nil, false
}
