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

We tested 50 different random permutations of a mix of ~20% Female papers and 80% Male papers in a SERP.  In each of these permutations, we ran the algorithm's "auction" 10,000 times and recorded the before/after of each of the individual papers' positions within the SERP.  

The objective of the algorithm is to ensure all papers move closer to the average position, and with an initial weighting multiplier of just the difference in the groups average pos vs. overall average position, this looks to be a positive improvement against our metric.  The results of this simulation are available in `results.csv`

#### Average change in rank after algorithm, by group
![Ranks](https://user-images.githubusercontent.com/12177307/155887174-bec3a015-49a5-45b6-9859-f212bcb3ad99.png)
 
 
#### The most disadvantaged group re-ordering
![Screen Shot 2022-02-27 at 9 47 09 AM](https://user-images.githubusercontent.com/12177307/155887186-1a12d0b7-f3ef-48d8-9095-5fa13a7525f7.png)


#### The most advantaged group reordering
![Screen Shot 2022-02-27 at 9 47 28 AM](https://user-images.githubusercontent.com/12177307/155887189-bf8dbab0-f26a-491d-bb6d-ab6672864bf3.png)
