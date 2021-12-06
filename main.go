package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserLocation struct {
	Userid       int        `json:"userid"`
	Userlocation [2]float64 `json:"userlocation"`
}

type PartnerLocation struct {
	Partnerid       int        `json:"partnerid"`
	Partnerlocation [2]float64 `json:"partnerlocation"`
}

var Partners [5]PartnerLocation

var bogotalimit1 [2]float64 = [2]float64{4.784373420346989, -73.99595035691422}
var bogotalimit2 [2]float64 = [2]float64{4.490865002856506, -74.27598825255971}

func generatePartners(userlocation [2]float64) {
	var maxdistance float64 = 2000
	maxdeg := maxdistance * (0.00001 / 1.11)
	maxpartnerx := userlocation[0] + maxdeg
	minpartnerx := userlocation[0] - maxdeg
	maxpartnery := userlocation[1] + maxdeg
	minpartnery := userlocation[1] - maxdeg

	for i := 0; i < 5; i++ {
		x := minpartnerx + rand.Float64()*(maxpartnerx-minpartnerx)
		y := minpartnery + rand.Float64()*(maxpartnery-minpartnery)

		Partners[i].Partnerid = i
		Partners[i].Partnerlocation = [2]float64{x, y}
	}
}

func manhattanDistance(coord1 [2]float64, coord2 [2]float64) float64 {
	var degdistance float64 = 0
	var distance float64 = 0
	for i := 0; i < 2; i++ {
		degdistance += math.Abs(coord1[i] - coord2[i])
	}

	distance = degdistance * (1.11 / 0.00001)

	return distance
}

func closestPartners(userlocation [2]float64) [5]PartnerLocation {
	generatePartners(userlocation)
	orderedpartners := Partners

	n := len(Partners)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			distance1 := manhattanDistance(userlocation, orderedpartners[j].Partnerlocation)
			distance2 := manhattanDistance(userlocation, orderedpartners[j+1].Partnerlocation)
			if distance1 > distance2 {
				temp := orderedpartners[j]
				orderedpartners[j] = orderedpartners[j+1]
				orderedpartners[j+1] = temp
			}
		}
	}

	return orderedpartners
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: home")

	fmt.Fprintf(w, "UN-Uber – interested_ms")
	fmt.Fprintf(w, "Developed by Gerson Nicolás Pineda Lara for Arquitectura de Software 2021-II")
}

func returnAllPartners(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPartners")

	json.NewEncoder(w).Encode(Partners)
}

func returnClosestsPartners(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: closestPartner")

	var userlocation []UserLocation
	err := json.NewDecoder(r.Body).Decode(&userlocation)
	if err != nil {
		fmt.Println(err)
	}

	// Checking coordinates limits for Bogotá
	if userlocation[0].Userlocation[0] < bogotalimit2[0] ||
		userlocation[0].Userlocation[0] > bogotalimit1[0] ||
		userlocation[0].Userlocation[1] > bogotalimit1[1] ||
		userlocation[0].Userlocation[1] < bogotalimit2[1] {
		fmt.Fprintf(w, "Coordinates aren't from Bogotá")
	}
	closestspartners := closestPartners(userlocation[0].Userlocation)
	fmt.Fprintf(w, "%+v\n", closestspartners)
}

func returnSingleClosestPartner(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnSingleClosestPartner")

	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["partnerid"])
	if err != nil {
		fmt.Println(err)
	}

	for _, partner := range Partners {
		if partner.Partnerid == key {
			json.NewEncoder(w).Encode(partner)
		}
	}
}

func handleRequests() {
	myrouter := mux.NewRouter().StrictSlash(true)
	myrouter.HandleFunc("/", home)
	myrouter.HandleFunc("/all-partners", returnAllPartners).Methods("GET")
	myrouter.HandleFunc("/partner/{partnerid}", returnSingleClosestPartner).Methods("GET")
	myrouter.HandleFunc("/closests-partners", returnClosestsPartners).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myrouter))
}

func main() {
	fmt.Println("Trying")
	Partners = [5]PartnerLocation{
		{Partnerid: 1, Partnerlocation: [2]float64{4.64087273964227, -74.08795803227153}},
		{Partnerid: 2, Partnerlocation: [2]float64{4.6440701368879855, -74.0851685349478}},
		{Partnerid: 3, Partnerlocation: [2]float64{4.642439316245012, -74.08092019197085}},
		{Partnerid: 4, Partnerlocation: [2]float64{4.633499379247723, -74.07992241027031}},
		{Partnerid: 5, Partnerlocation: [2]float64{4.62881549146964, -74.08204671977938}}}

	/*var userlocation UserLocation = UserLocation{Userid: 1,
		Userlocation: [2]float64{4.638913870896805, -74.08381574852237}}

	postbody, _ := json.Marshal(userlocation)
	responsebody := bytes.NewBuffer(postbody)*/

	handleRequests()
}
