const title = document.querySelector("#title");
const category = document.querySelector("#category");
const content = document.querySelector("#content");

const dailylog = document.querySelector("#dailylog-template");
const leetcode = document.querySelector("#leetcode-template");
const projecteuler = document.querySelector("#projecteuler-template");

dailylog.addEventListener("click", function () {
    title.value = "Daily Log 2021";
    category.value = "Daily Log";
    simplemde.value(`## LeetCode
- LeetCode Challenge 2021 []()
- [Post]()

## Project Euler
- Problem []()
- [Post]()`
    ); // Returns HTML from a custom parser

})
leetcode.addEventListener("click", function () {
    title.value = "[LeetCode]";
    category.value = "LeetCode";
    simplemde.value("# LeetCode Challenge 2021 \n ## []() \n \n ## My solution \n ``` \n ```")
})
projecteuler.addEventListener("click", function () {
    title.value = "[Project Euler] P";
    category.value = "Project Euler";
    simplemde.value(`# Problem
## []()

![](http://)

# My solution in [Jupyter Notebook]()`
    ); // Returns HTML from a custom parser

})