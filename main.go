package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// Objeto para la ubicación del usuario
type UserLocation struct {
	Userid       int        `json:"userid"`
	Userlocation [2]float64 `json:"userlocation"`
}

// Objeto para los socios
type PartnerLocation struct {
	Partnerid       int        `json:"partnerid"`
	Partnerlocation [2]float64 `json:"partnerlocation"`
}

// Lista de socios que se "buscarán"
var Partners [5]PartnerLocation

// Rectángulo con los límites de la ciudad en los extremos nororiental y suroccidental
var bogotalimit1 [2]float64 = [2]float64{4.784373420346989, -73.99595035691422}
var bogotalimit2 [2]float64 = [2]float64{4.490865002856506, -74.27598825255971}

// Función que verifica que la variable de entorno del puerto esté activa.
// Si ya lo está, asigna el valor pasado en los argumentos. Si no, crea la variable
// y asigna el valor pasado en los argumentos.
func getEnv(fallback string) string {
	value, exists := os.LookupEnv("PORT")
	if !exists {
		value = fallback
	}
	return ":" + value
}

// Se asume que las coordenas en grados son equivalentes a unas rectangulares dadas
// las dimensiones de la ciudad. Sin embargo, existe una equivalencia entre grados
// y metros.

func coordinatesChecker(coord [2]float64) bool {
	valid := true

	if len(coord) == 0 ||
		coord[0] < bogotalimit2[0] || coord[0] > bogotalimit1[0] ||
		coord[1] > bogotalimit1[1] || coord[1] < bogotalimit2[1] {
		valid = false
	}

	return valid
}

// Función que genera la cantidad de socios indicados en la variable global
// Partners. Los genera dentro de una circunferencia con centro en la ubicación del
// usuario y radio maxdistance (en metros).
func generatePartners(userlocation [2]float64) {
	// Máxima distancia al usuario en metros
	var maxmdistance float64 = 2000
	// La distancia se convierte a grados
	var maxddistance float64 = maxmdistance * (0.00001 / 1.11)

	for i := range Partners {
		// Se genera un ángulo aleatorio
		angle := 2 * math.Pi * rand.Float64()
		// Se genera un radio aleatorio dentro del rango. Se obtiene la
		// raíz cuadrada para que la distribución sea uniforme.
		rand_rad := maxddistance * math.Sqrt(rand.Float64())
		// Se calculan las coordenadas
		x := rand_rad*math.Cos(angle) + userlocation[0]
		y := rand_rad*math.Sin(angle) + userlocation[1]

		Partners[i].Partnerid = i
		Partners[i].Partnerlocation = [2]float64{x, y}
	}
}

// Función que calcula la distancia Manhattan entre dos puntos.
func manhattanDistance(coords1 [2]float64, coords2 [2]float64) float64 {
	var degdistance float64 = 0
	var distance float64 = 0
	// Se calcula la distancia en grados pero se convierte a metros
	for i := 0; i < 2; i++ {
		degdistance += math.Abs(coords1[i] - coords2[i])
	}

	distance = degdistance * (1.11 / 0.00001)

	return distance
}

// Función que organiza los socios generados de más a menos cercano según la
// distancia Manhattan y retorna el más cercano. Usa bubble sort.
func closestPartners(userlocation [2]float64) PartnerLocation {
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

	return orderedpartners[0]
}

// Función del endpoint sin ruta
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: home")

	fmt.Fprintf(w, "UN-Uber – interested_ms\n")
	fmt.Fprintf(w, "Developed by Gerson Nicolás Pineda Lara for Arquitectura de Software 2021-II")
}

// Función del endpoint /all-partners
func returnAllPartners(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPartners")

	json.NewEncoder(w).Encode(Partners)
}

// Función del endpoint /closest-partner
func returnClosestPartner(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: closestPartner")

	// Fue necesario declarar un arreglo a pesar de que solo se necesita 1
	var userlocation []UserLocation
	err := json.NewDecoder(r.Body).Decode(&userlocation)
	if err == io.EOF {
		fmt.Println("Request without information")
		json.NewEncoder(w).Encode("Bad user location information")
	} else if err != nil {
		fmt.Println(err)
	} else if !coordinatesChecker(userlocation[0].Userlocation) {
		// Se verifica que las coordenas entén en Bogotá
		fmt.Println("Coordinates aren't from Bogotá")
		json.NewEncoder(w).Encode("Bad coordinates")
	} else {
		closestspartners := closestPartners(userlocation[0].Userlocation)
		json.NewEncoder(w).Encode(closestspartners)
	}
}

// Función del endpoint /partner/{partnerid}
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

// Función enrutadora
func handleRequests() {
	myrouter := mux.NewRouter().StrictSlash(true)
	myrouter.HandleFunc("/", home)
	myrouter.HandleFunc("/all-partners", returnAllPartners).Methods("GET")
	myrouter.HandleFunc("/partner/{partnerid}", returnSingleClosestPartner).Methods("GET")
	myrouter.HandleFunc("/closest-partner", returnClosestPartner).Methods("POST")
	log.Fatal(http.ListenAndServe(getEnv("10000"), myrouter))
}

func main() {
	handleRequests()
}
