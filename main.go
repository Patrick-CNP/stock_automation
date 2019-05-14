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

func es(stock_symbol string) (bool, error){

	/*
	ES index prefix and field to use
	 */
	index_name_prefix := "/" + strings.ToLower(strings.Replace(stock_symbol, ":", "_", -1))
	index_field := "/stock/"
	var index_name string
	shouldWatch := false

	/*
		Time Series
	 */
	 ts,err := getTimeSeries(stock_symbol)
	 if err != nil {
		 time.Sleep(60 * time.Second)
		 return shouldWatch, err
	 }
	 index_name = index_name_prefix + "-ts" + index_field

	 /*
	 	Simple Moving Av 50 Day Av
	  */
	sma50, err := getSimpleMovingAv(stock_symbol, 50)
	if err != nil {
		time.Sleep(60 * time.Second)
		return shouldWatch, err
	}
	index_name = index_name_prefix + "-sma50" + index_field

	/*
		Simple Moving Av 15 Day Av
 	*/
	sma15, err := getSimpleMovingAv(stock_symbol, 15)
	if err != nil {
		time.Sleep(60 * time.Second)
		return shouldWatch, err
	}
	index_name = index_name_prefix + "-sma15" + index_field

	/*
		Convert the data to an SMAData object to send to ES
		The raw data is too much for ES to handle with lots of
		symbols.
	 */
	 smaData,err := getSMAData(*ts,*sma50, *sma15)

	 for _, sma := range(smaData.Data){
		 dateLayout := "2006-01-02"
	 	t, err := time.Parse(dateLayout, sma.Date)
	 	if err != nil{
	 		fmt.Println("Unable to parse date, skipping " + sma.Date )
	 		continue
		}
	 	fmt.Println("Checking time: " + t.String())
	 	fmt.Println("with :" + time.Now().Add(-72*time.Hour).String())
	 	if t.After(time.Now().Add(-96*time.Hour)){
	 		if sma.SMA50Day >= sma.Close && sma.Close > sma.SMA15Day{
				fmt.Println("Found a watcher: " + stock_symbol)
				shouldWatch = true
			} else if sma.Close >= sma.SMA50Day && sma.SMA50Day > sma.SMA15Day{
				fmt.Println("Found a watcher: " + stock_symbol)
				shouldWatch = true
			}

		}

	 }

	 handle(err)
	 err = smaData.Send(index_name, 10)
	 handle(err)

	 return shouldWatch,nil
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
	var watching []string

	for _,group := range s.Groups{
		for groupName,symbols := range group {
			if groupName == section {
				watching = nil
				for _,ea := range symbols {
					fmt.Println(groupName)

					/*
					Run the stock
					*/
					shouldWatch, err := es(ea.Symbol)
					if err != nil {
						continue
					}
					if shouldWatch {
						watching = append(watching, ea.Symbol)
					}

					/*
					Index
					*/
					err = createIndex(ea.Symbol)
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
				watchDash := groupName + "_Watching"
				fmt.Println(watchDash)
				fmt.Println(watching)
				err = createDashBoard(watchDash, &watching)
				handle(err)
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
	case "all":
		from_yaml(os.Args[2], "fav")
		from_yaml(os.Args[2], "ftse100")
		from_yaml(os.Args[2], "nasdaq")
		dashboards_from_yaml(os.Args[2])
	default:
		from_yaml(os.Args[1], os.Args[2])
	}
}
