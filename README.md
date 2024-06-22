# ProjectDGT

## Overview

ProjectDGT is a powerful tool designed to maintain your Discord status online 24/7. This project offers flexibility in deployment, with Docker support for seamless setup and management, as well as standalone executable files for those who prefer a more traditional approach. The core functionality is implemented in Go, leveraging the Gorilla WebSocket package for efficient real-time communication.

## Features

- Set and maintain Discord online status continuously
- Docker support for easy deployment and cross-platform compatibility
- Standalone executable files available for Linux (AMD64, ARM64) and Windows (AMD64, ARM64)
- Real-time communication using WebSockets for responsive status updates
- Flexible configuration options: environment variables or interactive prompt
- User-friendly interface with clear prompts and instructions
- Optional emoji support (configurable)

## Getting Started

### Prerequisites

- For Docker deployment:
  - Docker installed on your machine
  - Basic understanding of Docker
- For standalone execution:
  - A compatible system: Linux (AMD64 or ARM64) or Windows (AMD64 or ARM64)

### Installation

#### Option 1: Docker (Recommended)

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

#### Option 2: Standalone Executable

1. Visit the [Releases](https://github.com/Milou4Dev/ProjectDGT/releases) page of the ProjectDGT repository.
2. Download the appropriate executable for your system:
   - Linux: Choose between `amd64` or `arm64` versions
   - Windows: Choose between `amd64` or `arm64` versions
3. Run the executable file directly on your system.

## Configuration

ProjectDGT offers two methods of configuration:

1. Environment Variables (Advanced)
2. Interactive Prompt (Beginner-Friendly)

### Method 1: Environment Variables

ProjectDGT uses environment variables for configuration. The following variables are necessary for the application to function properly:

- `TOKEN`: Your Discord token
- `STATUS`: The status you want to set (e.g., "online", "idle", "dnd")
- `CUSTOM_STATUS`: Your custom status message
- `USE_EMOJI`: Set to "true" if you want to use an emoji in your status, "false" otherwise
- `EMOJI_NAME`: The name of the emoji (if `USE_EMOJI` is "true")
- `EMOJI_ID`: The ID of the emoji (if `USE_EMOJI` is "true")

To set these environment variables:

1. For Docker deployment:
   - You can pass these variables when running the Docker container using the `-e` flag:
     ```
     docker run --rm -it --init \
       -e TOKEN=your_token \
       -e STATUS=your_status \
       -e CUSTOM_STATUS=your_custom_status \
       -e USE_EMOJI=false \
       -e EMOJI_NAME=your_emoji_name \
       -e EMOJI_ID=your_emoji_id \
       projectdgtimg
     ```

2. For standalone executable:
   - Set these variables in your system environment before running the executable.
   - On Windows, you can use the `setx` command in Command Prompt:
     ```
     setx TOKEN your_token
     setx STATUS your_status
     ```
   - On Linux or macOS, you can export these variables in your shell:
     ```
     export TOKEN=your_token
     export STATUS=your_status
     ```

Note: Make sure to replace the placeholder values with your actual Discord token, desired status, and other configuration details.

### Method 2: Interactive Prompt

If you're not familiar with setting environment variables, don't worry! ProjectDGT provides an interactive setup process:

1. Simply run the application (either via Docker or the standalone executable).
2. The program will detect that environment variables are not set.
3. You'll be prompted to enter the necessary configuration details:
   - Your Discord token
   - Desired status (e.g., "online", "idle", "dnd")
   - Custom status message
   - Whether to use an emoji (yes/no)
   - If yes to emoji, the emoji name and ID

This method is perfect for beginners or those who prefer a guided setup process.

## Usage

1. Choose your preferred configuration method:
   - If using environment variables, ensure all necessary variables are set as described in the "Environment Variables" section.
   - If preferring the interactive method, simply start the application and follow the prompts.
2. Start ProjectDGT using either the Docker container or the standalone executable.
3. If you chose the environment variables method, the application will automatically use the provided settings to set your Discord status.
4. If you chose the interactive method, you'll be guided through the setup process.
5. For emoji usage:
   - If you're familiar with Discord emoji configuration, you can opt to use an emoji when prompted or set `USE_EMOJI` to "true".
   - If you're unsure about emoji configuration, it's recommended to choose not to use an emoji for a simpler setup.

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

- **LICENSE:** Contains the MIT license information for the project.
- **README.md:** Provides a comprehensive overview and instructions.
- **dockerfile:** Contains instructions for creating a Docker image.
- **go.mod:** Tracks the project dependencies.
- **go.sum:** Verifies the integrity of the dependencies.
- **main.go:** The main source code file of the project.

## Contributing

Contributions are welcome and appreciated! Here's how you can contribute:

1. Fork the repository
2. Create a new branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please ensure your code adheres to the project's coding standards and includes appropriate tests.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for full details.

## Acknowledgements

- [Gorilla WebSocket](https://github.com/gorilla/websocket) for the robust WebSocket implementation.
- The Go community for providing excellent libraries and resources.
- All contributors who have helped improve ProjectDGT.

## Support

If you encounter any issues or have questions, please [open an issue](https://github.com/Milou4Dev/ProjectDGT/issues) on the GitHub repository. We're here to help!

## Disclaimer

This project is intended for educational and personal use only. Please ensure you comply with Discord's terms of service when using this tool.
