# Santa (WIP)

Santa is a easy to use tool to help setup and generate a little bit of boilerplate for your Advent of Code.

## Commands

```
  config [--session <session_token>]    Set config options.
  init [--year <year>] [--force]        Initialize new AoC project directory
  new [day]                             Inside a project directory, create a day directory, downloads the input and creates initial solutions files.
```

## Roadmap

upgrade
  - [ ] automatically upgrade Santa.

config
  - [x] store session (--session <session>)
  - [ ] switch between global and local configs

init
  - [x] choose year (--year <year>)
  - [x] force init on non-empty directory (--force)
  - [ ] automatically identify if target path is already a Santa project

new
 - [x] choose day (--day <day>)
 - [ ] choose programming language
     - [ ] Go
     - [ ] Javascript
     - [ ] Typescript (may require bun or deno?)
     - [ ] Python
 - [ ] run solutions against input (directly on the aoc endpoints?)
