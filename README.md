# appsflyer-go

A Go client library for the [AppsFlyer Server-to-Server Events API](https://support.appsflyer.com/hc/en-us/articles/207034486-Server-to-Server-Events-API)

## Installation

```bash
go get github.com/carpenterscode/appsflyer-go
```

## Configuration

Set up the client with a JSON file of this format

```json
[
    {
        "appId": "com.company.android",
        "devKey": "aaaaaaaaaaaaaaaaaaaaaa",
        "platform": "android"
    },
    {
        "appId": "id111111111",
        "bundleId": "com.company.ios",
        "devKey": "aaaaaaaaaaaaaaaaaaaaaa",
        "platform": "ios"
    }
]
```

## Usage

```go
import (
        "time"

        af "github.com/carpenterscode/appsflyer-go"
)

func main() {
        tracker, trackErr := af.NewTracker()
        if trackErr != nil {
                panic(trackErr)
        }
        tracker.SetConfig("appsflyer.json")

        startTrial(tracker)

        subscribe(tracker)

        cancelSubscription(tracker)
}

func startTrial(tracker af.Tracker) {

        // User starts a trial
        evt := af.NewEvent("1111111111111-1111111", af.IOS).
                SetName(af.StartTrial).
                SetAdvertisingID("AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA").
                SetDeviceIP("1.2.3.4").
                SetPrice(59.99, "USD").
                SetDateValue("expiry", validDate).
                SetEventTime(time.Now())

        if err := tracker.Send(evt); err != nil {
                panic(err)
        }
}

func subscribe(tracker af.Tracker) {

        // User ends trial and pays for first subscription period
        evt := af.NewEvent("1111111111111-1111111", af.IOS).
                SetName(af.Subscribe).
                SetAdvertisingID("AAAAAAAA-AAAA-AAAA-AAAA-AAAAAAAAAAAA").
                SetDeviceIP("1.2.3.4").
                SetRevenue(59.99, "USD").
                SetDateValue("expiry", validDate).

        if err := tracker.Send(evt); err != nil {
                panic(err)
        }
}

func cancelSubscription(tracker af.Tracker) {

        // User cancels a subscription
        evt := af.NewEvent("1111111111111-1111111", af.Android)
        evt.SetName(af.EventName("cancel_subscription"))
        evt.SetRevenue(-59.99, "USD")

        if err := tracker.Send(evt); err != nil {
                panic(err)
        }
}
```

## Documentation

Official AppsFlyer docs

-   [Server-to-Server guide](https://support.appsflyer.com/hc/en-us/articles/207034486-Server-to-Server-Events-API)
-   [Subscription Tracking guide](https://support.appsflyer.com/hc/en-us/articles/360001279189-Subscription-Tracking-Guide#ServertoServer)

## License

[MIT](LICENSE)
