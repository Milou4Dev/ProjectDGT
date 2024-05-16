## Overview
Coming soon.

## Docker Support
Currently, only Docker is supported.

### Build and Run
To build and run the Docker container, use the following commands:
```
docker build --pull --no-cache --tag projectdgtimg . && \
docker run --rm --name projectdgt --init -it projectdgtimg
```
