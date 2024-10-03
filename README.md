# dkrTag: CLI Tool for Searching Docker Tags

`dkrTag` is a tool built in Go that fetches Docker tags from a specific repository using the [DockerHub APIs](https://docs.docker.com/reference/api/hub/latest/) and provides a UI interface to filter those tags. It features a user-friendly terminal interface built with [Charmbracelet's Bubble Tea](https://github.com/charmbracelet/bubbletea) framework and the [Bubbles List](https://github.com/charmbracelet/bubbles) component for easy interaction.

## Installation

1. Download the `dkrTag` binary from the [releases page](https://github.com/rajhawaldar/dkrTag/releases).
2. Install via Go:

    ```bash
    go install github.com/rajhawaldar/dkrTag@latest
    ```

## Usage

The syntax for using the tool:
```bash
dkrTag --repository <repository-name> [--namespace <namespace-name>]
```

Example:
```bash
dkrTag --repository nginx 
```

> [!TIP]
> If you're logged into Docker CLI using the ```docker login``` command, you can fetch tags from private repositories as well.

Example: 
```bash
dkrTag --repository webapp --namespace rajhawaldar
```
