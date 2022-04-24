package controller

import (
	"api/src/helpers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

var intVar int

type Motd struct {
	Msg string `json:"msg"`
	Url string `json:"url"`
}

type Usd struct {
	USD float64 `json:"USD"`
}

type Count struct {
	Total int
}

type ResponseConvert struct {
	Motd    Motd
	Success bool   `json:"success"`
	Base    string `json:"base"`
	Date    string `json:"date"`
	Rates   Usd
}

type OutCount struct {
	Komoditas []string `json:"komoditas"`
}

type TempOutputtgl struct {
	Tgl       string   `json:"tgl"`
	Komoditas []string `json:"komoditas"`
}

type TempOutput struct {
	Area_provinsi string    `json:"area_provinsi"`
	Price         []float64 `json:"price"`
}

type Output struct {
	Area_provinsi string `json:"area_provinsi"`
	Price         Price
}

type Price struct {
	Total  float64 `json:"total"`
	Avg    float64 `json:"avg"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Median float64 `json:"median"`
}

type Response struct {
	UUid          string `json:"uuid"`
	Komoditas     string `json:"komoditas"`
	Area_provinsi string `json:"area_provinsi"`
	Area_kota     string `json:"area_kota"`
	Size          string `json:"size"`
	Price         string `json:"price"`
	Tgl           string `json:"tgl_parsed"`
	Created_at    string `json:"timestamp"`
}

var Currency float64

type FetchController struct{}

var Cache = cache.New(5*time.Minute, 5*time.Minute)

func (f FetchController) Resource(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//Resources from efishery
		response, err := http.Get("https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list")
		if err != nil {
			error.ApiError(w, http.StatusInternalServerError, err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			error.ApiError(w, http.StatusNotFound, err.Error())
		}

		var responseObject []Response
		json.Unmarshal(responseData, &responseObject)
		for _, d := range responseObject {

			intVar, _ := strconv.ParseFloat(d.Price, 64)
			val := ConvertPrice(intVar)
			outd := struct {
				UUid          string  `json:"uuid"`
				Komoditas     string  `json:"komoditas"`
				Area_provinsi string  `json:"area_provinsi"`
				Area_kota     string  `json:"area_kota"`
				Size          string  `json:"size"`
				Price         string  `json:"price"`
				Tgl           string  `json:"tgl_parsed"`
				Created_at    string  `json:"timestamp"`
				Price_USD     float64 `json:"price_usd"`
			}{
				UUid:          d.UUid,
				Komoditas:     d.Komoditas,
				Area_provinsi: d.Area_provinsi,
				Area_kota:     d.Area_kota,
				Size:          d.Size,
				Price:         d.Price,
				Tgl:           d.Tgl,
				Created_at:    d.Created_at,
				Price_USD:     val,
			}
			helpers.RespondWithJSON(w, outd)
		}
	}
}

func ConvertPrice(in float64) float64 {

	// Check caching
	Currency, found := Cache.Get("currency")
	if found {
		Nil := Currency.(float64) * in
		return Nil
	}

	// if not caching
	resp, err := http.Get("https://api.exchangerate.host/latest?symbols=USD")
	if err != nil {
		log.Fatalln(err)
	}
	responseDataConv, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ugh: ", err)

	}
	var i ResponseConvert

	if err := json.Unmarshal([]byte(responseDataConv), &i); err != nil {
		fmt.Println("ugh: ", err)
	}

	Cache.Set("currency", i.Rates.USD, cache.NoExpiration)
	Nil := in * i.Rates.USD
	return Nil
}

func (f FetchController) TestResource(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var polo float64
		Cache.Set("currency", 90000, cache.NoExpiration)
		foo, found := Cache.Get("currency")
		if found {
			Nil := foo.(float64) * polo
			helpers.RespondWithJSON(w, Nil)
		}
	}
}

func (u FetchController) ResourcesUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		token := bearerToken[1]
		data, err := helpers.VerifyJwtToken(token)
		if err != nil {
			error.ApiError(w, http.StatusForbidden, err.Error())
			return
		}

		helpers.RespondWithJSON(w, data)
	}
}

func (u FetchController) Aggregate(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		token := bearerToken[1]
		data, err := helpers.VerifyJwtToken(token)
		if err != nil {
			error.ApiError(w, http.StatusForbidden, err.Error())
			return
		}

		if data.Role != "Admin" {
			error.ApiError(w, http.StatusInternalServerError, "You are not Admin!")
		}

		response, err := http.Get("https://stein.efishery.com/v1/storages/5e1edf521073e315924ceab4/list")
		if err != nil {
			error.ApiError(w, http.StatusInternalServerError, err.Error())
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			error.ApiError(w, http.StatusNotFound, err.Error())
		}

		var responseObject []Response
		json.Unmarshal(responseData, &responseObject)

		var outputlist = make(map[string]TempOutput)
		for _, inp := range responseObject {
			if inp.Price != "" {
				intVar, _ := strconv.ParseFloat(inp.Price, 64)
				outputlist[inp.Area_provinsi] = TempOutput{inp.Area_provinsi, append(outputlist[inp.Area_provinsi].Price, intVar)}
			}
		}

		var outputs = make(map[string]Output)
		for _, outp := range outputlist {

			var Total float64
			var min = outp.Price[0]
			var max = outp.Price[0]

			median, _ := stats.Median(outp.Price)
			roundedMedian, _ := stats.Round(median, 0)

			for _, value := range outp.Price {
				if value < min {
					min = value
				}
				if value > max {
					max = value
				}
				Total += value
			}
			len := len(outp.Price)
			avg := Total / float64(len)

			acas := Price{Total: Total, Avg: avg, Min: min, Max: max, Median: roundedMedian}
			outputs[outp.Area_provinsi] = Output{outp.Area_provinsi, acas}
		}

		helpers.RespondWithJSON(w, outputs)
	}

}
