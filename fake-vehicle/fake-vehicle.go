package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const APIKEY = "5b3ce3597851110001cf6248f5b8fac31f6b4ec7b7ebc6c9c731ede9"
const MANAGER_LOGIN = "admin2"
const MANAGER_PASSWORD = "qwerty"

type GeoPoint struct {
	TrackTime string  `json:"track_time"`
	GeoX      float64 `json:"geo_x"`
	GeoY      float64 `json:"geo_y"`
	// VehicleID uint    `json:"vehicle_id"`
}

type OpenStreetMapRoute struct {
	Type     string `json:"type"`
	Features []struct {
		Bbox       []float64 `json:"bbox"`
		Type       string    `json:"type"`
		Properties struct {
			Segments []struct {
				Distance float64 `json:"distance"`
				Duration float64 `json:"duration"`
				Steps    []struct {
					Distance    float64 `json:"distance"`
					Duration    float64 `json:"duration"`
					Type        int     `json:"type"`
					Instruction string  `json:"instruction"`
					Name        string  `json:"name"`
					WayPoints   []int   `json:"way_points"`
					ExitNumber  int     `json:"exit_number,omitempty"`
				} `json:"steps"`
			} `json:"segments"`
			Summary struct {
				Distance float64 `json:"distance"`
				Duration float64 `json:"duration"`
			} `json:"summary"`
			WayPoints []int `json:"way_points"`
		} `json:"properties"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
			Type        string      `json:"type"`
		} `json:"geometry"`
	} `json:"features"`
	Bbox     []float64 `json:"bbox"`
	Metadata struct {
		Attribution string `json:"attribution"`
		Service     string `json:"service"`
		Timestamp   int64  `json:"timestamp"`
		Query       struct {
			Coordinates [][]float64 `json:"coordinates"`
			Profile     string      `json:"profile"`
			Format      string      `json:"format"`
		} `json:"query"`
		Engine struct {
			Version   string    `json:"version"`
			BuildDate time.Time `json:"build_date"`
			GraphDate time.Time `json:"graph_date"`
		} `json:"engine"`
	} `json:"metadata"`
}

type ReverseGeocode struct {
	Geocoding struct {
		Version     string `json:"version"`
		Attribution string `json:"attribution"`
		Query       struct {
			Sources              []string `json:"sources"`
			Size                 int      `json:"size"`
			Private              bool     `json:"private"`
			PointLat             float64  `json:"point.lat"`
			PointLon             float64  `json:"point.lon"`
			BoundaryCircleRadius int      `json:"boundary.circle.radius"`
			BoundaryCircleLat    float64  `json:"boundary.circle.lat"`
			BoundaryCircleLon    float64  `json:"boundary.circle.lon"`
			Lang                 struct {
				Name      string `json:"name"`
				Iso6391   string `json:"iso6391"`
				Iso6393   string `json:"iso6393"`
				Via       string `json:"via"`
				Defaulted bool   `json:"defaulted"`
			} `json:"lang"`
			QuerySize int `json:"querySize"`
		} `json:"query"`
		Engine struct {
			Name    string `json:"name"`
			Author  string `json:"author"`
			Version string `json:"version"`
		} `json:"engine"`
		Timestamp int64 `json:"timestamp"`
	} `json:"geocoding"`
	Type     string `json:"type"`
	Features []struct {
		Type     string `json:"type"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			ID            string  `json:"id"`
			Gid           string  `json:"gid"`
			Layer         string  `json:"layer"`
			Source        string  `json:"source"`
			SourceID      string  `json:"source_id"`
			Name          string  `json:"name"`
			Confidence    float64 `json:"confidence"`
			Distance      float64 `json:"distance"`
			Accuracy      string  `json:"accuracy"`
			Country       string  `json:"country"`
			CountryGid    string  `json:"country_gid"`
			CountryA      string  `json:"country_a"`
			Region        string  `json:"region"`
			RegionGid     string  `json:"region_gid"`
			RegionA       string  `json:"region_a"`
			County        string  `json:"county"`
			CountyGid     string  `json:"county_gid"`
			CountyA       string  `json:"county_a"`
			Localadmin    string  `json:"localadmin"`
			LocaladminGid string  `json:"localadmin_gid"`
			Locality      string  `json:"locality"`
			LocalityGid   string  `json:"locality_gid"`
			Continent     string  `json:"continent"`
			ContinentGid  string  `json:"continent_gid"`
			Label         string  `json:"label"`
			Addendum      struct {
				Osm struct {
					Brand string `json:"brand"`
				} `json:"osm"`
			} `json:"addendum"`
		} `json:"properties"`
		Bbox []float64 `json:"bbox"`
	} `json:"features"`
	Bbox []float64 `json:"bbox"`
}

// radius passed in meters
func generateRandomPoint(center GeoPoint, radius float64) GeoPoint {
	var x0 = center.GeoX
	var y0 = center.GeoY

	// Convert Radius from meters to degrees.
	var rd = radius / 111300

	rand.Seed(time.Now().UnixNano())
	var u = 0.00000001 + rand.Float64()*0.999999999
	var v = 0.00000001 + rand.Float64()*0.999999999

	var w = rd * math.Sqrt(u)
	var t = 2 * math.Pi * v
	var x = w * math.Cos(t)
	var y = w * math.Sin(t)

	// Result point.
	return GeoPoint{GeoY: y + y0, GeoX: x + x0}
}

func getORSRoute(center, endpoint GeoPoint) (OpenStreetMapRoute, error) {

	var ORSRoute OpenStreetMapRoute
	centerCoordinatesString := fmt.Sprintf("%f", center.GeoX) + "," + fmt.Sprintf("%f", center.GeoY)
	endPointCoordinatesString := fmt.Sprintf("%f", endpoint.GeoX) + "," + fmt.Sprintf("%f", endpoint.GeoY)
	openRouteUrl := "https://api.openrouteservice.org/v2/directions/driving-car?api_key=" + APIKEY + "&start=" + centerCoordinatesString + "+&end=" + endPointCoordinatesString

	resp, err := http.Get(openRouteUrl)
	if err != nil {
		log.Println(err)
		return ORSRoute, err
	}
	if resp.StatusCode == 404 {
		log.Println("404 from OpenRoute, could not make route. Trying with another endpoint")
		return ORSRoute, fmt.Errorf("Cant get route")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	errUnmarshal := json.Unmarshal(body, &ORSRoute)
	if errUnmarshal != nil {
		log.Println("Cannot unmarshal response body")
		return ORSRoute, errUnmarshal
	}
	return ORSRoute, nil
}

func getReverseGeolocationAddress(gp GeoPoint) (string, error) {
	var revGeocode ReverseGeocode
	longitudeStr := fmt.Sprintf("%f", gp.GeoX)
	latitudeStr := fmt.Sprintf("%f", gp.GeoY)

	openRouteUrl := "https://api.openrouteservice.org/geocode/reverse?api_key=" + APIKEY + "+&point.lon=" + longitudeStr + "&point.lat=" + latitudeStr + "&boundary.circle.radius=1&size=1&sources=openstreetmap"

	resp, err := http.Get(openRouteUrl)
	if err != nil {
		log.Println(err)
		return "address not found", err
	}
	if resp.StatusCode == 404 {
		log.Println("404 from OpenRoute, could not reverse geolocation.")
		return "address not found", fmt.Errorf("Cant get address")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	errUnmarshal := json.Unmarshal(body, &revGeocode)
	if errUnmarshal != nil {
		log.Println("Cannot unmarshal response body")
		return "address not found", errUnmarshal
	}

	addressLabel := revGeocode.Features[0].Properties.Label
	return addressLabel, nil
}

func main() {

	vehicleID := flag.Int("vehicleID", -1, "vehicle ID")
	vehicleMS := flag.Float64("vehicleMaxSpeed", 80, "vehicle max speed")
	vehicleACC := flag.Float64("vehicleMaxAcc", 2.0, "vehicle max acceleration")
	vehicleStartX := flag.Float64("vehicleStartX", -83.16511660044354, "vehicle start longitude")
	vehicleStartY := flag.Float64("vehicleStartY", 42.53186889829862, "vehicle start latitude")
	withinRadius := flag.Float64("withinRadius", 15000, "radius where to generate random ride endpoint")
	flag.Parse()

	center := GeoPoint{GeoX: *vehicleStartX, GeoY: *vehicleStartY}
	if *vehicleID == -1 {
		fmt.Println("You should specify vehicle ID, it is mandatory")
		os.Exit(1)
	}

	fmt.Println("INITIATING FAKE VEHICLE AGENT USING THESE PARAMETERS:\n-----------------------------------------------------------")
	fmt.Println("VEHICLE_ID: ", *vehicleID)
	if *vehicleMS == 80 {
		fmt.Println("DEFAULT_MAX_SPEED: ", *vehicleMS)
	} else {
		fmt.Println("MAX_SPEED: ", *vehicleMS)
	}

	if *vehicleACC == 2.0 {
		fmt.Println("DEFAULT_MAX_ACCELERATION: ", *vehicleACC)
	} else {
		fmt.Println("MAX_ACCELERATION: ", *vehicleACC)
	}

	addressStart, err := getReverseGeolocationAddress(center)
	if err != nil {
		log.Println(err)
	} else if addressStart == "Taco Bell, Clawson, MI, USA" {
		fmt.Println("STARTING_WITH_DEFAULT ADDRESS: ", addressStart)
	} else {
		fmt.Println("STARTING_ADDRESS: ", addressStart)
	}

	fmt.Print("\nSEARCHING FOR RANDOM FINISH ADDRESS IN " + fmt.Sprintf("%f", *withinRadius) + " meters radius...\n\n")

	var orsRoute OpenStreetMapRoute
	var errGetRoute error
	for i := 0; i < 5; i++ {
		endpoint := generateRandomPoint(center, *withinRadius)
		orsRoute, errGetRoute = getORSRoute(center, endpoint)
		if errGetRoute == nil {
			fmt.Print("ROUTE FOUND SUCCESSFULLY.\n")
			break
		} else {
			fmt.Println("FAILED IN BUILDING ROUTE. TRYING ANOTHER FINSH POINT...")
			continue
		}
	}

	trackDistance := orsRoute.Features[0].Properties.Summary.Distance
	geopointsCount := len(orsRoute.Features[0].Geometry.Coordinates)
	avgGeopointsDistance := trackDistance / float64(geopointsCount)
	magicAvgSpeed := (*vehicleMS - (35**vehicleMS)/100) * (*vehicleACC / 2.43)
	avgPointPassTime := avgGeopointsDistance / magicAvgSpeed
	trafficJamIndex := 8.3
	gpsWriterNapTime := avgPointPassTime * trafficJamIndex

	finishSlice := orsRoute.Features[0].Geometry.Coordinates[geopointsCount-1]
	finishGeoPoint := GeoPoint{GeoX: finishSlice[0], GeoY: finishSlice[1]}
	addressFinish, err := getReverseGeolocationAddress(finishGeoPoint)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("FINISH ADDRESS IS: " + addressFinish + "\n")
	}

	fmt.Println("DEBUG SECTION:\n-----------------------------------------------------------")
	fmt.Println("TRACK_DISTANCE", trackDistance)
	fmt.Println("GEOPOINTS_COUNT", geopointsCount)
	fmt.Println("AVG_GEOPOINTS_DISTANCE", avgGeopointsDistance)
	fmt.Println("AVG_SPEED", magicAvgSpeed)
	fmt.Println("NAP_TIME", gpsWriterNapTime)
	fmt.Println("-----------------------------------------------------------")
	fmt.Println("STARTING MOVEMENT:")

	carParkApiUrl := "http://" + MANAGER_LOGIN + ":" + MANAGER_PASSWORD + "@localhost:8888/api/manager/1/vehicle/" + strconv.Itoa(25169) + "/checkpoint"

	startPoint := GeoPoint{
		GeoX:      center.GeoX,
		GeoY:      center.GeoY,
		TrackTime: TimeStampToUTCPFC3339String(time.Now()),
	}
	postGeoPoint(startPoint, carParkApiUrl)

	for _, coordinate := range orsRoute.Features[0].Geometry.Coordinates {

		geoPoint := GeoPoint{
			GeoX: coordinate[0],
			GeoY: coordinate[1],
			// VehicleID: uint(*vehicleID),
			TrackTime: TimeStampToUTCPFC3339String(time.Now()),
		}
		postGeoPoint(geoPoint, carParkApiUrl)
		time.Sleep(time.Second * time.Duration(gpsWriterNapTime))
	}

}

func TimeStampToUTCPFC3339String(timeStamp time.Time) string {

	RFC3339local := "2006-01-02T15:04:05Z"
	utcLoc, _ := time.LoadLocation("UTC")
	timeStamp = timeStamp.In(utcLoc)
	timeStampStr := timeStamp.Format(RFC3339local)

	return timeStampStr

}

func postGeoPoint(gp GeoPoint, url string) {
	geoPointJson, _ := json.Marshal(gp)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(geoPointJson))
	if err != nil {
		log.Println(err)
	}
	fmt.Sprintln(resp)
	resp.Body.Close()
}
