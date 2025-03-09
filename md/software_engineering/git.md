# Git
```bash
git config --global user.name "Mona lisa"
```

```bash
git init
git add .
git commit -m "First commit"
```

```bash
# remove .DS_Store file from GitHub that MAC OS X creates

# find and remove .DS_Store
find . -name .DS_Store -print0 | xargs -0 git rm -f --ignore-unmatch

# create .gitignore file, if needed
touch .gitignore
echo .DS_Store > .gitignore

# push changes to GitHub
git add .gitignore
git commit -m '.DS_Store removed'
git push origin master

```

```bash
# Create a new branch from the top of the main branch.
git checkout -b my-feature-branch main
```

```xml
# Delete Local branch
git branch -d <branchName>
git branch -D <branchName>
```

```bash
# Delete Remote Branch
git push origin --delete <branchName>
```

## Commit Messages
- ✨ :sparkles: - New feature
- 🔧 :wrench: - Fixes or improvements
- 🐛 :bug: - Bug fixes
- 📚 :books: - Documentation changes
- 🔥 :fire: - Removing code or features
- ✅ :white_check_mark: - Tests
- 📝 :memo: - Writing or updating documentation
- 🚀 :rocket: - Deployments
- 📦 :package: - Package updates
- 💄 :lipstick: - UI improvements