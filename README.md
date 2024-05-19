## Overview
Coming soon.

## Docker Support
Currently, only Docker is supported.

### Build and Run
To build and run the Docker container, use the following commands:
```
docker buildx build --pull --no-cache --tag projectdgtimg . --load && docker run --rm -it --init projectdgtimg
```
