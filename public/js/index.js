let index = {
    init: function () {
        $("#btn-search").on("click", () => {
            this.search();
        });
    },
    search: function () {
        let searchQuery = $("#search").val();
        location.href = `/post?search=${searchQuery}`
    }
}

index.init();