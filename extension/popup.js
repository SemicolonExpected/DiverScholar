let toggleButton = document.getElementById("tog");

toggleButton.addEventListener("change", async() => {
    let [tab] = await chrome.tabs.query({
        active: true,
        currentWindow: true
    });

    toggleButton.checked
    chrome.scripting.executeScript({
        target: {tabId: tab.id},
        function: swapSERP,
        args: [toggleButton.checked]
    });
})

function swapSERP(switched) {
    let original = document.getElementById("originalSERP");
    let reordered = document.getElementById("reorderedSERP");

    if (switched) {
        original.hidden = true;
        reordered.hidden = false;
    } else {
        original.hidden = false;
        reordered.hidden = true;
    }
}