# ProjectDGT

## Overview

ProjectDGT is a tool designed to keep your Discord status online 24/7. This project leverages Docker for easy deployment
and management. The main functionality is implemented in Go, and it uses the Gorilla WebSocket package for real-time
communication.

## Features

- Set and maintain Discord online status
- Docker support for easy deployment
- Real-time communication using WebSockets
- Environment variable support for hosting on platforms like Railway and Replit

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
   docker buildx build --pull --no-cache -t projectdgtimg . && \
   docker run --rm -it --init projectdgtimg
   ```

### Environment Variables

To configure the environment variables, you can set the `USE_CONFIG` variable to `on` and fill in the other required
variables as specified in the code.

## Usage

Once the Docker container is running, follow the on-screen prompts to set your Discord online status. Note: If you are
not familiar with configuring and using emojis, it is recommended to turn them off to avoid potential issues.

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