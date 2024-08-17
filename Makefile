# Define the binary name and service name
BINARY_NAME = gitmonitor
SERVICE_NAME = gitmonitor
SERVICE_FILE = /etc/systemd/system/$(SERVICE_NAME).service

# Define the build directory and source directory
BUILD_DIR = ./build
SRC_DIR = .

# Define the Go build command
GO_BUILD = go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go

# Define the service file content
SERVICE_CONTENT = "[Unit]\n\
Description=Git Monitor Service\n\
After=network.target\n\n\
[Service]\n\
ExecStart=$(BUILD_DIR)/$(BINARY_NAME)\n\
Restart=always\n\
User=$(USER)\n\
Group=$(USER)\n\
Environment=GO_ENV=production\n\
WorkingDirectory=$(SRC_DIR)\n\n\
[Install]\n\
WantedBy=multi-user.target\n"

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
	@echo -e $(SERVICE_CONTENT) | sudo tee $(SERVICE_FILE) > /dev/null
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
