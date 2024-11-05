# GAS

Github Account Switcher is a simple tool to switch between multiple GitHub accounts on the same machine. If you use your personal and work accounts on the same machine, you've probably felt the pain of switching between the two. This tool aims to make that process easier.

## Installation

GAS installation is simple. Just run the following command:

```bash
curl -s https://raw.githubusercontent.com/style77/gas/master/scripts/install.sh | bash
```

GAS installs itself in the `~/.gas` directory. It also adds the following line to your `.bashrc` or `.zshrc` file:

```bash
export PATH=$PATH:$HOME/.gas
```

This allows you to run the `gas` command from anywhere in your terminal.

(If you use windows, GAS will be installed in the `C:\Users\<username>\.gas` directory.)

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

Select the account you want to switch to from the list.

## Roadmap

- [x] Interactive add account
- [x] Switch between accounts interactively and by name/id
- [ ] Add account from command line
- [ ] Remove account
- [ ] List accounts
- [ ] Add support for repo-specific accounts

## License
This project is licensed under the MIT License - see the LICENSE file for details.
