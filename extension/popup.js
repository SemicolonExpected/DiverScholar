document.addEventListener('DOMContentLoaded', function () {
    var checkbox = document.querySelector('input[type="checkbox"]');
    chrome.storage.local.get('enabled', function (result) {
        if (result.enabled != null) {
            checkbox.checked = result.enabled;
        }
    });
    checkbox.addEventListener('click', function () {
        console.log(checkbox.checked);
        chrome.storage.local.set({ 'enabled': checkbox.checked }, function () {
            console.log("confirmed");
        });
    });
});

// var toggle = document.getElementById('tog');
// var enabled = toggle.checked;
// var background = chrome.extension.getBackgroundPage();

// chrome.storage.local.get('enabled', data => {
//     //enabled = !!data.enabled;

//     if(toggle.checked){
//         enabled = True;
//     }
//     else{
//         enabled = False
//     }
// });

// myButton.onclick = () => {
//     chrome.storage.local.set({enabled:enabled});
// };

// let changeColor = document.getElementById("changeColor");

// chrome.storage.sync.get("color", ({color}) => {
//     changeColor.style.backgroundColor = color;
// })

// changeColor.addEventListener("click", async () => {
//     let [tab] = await chrome.tabs.query({
//         active: true,
//         currentWindow: true
//     });
//     chrome.scripting.executeScript({
//         target: {tabId: tab.id},
//         function: setPageBackgroundColor,
//     });
// });

// function setPageBackgroundColor() {
//     chrome.storage.sync.get("color", ({color}) => {
//         document.body.style.backgroundColor = color;
//     });
// }