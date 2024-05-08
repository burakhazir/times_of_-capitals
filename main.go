package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// City struct'ı
type City struct {
	Name        string `json:"name"`
	Country     string `json:"country"`
	TimeZone    int    `json:"time_zone"`
	Continental string `json:"continental"`
}

// CityClock struct'ı, her bir şehrin güncel saat dilimini depolar
type CityClock struct {
	CityName    string
	CurrentTime time.Time
}

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/cityClocks", getCityClocks)
	fmt.Println("8080 portuna bağlanıldı...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Serve HTML dosyasını
func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Şehir saat bilgilerini döndür
func getCityClocks(w http.ResponseWriter, r *http.Request) {
	// JSON dosyasından tüm şehirleri al
	cities := getAllCities()

	// Tüm şehirlerin saat bilgisini almak için döngü kullan
	var cityClocks []CityClock
	for _, city := range cities {
		gmtTime := clock()                        // GMT saatini al
		cityTime := getCityTime(city, gmtTime)    // Şehir saati bilgisini al
		cityClocks = append(cityClocks, cityTime) // Şehir saatini depola
	}

	// JSON formatında şehir saatlerini gönder
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cityClocks)
	fmt.Println("working...")
}

// Tüm şehirleri JSON dosyasından al
func getAllCities() []City {
	jsonData, err := ioutil.ReadFile("file.json")
	if err != nil {
		log.Fatalf("JSON dosyasını okurken bir hata oluştu: %v", err)
	}

	// JSON verilerini depolamak için bir yapı oluştur
	var data struct {
		Cities []City `json:"Index"`
	}

	// JSON verilerini çözümle
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatalf("JSON verilerini çözümleme işlemi sırasında bir hata oluştu: %v", err)
	}

	return data.Cities
}

// GMT saatini al
func clock() time.Time {
	return time.Now().UTC()
}

// Verilen şehrin saatini belirlemek için gmtTime parametresini kullan
func getCityTime(city City, gmtTime time.Time) CityClock {
	// Şehir saat dilimini belirle
	cityTime := gmtTime.Add(time.Duration(city.TimeZone) * time.Minute)
	return CityClock{CityName: city.Name, CurrentTime: cityTime}
}
