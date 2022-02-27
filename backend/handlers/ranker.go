package handlers

import (
	"encoding/json"
	"google.golang.org/appengine/v2"
	"log"
	"net/http"
)

func init() {
	http.HandleFunc("/api/ranker", RankHandler)
}

type RankRequest struct {
	UserCookie User    `json:"user_cookie,omitempty"`
	SearchURL  string  `json:"search_url,omitempty"`
	Papers     []Paper `json:"papers,omitempty"`
}

type RankResponse struct {
	Ordering []int   `json:"ordering,omitempty"`
	Classes  []Group `json:"classes,omitempty"`
}

func RankHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	h := w.Header()
	h.Set("Access-Control-Allow-Origin", "*")

	var req RankRequest

	decErr := json.NewDecoder(r.Body).Decode(&req)
	if decErr != nil {
		log.Println("could not parse request:", decErr)
		// Silent fail
		return
	}

	// 1. Get User or Add
	// 2. Check if URL is in cache for User -> return if it is
	// 3. Check which papers have been scored -> score those that haven't been
	// 4. Re-weight the orders
	// 5. Auction
	// 6. Return results

	userKey, userData, err := upsertUser(ctx, &req)
	_ = userData
	_ = err

	searchKey, userSearch, err := upsertSearch(ctx, &req, userKey)
	_ = searchKey
	_ = err

	assignedGroups, err := scorePapersCached(ctx, userSearch)
	_ = err

	var resp RankResponse
	resp.Ordering = MCRank(assignedGroups)
	resp.Classes = assignedGroups

	result, _ := json.Marshal(resp)
	w.Write(result)
	return

}
