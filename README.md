<div align="center">

# fyle

<!-- ![Build Status](https://img.shields.io/github/actions/workflow/status/fmich7/http/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fmich7/fyle)](https://goreportcard.com/report/github.com/fmich7/fyle)
![Test Coverage](https://img.shields.io/badge/test--coverage-90%25-blue) -->

**File upload service with built-in encryption 📂🕵🏻**

[Description](#-description) • [Features](#-features) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [Testing](#%EF%B8%8F-testing)

</div>

## 📖 Description

_Fyle_ is a _file storage_ service written in **Go** that supports securely **file managment** with **AES-GCM** and a **nonce** for **unique encryption**. It uses **JWT authentication** for user access and streams data in chunks to keep memory usage low. Sensitive data like credentials and encryption keys are safely stored in the **system's keyring**.

## ✨ Features

- **End-to-end encryption**
  - Generate encryption key from user credentials and salt (from DB)
  - Store encryption key in system's keyring
  - Uniquely encrypt files with AES-GCM, nonce and the key
- **Streaming files**
  - Supports sending and receiving files in chunks (4kb chunks)
  - Encrypt and decrypt on the fly
- **JWT Authentication**
  - Token-based authentication using JWT (stored in keyring)
  - Ensures that only authorized users can access their files
- **CLI Client**
  - Easily interact with the server using a CLI client (see [📄Usage](#-features))
- **Testing**
  - Almost the whole codebase is covered with tests
  - Uses afero for filesystem mocking

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
fyle signup <username> <password>
```

> **[Optional]** Modify server configuration in `cmd/server/example_server.env`

## 📄 Usage

```bash
fyle ls                               # List all your files on the server
fyle upload <file_path>				  # Upload file to the server
fyle download <file_id> <output_path> # Download file from provided path
fyle signup <username> <password>     # Sign up a new user
fyle login <username> <password>      # Login to the server
```

> Run `fyle --help` to see all available commands and options.

## 💡 Documentation

Explore the full documentation for this package on

> [pkg.go.dev/github.com/fmich7/fyle](https://pkg.go.dev/github.com/fmich7/fyle#section-documentation)

## 🛠️ Testing

```bash
# go test -v -timeout 5s -race ./...
make test
# Generate test coverage file
make coverage
```
