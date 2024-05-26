import { ScrollIndicator } from "./scroll-indicator.js";
import { MasonryGrid } from "./masonry-grid.js";
import { VideoPlayer } from "./video-player.js";
import { SyncMenu } from "./syncMenu.js";
import { PlayerIframe } from "./player-iframe.js";
import { Routes } from "./routes.js";

let $grid = null;

$(document).ready(function () {

    // Initialize scroll indicator
    let scrollIndicator = new ScrollIndicator();

    // Initialize masonry grid
    let grid = new MasonryGrid();

    // Initialize video player
    let videoPlayer = new VideoPlayer();

    // Initialize Sync menu
    let syncplayMenu = new SyncMenu();

    // Add listener for "menu" button click to toggle sidebar
    $("#sidebar-menu-button").click(function () {
        $("#sidebar").toggleClass("closed");
        setTimeout(function () {
            grid.RefreshGrid();
            scrollIndicator.RefreshProgress();
        }, 300);
    })

    // Add listener for "back" button
    $(".back-button").click(function () {
        console.log("Going back...")
        window.history.back();
    })

    // Add listener for "play" to open iframe
    $(".iframe-browser-button").click((e) => {
        // get data from button
        let itemID = e.currentTarget.dataset.itemid;
        if (!itemID) {
            itemID = 0;
        }

        let playerIframe = new PlayerIframe();
        playerIframe.AddIframe(itemID);
    });

    // Add listener for "favorite" item
    $(".favorite-button").click((e) => {
        // get data from button
        let itemID = e.currentTarget.dataset.itemid;
        console.log(itemID);
        if (!itemID) {
            itemID = 0;
        }
        console.log(Routes)
        console.log(Routes.ApiRoutes.ToggleFavoriteItem)
        fetch(Routes.ApiRoutes.ToggleFavoriteItem, {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ "itemID": parseInt(itemID) }),
        }).then((response) => {
            response.json().then((data) => {
                if (data.message == "success") {
                    e.target.classList.toggle("active");
                }
            })
        })
    });

    // Add listener for "watched" item
    $(".watched-button").click((e) => {
        // get data from button
        let itemID = e.currentTarget.dataset.itemid;
        console.log(itemID);
        if (!itemID) {
            itemID = 0;
        }
        fetch(Routes.ApiRoutes.ToggleWatchedItem, {
            method: "POST",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ "itemID": parseInt(itemID) }),
        }).then((response) => {
            response.json().then((data) => {
                if (data.message == "success") {
                    e.target.classList.toggle("active");
                }
            })
        })
    });

})
