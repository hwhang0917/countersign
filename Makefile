DIST_DIR=dist
BINARY_NAME=countersign

build:
	mkdir -p $(DIST_DIR)
	go build -o ./$(DIST_DIR)/$(BINARY_NAME) main.go

clean:
	go clean
	rm -rf $(DIST_DIR)
