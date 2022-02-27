## Backend

Our project is fully hosted within Google Cloud Platform, and we have chosen to use AppEngine for the API (utilizing the legacy Datastore APIs as our persistence and cache layer)

Details on GCP usage can be found within the `app.yaml` and the `test_rank.http` file

### Algorithm 
https://github.com/SemicolonExpected/DiverScholar/blob/9391740d2f2205572190c62ca5e7fa3f9e6d3105/backend/handlers/auction.go#L40

1. When a request comes in from the Chrome Extension, we break the Paper into its component authors
2. Authors' First Names are then passed to genderize.io for "scoring." This yields a gender classification, with a probability/confidence
3. We score a paper for the balance of representation within its Authors, including a special flag if the first position is identified as female
4. Combining #2 and #3 we are able to simply cluster the papers into groups (see `ScorePaper` function)
5. We calculate the existing SERP's average ranking for each of the clusters, and use this as a baseline for how to augment their "bid" in the auction
6. After applying the necessary adjustments to a positional probability, we then run a randomized "auction" to select which paper should move to which position


### Simulation Results

