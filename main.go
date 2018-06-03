package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/rs/cors"
	rpio "github.com/stianeikeland/go-rpio"
)

const (
	SECRET = "42isTheAnswer"
)

type JWTData struct {
	// Standard claims are the standard jwt claims from the IETF standard
	// https://tools.ietf.org/html/rfc7519
	jwt.StandardClaims
	CustomClaims map[string]string `json:"custom,omitempty"`
}

var (
	bombPin = rpio.Pin(15)
)

type pins struct {
	trigPin rpio.Pin
	echoPin rpio.Pin
	mux     sync.Mutex
}

var bomba_state = false

var cont int64

type Conn struct {
	Status bool `json:"status,omitempty"`
}
type Level struct {
	Level int `json:"level,omitempty"`
}

var level Level

func Swap(arrayzor []float64, i, j int) {
	tmp := arrayzor[j]
	arrayzor[j] = arrayzor[i]
	arrayzor[i] = tmp
}

func bubbleSort(arrayzor []float64) {
	swapped := true
	for swapped {
		swapped = false
		for i := 0; i < len(arrayzor)-1; i++ {
			if arrayzor[i+1] < arrayzor[i] {
				Swap(arrayzor, i, i+1)
				swapped = true
			}
		}
	}
}

func getMeasurement() float64 {

	pinsObj.trigPin.Low()

	time.Sleep(time.Millisecond * 200)

	pinsObj.trigPin.High()

	time.Sleep(time.Microsecond * 20)

	pinsObj.trigPin.Low()

	var wg sync.WaitGroup
	var start time.Time
	var d time.Duration

	wg.Add(1)
	go func() {
		for pinsObj.echoPin.Read() == rpio.Low {
		}
		start = time.Now()
		for pinsObj.echoPin.Read() == rpio.High {
		}
		d = time.Since(start)
		wg.Done()
	}()

	wg.Wait()

	res := d.Seconds() * 17000

	fmt.Println("res", res)

	return res
}

func getDistance() float64 {
	ds := make([]float64, 0)
	var avg, sum float64
	for i := 0; i < 10; i++ {
		ds = append(ds, getMeasurement())
		time.Sleep(time.Millisecond * 100)
	}

	bubbleSort(ds)

	ds = ds[2:7]

	for _, v := range ds {
		sum += v
	}

	avg = sum / 5

	return avg
}

func convertToPCT(d float64) int {
	low := 20.0
	hight := 110.0

	if d < low {
		return 100
	} else if d > hight {
		return 0
	}

	rdistance := d - low
	total := hight - low

	return int(100 - ((rdistance * 100) / total))
}

func getlevel(w http.ResponseWriter, r *http.Request) {
	validateRequest(w, r)
	// if err := rpio.Open(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	fmt.Println(1)
	pinsObj.mux.Lock()
	d := getDistance()
	pinsObj.mux.Unlock()
	fmt.Println("d", d)
	// rpio.Close()
	fmt.Println(3)
	pctLevel := convertToPCT(d)
	fmt.Println("pctLevel", pctLevel)
	json.NewEncoder(w).Encode(Level{pctLevel})
	fmt.Println(5)
}

func GetBombStatus(w http.ResponseWriter, r *http.Request) {
	validateRequest(w, r)
	conn := Conn{bomba_state}
	json.NewEncoder(w).Encode(conn)
}

func TurnBombOn(w http.ResponseWriter, r *http.Request) {
	validateRequest(w, r)
	// if err := rpio.Open(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	bombPin.Output()
	// rpio.Close()
	bomba_state = true
	json.NewEncoder(w).Encode(Conn{true})
}

func TurnBombOff(w http.ResponseWriter, r *http.Request) {
	validateRequest(w, r)
	// if err := rpio.Open(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	bombPin.Input()
	// rpio.Close()
	bomba_state = false
	json.NewEncoder(w).Encode(Conn{false})
}

func login(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, "Login failed!", http.StatusUnauthorized)
		// http.Error(w, err.Error(), 1)
	}

	var userData map[string]string
	json.Unmarshal(body, &userData)

	// Demo - in real case scenario you'd check this against your database
	if userData["email"] == "admin@gmail.com" && userData["password"] == "admin123" {
		claims := JWTData{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
			},

			CustomClaims: map[string]string{
				"userid": "u1",
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(SECRET))
		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		json, err := json.Marshal(struct {
			Token string `json:"token"`
		}{
			tokenString,
		})

		if err != nil {
			log.Println(err)
			http.Error(w, "Login failed!", http.StatusUnauthorized)
		}

		w.Write(json)
	} else {
		http.Error(w, "Login failed!", http.StatusUnauthorized)
	}
}

var pinsObj = pins{
	trigPin: rpio.Pin(23),
	echoPin: rpio.Pin(24),
}

// main function to boot up everything
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/level", getlevel)
	mux.HandleFunc("/b_on", TurnBombOn)
	mux.HandleFunc("/b_off", TurnBombOff)
	mux.HandleFunc("/b_status", GetBombStatus)

	// Set pin to output mode
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	pinsObj.trigPin.Output()

	pinsObj.echoPin.Input()

	bombPin.Input()

	defer rpio.Close()

	fmt.Println("Server ready!!")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	handler := c.Handler(mux)

	go func() {
		for {
			fmt.Printf("%s", "Checkinkg bomb status ... ")
			if bomba_state {
				fmt.Printf("%s", "ON\n")
				// if err := rpio.Open(); err != nil {
				// 	fmt.Println(err)
				// 	os.Exit(1)
				// }
				pinsObj.mux.Lock()
				d := getMeasurement()
				pinsObj.mux.Unlock()
				pctLevel := convertToPCT(d)
				if pctLevel == 100 {
					bombPin.Input()
					bomba_state = false
				}
				// rpio.Close()
			} else {
				fmt.Printf("%s", "OFF\n")
			}
			time.Sleep(time.Second * 5)
		}
	}()

	log.Println("Listening for connections on port: ", 443)
	log.Fatal(http.ListenAndServeTLS("0.0.0.0:443", "certificate.crt", "private.key", handler))

}

func validateRequest(w http.ResponseWriter, r *http.Request) *JWTData {

	fmt.Println("host", r.Host)

	authToken := r.Header.Get("Authorization")
	authArr := strings.Split(authToken, " ")

	if len(authArr) != 2 {
		log.Println("Authentication header is invalid: " + authToken)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	jwtToken := authArr[1]

	claims, err := jwt.ParseWithClaims(jwtToken, &JWTData{}, func(token *jwt.Token) (interface{}, error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}
		return []byte(SECRET), nil
	})

	if err != nil {
		log.Println(err)
		http.Error(w, "Request failed!", http.StatusUnauthorized)
	}

	return claims.Claims.(*JWTData)
}
