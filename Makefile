# Define the binary name and service name
BINARY_NAME = gitmonitor
SERVICE_NAME = gitmonitor
SERVICE_FILE = ./$(SERVICE_NAME).service

# Define the build directory and source directory
SRC_DIR = $(shell pwd)
BUILD_DIR = $(SRC_DIR)/build

# Define the Go build command
GO_BUILD = go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/cmd/autodeploy/main.go

# Define the service file content
SERVICE_CONTENT = "[Unit]\nDescription=Git Monitor Service\nAfter=network.target\n\n[Service]\nExecStart=$(BUILD_DIR)/$(BINARY_NAME)\nRestart=always\nUser=$(USER)\nGroup=$(USER)\nEnvironment=GO_ENV=production\nWorkingDirectory=$(SRC_DIR)\n\n[Install]\nWantedBy=multi-user.target\n"

# Default target: build the binary and install the service
all: build install

# Build the Go application
build:
	@echo "Building the Go application..."
	$(GO_BUILD)
	@echo "Build completed!"

# Install the systemd service
install:
	@echo "Installing the systemd service..."
	@echo $(SERVICE_CONTENT)
	@sudo systemctl daemon-reload
	@sudo systemctl enable $(SERVICE_NAME)
	@sudo systemctl start $(SERVICE_NAME)
	@echo "Service installed and started!"

# Clean up the build files and remove the service
clean:
	@echo "Cleaning up..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@sudo systemctl stop $(SERVICE_NAME) || true
	@sudo systemctl disable $(SERVICE_NAME) || true
	@sudo rm -f $(SERVICE_FILE)
	@sudo systemctl daemon-reload
	@echo "Clean up completed!"

# Restart the service
restart:
	@sudo systemctl restart $(SERVICE_NAME)
	@echo "Service restarted!"

# Check the status of the service
status:
	@sudo systemctl status $(SERVICE_NAME)

.PHONY: all build install clean restart status
