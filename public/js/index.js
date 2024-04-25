$(document).ready(function () {
    // Add listener for "menu" button click to toggle sidebar
    $("#sidebar-menu-button").click(function () {
        console.log("clicked");
        $("#sidebar").toggleClass("closed");
    })
})