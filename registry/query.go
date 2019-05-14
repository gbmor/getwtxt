package registry

import (
	"sort"
	"strings"
	"time"
)

// QueryUser checks the user index for nicknames that contain the
// nickname provided as an argument. Entries are returned sorted
// by the date they were added to the index.
func (index UserIndex) QueryUser(name string) []string {
	var timekey = map[time.Time]string{}
	var keys TimeSlice
	var users []string
	imutex.RLock()
	for k, v := range index {
		if strings.Contains(v.Nick, name) {
			timekey[v.Date] = v.Nick + "\t" + k + "\t" + string(v.APIdate)
			keys = append(keys, v.Date)
		}
	}
	imutex.RUnlock()
	sort.Sort(keys)
	for _, e := range keys {
		users = append(users, timekey[e])
	}

	return users
}

// QueryTag returns all the known statuses that
// contain the provided tag.
func (index UserIndex) QueryTag(tag string) []string {
	var statusmap TimeMapSlice
	i := 0
	imutex.RLock()
	for _, v := range index {
		statusmap[i] = v.FindTag(tag)
		i++
	}
	imutex.RUnlock()

	return statusmap.SortByTime()
}

// FindTag takes a user's tweets and looks for a given tag.
// Returns the tweets with the tag as a []string.
func (userdata *Data) FindTag(tag string) TimeMap {
	var statuses TimeMap
	for k, e := range userdata.Status {
		parts := strings.Split(e, "\t")
		statusslice := strings.Split(parts[3], " ")
		for _, v := range statusslice {
			if v[1:] == tag {
				statuses[k] = e
				break
			}
		}
	}

	return statuses
}

// SortByTime returns a string slice of the query results
// sorted by time.Time
func (tm TimeMapSlice) SortByTime() []string {
	var unionmap TimeMap
	var times TimeSlice
	var data []string
	for _, e := range tm {
		for k, v := range e {
			unionmap[k] = v
		}
	}
	for k := range unionmap {
		times = append(times, k)
	}
	sort.Sort(times)
	for _, e := range times {
		data = append(data, unionmap[e])
	}

	return data
}

// GetStatuses returns the string slice containing a user's statuses
func (index UserIndex) GetStatuses(url string) TimeMap {
	imutex.RLock()
	status := index[url].Status
	imutex.RUnlock()
	return status
}