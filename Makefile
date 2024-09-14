# Directories
SERVER_DIR := server
CLIENT_DIR := client
DIST_DIR := dist
PUBLIC_DIR := $(DIST_DIR)/public
BINARY_NAME := countersign

# Go build flags
GOFLAGS := -v

# NPM command (use npm or yarn as needed)
NPM := npm

.PHONY: all build clean server client

all: clean build

build: client server

clean:
	rm -rf $(DIST_DIR)
	cd $(CLIENT_DIR) && $(NPM) run clean

server:
	mkdir -p $(DIST_DIR)
	cd $(SERVER_DIR) && go build $(GOFLAGS) -o ../$(DIST_DIR)/$(BINARY_NAME) main.go

client:
	cd $(CLIENT_DIR) && $(NPM) install
	cd $(CLIENT_DIR) && $(NPM) run build
	mkdir -p $(PUBLIC_DIR)
	rm -rf $(PUBLIC_DIR)/*
	mv $(CLIENT_DIR)/dist/* $(PUBLIC_DIR)/

run: build
	cd $(DIST_DIR) && ./$(BINARY_NAME)

dev:
	make -j2 dev-client dev-server

dev-client:
	cd $(CLIENT_DIR) && $(NPM) run dev

dev-server:
	cd $(SERVER_DIR) && go run main.go
