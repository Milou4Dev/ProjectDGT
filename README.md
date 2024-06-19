# ProjectDGT

## Overview

ProjectDGT is a tool designed to set your Discord online status. This project leverages Docker for easy deployment and
management. The main functionality is implemented in Go, and it uses the Gorilla WebSocket package for real-time
communication.

## Features

- Set Discord online status
- Docker support for easy deployment
- Real-time communication using WebSockets

## Getting Started

### Prerequisites

- Docker installed on your machine
- Basic understanding of Docker and Go

### Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/Milou4Dev/ProjectDGT.git
   cd ProjectDGT
   ```

2. **Build and Run the Docker Container:**
   ```sh
   docker buildx build --pull --no-cache --tag projectdgtimg . --load && docker run --rm -it --init projectdgtimg
   ```

## Usage

Once the Docker container is running, follow the on-screen prompts to set your Discord online status.

## Project Structure

```
ProjectDGT/
├── LICENSE
├── README.md
├── dockerfile
├── go.mod
├── go.sum
└── main.go
```

- **LICENSE:** Contains the license information for the project.
- **README.md:** This file, providing an overview and instructions.
- **dockerfile:** Instructions for creating a Docker image.
- **go.mod:** Tracks the dependencies of the project.
- **go.sum:** Verifies the integrity of the dependencies.
- **main.go:** The main source code file of the project.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgements

- [Gorilla WebSocket](https://github.com/gorilla/websocket) for the WebSocket implementation.