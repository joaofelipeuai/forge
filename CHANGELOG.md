# Changelog - Forge Framework

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-07-16

### 🎉 Initial Release

#### Added
- **Core Framework**
  - HTTP server with Express.js-like API
  - Route handling with parameters (`:id`)
  - Query parameter support
  - JSON/HTML/String response methods
  - Middleware pipeline system

- **Built-in Middleware**
  - Recovery middleware for panic handling
  - Logger middleware with request timing
  - CORS middleware with configurable headers
  - Rate limiting with IP-based tracking
  - JWT authentication with token validation

- **Advanced Features**
  - WebSocket support with broadcasting
  - Template engine with hot reload
  - File upload with validation
  - Hot reload for development
  - Graceful shutdown support

- **Security & Validation**
  - Input validation system
  - JWT token generation and validation
  - File type validation for uploads
  - Rate limiting protection

#### Fixed
- **Critical Security Issues**
  - Race condition in WebSocket broadcaster (thread safety)
  - Memory leak in rate limiter (goroutine lifecycle)
  - Input validation vulnerabilities

#### Technical Improvements
- Thread-safe WebSocket operations
- Memory-efficient rate limiting
- Comprehensive input validation
- Organized project structure
- Complete documentation

### 🏗️ Project Structure
- Followed Go Standard Project Layout
- Organized middleware and examples
- Centralized dependencies (single go.mod)
- Comprehensive documentation

### 📚 Documentation
- Complete README with examples
- API documentation
- Middleware guide
- Getting started guide
- Technical debt documentation
- Contributing guidelines

### 🧪 Testing
- Unit tests for core functionality
- Integration tests for middleware
- Example applications
- Validation examples

### 🎯 Performance
- Zero external dependencies
- Optimized for high concurrency
- Efficient memory usage
- Fast route matching with regex

---

## [Unreleased]

### Planned Features
- Configurable timeouts
- Structured logging
- Metrics collection
- Response compression
- Database integration helpers
- GraphQL support

---

## Version History

- **v1.0.0** - Initial stable release with all core features
- **v0.x.x** - Development versions (not released)

---

## Migration Guide

This is the first stable release, so no migration is needed.

For future versions, migration guides will be provided here.

---

## Contributors

- [João Felipe Souza](https://github.com/joaofelipeuai) - Creator and maintainer

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.