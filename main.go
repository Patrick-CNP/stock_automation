package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type fav_symbols struct {
	Group []struct{
		Name string `yaml:"name"`
		Symbol string `yaml:"symbol"`
	} `yaml:"favourites"`
}

func (symbols *fav_symbols) symbols(yml string) (error){
	/*
	Read stock symbols from the yaml file provided
	and populate the provided struct with them
	 */
	yamlFile,err := ioutil.ReadFile(yml)
	if err != nil {
		return err
	}

	//fmt.Println(string(yamlFile))

	err = yaml.Unmarshal(yamlFile, symbols)
	if err != nil{
		return err
	}

	return nil
}

func es(stock_symbol string){

	/*
	ES index prefix and field to use
	 */
	index_name_prefix := "/" + strings.ToLower(strings.Replace(stock_symbol, ":", "_", -1))
	index_field := "/stock/"
	var index_name string

	/*
	Time Series
	 */
	 ts,err := getTimeSeries(stock_symbol)
	 handle(err)
	 index_name = index_name_prefix + "-ts" + index_field
	 //err = ts.Send(index_name, 10)
	 //handle(err)

	 /*
	 Simple Moving Av 50 Day Av
	  */
	sma50, err := getSimpleMovingAv(stock_symbol, 50)
	handle(err)
	index_name = index_name_prefix + "-sma50" + index_field
	//err = sma50.Send(index_name, 10)
	//handle(err)

	/*
	Simple Moving Av 15 Day Av
 	*/
	sma15, err := getSimpleMovingAv(stock_symbol, 15)
	handle(err)
	index_name = index_name_prefix + "-sma15" + index_field
	//err = sma15.Send(index_name, 10)
	//handle(err)

	/*
	Convert the data to an SMAData object to send to ES
	The raw data is too much for ES to handle with lots of
	symbols.
	 */
	 smaData,err := getSMAData(*ts,*sma50, *sma15)
	 handle(err)
	 err = smaData.Send(index_name, 10)
	 handle(err)

}

func from_arg(){
	/*
	The symbol to pull data for
	*/
	stock_symbol := os.Args[1]
	stock_name := os.Args[2]

	/*
	run the stock
	 */
	es(stock_symbol)

	/*
	Index
	 */
	err := createIndex(stock_symbol)
	handle(err)

	/*
	Vis
	 */
	err = createSMAVisualisation(stock_symbol, stock_name)
	handle(err)

}

func from_yaml(){
	var s fav_symbols
	err := s.symbols("symbols.yml")
	handle(err)

	for _,g := range s.Group{
		fmt.Println(g.Symbol)
		/*
		run the stock
	 	*/
		es(g.Symbol)
		/*
		Index
 		*/
		err := createIndex(g.Symbol)
		handle(err)

		/*
		Vis
		 */
		err = createSMAVisualisation(g.Symbol, g.Name)
		handle(err)

		// pause to let ES catchup
		time.Sleep(60 * time.Second)
	}
}

func newIndex(){

	var s fav_symbols
	err := s.symbols("symbols.yml")
	handle(err)

	for _,g := range s.Group {
		fmt.Println(g.Name)
		err := createIndex(g.Symbol)
		handle(err)
	}
}

func newVis(){

	var s fav_symbols
	err := s.symbols("symbols.yml")
	handle(err)

	for _,g := range s.Group {
		fmt.Println(g.Name)
		err := createSMAVisualisation(g.Symbol, g.Name)
		handle(err)
	}
}

func main()  {
	//from_arg()
	from_yaml()
}
