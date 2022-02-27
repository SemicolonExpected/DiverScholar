package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spaolacci/murmur3"
	"google.golang.org/appengine/v2/datastore"
	"google.golang.org/appengine/v2/urlfetch"
	"io/ioutil"
	"log"
	"time"
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

func scorePapersCached(ctx context.Context, search *Search) ([]Group, error) {
	var paperDsKeys = make([]*datastore.Key, len(search.Papers))
	var paperKeyToPos = make(map[int64]int)

	hasher := murmur3.New64()

	log.Println("got paper count:", len(search.Papers))
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

	var cachedPapers []*Paper

	getErr := datastore.GetMulti(ctx, paperDsKeys, &cachedPapers)
	if getErr != nil {
		log.Println("scorePapersCached didn't find any cached papers:", getErr)
	}
	log.Println("cache paper len:", len(cachedPapers))
	var scores = make([]Group, len(search.Papers))
	var toSkip = make(map[int64]bool)
	for _, paper := range cachedPapers {
		scores[paperKeyToPos[paper.Key]] = ScorePaper(paper)
		toSkip[paper.Key] = true
	}

	log.Println("papercache len after:", len(paperKeyToPos), len(scores))

	var toLoad []*Paper
	var withKeys []*datastore.Key
	// remaining we need to go fetch, score, and load...
	for paperKey, paperPos := range paperKeyToPos {
		if _, ok := toSkip[paperKey]; ok {
			log.Println("\tskipping:", paperKey, toSkip[paperKey])
			continue
		}
		paperEntry := Paper{
			Key:     paperKey,
			URL:     search.Papers[paperPos].URL,
			Title:   search.Papers[paperPos].Title,
			Authors: make([]Author, len(search.Papers[paperPos].Authors)),
		}

		// check each of the authors
		for aPos, author := range search.Papers[paperPos].Authors {
			log.Println("\tchecking author", author.FullName)
			cachedAuthor, _ := getAuthorCached(ctx, author)
			log.Println("\t\tgot:", cachedAuthor.FullName)
			paperEntry.Authors[aPos] = cachedAuthor
		}

		// Score it
		scores[paperPos] = ScorePaper(&paperEntry)

		// Prepare to load cache
		toLoad = append(toLoad, &paperEntry)
		withKeys = append(withKeys, paperDsKeys[paperPos])
	}

	// Load before exit
	_, putErr := datastore.PutMulti(ctx, withKeys, toLoad)
	if putErr != nil {
		log.Println("scorePapersCached could not load new:", putErr)
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
		var nameErr error
		cachedAuthor.GenderedName, nameErr = getNameCached(ctx, author.FirstName)
		if nameErr != nil {
			return Author{}, nameErr
		}
		_, putErr := datastore.Put(ctx, authKey, &cachedAuthor)
		if putErr != nil {
			log.Println("could not insert author", putErr)
			return Author{}, putErr
		}
	}
	return cachedAuthor, nil
}

func getNameCached(ctx context.Context, firstName string) (Name, error) {
	if len(firstName) < 3 {
		log.Println("getNameCached first name too short:", firstName)
		return Name{}, nil
	}

	hasher := murmur3.New64()
	hasher.Reset()
	_, hashErr := hasher.Write([]byte(firstName))
	if hashErr != nil {
		log.Println("getNameCached could not hash first name:", hashErr)
		return Name{}, hashErr
	}
	firstHash := int64(hasher.Sum64())
	nameKey := datastore.NewKey(ctx, "NAME", "", firstHash, nil)
	var cachedName Name
	getErr := datastore.Get(ctx, nameKey, &cachedName)
	if getErr != nil {
		client := urlfetch.Client(ctx)
		client.Timeout = time.Millisecond * 1500

		url := fmt.Sprintf(
			"https://api.genderize.io?name=%s",
			firstName,
		)

		resp, err := client.Get(url)
		defer resp.Body.Close()
		if err != nil {
			log.Println("getNameCached could not retrieve name score:", err)
			return Name{}, err
		}
		log.Println("genderize response:", resp.StatusCode, resp.Status)

		respBody, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(respBody, &cachedName)
		if err != nil {
			return Name{}, err
		}

		// load to cache
		_, putErr := datastore.Put(ctx, nameKey, &cachedName)
		if putErr != nil {
			log.Println("getNameCached could not put the name into cache:", putErr)
		}
	}

	return cachedName, nil
}
