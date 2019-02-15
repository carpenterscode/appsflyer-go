package appsflyer

type deviceOS string

const (
	Android deviceOS = "android"
	IOS     deviceOS = "ios"
)

type EventName string

const (
	StartTrial EventName = "af_start_trial"
	Subscribe  EventName = "af_subscribe"
)

type EventParam string

const (
	ParamRevenue  EventParam = "af_revenue"
	ParamPrice    EventParam = "af_price"
	ParamCurrency EventParam = "af_currency"
)
