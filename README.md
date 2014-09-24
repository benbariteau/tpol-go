T'Pol - a command-specific shell
====

T'Pol is a "shell" for one specific command, allowing you to run subcommands or give it arguments without having to type the command first each time.

Use it like this:
```
$ tpol git
shell for /usr/bin/git
>git 
```
You'll be prompted with a "shell" where you type arguments and the environment's command will be run with those arguments

Example:
```bash
>echo stuff
stuff
```

Features:
* Subcommand history (per session)
* Bash completion (using linenoise)
* Command escape (just type `!command with args`)

Coming:
* Session-spanning (command specific) history
* Filename tab completion
* Readline support
