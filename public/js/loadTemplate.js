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
- [Post]()`
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
`)
})
projecteuler.addEventListener("click", function () {
    let n = title.value;
    // let n = parseInt(simplemde.value().split("\n")[0].slice(10));

    let url = "";
    if (n % 100 <= 50 && n % 10 != 0) {
        url = `https://github.com/dongsubkim/project_euler/blob/main/problem${parseInt(n / 100)}01-${parseInt(n / 100)}50/p${n}.ipynb`
    } else if (n % 100 == 0) {
        url = `https://github.com/dongsubkim/project_euler/blob/main/problem${parseInt(n / 100) - 1}51-${parseInt(n / 100)}00/p${n}.ipynb`
    } else {
        url = `https://github.com/dongsubkim/project_euler/blob/main/problem${parseInt(n / 100)}51-${parseInt(n / 100) + 1}00/p${n}.ipynb`
    }

    if (n < 100) {
        title.value = `[Project Euler] P0${n}. `; // ${title.value.slice(16)}`;
    } else {
        title.value = `[Project Euler] P${n}. `; // ${title.value.slice(16)}`;
    }
    category.value = "Project Euler";
    simplemde.value(`# Problem ${n}
## [](https://projecteuler.net/problem=${n})

![](http://)

# Check my solution in [Jupyter Notebook](${url})
`); // Returns HTML from a custom parser
})
programmers.addEventListener("click", function () {
    title.value = "[Programmers]";
    category.value = "Programmers";
    simplemde.value(`## []()

## My solution in 
\`\`\`

\`\`\` 

*Copied from https://github.com/dongsubkim/dailylog/issues/*`)
})
