# gau command is not working

Trying to use gau what the command is not working even though it 
is installed correctly.

When running gau you get this error
```bash
fatal: not a git repository (or any of the parent directories): .git
```

The message is not related to gau. Use `which`
```bash
which gau
```
Result
```bash
gau: aliased to git add --update
```

The gau command is an alias to a git command.

The alias can be removed using 

```bash
unalias gau
```
This only remove the alias for that session. 
To completely remove the alias the alias need to be remove from the config file.
Most of the time it should be in `~/.zshrc` or `.zprofile` but this gau alias is not there.

After some digging, I found that Oh My Zsh was the culprit. 
It uses a Git plugin that sets a bunch of aliases, including gau. 
The file is located at `~/.oh-my-zsh/plugins/git/git.plugin.zsh`, and simply uncommenting the line for the gau alias did the trick.