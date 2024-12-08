# Santa (WIP)

Santa is a easy to use tool to help setup and generate a little bit of boilerplate for your Advent of Code.

## Commands

```
  config [--session <session_token>]    Set config options.
  init [--year <year>] [--force]        Initialize new AoC project directory
  new [day]                             Inside a project directory, create a day directory, downloads the input and creates initial solutions files.
```

## Roadmap

command: upgrade
  - [x] automatically upgrade Santa.

command: config
  - [x] store session (--session <session>)
  - [ ] switch between global and local configs

command: init
  - [x] choose year (--year <year>)
  - [x] force init on non-empty directory (--force)
  - [ ] automatically identify if target path is already a Santa project

command: new
 - [x] create day folder with boilerplate code and input (--day <day>)
 - [ ] add test input and flag to run the code against it
 - [ ] choose programming language
 - [ ] run solutions against input (directly on the aoc endpoints?)

command: input
 - [ ] download input for specific day in current folder
 - [ ] prompt the user if they currently are in a santa day folder (then download this day input instead of requiring a flag)
 - [ ] download test input for specific day in current folder

languages:
 - [x] Go
 - [ ] Javascript
 - [ ] Typescript (may require bun or deno?)
 - [ ] Python
 - [ ] Zig
