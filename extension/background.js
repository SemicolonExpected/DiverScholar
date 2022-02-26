

let color = "#3aa757";

chrome.runtime.onInstalled.addListener(
    () => {
        chrome.storage.sync.set({color});
        console.log("default color is %cgreen", `color: ${color}`);
    }
);

chrome.tabs.onUpdated.addListener(function (tabId, changeInfo, tab) {

    // console.log(tabId);
    // console.log(changeInfo);
    // console.log(tab);

    if (changeInfo.status === "complete") {
        if (!tab.url.startsWith("https://arxiv.org/search/")) {
            return;
        }

        console.log(tab.url);
        chrome.scripting.executeScript({
            target: {tabId: tab.id},
            function: parseResult,
            function: populateResult,
        });
    }



})


function parseResult() {
    let SERP = document.querySelector("ol.breathe-horizontal"); //get all search results
    //console.log(SERP);
    let results = SERP.children;
    console.log(results);
    //console.log(results[i]); //this is in case they change their layout
    let titleAuthorPair = []
    for (var i = 0; i < results.length; i++) {
        let authors = [];
        let el = results[i].children;
        let title = el[1].innerHTML.trim();
        let el2 = el[2].children;
        for(var j = 1; j < el2.length; j++){
            authors.push(el2[j].innerHTML);
        }
        titleAuthorPair.push({"title": title, "authors": authors});
    }
    let url = window.location.href;
    let output = {"url": url, "entries" : titleAuthorPair};
    console.log(output);

    //SERP.innerHTML = "";

}

function populateResult() {
    // get result
    let SERP = document.querySelector("ol.breathe-horizontal");
    let results = SERP.children;

    let test = [0,1,3,2,4,5] //elissa should be 3
    for (var i = 0; i < test.length; i++) {
    	SERP.appendChild(results[test[i]].cloneNode(true))
    }
    for (var i = 0; i < test.length; i++) {
    	SERP.removeChild(SERP.firstElementChild)
    }
}