package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

/*
Standard post to ES
 */
func post(url string, b *[]byte) error {

	uri := "http://localhost:5601/" + url

	client := &http.Client{}
	client.Timeout = time.Second * 15

	body := bytes.NewBuffer(*b)
	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return err
	}

	fmt.Println(body)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("kbn-xsrf", "true")
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

	return nil
}

/*
Send a request to create an index based on stock symbol name
 */
func createIndex(symbol string) error {

	url := "api/saved_objects/index-pattern/" +
		strings.ToLower(strings.Replace(symbol, ":", "_", -1))

	body := []byte("{\"attributes\": {\"title\": \"" +
		strings.ToLower(strings.Replace(symbol, ":", "_", -1)) +
		"-*\"}}")

	err := post(url, &body)
	if err != err{
		return err
	}

	return nil
}

/*
Send a request to create a SMA visualisation for each symbol
 */
 func createSMAVisualisation(symbol string, stock_name string) error {

 	indexPattern := strings.ToLower(strings.Replace(symbol, ":", "_", -1))
 	url := "api/saved_objects/visualization/" + indexPattern

	title := stock_name + " Moving Average"

 	visState := `"visState": "{\"title\":\"` + title +
 		`\",\"type\":\"line\",\"params\":{\"type\":\"line\",\"grid\":{\"categoryLines\":false,\"style\":{\"color\":\"#eee\"}},\"categoryAxes\":[{\"id\":\"CategoryAxis-1\",\"type\":\"category\",\"position\":\"bottom\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\"},\"labels\":{\"show\":true,\"truncate\":100},\"title\":{}}],\"valueAxes\":[{\"id\":\"ValueAxis-1\",\"name\":\"LeftAxis-1\",\"type\":\"value\",\"position\":\"left\",\"show\":true,\"style\":{},\"scale\":{\"type\":\"linear\",\"mode\":\"normal\"},\"labels\":{\"show\":true,\"rotate\":0,\"filter\":false,\"truncate\":100},\"title\":{\"text\":\"Closing Price\"}}],\"seriesParams\":[{\"show\":\"true\",\"type\":\"line\",\"mode\":\"normal\",\"data\":{\"label\":\"Closing Price\",\"id\":\"1\"},\"valueAxis\":\"ValueAxis-1\",\"drawLinesBetweenPoints\":true,\"showCircles\":true},{\"show\":true,\"mode\":\"normal\",\"type\":\"line\",\"drawLinesBetweenPoints\":true,\"showCircles\":true,\"data\":{\"id\":\"2\",\"label\":\"50 Day SMA\"},\"valueAxis\":\"ValueAxis-1\"},{\"show\":true,\"mode\":\"normal\",\"type\":\"line\",\"drawLinesBetweenPoints\":true,\"showCircles\":true,\"data\":{\"id\":\"3\",\"label\":\"15 Day SMA\"},\"valueAxis\":\"ValueAxis-1\"}],\"addTooltip\":true,\"addLegend\":true,\"legendPosition\":\"right\",\"times\":[],\"addTimeMarker\":false},\"aggs\":[{\"id\":\"1\",\"enabled\":true,\"type\":\"avg\",\"schema\":\"metric\",\"params\":{\"field\":\"Close\",\"customLabel\":\"Closing Price\"}},{\"id\":\"2\",\"enabled\":true,\"type\":\"avg\",\"schema\":\"metric\",\"params\":{\"field\":\"SMA50Day\",\"customLabel\":\"50 Day SMA\"}},{\"id\":\"3\",\"enabled\":true,\"type\":\"avg\",\"schema\":\"metric\",\"params\":{\"field\":\"SMA15Day\",\"customLabel\":\"15 Day SMA\"}},{\"id\":\"4\",\"enabled\":true,\"type\":\"date_histogram\",\"schema\":\"segment\",\"params\":{\"field\":\"Date\",\"useNormalizedEsInterval\":true,\"interval\":\"d\",\"time_zone\":\"Europe/London\",\"drop_partials\":false,\"customInterval\":\"2h\",\"min_doc_count\":1,\"extended_bounds\":{},\"customLabel\":\"Date\"}}]}",`
	 body := []byte(`{
    "attributes": {
      "title": "` + title + `",` + visState +
      `"uiStateJSON": "{\"vis\":{\"colors\":{\"50 Day SMA\":\"#629E51\",\"Closing Price\":\"#447EBC\",\"15 Day SMA\":\"#BF1B00\"}}}",
      "description": "",
      "version": 1,
      "kibanaSavedObjectMeta": {
        "searchSourceJSON": "{\"index\":\"` + indexPattern + `\",\"query\":{\"query\":\"\",\"language\":\"lucene\"},\"filter\":[]}"
      }
	}
  }`)

	 err := post(url, &body)
	 if err != err{
		 return err
	 }

	 return nil

 }
