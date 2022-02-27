# DiverScholar

In recent years, studies have shown that publishing in academia suffers from systemic biases towards already marginalized groups. In one study it was shown that in the worlds largest multidisciplinary journal Nature's articles quote men TWICE as much as women, are more likely to quote those with names commonly used in English speaking cultures. In fact, Nature's journalists were found to quote those with East Asian names LESS! ([Nature 2021](https://www.nature.com/articles/d41586-021-01676-7)) Another study found that in elite medical journals, papers authored women were half as likely to be cited than similar papers authored by men. Not only this, but women who coauthored with women had the fewest median citations whereas men coauthoring with men had the most. In their paper, Chatterjee and Werner attribute this bias to lack of visibility of women in academic medicine.([Chatterjee 2021](https://jamanetwork.com/journals/jamanetworkopen/fullarticle/2781617)) Harvard medical researcher, Julie Silvers believes that both the tendency of men to self promote more often than women, as well as implicit bias against women could play a part in both lack of visibility and lack of citations. ([Reardon 2021](https://www.nature.com/articles/d41586-021-02102-8#ref-CR5))
In an industry where citations and visibility are the keys to career progression, this bias already can hinder the prospects of marginalized groups. This is compounded upon by the fact that impact positively correlates with the probability of being cited—that is to say the papers with the most citations are more likely to be cited leading to a feedback loop where the cited get more cited.

## What it does

We hope to mitigate that by increasing the visibility of papers that are either led by women, or have a group of authors with strong female representation by both highlighting as well as boosting papers who meet the aforementioned criteria in academic search results. To do this, we devised a browser extension that would process search results from arxiv and slightly boost the rankings of papers that are women led or with a good ratio of women to men. Because there is a correlation between the search engine ranking of a result and engagement with the result ([Dean 2019](https://backlinko.com/google-ctr-stats)), we hope this will lead to more research done by women being seen in science.

## How we built it
Our extension is built off javascript, python, and golang. The extension scrapes search results from academic article archive *arxiv.org* and passes it to our back end. The backend uses takes the original rankings and reorders them using a fairness algorithm that instead of aiming for the top spot, we tried to make it so that both men and women on average rank similarly. This leads to more equitable statistics in the long run.

## Challenges we ran into
Currently we are limited in the amount of articles a day we can classify as being women led or having good female representation by the classifer we use to categorize gender: [genderize.io](https://genderize.io/). To mitigate this we implemented a cache so that if a name has been classified before it wouldn't have to be classifed again. This will overtime immensely cut down on the need to classify and therefore not have to worry about hitting the classification limit. Similarly we also have an article specific cache so not only do we not need to classify the author names, but we dont have to do any probability calculations again either!
We are also limited by the fact that since our weights are based off results of a classifier there will always be a margin of error without scholars self-disclosing as being of a certain gender. 

## Accomplishments that we're proud of
- Building an app in under 2 days
- Having participated in a Hackathon in general!

## What we learned
- If you give them the permissions, Chrome extensions have a lot of functionality.

## What's next for DiverScholar

We hope to be able to expand our functionality to other academic search engines and libraries e.g. Google Scholar, ACM. 
We also hope to be able incorporate more groups and different ways scholars may be disadvantaged into our algorithm and ultimately uplift all underrepresented groups.
We also hope that this inspires third parties to make adjustments to search results or recommendations pages, or news stories, or timelines -- we think this aggregate approach might have some larger potential in all aspects of web content

## Installation
Clone this repository and then in google chrome: go to [chrome://extensions/ ](chrome://extensions/)

Click Load unpacked and point to where the extension folder of this repo is located on your computer.

Search something on arxiv.org and watch it go!


## Contributions
Victoria:   Idea, Research and Development, Front End, Documentation

Daniel:     Fairness Algorithm, Backend, Documentation