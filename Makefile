app_name = toshiki-rssbot

VERSION=$(shell git describe --tags --always)
DATA=$(shell date)
COMMIT=$(shell git rev-parse --short HEAD)
test:
	go test ./... -v

all: build

build: get
	go build -trimpath -ldflags \
	"-s -w -buildid= \
	-X 'github.com/andatoshiki/toshiki-rssbot/internal/config.commit=$(COMMIT)' \
	-X 'github.com/andatoshiki/toshiki-rssbot/internal/config.date=$(DATA)' \
	-X 'github.com/andatoshiki/toshiki-rssbot/internal/config.version=$(VERSION)'" -o $(app_name)

get:
	go mod download

run:
	go run .

clean:
	rm toshiki-rssbot

format:
	gofmt -s -w .

token:
	bash ./telegraph.sh