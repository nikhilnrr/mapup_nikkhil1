package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"
)

type SortRequest struct {
	ToSort [][]int `json:"to_sort"`
}

type SortResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNS       int64   `json:"time_ns"`
}

func sortHandler(w http.ResponseWriter, r *http.Request, sortFunc func([][]int) [][]int) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request SortRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	startTime := time.Now()
	sortedArrays := sortFunc(request.ToSort)
	duration := time.Since(startTime).Nanoseconds()

	response := SortResponse{
		SortedArrays: sortedArrays,
		TimeNS:       duration,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sortSequential(toSort [][]int) [][]int {
	sortedArrays := make([][]int, len(toSort))
	for i, arr := range toSort {
		sorted := make([]int, len(arr))
		copy(sorted, arr)
		sort.Ints(sorted)
		sortedArrays[i] = sorted
	}
	return sortedArrays
}

func sortConcurrent(toSort [][]int) [][]int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	sortedArrays := make([][]int, len(toSort))

	for i, arr := range toSort {
		wg.Add(1)
		go func(i int, arr []int) {
			defer wg.Done()

			sorted := make([]int, len(arr))
			copy(sorted, arr)
			sort.Ints(sorted)

			mu.Lock()
			sortedArrays[i] = sorted
			mu.Unlock()
		}(i, arr)
	}

	wg.Wait()
	return sortedArrays
}

func main() {
	http.HandleFunc("/process-single", func(w http.ResponseWriter, r *http.Request) {
		sortHandler(w, r, sortSequential)
	})

	http.HandleFunc("/process-concurrent", func(w http.ResponseWriter, r *http.Request) {
		sortHandler(w, r, sortConcurrent)
	})

	port := 8000
	addr := ":" + strconv.Itoa(port)
	println("Server listening on", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
