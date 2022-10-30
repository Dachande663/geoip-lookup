package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type LookupResponseSuccess struct {
	Status string  `json:"status"`
	Info   *IpInfo `json:"info"`
}

type LookupResponseError struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type StatusResponse struct {
	Status            string `json:"status"`
	StatUptime        uint64 `json:"stat_uptime"`
	StatTotalRequests uint64 `json:"stat_total_requests"`
	StatTotalLookups  uint64 `json:"stat_total_lookups"`
	StatTotalFound    uint64 `json:"stat_total_found"`
	StatTotalMissing  uint64 `json:"stat_total_missing"`
	DbLastUpdated     string `json:"db_last_updated"`
	DbRecordCount     uint   `json:"db_record_count"`
}

func indexRoute(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"status":"error","error":"Route not found"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"name":"geoip","version":"1.0.0"}`)
}

func statusRoute(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(StatusResponse{
		Status:            "up",
		StatUptime:        uint64(time.Now().Sub(startTime).Seconds()),
		StatTotalRequests: totalRequests,
		StatTotalLookups:  totalLookups,
		StatTotalFound:    totalLookupsFound,
		StatTotalMissing:  totalLookups - totalLookupsFound,
		DbLastUpdated:     time.Unix(int64(ipDatabase.Metadata.BuildEpoch), 0).Format(time.RFC3339),
		DbRecordCount:     ipDatabase.Metadata.NodeCount,
	})
}

func lookupRoute(w http.ResponseWriter, req *http.Request) {
	atomic.AddUint64(&totalLookups, 1)
	inputIp := strings.TrimPrefix(req.URL.Path, "/lookup/")
	ipInfo, err := lookupIp(inputIp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(LookupResponseError{Status: "error", Error: err.Error()})
		return
	}
	atomic.AddUint64(&totalLookupsFound, 1)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LookupResponseSuccess{Status: "success", Info: ipInfo})
}

func authorize(r *http.Request) bool {
	if config["HTTP_AUTH"] == "" {
		return true
	}
	expectedHeader := []byte("Bearer " + config["HTTP_AUTH"])
	providedHeader := []byte(r.Header.Get("Authorization"))
	return subtle.ConstantTimeCompare(expectedHeader, providedHeader) == 1
}

func appMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&totalRequests, 1)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-App", "geoip")
		if !authorize(r) {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(LookupResponseError{Status: "error", Error: "Unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
