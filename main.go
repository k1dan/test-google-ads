package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	WrongRequestBody = "WRONG_REQUEST"
	WrongHeaders	 = "WRONG_HEADERS"
)

// ErrorResponse represents error that is returned from Auth API
type ErrorResponse struct {
	Err Error `json:"error"`
}

// Error represents error body that is returned from Auth API
type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// CampaignErrResponse represents error that is returned from Auth API for retrieve campaign request
type CampaignErrResponse struct {
	Err struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

//Payload is a type for Payload
type Payload struct {
	Query string `json:"query"`
}

//CampaignsStats is a type for Campaigns Stats
type CampaignsStats struct {
	Results   []Result `json:"results"`
	FieldMask string   `json:"fieldMask"`
	RequestID string   `json:"requestId"`
}

//Campaign is a type for Campaign
type Campaign struct {
	ResourceName           string `json:"resourceName"`
	AdvertisingChannelType string `json:"advertisingChannelType"`
	Name                   string `json:"name"`
	ID                     string `json:"id"`
}

//Metrics is a type for Metrics
type Metrics struct {
	Clicks      string `json:"clicks"`
	Spend       string  `json:"costMicros"`
	Impressions string `json:"impressions"`
}

//Segments is a type for Segments
type Segments struct {
	Date string `json:"date"`
}

//Customer is a type for Customer
type Customer struct {
	ResourceName string `json:"resourceName"`
	Currency string `json:"currencyCode"`
}

//Result is a type for Result
type Result struct {
	Campaign Campaign `json:"campaign"`
	Metrics  Metrics  `json:"metrics"`
	Segments Segments `json:"segments"`
	Customer Customer `json:"customer"`
}

func searchStream(w http.ResponseWriter, r *http.Request) {
	
	var p Payload

	developerToken := r.Header.Get("developer-token")
	loginCustomerID := r.Header.Get("login-customer-id")
	authorization := r.Header.Get("Authorization")

	if developerToken == "" || loginCustomerID == "" || authorization == "" {
		respondWithError(w, http.StatusUnauthorized, "headers are missing", WrongHeaders)
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "unable to read request body", WrongRequestBody)
	}

	defer r.Body.Close()

	customerID := mux.Vars(r)["id"]
	dateStart, dateEnd := parseDates(p.Query)

	campaignStats := createResults(customerID, dateStart, dateEnd)

	respondWithJSON(w, http.StatusOK, campaignStats)


}

func createResults(customerID, dateStart, dateEnd string) CampaignsStats {

	mockCampaigns := []Campaign{
		{
			ResourceName: "customers/" + customerID + "/campaigns/" + "14344919920",
			AdvertisingChannelType: "DISPLAY",
			Name: "customer " + customerID + " Mock Campaign 1",
			ID: customerID,
		},
		{
			ResourceName: "customers/" + customerID + "/campaigns/" + "14344919921",
			AdvertisingChannelType: "SEARCH",
			Name: "customer " + customerID + "Mock Campaign 2",
			ID: customerID,
		},
		{
			ResourceName: "customers/" + customerID + "/campaigns/" + "143449199212",
			AdvertisingChannelType: "SEARCH",
			Name: "customer " + customerID + "Mock Campaign 3",
			ID: customerID,
		},
		{
			ResourceName: "customers/" + customerID + "/campaigns/" + "14344919923",
			AdvertisingChannelType: "DISPLAY",
			Name: "customer " + customerID + "Mock Campaign 4",
			ID: customerID,
		},
	}

	var results []Result

	for _, v := range mockCampaigns {

		result := Result{
			Customer: Customer{
				ResourceName: "customers/" + customerID,
				Currency: "EUR",
			},
			Campaign: v,
			Metrics: Metrics{
				Clicks: strconv.Itoa(getRandom(0, 100)),
				Spend: strconv.Itoa(getRandom(1000, 1000000)),
				Impressions: strconv.Itoa(getRandom(0, 2000)),
			},
			Segments: Segments{
				Date: dateEnd,
			},
		}

		results = append(results, result)
	}

	return CampaignsStats{
		Results: results,
		FieldMask: "mock field mask",
		RequestID: "mock request id",
	}

}


func getRandom(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min + 1) + min
}

func parseDates(s string) (string, string){
	r := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	m := r.FindAllStringSubmatch(s, 2)

	return m[0][0], m[1][0]
}

func respondWithError(w http.ResponseWriter, code int, message, status string,) {

	errorResponse := ErrorResponse{
		Err: Error{
			Code: code,
			Message: message,
			Status: status,
		},
	}

	respondWithJSON(w, code, errorResponse)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/customers/{id}/googleAds:searchStream", searchStream)

	log.Fatal(http.ListenAndServe(":8080", router))

}


