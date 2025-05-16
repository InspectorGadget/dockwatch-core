# DockWatch Core

## Overview
DockWatch Core is a lightweight, open-source monitoring tool designed to keep track of your Docker containers. It provides real-time insights into the performance and status of your containers, making it easier to manage and optimize your Docker environment.

DockWatch has to be paird with the [DockWatch Frontend](https://github.com/InspectorGadget/dockwatch-ui)

## Features
- Real-time monitoring of Docker containers
- Resource usage statistics (CPU, memory, disk I/O)
- User-friendly web interface

## Installation
1. Clone the repository:
   ```bash
    git clone git@github.com:InspectorGadget/dockwatch-core.git
    cd dockwatch-core
    ```
2. Review the environment variables in the `docker-compose.yml` file. You can customize the configuration as needed.
    - Alternatively, you can create a `config.json` file in the root directory with the following structure:
        ```json
        {
            "DOCKER_HOST": "unix:///var/run/docker.sock"
        }
        ```
3. Run the docker-compose command to start the application:
   ```bash
    docker-compose up -d
    ```
4. Access the WebSocket at `http://localhost:8080/socket`.