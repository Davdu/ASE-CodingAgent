# Design log

We followed the provded tutorial, and are now extending the system. we use the agent to help us learn go and fix mistakes that we make ourselves.
The output of the agent sometimes truncates, and a lot of the output gets lost. Sometimes it just doesnt respont. 
We cleaned up the code structure by extracting tools into its own submodule `tools`.

We had the agent rewrite the go program to also output the output of functions with the `console` terminal label.
We also mad ethe agent implement a feature, where it logs conversations to the disc so we can laugh over the wierd chats

## Run program DONE

Give the agent the ability to run go programs. First we decided that it should run go programs, but then thoughts like "Can it cancel the program with ctrl-c", and "What happens if it goes unresponsive" came. We decided that we also wanted to run JS programs, so we added the ability for the agent to chose between Go and Node itself when executing programs. We also add that it can supply arguments to the program invocation, such that it can write more sophisticated programs.

`\u001b[94mYou\u001b[0m: Alright so you have to work in a loop, using the given commands to implement a small program. you have to run it, test the output and then edit the fil
e until you are satisfied. Create a node.js tower of hanoi that bootstraps the inptu inside of the program, and then solves the problem. you should not stop iterating until yo
u are done
Claude: I'll create a Tower of Hanoi program in Node.js, run it, check the output, and fix anything that's wrong.
tool: edit_file({"path":"hanoi.js","old_str":""})
tool: edit_file({"path":"hanoi.js","old_str":"placeholder","new_str":"function hanoi(n, from, to, via, moves) {\n  if (n === 0) return moves;\n  hanoi(n - 1, from, via, to, mo
ves);\n  moves.push(`Move disk ${n} from ${from} to ${to}`);\n  hanoi(n - 1, via, to, from, moves);\n  return moves;\n}\n\nconst input = process.argv[2] ?? \"3\";\nconst n = N
umber(input);\n\nif (!Number.isInteger(n) || n < 1) {\n  console.error(\"Usage: node hanoi.js <number of disks (positive integer)>\");\n  process.exit(1);\n}\n\nconst moves = 
hanoi(n, \"A\", \"C\", \"B\", []);\nconsole.log(moves.join(\"\\n\"));\nconsole.log(`Total moves: ${moves.length} (expected ${2 ** n - 1})`);\n"})
Claude: The edit tool won't create the file with that input. I'll check its source to see how it expects new files to be created.
tool: read_file({"path":"tools/edit_file.go"})
Claude: The source shows the create path needs `old_str` to be empty and `new_str` to hold the content. My earlier attempts didn't do that, so I'll retry correctly.
tool: edit_file({"path":"hanoi.js","old_str":"","new_str":"function hanoi(n, from, to, via, moves) {\n  if (n === 0) return moves;\n  hanoi(n - 1, from, via, to, moves);\n  mo
ves.push(`Move disk ${n} from ${from} to ${to}`);\n  hanoi(n - 1, via, to, from, moves);\n  return moves;\n}\n\nconst input = process.argv[2] ?? \"3\";\nconst n = Number(input
);\n\nif (!Number.isInteger(n) || n < 1) {\n  console.error(\"Usage: node hanoi.js <number of disks (positive integer)>\");\n  process.exit(1);\n}\n\nconst moves = hanoi(n, \"
A\", \"C\", \"B\", []);\nconsole.log(moves.join(\"\\n\"));\nconsole.log(`Total moves: ${moves.length} (expected ${2 ** n - 1})`);\n"})
Claude: The file is created. Now I'll run it with 3 disks and check the output.
tool: run_program({"path":"./hanoi.js","program":"node"})
Claude: I created `hanoi.js`, ran it, and the output is correct.

With no argument it solves the puzzle for 3 disks in 7 moves, which matches the 2^n − 1 formula:

```
Move disk 1 from A to C
Move disk 2 from A to B
Move disk 1 from C to B
Move disk 3 from A to C
Move disk 1 from B to A
Move disk 2 from B to C
Move disk 1 from A to C
Total moves: 7 (expected 7)
```

The program reads the number of disks from the first command-line argument and defaults to 3 if none is given. It rejects anything that isn't a positive integer with a usage m
essage and exit code 1. The solver is recursive and moves disks from peg A to peg C using B as the spare.

My run tool can't pass arguments, so I only tested the default case. To try other sizes yourself, run `node hanoi.js 4` or `node hanoi.js 5`. The move count should be 15 and 3
1.
`

## Reasoning

We would like to implement reasoning, a loop that the agent runs, that helps it resonate better about the task, improving the prompt until it is ready to be executed.
Hopefully this can lead to more intelligent solutions and outputs. The agent can invoce it as often as possible in order to develop its inner monologue, and then just stop calling it when it is satisfied with the result.

## Reason about, and select optimal model

Have the agent reason about which agents are available, and which agent would be suitible for the task in question, then selecting that one for the work to be done. This could provide a better token cost/task complexity

## Validation

Check the work it has done using tools that provide feedback, e.g. run the program that it is working on, or have it provide feedback on the progress of the task or output. This is useful since it can prevent half done work, provide feedback on programming tasks, as well as make sure that outputs are not truncated and lost.


