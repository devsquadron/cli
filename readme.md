### Installing production from source

- you will need [go](https://go.dev/doc/install) to build from source

- clone the repository
```
git clone https://github.com/devsquadron/project-manager.git
cd project-manager
```

- if you already have the repo cloned, switch to `mainline` and pull the latest
```
git checkout mainline
git pull
git pull --tags # optional
```

- make the executables, install globally, and install zsh completions
```
make prod
sudo make install
make zsh-completions
```
NOTE: you may need to restart shell to get full completions

- run the following to get started contributing!
```
ds list --tag dev
```
