let $grid = null;

$(document).ready(function () {

    // generate masonry grid
    $grid = $('.masonry-container').isotope({
        // options
        itemSelector: '.masonry-item',
        percentPosition: true,
        masonry: {
            columnWidth: '.masonry-item',
            horizontalOrder: true
        }
    });

    // Add listener for "menu" button click to toggle sidebar
    $("#sidebar-menu-button").click(function () {
        console.log("clicked");
        $("#sidebar").toggleClass("closed");
        setTimeout(function () {
            RefreshGrid();
        }, 300);
    })

    $grid.imagesLoaded().progress(function () {
        console.log("images loaded")
        RefreshGrid()
    });
})


function RefreshGrid() {
    if ($grid == null) {
        return;
    }
    $grid.isotope('layout');
}