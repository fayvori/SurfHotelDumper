package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchIdResponse struct {
	Status   string `json:"status"`
	SearchId int    `json:"searchId"`
}

type HotelOptions struct {
	HalfBoard         bool `json:"halfBoard"`
	AllInclusive      bool `json:"allInclusive"`
	FullBoard         bool `json:"fullBoard"`
	UltraAllInclusive bool `json:"ultraAllInclusive"`
	Refundable        bool `json:"refundable"`
	Breakfast         bool `json:"breakfast"`
	Available         int  `json:"available"`
}

type HotelRoom struct {
	Type           string       `json:"type"`
	AgencyName     string       `json:"agencyName"`
	Desc           string       `json:"desc"`
	FullBookingURL string       `json:"fullBookingURL"`
	Tax            int          `json:"tax"`
	AgencyId       string       `json:"agencyId"`
	Options        HotelOptions `json:"options"`
	Price          int          `json:"price"`
	BookingURL     string       `json:"bookingURL"`
	Total          int          `json:"total"`
	InternalTypeId string       `json:"internalTypeId"`
}

type HotelLocation struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type HotelPhoto struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"json:"id"`
	HotelId int                `bson:"hotelId,omitempty"json:"hotelId"`
	Photos  []int              `bson:"photos,omitempty"json:"photos"`
}

type Hotel struct {
	MaxPricePerNight int            `json:"maxPricePerNight"`
	Location         HotelLocation  `json:"location"`
	FullUrl          string         `json:"fullUrl"`
	MinPriceTotal    int            `json:"minPriceTotal"`
	Id               int            `json:"id"`
	PhotosByRoomType map[string]int `json:"photosByRoomType"`
	Url              string         `json:"url"`
	PhotoCount       int            `json:"photoCount"`
	Popularity       int            `json:"popularity"`
	MaxPrice         int            `json:"maxPrice"`
	Amenities        []int          `json:"amenities"`
	Address          string         `json:"address"`
	GuestScore       int            `json:"guestScore"`
	Name             string         `json:"name"`
	Price            int            `json:"price"`
	Rating           int            `json:"rating"`
	Rooms            []HotelRoom    `json:"rooms"`
	Stars            int            `json:"stars"`
	Distance         float64        `json:"distance"`
	PhotoHotel       []int          `json:"photoHotel"`
	Iata             string         `json:"iata""`
}

type HotelResponse struct {
	Result []Hotel `json:"result"`
}

type Pagination struct {
	Page          int    `json:"page"`
	PageSize      int    `json:"page_size"`
	Total         int    `json:"total"`
	MinPrice      int    `json:"minPrice"`
	MaxPrice      int    `json:"maxPrice"`
	TotalFiltered int    `json:"total_filtered"`
	First         string `json:"first"`
	Last          string `json:"last"`
	Prev          string `json:"prev"`
	Next          string `json:"next"`
	Query         struct {
		Page     []string `json:"page"`
		PageSize []string `json:"page_size"`
	} `json:"query"`
}

type PaginationHotels struct {
	Pagination `json:"pagination"`
	Result     []Hotel `json:"result"`
}
