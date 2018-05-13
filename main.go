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
	// Use mcu pin 10, corresponds to physical pin 19 on the pi
	bombPin = rpio.Pin(17)
	trigPin = rpio.Pin(23)
	echoPin = rpio.Pin(24)
)

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

	trigPin.Low()

	time.Sleep(time.Millisecond * 200)

	trigPin.High()

	time.Sleep(time.Microsecond * 20)

	trigPin.Low()

	var wg sync.WaitGroup
	var start time.Time
	var d time.Duration

	wg.Add(1)
	go func() {
		for echoPin.Read() == rpio.Low {
		}
		start = time.Now()
		for echoPin.Read() == rpio.High {
		}
		d = time.Since(start)
		wg.Done()
	}()

	wg.Wait()

	res := d.Seconds() * 17000

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
	json.NewEncoder(w).Encode(Level{int(getMeasurement())})
	// json.NewEncoder(w).Encode(Level{34})
}

func GetBombStatus(w http.ResponseWriter, r *http.Request) {
	conn := Conn{bomba_state}
	json.NewEncoder(w).Encode(conn)
}

func TurnBombOn(w http.ResponseWriter, r *http.Request) {
	bombPin.Output()
	bomba_state = true
	json.NewEncoder(w).Encode(Conn{true})
}

func TurnBombOff(w http.ResponseWriter, r *http.Request) {
	bombPin.Input()
	bomba_state = false
	json.NewEncoder(w).Encode(Conn{false})
}

func login(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Header", r.Header)
	fmt.Println("Host", r.Host)
	fmt.Println("PostForm", r.PostForm)
	fmt.Println("RemoteAddr", r.RemoteAddr)

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

// main function to boot up everything
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/level", getlevel)
	mux.HandleFunc("/b_on", TurnBombOn)
	mux.HandleFunc("/b_off", TurnBombOff)
	mux.HandleFunc("/b_status", GetBombStatus)
	// mux.Headers("", "")
	// mux.Headers("Access-Control-Allow-Origin", "*")

	// Set pin to output mode

	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	trigPin.Output()

	echoPin.Input()

	bombPin.Output()
	bombPin.High()

	// Open and map memory to access gpio, check for errors

	// var wg sync.WaitGroup

	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for {

	// 		d := getDistance()
	// 		level.mu.Lock()
	// 		level.Level = convertToPCT(d)
	// 		level.mu.Unlock()
	// 		fmt.Println(time.Now().Format(time.RFC3339), " level Readed ", level)
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	fmt.Println("Server ready!!")
	// err := http.ListenAndServeTLS("0.0.0.0:443", "certificate.crt", "private.key", handlers.CORS()(router))
	// err := http.ListenAndServe("0.0.0.0:8082", handlers.CORS()(router))
	// if err != nil {
	// 	fmt.Println("Error", err)
	// }

	// mux.HandleFunc("/", hello)
	// mux.HandleFunc("/login", login)
	// mux.HandleFunc("/account", account)

	// "http://localhost:8100", "https://tinaco2.tk", "http://localhost:4200", "http://192.168.1.65", "http://192.168.0.14"

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Authorization", "Cache-Control", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})

	// Insert the middleware
	handler := c.Handler(mux)

	// handler := cors.Default().Handler(mux)

	log.Println("Listening for connections on port: ", 8082)
	log.Fatal(http.ListenAndServeTLS("0.0.0.0:443", "certificate.crt", "private.key", handler))

	// wg.Wait()

	// rpio.Close()

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
