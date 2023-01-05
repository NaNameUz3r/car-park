package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/XANi/loremipsum"
)

type intslice []int

func (i *intslice) String() string {
	return fmt.Sprintf("%d", *i)
}

func (i *intslice) Set(value string) error {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		*i = append(*i, -1)
	} else {
		*i = append(*i, tmp)
	}
	return nil
}

var USAGE = "\nUSAGE: \n -e <ENTERPRISE_ID> -c <CARS AMOUNT FOR THIS ENTERPRISE> ... \n This pattern is repeatable. Pass arguments sequentially."

var loremIpsumGenerator = loremipsum.New()
var loremIpsum = strings.Split(loremIpsumGenerator.Paragraphs(4), " ")

func checkArgs() error {
	var err error

	neAgrs := fmt.Errorf("Not enough arguments.")
	wrongArgs := fmt.Errorf("Wrong argument combination.")
	if len(os.Args[1:]) <= 3 {
		err = neAgrs
		return err
	}

	e_counter := 0
	c_counter := 0
	for i, a := range os.Args {

		if a == "-e" {
			e_counter += 1
		}
		if a == "-c" {
			c_counter += 1
		}

		if a == "-c" && os.Args[i-2] != "-e" {
			err = wrongArgs
		}
		if a == "-e" && i > 1 && os.Args[i-2] != "-c" {
			err = wrongArgs
		}
	}

	if e_counter != c_counter {
		err = wrongArgs
	}
	return err
}

type Vehicle struct {
	Description      string `json:"description"`
	Price            uint   `json:"price"`
	Mileage          uint   `json:"mileage"`
	ManufacturedYear uint   `json:"manufactured"`

	EnterpriseID uint     `json:"enterprise_id"`
	CarModel     CarModel `json:"car_model" binding:"required"`
	CarModelID   uint     `json:"carmodel_id"`
}

type CarModel struct {
	Brand            string  `json:"brand"`
	CarType          string  `json:"carmodel"`
	FuelTankCapacity uint    `json:"fueltankcapacity"`
	Seats            uint    `json:"seats"`
	EngineVolume     float64 `json:"enginevolume"`
}

type Driver struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Salary    float64 `json:"salary"`

	VehicleID uint `json:"vehicle_id"`
	IsActive  bool `json:"is_active"`
}

var carBrands = []string{"BMW", "MERCEDES", "LADA"}
var carTypes = []string{"Sedan", "Truck", "Hatchback", "Sportcar"}
var Names = []string{"John", "Will", "Sergei", "Ivan", "Tendzin", "Hao", "Alex"}
var LastNames = []string{"Peterson", "Prokopenko", "Gyaltco", "Smith", "Chi", "Pots"}

func randomSubstring(s []string, length int) string {
	resultString := ""
	stringLength := len(s)
	for i := 0; i < length; i++ {
		rand.Seed(time.Now().UnixNano())
		resultString = resultString + " " + s[rand.Int()%stringLength]
	}
	resultString = strings.Title(resultString)
	resultString = strings.TrimSpace(resultString)
	return resultString
}

func main() {
	var enterprises intslice
	var cars intslice
	flag.Var(&enterprises, "e", "Enterprise ID.")
	flag.Var(&cars, "c", "Cars amount to generate.")
	flag.Parse()

	err := checkArgs()
	if err != nil {
		fmt.Println(err.Error(), USAGE)
		os.Exit(1)
	}

	PostVehicleUrl := "http://superadmin:superadminqwerty@localhost:8888/api/manager/666/vehicle/create"
	PostDrivereUrl := "http://superadmin:superadminqwerty@localhost:8888/api/save/drivers"

	for i, e := range enterprises {

		for car := 0; car < cars[i]; car++ {
			rand.Seed(time.Now().UnixNano())
			randomBrand := rand.Int() % len(carBrands)
			randomCarType := rand.Int() % len(carTypes)
			vehicle := Vehicle{
				Description:      randomSubstring(loremIpsum, 5),
				Price:            uint(rand.Intn(500000-85000+1) + 85000),
				Mileage:          uint(rand.Intn(350000-2500+1) + 2500),
				ManufacturedYear: uint(rand.Intn(2022-1980+1) + 1980),
				EnterpriseID:     uint(e),
				CarModel: CarModel{
					Brand:            carBrands[randomBrand],
					CarType:          carTypes[randomCarType],
					FuelTankCapacity: uint(rand.Intn(250-15+1) + 15),
					Seats:            uint(rand.Intn(10-1+1) + 1),
					EngineVolume:     2.4 + rand.Float64()*15.5 - 2.4,
				},
			}

			vehicleJSON, _ := json.Marshal(vehicle)
			resp, err := http.Post(PostVehicleUrl, "application/json", bytes.NewBuffer(vehicleJSON))
			if err != nil {
				log.Println(err)
			}
			fmt.Sprintln(resp)

			if car%8 == 0 {
				type Response struct {
					ID      uint
					Message string
				}

				var r Response
				b, _ := io.ReadAll(resp.Body)
				json.Unmarshal([]byte(b), &r)
				vehicleID := r.ID

				randomName := rand.Int() % len(Names)
				randomLastName := rand.Int() % len(LastNames)
				driver := Driver{
					FirstName: Names[randomName],
					LastName:  LastNames[randomLastName],
					Salary:    float64(rand.Intn(50000-35000+1) + 35000),
					VehicleID: vehicleID,
					IsActive:  true,
				}

				jsonDriver, _ := json.Marshal(driver)
				resp2, err := http.Post(PostDrivereUrl, "application/json", bytes.NewBuffer(jsonDriver))
				if err != nil {
					log.Println(err)
				}
				fmt.Println(resp2)
			}
			resp.Body.Close()

		}
	}

}
