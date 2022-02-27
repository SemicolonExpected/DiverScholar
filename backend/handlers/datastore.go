package handlers

import (
	"context"
	"encoding/json"
	"github.com/spaolacci/murmur3"
	"google.golang.org/appengine/v2/datastore"
	"log"
)

func upsertUser(ctx context.Context, req *RankRequest) (*datastore.Key, *User, error) {
	userKey := datastore.NewKey(ctx, "USER", req.UserCookie.Key, 0, nil)

	// get
	var des string
	getErr := datastore.Get(ctx, userKey, &des)
	if getErr != nil {
		log.Println("get err:", getErr)
		_, putErr := datastore.Put(ctx, userKey, &req.UserCookie)
		if putErr != nil {
			log.Println("put err:", putErr)
			return nil, nil, putErr
		}
	}

	return userKey, &req.UserCookie, nil
}

func upsertSearch(ctx context.Context, req *RankRequest, user *datastore.Key) (
	*datastore.Key, *Search, error) {

	hasher := murmur3.New64()
	_, hashErr := hasher.Write([]byte(req.SearchURL))
	if hashErr != nil {
		log.Println("could not hash search url:", hashErr, "(", req.SearchURL, ")")
		return nil, nil, hashErr
	}
	urlHash := int64(hasher.Sum64())
	userSearch := Search{
		Key:    urlHash,
		URL:    req.SearchURL,
		Papers: req.Papers,
	}

	payload, err := json.Marshal(userSearch)
	if err != nil {
		log.Println("could not marshall...", err)
	}

	searchKey := datastore.NewKey(ctx, "SEARCH", "", urlHash, user)

	plist := datastore.PropertyList{
		datastore.Property{Name: "payload", NoIndex: true, Value: payload},
	}

	_, putErr := datastore.Put(ctx, searchKey, &plist)
	if putErr != nil {
		log.Println("could not insert search for user:",
			putErr, "(", user.String(), ", ", req.SearchURL)
		return nil, nil, putErr
	}

	return searchKey, &userSearch, nil

}

func scorePaperCached(ctx context.Context, search *Search) ([]int, error) {
	var paperDsKeys = make([]*datastore.Key, len(search.Papers))
	var paperKeyToPos = make(map[int64]int)

	hasher := murmur3.New64()

	for pos, paper := range search.Papers {
		hasher.Reset()
		_, hashErr := hasher.Write([]byte(paper.URL))
		if hashErr != nil {
			log.Println("couldn't hash", paper.URL)
			continue
		}

		pHash := int64(hasher.Sum64())
		paperDsKeys[pos] = datastore.NewKey(ctx, "PAPER", "", pHash, nil)
		paperKeyToPos[pHash] = pos
	}

	var cachedPapers []Paper

	getErr := datastore.GetMulti(ctx, paperDsKeys, &cachedPapers)
	if getErr != nil {
		log.Println("couldn't get any:", getErr)

	}

	var scores = make([]int, len(search.Papers))

	for _, paper := range cachedPapers {
		scores[paperKeyToPos[paper.Key]] = scorePaper(&paper)
		delete(paperKeyToPos, paper.Key)
	}

	// remaining we need to go load...
	for paperKey, paperPos := range paperKeyToPos {
		paperEntry := Paper{
			Key:     paperKey,
			URL:     search.Papers[paperPos].URL,
			Title:   search.Papers[paperPos].Title,
			Authors: make([]Author, len(search.Papers[paperPos].Authors)),
		}

		// check each of the authors
		for aPos, author := range search.Papers[paperPos].Authors {
			cachedAuthor, _ := getAuthorCached(ctx, author)
			paperEntry.Authors[aPos] = cachedAuthor
		}
		scores[paperPos] = scorePaper(&paperEntry)
	}
	return scores, nil
}

func getAuthorCached(ctx context.Context, author Author) (Author, error) {
	hasher := murmur3.New64()
	hasher.Reset()
	_, hashErr := hasher.Write([]byte(author.AuthorLink))
	if hashErr != nil {
		return Author{}, hashErr
	}
	authHash := int64(hasher.Sum64())
	var cachedAuthor Author
	authKey := datastore.NewKey(ctx, "AUTHOR", "", authHash, nil)
	getErr := datastore.Get(ctx, authKey, &cachedAuthor)
	if getErr != nil {
		cachedAuthor.Key = authHash
		cachedAuthor.AuthorLink = author.AuthorLink
		cachedAuthor.FullName = author.FullName
		cachedAuthor.FirstName = author.FirstName
		cachedAuthor.Score = *getName(ctx, author.FirstName)

		_, putErr := datastore.Put(ctx, authKey, &cachedAuthor)
		if putErr != nil {
			log.Println("could not insert author", putErr)
			return Author{}, putErr
		}
	}
	return cachedAuthor, nil
}

func getNameCahced(ctx context.Context) {

}
