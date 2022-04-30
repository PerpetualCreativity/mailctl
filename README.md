# mailctl (in alpha)

**modern console-based email app** (not TUI!)

Thanks for checking out `mailctl`!

# installation

Currently, as `mailctl` is in alpha, the only way to install it is through `go`:

```sh
go install github.com/PerpetualCreativity/mailctl@latest
```

# config

Run `mailctl init` to put a sample configuration file at `~/.mailctl.yml`. Alternatively, you can [view the sample in GitHub](/cmd/mailctl.yml).

# usage

The name of the executable is `mailctl` (you may wish to alias this). You can specify which account to use for any command with `-a 1`

| command                            | supported? | function                                                                                                                                   |
|------------------------------------|------------|--------------------------------------------------------------------------------------------------------------------------------------------|
| `list [folder] [number]`           | yes        | Lists emails in folder (if folder is not specified, a select menu will appear). By default, displays last 10 messages.                     |
| `view id [folder]`                 | yes        | View an email or a draft. Requires email or draft id. Uses Inbox folder by default.                                                        |
| `new`                              | yes        | Write a new email.                                                                                                                         |
| `reply id [folder]`                | yes        | Write a reply draft to an email. Requires email id.                                                                                        |
| `edit id [folder]`                 | yes        | Edit a draft. Requires draft id.                                                                                                           |
| `send id`                          | yes        | Send a draft. Requires draft id.                                                                                                           |
| `move id currentFolder [toFolder]` | yes        | Move an email or draft to a different folder. Requires email or draft id and current folder. Will prompt for destination if not specified. |
| `search`                           | no         | Search for an email. Syntax not defined yet.                                                                                               |

