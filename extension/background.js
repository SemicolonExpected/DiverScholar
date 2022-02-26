

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
        });
    }



})


function parseResult() {
    let SERP = document.querySelector("ol.breathe-horizontal");
    console.log(SERP);
    let results = SERP.children;
    console.log(results);
    let current = [];
    for (var i = 0; i < results.length; i++) {
        current.push(results[i]);
    }
    console.log(current[0]);

    SERP.innerHTML = "";


}