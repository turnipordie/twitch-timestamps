package tsprocessor

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type Comment struct {
	CreatedAt string `json:"created_at"`
}

type Video struct {
	CreatedAt string `json:"created_at"`
}

type Data struct {
	Comments []Comment `json:"comments"`
	Video    Video
}

func Process(jsonData []byte, avgMult float64) []string {
	var data Data
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return []string{}
	}

	vodStartMs, _ := time.Parse(time.RFC3339, data.Video.CreatedAt)

	timestamps := make([]int, 0, len(data.Comments))
	for _, comment := range data.Comments {
		t, err := time.Parse(time.RFC3339, comment.CreatedAt)
		if err != nil {
			println("invalid ts " + comment.CreatedAt)
			continue
		}

		timestamp := t.Sub(vodStartMs) / time.Second
		timestamps = append(timestamps, int(timestamp))
	}

	_, frequencyPerMin := FrequencyMap(timestamps)
	return FindPeaks(frequencyPerMin, int(float64(getAverage(frequencyPerMin))*avgMult))
}

func getAverage(m map[int]int) int {
	avg := 0
	for _, v := range m {
		avg += v
	}
	avg /= len(m)
	return avg
}

func FrequencyMap(tsList []int) (map[int]int, map[int]int) {
	frequencyPerSec := make(map[int]int)
	frequencyPerMin := make(map[int]int)
	for _, second := range tsList {
		frequencyPerSec[second]++
		frequencyPerMin[(int)(second/60)]++
	}
	return frequencyPerSec, frequencyPerMin
}

func FindPeaks(frequencyMap map[int]int, baseline int) []string {
	keys := make([]int, 0)
	for k := range frequencyMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	var timestamps []string
	for _, k := range keys {
		msgCount := frequencyMap[k]
		if msgCount > baseline {
			hour := k / 60
			minute := k % 60
			timestamps = append(timestamps, fmt.Sprintf("%02d:%02d", hour, minute))
		}
	}

	return timestamps
}
