# 🔑 nspm

Secure monolithic password manager written in Go

# ❓ About :

nspm is a secure, cross-platform command line password manager written in Go.

# 💡 Features:

1. Database encryption using AES-256
2. Categorization of entries
3. Unlimited number of entries across unlimited categories
4. Ability to move entries between categories
5. Password generator for creating strong passwords for entries

# 🖥️ Screenshots:
![image](https://github.com/regularenthropy/nanopm/assets/89523758/8dabb102-324e-4873-8d9b-4db841249817)



# 🛠️ Installation:
You can download nspm from the releases page on Github or compile it manually. To compile it manually, follow these steps:
```
Linux:  
    sudo make install
Windows:
    go build -a -gcflags=all="-l -B" -ldflags="-w -s"
```
# 📖 Manual:

1. Initialising nspm:

- By default, nspm creates a `.nspm` folder in the `$HOME` directory to hold all databases (except those passed as an argument).

- If no databases are found, nspm will prompt you to pass the `-n` argument to initialize a new database.

2. There are two different ways to select a database in nspm:
   - Run nspm with the `-f` flag and pass a database path as an argument.
   - Select a database from the menu that will be shown if nspm finds any databases in `$HOME/.nspm`.

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
