fpath="/usr/local/share/zsh/site-functions"
binpath="/usr/local/bin"

build: clean prep
	go build -o ./ds cli/main.go

prep:
	go mod tidy
	rm -f constants/Version.txt
	go generate ./...

zsh-completions:
	sudo mkdir -p "${fpath}"
	ds --completion zsh > "./_ds"
	sudo mv "./_ds" "${fpath}/_ds"

all: clean prep build git-hooks

prod: clean prep
	go build -tags prod -o ./ds cli/main.go

prod-all: prod git-hooks-prod

git-hooks: 
	go build -o ./.git/hooks/prepare-commit-msg ./hooks/prepare-commit-msg/main.go
	go build -o ./.git/hooks/commit-msg ./hooks/commit-msg/main.go
	# go build -o ./.git/hooks/post-checkout ./hooks/post-checkout/main.go

git-hooks-prod: 
	go build -tags prod -o ./.git/hooks/prepare-commit-msg ./hooks/prepare-commit-msg/main.go
	go build -tags prod -o ./.git/hooks/commit-msg ./hooks/commit-msg/main.go
	# go build -tags prod -o ./.git/hooks/post-checkout ./hooks/post-checkout/main.go

install:
	mv ./ds "${binpath}/ds"

clean:
	rm -rf ./bin
	rm -f ./ds
	rm -f ./.git/hooks/prepare-commit-msg
	rm -f ./.git/hooks/commit-msg
	rm -f ./.git/hooks/post-checkout

uninstall:
	rm -f "${binpath}/ds"
	rm -f "${fpath}/_ds"

release:
	go run ./scripts/release/main.go
