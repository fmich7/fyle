<div align="center">

# fyle

<!-- ![Build Status](https://img.shields.io/github/actions/workflow/status/fmich7/http/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fmich7/fyle)](https://goreportcard.com/report/github.com/fmich7/fyle)
![Test Coverage](https://img.shields.io/badge/test--coverage-90%25-blue) -->

**File upload service with built-in encryption 📂🕵🏻**

[Description](#-description) • [Features](#-features) • [Quick Start](#-quick-start) • [Usage](#-usage) • [Documentation](#-documentation) • [Testing](#%EF%B8%8F-testing)

</div>

## 📖 Description

_Fyle_ is a _file storage_ service written in **Go** that supports secure **file managment** with **AES-GCM** for **unique encryption**. It uses **JWT authentication** for user access and streams data in chunks to keep memory usage low. Sensitive data like credentials and encryption keys are safely stored in the **system's keyring**. Access the server through the **CLI client**. Application is tested with **unit and integration tests** with high codebase test coverage.

## ✨ Features

- **End-to-end encryption** Encrypt files with AES-GCM to keep user data secure and private
- **Security** Different encryptions for each uploaded file
- **Keyring** Safely store sensitive data like credentials and encryption keys in the system's keyring
- **Streaming files** Upload and download files in chunks to keep memory usage low
- **JWT Authentication** Provides user authentication on the server
- **CLI Client** Easily interact with server through Cli Client (see [📄Usage](#-usage))
- **Testing** Unit and integration tests, high codebase test coverage

## 🚀 Quick Start

### 1. Install CLI client

```bash
go install github.com/fmich7/fyle/cmd/fyle
```

### 2. Run server and PostgreSQL

```bash
docker-compose up
```

### 3. Use CLI client

```bash
fyle signup [username] [password]
```

> **[Optional]** Modify server configuration in `cmd/server/example_server.env`

## 📄 Usage

```bash
fyle ls                                # List all your files on the server
fyle upload [localPath] [serverPath]   # Upload file to the server
fyle download [serverPath] [localPath] # Download file from provided path
fyle signup [username] [password]      # Sign up a new user
fyle login [username] [password]       # Login to the server
```

> Run `fyle --help` to see all available commands and options.

## 💡 Documentation

Explore the full documentation for this package on

> [pkg.go.dev/github.com/fmich7/fyle](https://pkg.go.dev/github.com/fmich7/fyle#section-documentation)

## 🛠️ Testing

The application has been tested on the following platforms:

- Linux (Ubuntu 22.04, WSL2)
- Windows 10
- Go version 1.23.3

To run end-to-end integration test that verifies the entire workflow:

```bash
go test -v ./tests/integration_test.go
```

Run all tests:

```bash
#	go test -v -timeout 10s -race ./pkg/... ./tests/...
make test
# Generates test coverage file
#	go test -timeout 10s -race ./pkg/... ./tests/... -coverprofile=cover.out
make coverage
```
