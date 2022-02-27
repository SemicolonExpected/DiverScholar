let toggleButton = document.getElementById("tog");

toggleButton.addEventListener("change", async() => {
    console.log("here");
    let [tab] = await chrome.tabs.query({
        active: true,
        currentWindow: true
    });
    chrome.scripting.executeScript({
        target: {tabId: tab.id},
        function: swapSERP,
    });
})

function swapSERP() {
    console.log("here2");

    let original = document.getElementById("originalSERP");
    let reordered = document.getElementsByClassName("reorderedSERP");

    original.hidden = !original.hidden;
    reordered.hidden = !reordered.hidden;
}