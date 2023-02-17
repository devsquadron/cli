### Installing production from source

- you will need [go](https://go.dev/doc/install) to build from source

- clone the repository
```
git clone https://github.com/devsquadron/ds.git
cd ds
```

- if you already have the repo cloned, switch to `mainline` and pull the latest
```
git checkout mainline
git pull
git pull --tags # optional but get's version number for debugging purposes
```

- make the executables, install globally, and install zsh completions
```
make prod
sudo make install
make zsh-completions
```
NOTE: you may need to restart shell to get full completions

- use the [ds list](https://developersquadron.com/getting-started/listing-tasks/) command to start understanding the queue
