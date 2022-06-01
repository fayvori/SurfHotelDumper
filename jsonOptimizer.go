package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Country struct {
	Id        string `json:"id"`
	Name      interface{}
	Code      string      `json:"code"`
	CountryId string      `json:"countryId"`
	Longitude string      `json:"longitude"`
	StateCode interface{} `json:"state_code"`
	Latitude  string      `json:"latitude"`
}

type Marshaled struct {
	Countries []string `json:"countries"`
}

func main() {

	var countries []Country
	dataFromFile, err := ioutil.ReadFile("countries.json")

	if err != nil {
		log.Fatal("error while tryna open file")
	}

	err = json.Unmarshal(dataFromFile, &countries)

	if err != nil {
		log.Fatal("unmarshalling error")
	}

	var iatasCodes []string

	for _, v := range countries {
		str := fmt.Sprintf("%s", v.Code)

		if str == "" {
			continue
		}

		iatasCodes = append(iatasCodes, str)
	}

	var toMarshal Marshaled
	toMarshal.Countries = iatasCodes

	marshaled, err := json.Marshal(toMarshal)

	f, err := os.Create("optimizedCountries.json")
	defer f.Close()

	_, err = f.Write(marshaled)
	if err != nil {
		log.Fatal("failed to write to file")
	}

	fmt.Println(string(marshaled))

}
