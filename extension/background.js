
chrome.runtime.onInstalled.addListener(
    () => {
        let details = {name: "searcherID", url: "https://dwio-hack-nyu-2022.uc.r.appspot.com/"};
        chrome.cookies.get(
            details,
            (cookie) => {
                if (cookie == null) {
                    let searcherID = makeid(32);
                    chrome.storage.sync.set({searcherID});
                    details.value = searcherID;
                    chrome.cookies.set(details, (cookiez) => {
                        console.log("set cookie:", cookiez);
                    });
                } else {
                    console.log("got one:", cookie);
                    let searcherID = cookie.value;
                    chrome.storage.sync.set({searcherID});
                }
            }
        );

     }
 );



chrome.tabs.onUpdated.addListener(function (tabId, changeInfo, tab) {

    if (changeInfo.status === "complete") {
        if (!tab.url.startsWith("https://arxiv.org/search/")) {
            return;
        }
        console.log(tab.url);

        // Parse in the target page, and store the original SERP
        chrome.scripting.executeScript({
                target: {tabId: tab.id},
                function: parseResult
            },
            (injectionResults) => {
                console.log(injectionResults);
                let originalSERP = injectionResults[0].result;
                chrome.storage.local.set({originalSERP});
                console.log(originalSERP);
            });

        // Send the requests to the API

        chrome.scripting.executeScript({
            target: {tabId: tab.id},
            function: callRanker,
            args: ["helloz"]
        });

        console.log("after");
    }
});




function parseResult() {
    let SERP = document.querySelector("ol.breathe-horizontal"); //get all search results
    if (SERP == null) {
        return;
    }
    SERP.id = "originalSERP";
    let reorderedSERP = SERP.cloneNode(true);
    reorderedSERP.id = "reorderedSERP";
    reorderedSERP.hidden = false;
    reorderedSERP.style.backgroundColor = "#990000";

    SERP.parentNode.insertBefore(reorderedSERP, SERP.nextSibling);
    SERP.hidden = true;

    let papers = [];
    for (let i = 0; i < reorderedSERP.children.length; i++) {
        let li = reorderedSERP.children[i];
        let authBlock = li.querySelector("p.authors")
            .getElementsByTagName("a");
        let titleBlock = li.querySelector("p.title");

        let authors = []
        for (let a = 0; a < authBlock.length; a++) {
            let nameArr = authBlock[a].text.trim().split(" ");
            authors.push({
                author_link: authBlock[a].href,
                full_name: authBlock[a].text.trim(),
                first_name: nameArr[0],
            })
        }

        papers.push({
            URL: document.URL,
            Title: titleBlock.textContent.trim(),
            authors: authors,
        })
    }

    return papers;
    // Added some more fields to the parsed <li>s
    /*
    let titleAuthorPair = []
    for (var i = 0; i < results.length; i++) {
        let authors = [];
        let el = results[i].children;
        let title = el[1].innerHTML.trim();
        let el2 = el[2].children;
        for(let j = 1; j < el2.length; j++){
            authors.push(el2[j].innerHTML);
        }
        titleAuthorPair.push({"title": title, "authors": authors});
    }
    let url = window.location.href;
    let output = {"url": url, "entries" : titleAuthorPair};
    console.log(output);
    */
}

function callRanker(paperList) {

}

function populateResult() {
    // get result
    let SERP = document.querySelector("ol.breathe-horizontal");
    let results = SERP.children;

    let test = [0,1,3,2,4,5] //elissa should be 3
    let test1 = [0,0,1,0,0,0]
    for (var i = 0; i < test.length; i++) {
    	SERP.appendChild(results[test[i]].cloneNode(true))
    }
    for (var i = 0; i < test.length; i++) {
    	SERP.removeChild(SERP.firstElementChild)
    }
    for (var i = 0; i < test.length; i++) {
    	if(test1[i] == 1){
    		temp = document.createElement('span')
    		temp.innerHTML = "<b>&nbsp;<font color='gold' size = '+1'>&#9728;</font> Women Led or Strong Female Representation <font color='gold' size = '+1'>&#9728;</font>&nbsp;</b>";
    		temp.style.backgroundColor = "lavender";
    		temp.style.borderRadius = "20px";
    		results[i].children[2].appendChild(temp)
    	}
    }
}

function makeid(length) {
    var result           = '';
    var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for ( var i = 0; i < length; i++ ) {
        result += characters.charAt(Math.floor(Math.random() *
            charactersLength));
    }
    return result;
}