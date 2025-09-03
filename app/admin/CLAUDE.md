# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go microservice project based on the Kratos framework. It follows a clean architecture pattern with layered components:

1. **API Layer**: Defines service interfaces using Protocol Buffers (proto files)
2. **Service Layer**: Implements the gRPC and HTTP services that handle requests
3. **Biz Layer**: Contains business logic and use cases
4. **Data Layer**: Manages data access and persistence
5. **Server Layer**: Configures and initializes HTTP/gRPC servers

## Common Commands

### Development Setup
```bash
# Initialize development environment (install dependencies)
make init

# Generate code from proto files
make api

# Generate all files (API + wire dependency injection)
make all

# Run the application
make run
```

### Building
```bash
# Build the application
make build

# Clean build artifacts
make clean
```

### Code Generation
```bash
# Generate code from proto files using buf
make api

# Generate wire dependency injection
make wire

# Generate all files
make all
```

## Project Structure

- `api/` - Contains proto files and generated code
- `app/admin/` - Main application with cmd, configs, and internal packages
- `app/admin/cmd/` - Application entry point with main.go and wire.go
- `app/admin/internal/` - Internal application components:
  - `biz/` - Business logic layer
  - `conf/` - Configuration files and structs
  - `data/` - Data access layer
  - `server/` - Server configuration (HTTP/gRPC)
  - `service/` - Service implementation layer
- `pkg/` - Shared packages across applications
- `third_party/` - Third-party proto dependencies

## Architecture Patterns

The project follows a layered architecture with dependency injection using Google Wire:

1. **Entry Point**: `app/admin/cmd/main.go` initializes configuration, logging, and starts the application
2. **Dependency Injection**: Google Wire is used to manage dependencies between layers
3. **Service Layer**: Implements gRPC services and HTTP handlers
4. **Business Logic**: Usecases in the biz layer contain core business logic
5. **Data Access**: Repository pattern in the data layer for data operations

## Code Generation

The project uses buf (https://buf.build) for Protocol Buffer management and code generation:

- Proto files are in `api/proto/`
- Generated code goes to `api/gen/go/`
- API documentation is generated to `docs/`
- 新增或修改`api/proto`下的proto文件后在`app/admin/`目录下执行`make api`命令来生成对应的代码，对应的代码在`api/gen/go/`下。同时生成`openapi.yaml`在`docs/`目录下
- 新增或修改`internal/conf`下的proto文件后在`app/admin/`目录下执行`make config`命令来生成对应的代码
- 新增或修改`internal/data/ent/shema`下的文件后在`app/admin/`目录下执行`make ent`命令来生成对应的代码

## Configuration

- Application configuration files are in `app/admin/configs/`
- Configuration structure is defined in `app/admin/internal/conf/conf.proto`