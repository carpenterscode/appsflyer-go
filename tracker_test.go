package appsflyer

import (
	"encoding/json"
	"testing"
	"time"
)

const (
	advertisingID    = "AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA"
	appsFlyerID      = "1111111111111-1111111"
	iosEventJSON     = `{"appsflyer_id":"1111111111111-1111111","bundle_id":"com.company.app","ip":"1.2.3.4","eventName":"af_start_trial","eventTime":"2019-02-15 06:30:36.869","eventValue":"{\"af_currency\":\"USD\",\"af_revenue\":\"59.99\",\"expiry\":\"2014/11/11\"}","idfa":"AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA","af_events_api":"true"}`
	androidEventJSON = `{"advertising_id":"AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA","appsflyer_id":"1111111111111-1111111","eventName":"af_start_trial","eventValue":"","af_events_api":"true"}`
)

var validDate = time.Date(2014, time.November, 11, 19, 0, 0, 0, time.Local)
var validTime = time.Date(2019, time.February, 15, 6, 30, 36, 869*1000000, time.UTC)

func TestEventWithOptionalParams(t *testing.T) {
	evt := NewEvent(appsFlyerID, IOS)
	evt.SetName(StartTrial)
	evt.SetAdvertisingID(advertisingID)
	evt.SetBundleID("com.company.app")
	evt.SetDeviceIP("1.2.3.4")
	evt.SetRevenue(59.99, "USD")
	evt.SetDateValue("expiry", validDate)
	evt.SetEventTime(validTime)

	if buf, err := json.Marshal(evt); err != nil {
		t.Errorf("Should have read event. %s", err)
	} else if string(buf) != iosEventJSON {
		t.Error(string(buf))
		t.Error("Should become valid JSON")
	}
}

func TestEventWithoutOptionalParams(t *testing.T) {
	evt := NewEvent(appsFlyerID, Android)
	evt.SetName(StartTrial)
	evt.SetAdvertisingID(advertisingID)

	if buf, err := json.Marshal(evt); err != nil {
		t.Errorf("Should have read event. %s", err)
	} else if string(buf) != androidEventJSON {
		t.Error("Should become valid JSON")
	}
}

func TestConfigureTracker(t *testing.T) {
	if _, trackerErr := NewTracker(); trackerErr != nil {
		t.Error(trackerErr)
	}
	// TODO: Read a test JSON file
}
