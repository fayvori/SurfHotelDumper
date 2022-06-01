package main

import (
	"SurfHotelsDumper/constants"
	"SurfHotelsDumper/hasher"
	"SurfHotelsDumper/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var (
	client = resty.New()
	Ctx    = context.TODO()
)

func ReverseSlice(data interface{}) {
	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		panic(errors.New("data must be a slice type"))
	}
	valueLen := value.Len()
	for i := 0; i <= int((valueLen-1)/2); i++ {
		reverseIndex := valueLen - 1 - i
		tmp := value.Index(reverseIndex).Interface()
		value.Index(reverseIndex).Set(value.Index(i))
		value.Index(i).Set(reflect.ValueOf(tmp))
	}
}

func AddPhotosToHotelDbResponse(hotels *models.HotelResponse) {
	var hotelsIds []string

	if len(hotels.Result) > 200 {
		for i := 0; i < 200; i++ {
			hotelsIds = append(hotelsIds, strconv.Itoa(hotels.Result[i].Id))
		}
	} else {
		for i := 0; i < len(hotels.Result); i++ {
			hotelsIds = append(hotelsIds, strconv.Itoa(hotels.Result[i].Id))
		}
	}

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"id": strings.Join(hotelsIds, ","),
		}).
		Get("https://yasen.hotellook.com/photos/hotel_photos")

	if err != nil {
		fmt.Println(err.Error())
	}

	var photos map[string][]int
	err = json.Unmarshal(resp.Body(), &photos)

	for i := 0; i < len(hotels.Result); i++ {
		id := strconv.Itoa(hotels.Result[i].Id)
		id1, _ := strconv.Atoi(id)

		if hotels.Result[i].Id == id1 {
			hotels.Result[i].PhotoHotel = photos[id]
		}
	}
}

const (
	startDate string = "2022-06-04"
	endDate   string = "2022-06-16"
	currency  string = "RUB"
	language  string = "ru"
)

type MarshaledIatas struct {
	Countries []string `json:"countries"`
}

func main() {
	const uri = "mongodb://localhost:27017"
	connect, err := mongo.Connect(Ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Printf(err.Error())
	}

	defer func() {
		if err := connect.Disconnect(Ctx); err != nil {
			panic(err)
		}
	}()

	if err := connect.Ping(Ctx, readpref.Primary()); err != nil {
		log.Fatalf(err.Error())
	}

	client.
		SetRetryCount(400)

	coll := connect.Database("surf-hotelDumper").Collection("hotels")

	var iatas MarshaledIatas

	f, err := ioutil.ReadFile("optimizedCountries.json")
	err = json.Unmarshal(f, &iatas)

	for _, iata := range iatas.Countries {
		hotelsHash := hasher.Md5HotelHasher(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s:%s:%s",
			constants.TOKEN,
			constants.MARKER,
			"1",
			startDate,
			endDate,
			currency,
			constants.CUSTOMER_IP,
			iata,
			"ru",
			"1",
		))

		fmt.Println(iata)

		respSearchId, err := client.R().
			EnableTrace().
			Get(fmt.Sprintf("%s/start.json?iata=%s&checkIn=%s&checkOut=%s&adultsCount=%s&customerIP=%s&lang=%s&currency=%s&waitForResult=%s&marker=%s&signature=%s",
				constants.HOTELLOOK_ADDR,
				iata,
				startDate,
				endDate,
				"1",
				constants.CUSTOMER_IP,
				language,
				currency,
				"1",
				constants.MARKER,
				hotelsHash,
			))

		fmt.Println(string(respSearchId.Body()))

		if err != nil {
			log.Printf(err.Error())
		}

		var hotels models.HotelResponse
		err = json.Unmarshal(respSearchId.Body(), &hotels)
		if err != nil {
			log.Printf(err.Error())
		}

		//fmt.Printf("im here")
		//time.Sleep(30 * time.Second)
		//
		//hotelHash := hasher.Md5HotelHasher(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
		//	constants.TOKEN,
		//	constants.MARKER,
		//	strconv.Itoa(constants.HOTELS_LIMIT),
		//	"0",
		//	"1",
		//	strconv.Itoa(searchId.SearchId),
		//	"0",
		//	"popularity",
		//))
		//
		//resp, err := client.R().
		//	SetQueryParams(map[string]string{
		//		"searchId":   strconv.Itoa(searchId.SearchId),
		//		"limit":      strconv.Itoa(constants.HOTELS_LIMIT),
		//		"sortBy":     "popularity",
		//		"sortAsc":    "0",
		//		"roomsCount": "1",
		//		"offset":     "0",
		//		"marker":     constants.MARKER,
		//		"signature":  hotelHash,
		//	}).
		//	Get(fmt.Sprintf("%s/getResult.json", constants.HOTELLOOK_ADDR))
		//
		//var hotels models.HotelResponse
		//err = json.Unmarshal(resp.Body(), &hotels)
		//
		//if err != nil {
		//	log.Printf(err.Error())
		//}

		fmt.Println(len(hotels.Result))
		if len(hotels.Result) > 0 {
			AddPhotosToHotelDbResponse(&hotels)

			ReverseSlice(hotels.Result)
			for _, v := range hotels.Result {
				// set iata for searching
				v.Iata = iata

				// set room total to 0 for a skeleton app condition
				v.Rooms[0].Total = 0

				result, err := coll.InsertOne(constants.Ctx, v)
				if err != nil {
					log.Printf(err.Error())
				}

				fmt.Printf("Inserted %d\n", result.InsertedID)
			}
		}
	}
}
