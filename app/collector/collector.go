package collector

/*
 Stock struc to store data for SMA. Data is moved here before being send to ES
to ease the load on ES. Un required data points are stripped out before uploading to ES
 */
type SMAData struct {
	Data []SMA
}

type SMA struct {

	/*
	From Time Series
	*/
	Date   string
	Open   float32
	Close  float32
	High   float32
	Low    float32
	Volume float32

	/*
	From Simple Moving Av - 50 Day
	*/
	SMA50Day float32

	/*
	From Simple Moving Av - 15 Day
	*/
	SMA15Day float32
}

/*
	Common MetaData
 */
type MetaData struct {
	Symbol string `json:"1. Symbol"`
	Indicator string `json:"2: Indicator"`
	LastRefreshed string `json:"3: Last Refreshed"`
	Interval string `json:"4: Interval"`
	TimePeriod int `json:"5: Time Period"`
	SeriesType string `json:"6: Series Type"`
	TimeZone string `json:"7: Time Zone"`
} //`json:"Meta Data"`

/*
Time Series Data Struct
 */
type TimeSeries struct {
	MetaData struct {
		Information string `json:"1. Information"`
		Symbol string `json:"2. Symbol"`
		LastRefreshed string `json:"3. Last Refreshed"`
		Interval string `json:"4. Interval String"`
		OutputSize string `json:"5. Output Size"`
		Timezone string `json:"6. Time Zone"`
	} `json:"Meta Data"`

	Data map[string]struct {
		Open float32 `json:"1. open,string"`
		High float32 `json:"2. high,string"`
		Low float32 `json:"3. low,string"`
		Close float32 `json:"4. close,string"`
		Volume float32 `json:"5. volume,string"`
	} `json:"Time Series (Daily)"`
}

type ESTimeSeries struct {
	Date string
	Open float32
	High float32
	Low float32
	Close float32
	Volume float32
}

/*
	SimpleMovingAv Struc
 */
type SimpleMovingAv struct {

	Meta MetaData `json:"Meta Data"`

	Data map[string]struct {
		Value float32 `json:"SMA,string"`
	} `json:"Technical Analysis: SMA"`
}

type ESSimpleMovingAv struct {
	Date string
	SMA float32
}

/*
	ExponentialMovingAv Struc
 */
type ExponentialMovingAv struct {

	Meta MetaData `json:"Meta Data"`

	Data map[string]struct {
		Value float32 `json:"EMA,string"`
	} `json:"Technical Analysis: EMA"`
}

type ESExponentialMovingAv struct {
	Date string
	EMA float32
}

/*
	Any object implementing the collector inetrface must
	contain the following methods.
 */
type collector interface {
	getTimeSeries(symbol string) (*TimeSeries, error)
	getSimpleMovingAv(symbol string, window int) (*SimpleMovingAv, error)
	getExponentialMovingAv(symbol string, window int) (*ExponentialMovingAv, error)
}

/*
	Select the collector
 */
func newCollector(collectorType string) (*collector, error){
	/*
		Return a stock collector object based on the
		collectorType string
	 */
	 switch collectorType {
	 	case "alphavantage":

	 	}
}
