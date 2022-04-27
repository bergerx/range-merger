package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

type rangeObject struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type rangeList []rangeObject

func (rs rangeList) Len() int {
	return len(rs)
}

func (rs rangeList) Less(x, y int) bool {
	int1 := rs[x]
	int2 := rs[y]
	if int1.Start < int2.Start {
		return true
	}
	if int1.Start == int2.Start && int1.End < int2.End {
		return true
	}
	return false
}

func (rs rangeList) Swap(x, y int) {
	r := rs
	r[x], r[y] = r[y], r[x]
}

func (rs rangeList) MergeOverlappingRanges() rangeList {
	fmt.Fprintln(os.Stderr, "got a request")
	sort.Sort(rs)
	var mergedRangeList rangeList
	mergedRangeList = append(mergedRangeList, rs[0])
	for i, r := range rs {
		if i == 0 {
			continue
		}
		lastInMerged := mergedRangeList[len(mergedRangeList)-1]
		if lastInMerged.End < r.Start {
			mergedRangeList = append(mergedRangeList, r)
			continue
		}
		lastInMerged.End = r.End
		mergedRangeList[len(mergedRangeList)-1] = lastInMerged
	}
	return mergedRangeList
}

func overlaps(w http.ResponseWriter, req *http.Request) {
	var inputRanges rangeList
	all, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	err = json.Unmarshal(all, &inputRanges)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintln(w, err)
		return
	}
	merged := inputRanges.MergeOverlappingRanges()
	out, err := json.Marshal(merged)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintf(w, "%s\n", string(out))
}

func main() {
	http.HandleFunc("/api/overlaps", overlaps)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Http Listen error: %v", err)
	}
}
