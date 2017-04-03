package http

import (
	"encoding/json"
	gocache "github.com/pmylund/go-cache"
	"github.com/unrolled/render"
	"net/http"
	"time"
)

var (
	cache = gocache.New(5*time.Minute, 30*time.Second)
)

type Data struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func getJson(this interface{}, url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(this)
}

func loadData(this *Data, url string) (interface{}, bool) {
	if cached, found := cache.Get("data"); found {
		cache.Set("data_tmp", cached, gocache.NoExpiration)
		return cached, found
	}

	if cached_tmp, found_tmp := cache.Get("data_tmp"); found_tmp {
		go func() {
			getJson(this, url)
			cache.Set("data", this, gocache.DefaultExpiration)
		}()
		return cached_tmp, found_tmp
	}

	getJson(this, url)
	cache.Set("data", this, gocache.DefaultExpiration)

	return this, false
}

func setDefaultHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Vary", "Accept-Encoding")
}

func setCacheHeader(w http.ResponseWriter, found bool) {
	v := "MISS"
	if found {
		v = "HIT"
	}
	w.Header().Set("X-Cache", v)
}

func main() {
	render := render.New()
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setDefaultHeaders(w)

		data := new(Data)
		url := ""

		res, found := loadData(data, url)
		setCacheHeader(w, found)

		render.JSON(w, http.StatusOK, res)
	})

	http.ListenAndServe(":8000", mux)
}
