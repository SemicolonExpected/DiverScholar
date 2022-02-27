package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/appengine/v2"
	"google.golang.org/appengine/v2/urlfetch"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	Ordering []int `json:"ordering,omitempty"`
}

func RankHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	var req RankRequest

	decErr := json.NewDecoder(r.Body).Decode(&req)
	if decErr != nil {
		log.Println("could not parse request:", decErr)
		// Silent fail
		return
	}

	log.Printf("request: %+v", req)

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
	_ = userSearch
	_ = err

	ranks, err := scorePaperCached(ctx, userSearch)
	_ = ranks
	_ = err

	result, _ := json.Marshal(ranks)
	w.Write(result)
	return

}

func getName(ctx context.Context, first string) *Name {

	// TODO: cache the responses from the name server (also into datastore)

	client := urlfetch.Client(ctx)
	client.Timeout = time.Millisecond * 200

	resp, err := client.Get(fmt.Sprintf("https://api.genderize.io?name=%s", first))
	if err != nil {
		log.Fatal("bad client req to service:", err)
	}
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("bad http request")
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var result Name
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		log.Fatal("json body err")
	}

	return &result
}
