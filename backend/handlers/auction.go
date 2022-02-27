package handlers

func scorePaper(paper *Paper) int {
	result := 0
	for _, a := range paper.Authors {
		if a.Score.Confidence <= .6 {
			continue
		}
		if a.Score.Gender == "male" {
			result--
		} else if a.Score.Gender == "female" {
			result++
		}
	}

	return result
}

func mcRanks(scores []int) []int {

	// ranking algorithm

	return scores
}
