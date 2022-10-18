package gsm_test

import (
	"testing"

	"github.com/go-msvc/utils/gsm"
)

func Test1(t *testing.T) {
	if c, ok := gsm.CountryByMCC("655"); !ok {
		t.Fatalf("MCC 655 not found")
	} else {
		if !c.HasName("South Africa") {
			t.Fatalf("Not SA: %+v", c)
		}
	}

	//look for MCC and MNC as in first 5 digits of IMSI
	if n, ok := gsm.NetworkByMCCMNC("65501"); !ok {
		t.Fatalf("65501 not found")
	} else {
		if !n.HasName("Vodacom") {
			t.Fatalf("Not Vodacom: %+v", n)
		}
	}
	if c, s, ok := gsm.CountryByCC("27824526299"); !ok {
		t.Fatalf("27 not found")
	} else {
		if !c.HasName("South Africa") {
			t.Fatalf("CC not SA: %+v", c)
		}
		if s != "824526299" {
			t.Fatalf("Not remain(%s)", s)
		}
	}

	//load will not append but add to existing entries
	//if want to load fresh, need to add a reset func to discard init load data
	//but for now keep the init load because this info changes seldom
	if err := gsm.LoadCSVFile("./mcc-mnc-table.csv"); err != nil {
		t.Fatalf("failed to load: %+v", err)
	}
	t.Logf("loaded")
}
