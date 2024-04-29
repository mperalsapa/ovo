export class MasonryGrid {
    grid = null;
    constructor() {
        // generate masonry grid
        this.grid = $('.masonry-container').isotope({
            // options
            itemSelector: '.masonry-item',
            // percentPosition: true,
            masonry: {
                columnWidth: '.masonry-item',
                horizontalOrder: true,
                fitWidth: true,
                gutter: 25
            }
        });


        this.grid.imagesLoaded().progress(this.RefreshGrid.bind(this));
    }

    RefreshGrid() {
        if (this.grid == null) {
            console.log("No grid found")
            return;
        }
        console.log("Refreshing grid")
        this.grid.isotope('layout');
    }

}