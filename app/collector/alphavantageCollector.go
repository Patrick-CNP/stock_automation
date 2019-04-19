package collector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"stock_automation/app/utils"
	"strconv"
	"time"
)

/*
API Struct
 */
type Api struct {
	StockData struct{
		Url string `json:"url"`
		ApiKey string `json:"api_key"`
	} `json:"stock_data"`
}

/*
Get the api data
 */
func getAPIData(fileName string) (*Api){

	/*
	Open
	 */
	apiFile,err := os.Open(fileName)
	utils.HandleErr(err)

	defer apiFile.Close()

	/*
		Read
	 */
	byteValue,err := ioutil.ReadAll(apiFile)
	utils.HandleErr(err)

	/*
	Unmarshal
	 */
	apiData := new(Api)
	err = json.Unmarshal(byteValue, apiData)
	utils.HandleErr(err)

	return apiData

}

/*
Get the data set for SMA data
 */
func getSMAData(ts TimeSeries, sma50 SimpleMovingAv, sma15 SimpleMovingAv) (*SMAData, error){

	var smaData SMAData
	var sma SMA

	for date,data := range ts.Data {
		sma.Date = date
		sma.Open = data.Open
		sma.Close = data.Close
		sma.High = data.High
		sma.Low = data.Low
		sma.Volume = data.Volume
		sma.SMA50Day = sma50.Data[date].Value
		sma.SMA15Day = sma15.Data[date].Value

		smaData.Data = append(smaData.Data, sma)
	}

	return  &smaData, nil
}

/*
	http client for api calls to Data source
 */
var myClient = &http.Client{Timeout: 10 * time.Second}

/*
	API data for the stock data
 */
var apiData = getAPIData("api.json")

/*
 	Get the Time Series Data
 */
func getTimeSeries(symbol string) (*TimeSeries, error) {

	target := new(TimeSeries)

	url := apiData.StockData.Url + "/query?function=TIME_SERIES_DAILY&symbol=" +
		symbol +
		"&apikey=" +
		apiData.StockData.ApiKey

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

	url := apiData.StockData.Url + "/query?function=SMA&symbol=" +
		symbol +
		"&interval=daily&time_period=" +
		strconv.Itoa(window) +
		"&series_type=close&apikey=" +
		apiData.StockData.ApiKey

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

	url := apiData.StockData.Url + "/query?function=SMA&symbol=" +
		symbol +
		"&interval=daily&time_period=" +
		strconv.Itoa(window) +
		"&series_type=close&apikey=" +
		apiData.StockData.ApiKey

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

