package handlers

import (
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/sampleuv"
	"time"
)

func ScorePaper(paper *Paper) Group {
	result := 0

	fl := paper.Authors[0].GenderedName.Gender == "female" &&
		paper.Authors[0].GenderedName.Confidence > .6

	for _, a := range paper.Authors {
		if a.GenderedName.Confidence <= .6 {
			continue
		}
		if a.GenderedName.Gender == "male" {
			result--
		} else if a.GenderedName.Gender == "female" {
			result++
		}
	}

	if result == -1*len(paper.Authors) {
		return MaleOnly
	} else if !fl && 0 < result && result < 1 {
		return GenderDiverse
	} else if fl && 0 < result && result < 1 {
		return FemaleLedDiverse
	} else if fl {
		return FemaleLed
	} else {
		return Unknown
	}

}

func MCRank(scores []Group) []int {
	result := make([]int, len(scores))

	tally := getAvgPos(scores)
	weights := getWeightAdj(scores, tally)

	result = weightedTake(weights)

	return result
}

type metric struct {
	Sum  int
	Cnt  int
	Mult float64
}

func getAvgPos(scores []Group) map[Group]*metric {

	tally := map[Group]*metric{
		MaleOnly:         {0, 0, 1.0},
		FemaleLed:        {0, 0, 1.0},
		GenderDiverse:    {0, 0, 1.0},
		FemaleLedDiverse: {0, 0, 1.0},
		Unknown:          {0, 0, 1.0},
	}

	avg := float64(len(scores)+1) / float64(2)
	for pos, val := range scores {
		g := tally[val]
		g.Sum += pos
		g.Cnt += 1
	}

	for _, v := range tally {
		v.Mult = (float64(v.Sum) / float64(v.Cnt)) / avg
	}

	return tally
}

func getWeightAdj(scores []Group, tally map[Group]*metric) []float64 {
	sz := len(scores)
	weights := make([]float64, sz)
	running := float64(0)

	for pos, val := range scores {
		running += float64(sz - pos + 1)
		baseline := float64(sz-pos+1) / running
		weights[pos] = baseline*tally[val].Mult + .01
	}
	return weights
}

func weightedTake(weights []float64) []int {
	result := make([]int, len(weights))
	w := sampleuv.NewWeighted(
		weights,
		rand.New(rand.NewSource(uint64(time.Now().UnixNano()))),
	)
	for i := 0; i < len(weights); i++ {
		next, ok := w.Take()
		if !ok || next < 0 {
			result[i] = i
		} else {
			result[i] = next
		}
	}
	return result
}
