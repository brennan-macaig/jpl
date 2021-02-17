# JPL
_Build software with rocket speed._

## What?
A simple build tool which just runs a list of commands in order.
No garbage syntax, no crazy assumptions, no BS.

## Why?
I got frustrated with existing build tools:

- Make is impossible to read, and incredibly complex
- Gradle depends on Java (??!?!)
- Bash sucks

## How?

By default, JPL assumes you have a YAML file at the top-level of your project called
`jpl.yml` (you can change this). It reads the YAML file, and executes the commands,
one at a time.

For now, JPL has two "stages": test and build. You specify on the command line which
you want. For build, use `-b` or no-opt, and for test use `-t`. If `-t` is given, it
will only run the tests and then exit.

If a command returns non-zero, JPL will stop trying to process your build and exit.

## build.yml reference

### Structural Overview
Every `build.yml` file that JPL has the same basic structure:

- **version** - what version of the build file is this? (Only 1.0 for now)
- **conmfig** - some configuration options for JPL
- **variables** - variables you'd like to append to whatever was taken (or not)
from the OS environment
- **tasks** - The tasks that need to run (build and test only) 

Tasks are made out of collections of modules, whose reference can be found below.
However, it should be noted that at current only two modules are supported:
`execute` and `copy`. More modules will be added on request, or if they are important.

### Modules
Modules live inside tasks (tasks being only `build` and `test`), and are executed from
top of the file to the bottom. Modules must look something like this:
```yaml
- module: moduleName
  param: someParameter
  maybeAList:
    - list item 1
    - list item 2
  anotherParam: anotherParam1
```
Modules each have their own defined parameters (check their reference below for more detail).

### Task Ordering
Tasks (`build` and `test`) are run from top of the config to bottom. That means that
each module is read, line by line, and run from the first provided to the last.
No finicky garbage, no weird build ordering. Tasks, and modules in tasks, are run from
the top of the file to the bottom (just like source code).


### Module: Execute
Execute is likely where the bulk of the work that JPL will be doing is found.
Execute works similarly to a bash script, it runs the commands given one at a time,
from top of the file to the bottom.

```yaml
- module: execute # Module name, must be exactly "execute"
  commands: # Commands to run
    - command1 # Commands are run in order
    - command2 # From top to bottom
    - ... # You can have as many commands as you want
```

### Module: Copy
Copy is a very simple module: copy the files given to a new location. If multiple files
are given, the files will be shoved into `dest/sourceFileName`. Otherwise, if only
one file is given, it will be copied exactly to `dest`.

```yaml
- module: copy # Module name, must be exactly "copy"
  src: # Any number of files to copy
    - file1
    - file2
  dest: destination # Where to copy files to (for multiple files, it will be destination/fileName)
```