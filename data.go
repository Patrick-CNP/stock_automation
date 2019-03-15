package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
  Generate a key from https://www.alphavantage.co/
*/
var key = "<KEY_NEEDED>"

/*
Time Series Data Struc
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
	VolumeWeightedAveragePrice Struc
 */
type VolumeWeightedAveragePrice struct {

	MetaData struct {
		Symbol string `json:"1. Symbol"`
		Indicator string `json:"2: Indicator"`
		LastRefreshed string `json:"3: Last Refreshed"`
		Interval string `json:"4: Interval"`
		TimeZone string `json:"5: Time Zone"`
	} `json:"Meta Data"`

	Data map[string]struct {
		Value float32 `json:"VWAP,string"`
	} `json:"Technical Analysis: VWAP"`
}

type ESVolumeWeightedAveragePrice struct {
	Date string
	VWAP float32
}

/*
 Stock struc to be sent to ES
 */
type Stock struct {
	/*
	From Time Series
	 */
	Date string
	Open float64
	Close float64
	High float64
	Low float64
	Volume float64

	/*
	From Simple Moving Av
	 */
	 SimpleMovingAv float64

	 /*
	From ExponentialMoving Av
	 */
	ExponentialMovingAv float64
}

/*
	http client for api calls
 */
var myClient = &http.Client{Timeout: 10 * time.Second}

/*
 	Get the Time Series Data
 */
func getTimeSeries(symbol string) (*TimeSeries, error) {

	target := new(TimeSeries)

	url := "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY&symbol=" +
		symbol +
		"&outputsize=full&apikey=" +
		key

	r, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

/*
 	Get the Simple Moving Average Data
 */
func getSimpleMovingAv(symbol string, window int) (*SimpleMovingAv, error) {

	target := new(SimpleMovingAv)

	url := "https://www.alphavantage.co/query?function=SMA&symbol=" +
		symbol +
		"&interval=daily&time_period=" +
		strconv.Itoa(window) +
		"&series_type=close&outputsize=full&apikey=" +
		key

	r, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

/*
 	Get the Exponential Moving Average Data
 */
func getExponentialMovingAv(symbol string, window int) (*ExponentialMovingAv, error) {

	target := new(ExponentialMovingAv)

	url := "https://www.alphavantage.co/query?function=SMA&symbol=" +
		symbol +
		"&interval=daily&time_period=" +
		strconv.Itoa(window) +
		"&series_type=close&outputsize=full&apikey=" +
		key

	r, err := myClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

/*
	Process Time Series, Simple Moving Average & Exponential Moving Average
	to return a stock object
 */
 /*
func getStock(ts *TimeSeries, sma *SimpleMovingAv, ema *SimpleMovingAv) (*[]Stock, error) {

		Loop through Time Series, match the dates and
		pull the matching values from SMA and EMA


	 stock_data := new([]Stock)
	 var e error

	 for ts_date, ts_data := range ts.Data {

		s := new(Stock)

		s.Date = ts_date

		s.Volume,e = strconv.ParseFloat(ts_data.Volume, 64)
		 handle(e)

		s.Open,e = strconv.ParseFloat(ts_data.Open, 64)
		 handle(e)

		s.Close,e = strconv.ParseFloat(ts_data.Close, 64)
		 handle(e)

		s.High,e = strconv.ParseFloat(ts_data.High, 64)
		 handle(e)

		s.Low,e = strconv.ParseFloat(ts_data.Low, 64)
		 handle(e)

		s.SimpleMovingAv,e = strconv.ParseFloat(sma.Data[ts_date].Value, 64)
		 handle(e)

		s.ExponentialMovingAv,e = strconv.ParseFloat(ema.Data[ts_date].Value, 64)
		handle(e)

		*stock_data = append(*stock_data, *s)
	 }

	 return stock_data, nil
}
*/

func EsPut(index string, b *[]byte, sem *chan int ) error {
	uri := "http://localhost:9200" + index

	client := &http.Client{}
	client.Timeout = time.Second * 15

	body := bytes.NewBuffer(*b)
	req, err := http.NewRequest(http.MethodPut, uri, body)
	if err != nil {
		return err
	}

	fmt.Println(body)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do() failed with '%s'\n", err)
		return err
	}

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ioutil.ReadAll() failed with '%s'\n", err)
		return err
	}

	fmt.Printf("Response status code: %d, text:\n%s\n", resp.StatusCode, string(d))

	// release chanel
	<-*sem
	return nil
}

/*
Convert the input data sets to structs with format for ES
 */
type Converter interface {
	Convert(*interface{}) error
}

func (ts *TimeSeries) Convert(series *[]ESTimeSeries) error{

	var s ESTimeSeries

	for date,values := range ts.Data {
		s = ESTimeSeries{
			Date:date,
			Open:values.Open,
			Close:values.Close,
			Volume:values.Volume,
			High:values.High,
			Low:values.Low }
		*series = append(*series, s)
	}
	return nil
}


func (sma *SimpleMovingAv) Convert( esSMA *[]ESSimpleMovingAv ) error{
	var es ESSimpleMovingAv

	for date,values := range sma.Data{
		es = ESSimpleMovingAv{
			Date:date,
			SMA:values.Value }
		*esSMA = append(*esSMA, es)
		}

	return nil
}

func (ema *ExponentialMovingAv) Convert( esEMA *[]ESExponentialMovingAv ) error{
	var es ESExponentialMovingAv

	for date,values := range ema.Data{
		es = ESExponentialMovingAv{
			Date:date,
			EMA:values.Value }
		*esEMA = append(*esEMA, es)
	}

	return nil
}


/*
Send the different data types to ES
 */

type Sender interface {
	Send(index_name string, concurrency int) error
}

func (ts *TimeSeries) Send( index_name string, concurrency int ) error {

	/*
	Used for output of each Marshaled json
 	*/
	var b []byte
	var err error

	/*
	Convert to ES compatible Time Series
	 */
	 var esTimeSeries []ESTimeSeries
	 ts.Convert(&esTimeSeries)

	/*
	Channel to manage concurrent processes
	 */
	c := make(chan int, concurrency)

	indexCount := 0
	for _,d := range esTimeSeries {

		fmt.Println(d)

		b,err = json.Marshal(d)
		fmt.Println(string(b))
		if err != nil {
		   return err
		}

		c <- 0
		i := index_name + strconv.Itoa(indexCount)
		go EsPut(i, &b, &c)
		indexCount++
	}

	return nil
}

func (sma *SimpleMovingAv) Send( index_name string, concurrency int ) error {

	/*
	Used for output of each Marshaled json
 	*/
	var b []byte
	var err error

	/*
	Convert to ES compatible Time Series
	 */
	var esSMA []ESSimpleMovingAv
	sma.Convert(&esSMA)

	/*
	Channel to manage concurrent processes
	 */
	c := make(chan int, concurrency)

	indexCount := 0
	for _,d := range esSMA {

		fmt.Println(d)
		b,err = json.Marshal(d)

		/*
		Update the value key in the data to separate time windows
 		*/
		var kv string
		poss_kv := strings.Split(index_name, "-")
		if len(poss_kv) > 1{
			kv = poss_kv[len(poss_kv)-1]
		} else {
			kv = poss_kv[0]
		}
		/*
		Convert []byte to a string, replace the 'SMA' key with new value
		then convert back to []byte
		 */
		b := []byte(strings.Replace(string(b), "SMA", kv, 1))

		fmt.Println(string(b))
		if err != nil {
			return err
		}

		c <- 0
		i := index_name + strconv.Itoa(indexCount)
		go EsPut(i, &b, &c)
		indexCount++
	}

	return nil
}

func (ema *ExponentialMovingAv) Send( index_name string, concurrency int ) error {

	/*
	Used for output of each Marshaled json
 	*/
	var b []byte
	var err error

	/*
	Convert to ES compatible Time Series
	 */
	var esEMA []ESExponentialMovingAv
	ema.Convert(&esEMA)

	/*
	Channel to manage concurrent processes
	 */
	c := make(chan int, concurrency)

	indexCount := 0
	for _,d := range esEMA {

		fmt.Println(d)

		b,err = json.Marshal(d)
		fmt.Println(string(b))
		if err != nil {
			return err
		}

		c <- 0
		i := index_name + strconv.Itoa(indexCount)
		go EsPut(i, &b, &c)
		indexCount++
	}

	return nil
}
