package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

type ForecastResp struct {
	Timezone string `json:"timezone"`
	Daily    Daily  `json:"daily"`
}

type Daily struct {
	Time        []string  `json:"time"`
	TempMax     []float64 `json:"temperature_2m_max"`
	TempMin     []float64 `json:"temperature_2m_min"`
	WeatherCode []int     `json:"weather_code"`
}

func weatherFetch(ctx context.Context, lat float64, lon float64, days int) (ForecastResp, error) {
	base := "https://api.open-meteo.com/v1/forecast"

	q := url.Values{}

	q.Set("latitude", strconv.FormatFloat(lat, 'f', 4, 64))
	q.Set("longitude", strconv.FormatFloat(lon, 'f', 4, 64))
	q.Set("daily", "temperature_2m_max,temperature_2m_min,weather_code")
	q.Set("forecast_days", strconv.Itoa(days))
	q.Set("timezone", "Asia/Seoul")

	finalURL := base + "?" + q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, finalURL, nil)
	if err != nil {
		return ForecastResp{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "my-weather-wrapper/1.0")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return ForecastResp{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return ForecastResp{}, fmt.Errorf("open-meteo error: status=%d body=%s", resp.StatusCode, string(b))
	}

	var out ForecastResp
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return ForecastResp{}, err
	}
	return out, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	daysStr := r.URL.Query().Get("days")

	if latStr == "" || lonStr == "" {
		http.Error(w, "missing lat or lon", http.StatusBadRequest)
		return
	}
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "invalid lat", http.StatusBadRequest)
		return
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "invalid lon", http.StatusBadRequest)
		return
	}
	days := 7
	if daysStr != "" {
		n, err := strconv.Atoi(daysStr)
		if err != nil || n <= 0 || n > 16 {
			http.Error(w, "invalid days (1~16)", http.StatusBadRequest)
			return
		}
		days = n
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	latKey := strconv.FormatFloat(lat, 'f', 4, 64)
	lonKey := strconv.FormatFloat(lon, 'f', 4, 64)
	cacheKey := fmt.Sprintf("%s:%s:%d", latKey, lonKey, days)

	cached, err := rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("X-Cache", "HIT")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cached))
		return
	}

	if err != nil && err != redis.Nil {
		log.Printf("redis get error: %v", err)
	}

	data, err := weatherFetch(ctx, lat, lon, days)
	if err != nil {
		http.Error(w, "failed to fetch weather: "+err.Error(), http.StatusBadGateway)
		return
	}
	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "failed to make json", http.StatusInternalServerError)
		return
	}
	ttl := 30 * time.Minute
	if err := rdb.Set(ctx, cacheKey, b, ttl).Err(); err != nil {
		log.Printf("redis set error: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("redis connect failed: ", err)
	}

	http.HandleFunc("/weather", weatherHandler)
	fmt.Println("listening on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
