import { ScrollIndicator } from "./scroll-indicator.js";
import { MasonryGrid } from "./masonry-grid.js";
let $grid = null;

$(document).ready(function () {

    // Initialize scroll indicator
    let scrollIndicator = new ScrollIndicator();

    // Initialize masonry grid
    let grid = new MasonryGrid();

    // Add listener for "menu" button click to toggle sidebar
    $("#sidebar-menu-button").click(function () {
        console.log("clicked");
        $("#sidebar").toggleClass("closed");
        setTimeout(function () {
            grid.RefreshGrid();
            scrollIndicator.RefreshProgress();
        }, 300);
    })

    // Add listener for "back" button
    $("#back-button").click(function () {
        console.log("Going back...")
        window.history.back();
    })

})
