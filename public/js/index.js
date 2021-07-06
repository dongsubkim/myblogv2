let index = {
    init: function () {
        $("#btn-search").on("click", () => {
            this.search();
        });
    },
    search: function () {
        let searchQuery = $("serach").val();
        console.log(searchQuery);
        $.ajax({
            type: "GET",
            url: `/post?search=${searchQuery}`,
        }).done(function (resp) {
            console.log(resp)
        }).fail(function (error) {
            console.log(error)
        });
    }
}

index.init();