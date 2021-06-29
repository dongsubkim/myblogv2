for (const tb of document.querySelectorAll("table")) {
    tb.classList.add("pure-table", "pure-table-bordered", "mb-3");
}
const btns = document.querySelectorAll(".edit-btn")
console.log(btns)
for (var i = 0; i < btns.length; i++) {
    let btn = btns[i]
    let id = "#editForm" + i
    btn.addEventListener("click", () => {
        console.log(id)
        let editForm = document.querySelector(id)
        editForm.classList.toggle("d-none")
    })
}