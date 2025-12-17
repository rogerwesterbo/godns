# Binary Installation Guide

This guide explains how to install and run GoDNS using the pre-compiled binaries available on the [Releases](https://github.com/rogerwesterbo/godns/releases) page.

## Prerequisites

GoDNS requires a **Valkey** (or Redis) server for data storage.

- **Valkey/Redis**: Version 6.0 or higher recommended.

## Step 1: Download the Binary

1.  Go to the [Releases](https://github.com/rogerwesterbo/godns/releases) page.
2.  Download the archive for your operating system and architecture:
    - **Linux**: `godns-<version>-linux-amd64.tar.gz` (or `arm64`)
    - **macOS**: `godns-<version>-darwin-amd64.tar.gz` (or `arm64` for Apple Silicon)
    - **Windows**: `godns-<version>-windows-amd64.exe`

## Step 2: Extract and Install

### Linux / macOS

```bash
# Extract the archive
tar -xzf godns-*.tar.gz

# Move the binary to a location in your PATH (optional)
sudo mv godns /usr/local/bin/
```

### Windows

Extract the `.zip` file and place `godns.exe` in a folder of your choice.

## Step 3: Configuration

GoDNS is configured using environment variables. You can set these in your shell or create a `.env` file in the same directory as the binary.

### Minimal Configuration (No Auth)

If you just want to run the DNS server without Keycloak authentication:

```bash
export VALKEY_HOST=localhost
export VALKEY_PORT=6379
export AUTH_ENABLED=false
```

### Full Configuration

See the table below for common configuration options:

| Variable          | Default     | Description                                      |
| :---------------- | :---------- | :----------------------------------------------- |
| `DNS_SERVER_PORT` | `:53`       | Port for the DNS server (UDP/TCP)                |
| `HTTP_API_PORT`   | `:8080`     | Port for the REST API                            |
| `VALKEY_HOST`     | `localhost` | Hostname of the Valkey/Redis server              |
| `VALKEY_PORT`     | `6379`      | Port of the Valkey/Redis server                  |
| `AUTH_ENABLED`    | `true`      | Enable/Disable Keycloak authentication           |
| `LOG_LEVEL`       | `info`      | Logging level (`debug`, `info`, `warn`, `error`) |

## Step 4: Running the Server

### Linux / macOS

Note: Binding to port 53 usually requires root privileges.

```bash
# Run with sudo if using port 53
sudo VALKEY_HOST=localhost ./godns
```

Or if you installed it to `/usr/local/bin`:

```bash
sudo godns
```

### Windows

Open PowerShell or Command Prompt as Administrator (if using port 53) and run:

```powershell
.\godns.exe
```

## Running as a Service (Linux Systemd)

To run GoDNS as a background service on Linux, create a systemd unit file.

1.  Create the file `/etc/systemd/system/godns.service`:

```ini
[Unit]
Description=GoDNS Server
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/godns
Restart=on-failure
# Environment variables
Environment="VALKEY_HOST=localhost"
Environment="AUTH_ENABLED=false"
# Or load from file
# EnvironmentFile=/etc/godns/godns.env

[Install]
WantedBy=multi-user.target
```

2.  Reload systemd and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now godns
```

3.  Check status:

```bash
sudo systemctl status godns
```
