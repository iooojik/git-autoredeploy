# Git Auto Redeploy

Git Auto Redeploy is a Go application that monitors specified Git repositories for changes and automatically runs custom commands when changes are detected. It is designed to be lightweight and efficient, making it suitable for use in automated deployment pipelines.

## Features

- Monitors multiple Git repositories for updates.
- Executes custom commands when updates are detected.
- Configurable check intervals and commands.
- Supports YAML configuration files.
- Can be set up as a daemon to run on system startup.

## Requirements

- Go 1.23 or later
- Git

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/git-autoredeploy.git
    cd git-autoredeploy
    ```

2. **Build the application:**

    ```bash
    go build -o build/git-autoredeploy ./cmd/autodeploy/main.go
    ```

3. **Configuration:**

   Create a configuration file named `config.yaml` in the `configs` directory. Here's an example configuration:

    ```yaml
    projects:
      - name: project1
        repo: "git@github.com:username/project1.git"
        directory: "/path/to/project1"
        command: "bash deploy.sh"
      - name: project2
        repo: "git@github.com:username/project2.git"
        directory: "/path/to/project2"
        command: "python deploy.py"
    check_interval: 60
    ```

4. **Run the application:**

    ```bash
    ./build/git-autoredeploy -configDir ./configs
    ```

## Usage

- **Command-line flags:**
    - `-configDir`: Specifies the directory where the configuration file is located. Default is the current directory.

- **Daemon Setup:**

  To set up the application as a daemon on Linux:

    1. Create a systemd service file:

       ```ini
       [Unit]
       Description=Git Auto Redeploy Service
       After=network.target
  
       [Service]
       ExecStart=/path/to/git-autoredeploy -configDir /path/to/configs
       Restart=always
       User=yourusername
       Group=yourusergroup
       WorkingDirectory=/path/to/working/directory
  
       [Install]
       WantedBy=multi-user.target
       ```

    2. Reload systemd and start the service:

       ```bash
       sudo systemctl daemon-reload
       sudo systemctl enable git-autoredeploy
       sudo systemctl start git-autoredeploy
       ```
       
## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss any changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
