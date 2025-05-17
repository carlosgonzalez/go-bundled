[![CI](https://github.com/carlosgonzalez/go-bundled/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/carlosgonzalez/go-bundled/actions/workflows/ci.yml)

# Simple Golang Setup with Echo, GORM, Testify and Mockery

## Description

This project aims to provide a straightforward setup for beginners who want to explore Golang.

Whether you’re new to programming or just diving into Go, this repository will help you get started with essential concepts and tools.

## Features

### Routing with Echo

Learn how to set up routes and handle HTTP requests using the Echo framework.

### Using GORM as an ORM

Understand how to interact with databases using GORM, a powerful Go ORM library.

### Concurrency with Channels and Goroutines

Explore Go’s concurrency model by using channels and goroutines effectively.

## Project Structure

```
go-bundled/
┣ .github/
┃ ┣ workflows/
┃ ┗ dependabot.yml
┣ internal/
┃ ┣ handlers/
┃ ┣ middlewares/
┃ ┣ models/
┃ ┣ repositories/
┃ ┗ services/
┣ mocks/
┃ ┗ mock_UserRepositoryInterface.go
┣ pkg/
┃ ┗ validators/
┣ tmp/
┃ ┣ build-errors.log
┃ ┗ main
┣ .air.toml
┣ .gitignore
┣ .mockery.yaml
┣ README.md
┣ docker-compose.yml
┣ go.mod
┣ go.sum
┗ main.go

```

### Testing

Learn how to write unit tests for your Go code to ensure reliability with Testify and Mockery.

### CI

Integrate linters (such as golangci-lint) to maintain code quality.

# Installation

Make sure you have Go installed on your system. If not, download it from [here](https://go.dev/doc/install).

Clone this repository:

```bash
git clone git@github.com:carlosgonzalez/go-bundled.git
```

Navigate to the project directory:

```bash
cd go-bundled
```

Install dependencies:

```bash
go mod tidy
```

Mockery

```bash
brew install mockery
```

Install `air` if you want to have live-reload capabilities

```bash
go install github.com/cosmtrek/air@latest
```

# Usage

Run the application:

```bash
docker-compose up -d
```

if you have air installed, run:

```bash
air
```

otherwise simply run with

```bash
go run .
```

Generate mocks with Mockery

```bash
mockery --all
```

# Contributing

Contributions are welcome! If you find any issues or want to enhance the project, feel free to submit a pull request.
