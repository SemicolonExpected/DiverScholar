package main

import (
	"backend/handlers"
	"fmt"
	"log"
	"math/rand"
	"os"
)

var maleAuthor = handlers.Author{
	FirstName: "mark", GenderedName: handlers.Name{FirstName: "mark", Gender: "male", Confidence: .9},
}

var femaleAuthor = handlers.Author{
	FirstName: "mary", GenderedName: handlers.Name{FirstName: "mary", Gender: "female", Confidence: .9},
}

var unknownAuthor = handlers.Author{
	FirstName: "pat", GenderedName: handlers.Name{FirstName: "mary", Gender: "female", Confidence: .5},
}

var PaperList = []handlers.Paper{
	{Key: 0, URL: "a", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 1, URL: "b", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 2, URL: "c", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 3, URL: "d", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 4, URL: "e", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 5, URL: "f", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 6, URL: "g", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 7, URL: "h", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 8, URL: "i", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 9, URL: "j", Title: "a", Authors: []handlers.Author{maleAuthor, unknownAuthor}},
	{Key: 10, URL: "k", Title: "a", Authors: []handlers.Author{unknownAuthor, maleAuthor}},
	{Key: 11, URL: "l", Title: "a", Authors: []handlers.Author{unknownAuthor, unknownAuthor}},
	{Key: 12, URL: "m", Title: "a", Authors: []handlers.Author{maleAuthor, maleAuthor}},
	{Key: 13, URL: "n", Title: "a", Authors: []handlers.Author{maleAuthor, femaleAuthor}},
	{Key: 14, URL: "o", Title: "a", Authors: []handlers.Author{femaleAuthor, maleAuthor}},
	{Key: 15, URL: "p", Title: "a", Authors: []handlers.Author{femaleAuthor, femaleAuthor}},
}

func main() {
	f, err := os.Create("/Users/walt/hack/DiverScholar/backend/results.csv")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, _ = f.WriteString("permutation_id,simulation_id,paper_type,starting_pos,ending_pos,delta\n")

	for iter := 0; iter < 50; iter++ {
		rand.Shuffle(len(PaperList), func(i int, j int) {
			PaperList[i], PaperList[j] = PaperList[j], PaperList[i]
		})

		groupAssigned := make([]handlers.Group, len(PaperList))

		for i, p := range PaperList {
			s := handlers.ScorePaper(&p)
			groupAssigned[i] = s

		}

		for simul := 0; simul < 1000; simul++ {
			new_ranks := handlers.MCRank(groupAssigned)

			for i := 0; i < len(PaperList); i++ {
				_, err := fmt.Fprintf(f, "%d,%d,%s,%d,%d,%d\n",
					iter,
					simul,
					groupAssigned[i],
					i,
					new_ranks[i],
					new_ranks[i]-i)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

}
