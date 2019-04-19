package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Symbol struct {
	Name string `yaml:"name"`
	Symbol string `yaml:"symbol"`
}

type SymbolGroups struct {
	Groups []map[string][]Symbol `yaml:"groups"`
}

func (s *SymbolGroups) symbols(yml string) (error){
	/*
		Read stock symbols from the yaml file provided
		and populate the provided struct with them
	 */
	yamlFile,err := ioutil.ReadFile(yml)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, s)
	if err != nil{
		return err
	}

	return nil
}

func es(stock_symbol string) error{

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
	 if err != nil {
		 time.Sleep(60 * time.Second)
		 return err
	 }
	 index_name = index_name_prefix + "-ts" + index_field

	 /*
	 	Simple Moving Av 50 Day Av
	  */
	sma50, err := getSimpleMovingAv(stock_symbol, 50)
	if err != nil {
		time.Sleep(60 * time.Second)
		return err
	}
	index_name = index_name_prefix + "-sma50" + index_field

	/*
		Simple Moving Av 15 Day Av
 	*/
	sma15, err := getSimpleMovingAv(stock_symbol, 15)
	if err != nil {
		time.Sleep(60 * time.Second)
		return err
	}
	index_name = index_name_prefix + "-sma15" + index_field

	/*
		Convert the data to an SMAData object to send to ES
		The raw data is too much for ES to handle with lots of
		symbols.
	 */
	 smaData,err := getSMAData(*ts,*sma50, *sma15)
	 handle(err)
	 err = smaData.Send(index_name, 10)
	 handle(err)

	 return nil
}

func from_arg(){
	/*
		The symbol to pull data for
	*/
	stock_symbol := os.Args[2]
	stock_name := os.Args[3]

	/*
		Run the stock
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

func from_yaml(symbolsFile string, section string){
	var s SymbolGroups
	err := s.symbols(symbolsFile)
	handle(err)

	for _,group := range s.Groups{
		for groupName,symbols := range group {
			if groupName == section {
				for _,ea := range symbols {
					fmt.Println(groupName)
					/*
					Run the stock
					*/
					err = es(ea.Symbol)
					if err != nil {
						continue
					}

					/*
					Index
					*/
					err := createIndex(ea.Symbol)
					handle(err)

					/*
					Vis
					*/
					err = createSMAVisualisation(ea.Symbol, ea.Name)
					handle(err)

					/*
					Pause to let ES catchup
					*/
					time.Sleep(60 * time.Second)
				}
			}
		}

	}

	/*
		Create the dashboards
 	*/
	dashboards_from_yaml(symbolsFile)

}

func dashboards_from_yaml(symbolsFile string){
	var s SymbolGroups
	err := s.symbols(symbolsFile)
	handle(err)
	var groupSymbols []string

	for _,group := range s.Groups{
		for groupName,symbols := range group {
			for _,ea := range symbols {
				groupSymbols = append(groupSymbols, ea.Symbol )
			}
			err = createDashBoard(groupName, &groupSymbols)
			handle(err)
		}
		/*
			Empty the symbols array
		 */
		groupSymbols = nil
	}
}

func main()  {
	switch os.Args[1] {
	case "symbol":
		from_arg()
	case "dashonly":
		dashboards_from_yaml(os.Args[2])
	default:
		from_yaml(os.Args[1], os.Args[2])
	}
}
