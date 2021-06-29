const title = document.querySelector("#title");
const category = document.querySelector("#category");
const content = document.querySelector("#content");

const dailylog = document.querySelector("#dailylog-template");
const leetcode = document.querySelector("#leetcode-template");
const projecteuler = document.querySelector("#projecteuler-template");
const programmers = document.querySelector("#programmers-template");

dailylog.addEventListener("click", function () {
    title.value = "Daily Log 2021";
    category.value = "Daily Log";
    simplemde.value(`## LeetCode
- LeetCode Challenge 2021 []()
- [Post]()

## Project Euler
- Problem []()
- [Post]()

*Copied from https://github.com/dongsubkim/dailylog/issues/ *`
    ); // Returns HTML from a custom parser

})
leetcode.addEventListener("click", function () {
    title.value = "[LeetCode]";
    category.value = "LeetCode";
    simplemde.value(`# LeetCode Challenge 2021
## []()

## My solution in 
\`\`\`

\`\`\` 

*Copied from https://github.com/dongsubkim/dailylog/issues/ *`)
})
projecteuler.addEventListener("click", function () {
    title.value = "[Project Euler] P0";
    category.value = "Project Euler";
    simplemde.value(`# Problem 
## []()

![](http://)

# Check my solution in [Jupyter Notebook]()

*Copied from https://github.com/dongsubkim/dailylog/issues/ *`
    ); // Returns HTML from a custom parser
})
programmers.addEventListener("click", function () {
    title.value = "[Programmers]";
    category.value = "Programmers";
    simplemde.value(`## []()

## My solution in 
\`\`\`

\`\`\` 

*Copied from https://github.com/dongsubkim/dailylog/issues/ *`)
})
