shex - a command-specific shell
====

shex is a "shell" for one specific command.

Use it like this:
```bash
$ shex git
shell for /usr/bin/git
>git 
```
You'll be prompted with a "shell" where you type arguments and the environment's command will be run with those arguments

Example:
```bash
>echo stuff
stuff
```

This is a very bare-bones command right now. I hope to add stuff like color and readline support in the future.
