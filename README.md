# GAS

Github Account Switcher is a simple tool to switch between multiple GitHub accounts on the same machine. If you use your personal and work accounts on the same machine, you've probably felt the pain of switching between the two. This tool aims to make that process easier.

[![asciicast](https://asciinema.org/a/IhujYFPVonyFwRXvv5Z1cjhTS.png)](https://asciinema.org/a/IhujYFPVonyFwRXvv5Z1cjhTS)

> GAS is a work in progress. If you have any feedback or suggestions, please open an issue or a pull request.

> GAS is written for cross-platform compatibility, but installation scripts are only available for Unix-based systems at the moment. If you use Windows, you can install GAS manually by downloading the binary from the [releases page](https://github.com/style77/gas/releases) and adding it to your PATH manually.

## Installation

GAS installation is simple. Just run the following command:

```bash
curl -sSL https://raw.githubusercontent.com/style77/gas/master/scripts/install.sh | bash
```

GAS installs itself in the `/usr/local/bin` directory. It also adds the following line to your `.bashrc` or `.zshrc` file:

```bash
export PATH=$PATH:$HOME/.gas
```

This allows you to run the `gas` command from anywhere in your terminal.

> Keep in mind for ssh key generation you need to have `ssh-keygen` installed on your system and GAS does not install it for you.

### Config file

GAS stores your account details in the `~/.gas.yaml` file. You can edit this file directly to add or remove accounts.
If you remove this file, GAS will create a new once you run the `gas` command again.

## Usage

- Add a new account:

```bash
gas new
```

This will prompt you to enter your account details interactively.

- Switch between accounts:

```bash
gas switch
```

- Setup repo:

```bash
gas setup
```

This will set up remoteUrl for the current repo to use the selected account.

- Push/Commit to repo:

Since v1.1.0, you can use gas as your git command. This will confirm that you use the correct account before pushing or committing. GAS works as a wrapper around the git command, so you can use it as you would use git normally.

```bash
gas push
```

```bash
gas commit -m "your message"
```

Select the account you want to switch to from the list.

### Setting up different acronym

There is high possibility of 'GAS' being an acronym for different programs on your machine. e.g. GNU Assembler. 
You can set up a different acronym for GAS by changing the binary file name to something else e.g. `github-switcher` and then changing `PATH` in your `.bashrc` or `.zshrc` file to point to the new binary file name.

```bash
export PATH=$PATH:$HOME/github-switcher

# previously it was
# export PATH=$PATH:$HOME/gas
```

## Roadmap

- [x] Interactive add account
- [x] Switch between accounts interactively and by name/id
- [ ] Add account from command line
- [ ] Remove account
- [ ] List accounts
- [x] Add support for repo-specific accounts

## License
This project is licensed under the MIT License - see the LICENSE file for details.
