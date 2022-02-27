
        //it is enabled, do accordingly
    

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

    // console.log(tabId);
    // console.log(changeInfo);
    // console.log(tab);
   //chrome.storage.local.get('enabled', data => {
    //if (data.enabled) {
	    if (changeInfo.status === "complete") {
	        if (!tab.url.startsWith("https://arxiv.org/search/")) {
	            return;
	        }

	        console.log(tab.url);
	        chrome.scripting.executeScript({
	            target: {tabId: tab.id},
	            function: parseResult,
	        });
	        chrome.scripting.executeScript({
	            target: {tabId: tab.id},
	            function: populateResult,
	        });
	    }
	//} 
	//});

});


function callRanker(paperList) {
	const url = "https://dwio-hack-nyu-2022.uc.r.appspot.com/api/score";
	const http = new XMLHttpRequest();
	http.open("GET", url);
	http.send(paperList);

	Http.onreadystatechange = (e) => {
 	 console.log(Http.responseText)
	}

}

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
    //callRanker(output);

    //SERP.innerHTML = "";
    const posturl = "https://dwio-hack-nyu-2022.uc.r.appspot.com/api/ranker";
	const http = new XMLHttpRequest();
	http.open("POST", posturl);
	http.setRequestHeader("Content-Type", "application/json");
	http.send(JSON.stringify(output));

	Http.onreadystatechange = (e) => {
 	 	console.log(http.responseText)
	}
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