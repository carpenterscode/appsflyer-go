package appsflyer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	appsflyerURLFormat  = "https://api2.appsflyer.com/inappevent/%s"
	appsflyerDateFormat = "2006/01/02"
	appsflyerTimeFormat = "2006-01-02 15:04:05.000"
)

type Tracker struct {
	client    http.Client
	platforms map[deviceOS]app
}

func NewTracker() (*Tracker, error) {

	// An instance of http.DefaultTransport with 10 max idle conn
	client := http.Client{Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}}

	return &Tracker{client: client}, nil
}

func (t *Tracker) SetConfig(configPath string) error {
	configFile, configErr := os.Open(configPath)
	defer configFile.Close()
	if configErr != nil {
		return configErr
	}
	var config []app
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return err
	}

	t.platforms = make(map[deviceOS]app, len(config))
	for _, app := range config {
		t.platforms[app.Platform] = app
	}
	return nil
}

func (t Tracker) Send(evt *Event) error {

	if evt == nil {
		return fmt.Errorf("AppsFlyer tracker should not send a nil event")
	} else if string(evt.body.EventName) == "" {
		return fmt.Errorf("AppsFlyer event should have an event name")
	}

	appInfo := t.platforms[evt.platform]

	evt.body.BundleID = appInfo.BundleID

	appsFlyerURL := fmt.Sprintf(appsflyerURLFormat, appInfo.AppID)

	var body io.Reader
	if data, jsonErr := json.Marshal(evt); jsonErr != nil {
		return jsonErr
	} else {
		body = bytes.NewReader(data)
	}

	req, reqErr := http.NewRequest(http.MethodPost, appsFlyerURL, body)
	if reqErr != nil {
		return reqErr
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authentication", appInfo.DevKey)

	if resp, err := t.client.Do(req); err != nil {
		return err
	} else if resp.StatusCode != 200 {
		resp.Body.Close()
		return fmt.Errorf("AppsFlyer server response %s", resp.Status)
	} else {
		resp.Body.Close()
	}

	return nil
}

type app struct {
	AppID    string   `json:"appId"`
	BundleID string   `json:"bundleId"`
	DevKey   string   `json:"devKey"`
	Platform deviceOS `json:"platform"`
}

type Event struct {
	body
	platform deviceOS
	values   map[EventParam]string
}

type body struct {
	AdvertisingID string `json:"advertising_id,omitempty"`
	AppsFlyerID   string `json:"appsflyer_id"`
	BundleID      string `json:"bundle_id,omitempty"`
	DeviceIP      string `json:"ip,omitempty"`
	EventCurrency string `json:"eventCurrency,omitempty"`
	EventName     string `json:"eventName"`
	EventTime     string `json:"eventTime,omitempty"`
	EventValue    string `json:"eventValue"`
	IDFA          string `json:"idfa,omitempty"`
	UseEventsAPI  bool   `json:"af_events_api,string"`
}

func NewEvent(appsFlyerID string, platform deviceOS) *Event {
	return &Event{
		body: body{
			AppsFlyerID:  appsFlyerID,
			UseEventsAPI: true,
		},
		platform: platform,
		values:   make(map[EventParam]string),
	}
}

func (evt *Event) SetName(name EventName) *Event {
	evt.body.EventName = string(name)
	return evt
}

func (evt *Event) SetEventTime(eventTime time.Time) *Event {
	evt.body.EventTime = eventTime.Format(appsflyerTimeFormat)
	return evt
}

func (evt *Event) SetValue(param EventParam, val string) *Event {
	evt.values[param] = val
	return evt
}

func (evt *Event) SetDateValue(param EventParam, val time.Time) *Event {
	evt.values[param] = val.Format(appsflyerDateFormat)
	return evt
}

func (evt *Event) SetRevenue(revenue float64, currency string) *Event {
	evt.values[ParamCurrency] = currency
	evt.values[ParamRevenue] = fmt.Sprintf("%.2f", revenue)
	return evt
}

func (evt *Event) SetPrice(price float64, currency string) *Event {
	evt.values[ParamCurrency] = currency
	evt.values[ParamPrice] = fmt.Sprintf("%.2f", price)
	return evt
}

func (evt *Event) SetAdvertisingID(advertisingID string) *Event {
	if evt.platform == Android {
		evt.AdvertisingID = advertisingID
	} else {
		evt.IDFA = advertisingID
	}
	return evt
}

func (evt *Event) SetDeviceIP(deviceIP string) *Event {
	evt.DeviceIP = deviceIP
	return evt
}

func (evt Event) MarshalJSON() ([]byte, error) {
	if len(evt.values) > 0 {
		if data, err := json.Marshal(evt.values); err != nil {
			return nil, err
		} else {
			evt.EventValue = string(data)
		}
	}

	return json.Marshal(evt.body)
}
