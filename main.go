package main

import (
	"SurfHotelsDumper/constants"
	"SurfHotelsDumper/databases"
	"SurfHotelsDumper/hasher"
	"SurfHotelsDumper/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"strconv"
	"time"
)

var (
	client = resty.New()
	Ctx    = context.TODO()
)

func AddPhotos(hotel *models.Hotel) {
	stringId := strconv.Itoa(hotel.Id)
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"id": stringId,
		}).
		Get("https://yasen.hotellook.com/photos/hotel_photos")

	if err != nil {
		fmt.Println(err.Error())
	}

	var photos map[string][]int
	err = json.Unmarshal(resp.Body(), &photos)

	hotel.PhotoHotel = photos[stringId]
}

const (
	startDate string = "2022-05-26"
	endDate   string = "2022-06-10"
	currency  string = "RUB"
	language  string = "ru"
)

func main() {
	const uri = "mongodb://localhost:27017"
	mongo, err := mongo.Connect(Ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Printf(err.Error())
	}

	defer func() {
		if err := mongo.Disconnect(Ctx); err != nil {
			panic(err)
		}
	}()

	if err := mongo.Ping(Ctx, readpref.Primary()); err != nil {
		log.Fatalf(err.Error())
	}

	client.
		SetRetryCount(200)

	coll := mongo.Database("surf-hotelDumper").Collection("hotels")

	iatas := []string{"MOW", "AER", "LED"}

	for _, iata := range iatas {

		searchIdhash := hasher.Md5HotelHasher(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s:%s:%s",
			constants.TOKEN,
			constants.MARKER,
			"1",
			startDate,
			endDate,
			currency,
			constants.CUSTOMER_IP,
			iata,
			"ru",
			"0",
		))

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
				"0",
				constants.MARKER,
				searchIdhash,
			))

		if err != nil {
			log.Printf(err.Error())
		}

		var searchId models.SearchIdResponse
		err = json.Unmarshal(respSearchId.Body(), &searchId)
		if err != nil {
			log.Printf(err.Error())
		}

		time.Sleep(20 * time.Second)

		hotelHash := hasher.Md5HotelHasher(fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
			constants.TOKEN,
			constants.MARKER,
			strconv.Itoa(constants.HOTELS_LIMIT),
			"0",
			"1",
			strconv.Itoa(searchId.SearchId),
			"0",
			"popularity",
		))

		resp, err := client.R().
			SetQueryParams(map[string]string{
				"searchId":   strconv.Itoa(searchId.SearchId),
				"limit":      strconv.Itoa(constants.HOTELS_LIMIT),
				"sortBy":     "popularity",
				"sortAsc":    "0",
				"roomsCount": "1",
				"offset":     "0",
				"marker":     constants.MARKER,
				"signature":  hotelHash,
			}).
			Get(fmt.Sprintf("%s/getResult.json", constants.HOTELLOOK_ADDR))

		var hotels models.HotelResponse
		err = json.Unmarshal(resp.Body(), &hotels)

		if err != nil {
			log.Printf(err.Error())
		}

		if len(hotels.Result) > 0 {
			for _, v := range hotels.Result {
				//AddPhotos(&v)

				// set iata for search
				v.Iata = iata

				// set room total to 0 for a skeleton app condition
				v.Rooms[0].Total = 0

				// set hotel photos now for omitting request to hotellook and fastest speed
				v.PhotoHotel = []int{5162091534, 5162091993}

				result, err := coll.InsertOne(databases.Ctx, v)

				if err != nil {
					log.Printf(err.Error())
				}

				fmt.Printf("Inserted %d\n", result.InsertedID)
			}
		}
	}
}
