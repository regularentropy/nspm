# 🔑 nanopm

Secure monolithic password manager written in Go

# ❓ About :

nanopm is a secure, cross-platform command line password manager written in Go.

# 💡 Features:

1. Database encryption using AES-256
2. Categorization of entries
3. Unlimited number of entries across unlimited categories
4. Ability to move entries between categories
5. Password generator for creating strong passwords for entries

# 🖥️ Screenshots:
![image](https://github.com/regularenthropy/nanopm/assets/89523758/ace53eee-a396-4009-8eb8-1731be01e072)

# 🛠️ Installation:
You can download Nanopm from the releases page on Github or compile it manually. To compile it manually, follow these steps:
```
Linux:  
    sudo make install
Windows:
    go build -a -gcflags=all="-l -B" -ldflags="-w -s"
```
# 📖 Manual:

1. Initialising nanopm:

- By default, nanopm creates a `.nanopm` folder in the `$HOME` directory to hold all databases (except those passed as an argument).

- If no databases are found, Nanopm will prompt you to pass the `-n` argument to initialize a new database.

2. There are two different ways to select a database in nanopm:
   - Run nanopm with the `-f` flag and pass a database path as an argument.
   - Select a database from the menu that will be shown if nanopm finds any databases in `$HOME/.nanopm`.

# 👤 Authors:

regularenthropy - main developer  
Contributors are welcome!

# ✅ TODO
- [X] Add Windows support
- [X] Refactor codebase to look better
- [ ] Create a better name

License:
--------
GPLv3 - See [LICENSE](/LICENSE)
