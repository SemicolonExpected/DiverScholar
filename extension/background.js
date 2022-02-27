
let switched = true;

chrome.storage.sync.set({switched: switched});


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
                    chrome.cookies.set(details, (cookiez) => {});
                } else {
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

        // Parse in the target page, and store the original SERP
        chrome.scripting.executeScript({
                target: {tabId: tab.id},
                function: parseResult
            },
            (injectionResults) => {
                let originalSERP = injectionResults[0].result;
                chrome.storage.sync.get(['searcherID'],
                    function(sess) {
                        let req = originalSERP;
                        req.user_cookie.key = sess.searcherID;

                        callRanker(req, function(response) {
                            let rankings = JSON.parse(response);

                            console.log(rankings);

                            chrome.scripting.executeScript({
                                target: {tabId: tab.id},
                                function: populateResult,
                                args: [rankings]
                            })



                            chrome.storage.sync.set({rankings});
                        });

                    });
                chrome.storage.local.set({originalSERP});
            });
    }
});


function callRanker(paperList, callback) {
    const posturl = "https://dwio-hack-nyu-2022.uc.r.appspot.com/api/ranker";
    fetch(posturl, {
        body: JSON.stringify(paperList),
        method: "POST"
    }).then(r => r.text()).then(r => callback(r));

}


function parseResult() {
    let SERP = document.querySelector("ol.breathe-horizontal"); //get all search results
    if (SERP == null) {
        return;
    }
    SERP.id = "originalSERP";
    let reorderedSERP = SERP.cloneNode(true);
    reorderedSERP.id = "reorderedSERP";
    reorderedSERP.hidden = false;
    //reorderedSERP.style.backgroundColor = "#990000";

    SERP.parentNode.insertBefore(reorderedSERP, SERP.nextSibling);
    SERP.hidden = true;

    let papers = [];
    for (let i = 0; i < reorderedSERP.children.length; i++) {
        let li = reorderedSERP.children[i];
        let authBlock = li.querySelector("p.authors")
            .getElementsByTagName("a");
        let titleBlock = li.querySelector("p.title");
        let url = li.querySelector("div.is-marginless p.list-title")
            .getElementsByTagName("a");

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
            URL: url[0].href,
            Title: titleBlock.textContent.trim(),
            authors: authors,
        })
    }
    return {
        user_cookie: {
            key: "empty",
        },
        search_url: window.location.href,
        papers: papers
    };
}


function populateResult(ranks) {

    let SERP = document.getElementById("reorderedSERP");
    let results = SERP.children;

    let test = [0,1,3,2,4,5] //elissa should be 3
    let test1 = [0,0,1,0,0,0]
    for (var i = 0; i < ranks['ordering'].length; i++) {
    	SERP.appendChild(results[ranks['ordering'][i]].cloneNode(true))
    }
    for (var i = 0; i < ranks['ordering'].length; i++) {
    	SERP.removeChild(SERP.firstElementChild)
    }
    for (var i = 0; i < test.length; i++) {
    	if(ranks['classes'][i] == 1){
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