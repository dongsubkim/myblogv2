const title = document.querySelector("#title");
const category = document.querySelector("#category");
const content = document.querySelector("#content");

const dailylog = document.querySelector("#dailylog-template");
const leetcode = document.querySelector("#leetcode-template");
const projecteuler = document.querySelector("#projecteuler-template");
const programmers = document.querySelector("#programmers-template");

let d = new Date();
dailylog.addEventListener("click", function () {
    let month = d.getMonth() + 1
    if (month < 10) {
        month = "0" + month.toString();
    }
    let date = d.getDate();
    if (date < 10) {
        date = "0" + date.toString();
    }
    console.log(d);
    title.value = `Daily Log ${d.getFullYear()}-${month}-${date}`;
    category.value = "Daily Log";
    simplemde.value(`## LeetCode
- ${d.toLocaleString('en-US', { month: 'long' })} LeetCode Challenge 2021 []()
- [Post]()

## Project Euler
- Problem []()
- [Post]()`
    ); // Returns HTML from a custom parser

})
leetcode.addEventListener("click", function () {

    title.value = "[LeetCode] ";
    category.value = "LeetCode";
    simplemde.value(`# ${d.toLocaleString('en-US', { month: 'long' })} LeetCode Challenge 2021
## []()

## My solution in 
\`\`\`

\`\`\` 
`)
})
projecteuler.addEventListener("click", function () {
    let n = title.value;
    // let n = parseInt(simplemde.value().split("\n")[0].slice(10));
    let q = parseInt(n / 50);
    if (n % 50 == 0) {
        q--;
    }
    let url = `https://github.com/dongsubkim/project_euler/blob/main/problem${q * 50 + 1}-${(q + 1) * 50}/p${n}.ipynb`;

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
    title.value = "[Programmers] ";
    category.value = "Programmers";
    simplemde.value(`## []()

## My solution in 
\`\`\`

\`\`\` 

*Copied from https://github.com/dongsubkim/dailylog/issues/*`)
})
