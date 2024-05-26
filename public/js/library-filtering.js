export class LibraryFiltering {
    sortButton;
    filterButton;
    grid;
    constructor() {
        let grid = $(".masonry-container")
        if (!grid) {

            return;
        }

        let sortButton = document.getElementById('sort-button');
        if (sortButton) {
            this.sortButton = sortButton;
            this.sortButton.addEventListener('click', this.OrderItems.bind(this));
        }

        let filterButton = document.getElementById('filter-button');
        if (filterButton) {
            this.filterButton = filterButton;
            this.filterButton.addEventListener('click', this.FilterItems.bind(this));
        }
    }

    OrderItems() {
        // Get current Order from url
        let currentOrderBy = new URLSearchParams(window.location.search).get('order_by');
        if (currentOrderBy == null) {
            currentOrderBy = "title";
        }

        let currentOrder = new URLSearchParams(window.location.search).get('order');
        if (currentOrder == null) {
            currentOrder = "asc";
        }

        Swal.fire({
            html: `<div class="flex flex-col items-start gap-2">
            <h2>Order by</h2>
                <label>
                    <input type="radio" value="title" name="order_by" ${currentOrderBy == "title" ? "checked" : ""}>
                    Title
                </label>
                <label>
                    <input type="radio" value="meta_rating" name="order_by" ${currentOrderBy == "meta_rating" ? "checked" : ""} >
                    Rating
                </label>
                <label>
                    <input type="radio" value="duration" name="order_by" ${currentOrderBy == "duration" ? "checked" : ""} >
                    Runtime
                </label>
                <label>
                    <input type="radio" value="release_date" name="order_by" ${currentOrderBy == "release_date" ? "checked" : ""} >
                    Release Date
                </label>
                <label>
                    <input type="radio" value="created_at" name="order_by" ${currentOrderBy == "created_at" ? "checked" : ""} >
                    Scan Date
                </label>
            </div>
            <div class="flex flex-col items-start gap-2 mt-5">
            <h2>Order</h2>
                <label>
                    <input type="radio" value="asc" name="order" ${currentOrder == "asc" ? "checked" : ""}>
                    Ascending
                </label>
                <label>
                    <input type="radio" value="desc" name="order" ${currentOrder == "desc" ? "checked" : ""}>
                    Descending
                </label>
            </div>`,
            customClass: {
                confirmButton: 'button button-primary',
            },
            buttonsStyling: false,
            scrollbarPadding: false,
        }).then((result) => {
            if (result.isConfirmed) {
                let order_by = $("input[name='order_by']:checked").val();
                let order = $("input[name='order']:checked").val();

                if (!order_by || !order) {

                    return;
                }

                let url = new URL(window.location.href);
                url.searchParams.set('order_by', order_by);
                url.searchParams.set('order', order);
                window.location.href = url.toString();
            }
        })
    }

    FilterItems() {
        console.log("Filtering items...");
    }

}