GO ?= go
GOLINT ?= golangci-lint
SRCS := $(wildcard $(CURDIR)/*.go) $(wildcard $(CURDIR)/**/*.go) $(wildcard $(CURDIR)/**/**/*.go)

TARGET := snake-go

default: all

all: $(TARGET) $(SRCS)

$(TARGET): $(SRCS)
	$(GO) build -o $(TARGET) -gcflags="-e"

tidy:
	go mod tidy \;

clean:
	rm -rf $(TARGET)
