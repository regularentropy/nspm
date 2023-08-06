# 🔑 nanopm

Secure monolithic password manager written in Go

# ❓ About :

nanopm is a command line password manager written in Go.  
It works on all Linux distributions that have GoLang in repositories.  
Unfortunately, it doesn't currently support Windows, but I'm working on it.
Stay tuned :)

# 💡 Features:

1. Database encryption via AES-256
2. Categories for entries
3. Ability to have an unlimited number of entries between unlimited categories
4. Ability to move entries between categories
5. Ability to generate strong password for the entries

# 🛠️ Installation:
```
Linux:  
    sudo make install
Windows:
    Coming soon
```
# 📖 Manual:

1. Initialising nanopm:
   - By default, nanopm creates a .nanopm folder in the $HOME directory, which will hold all databases (except those passed as an argument).
   - If no databases are found, nanopm will ask you to pass the "-n" argument to initialize a new database.

3. There are two different ways to select a database in nanopm:
   - Run nanopm with the "-f" flag and pass a database path as an argument.
   - Select a database from the menu that will be shown if nanopm finds any databases in $HOME/.nanopm

# 👤 Authors:

regularenthropy - main developer  
Contributors are welcome!

# ✅ TODO
- [X] Add Windows support
- [ ] Refactor codebase to look better
- [ ] Create a better name

License:
--------
GPLv3 - See [LICENSE](/LICENSE)
