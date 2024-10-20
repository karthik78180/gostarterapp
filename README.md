# GoStarterApp - CLI Tool

## Overview

**GoStarterApp** is a command-line interface (CLI) tool designed to fetch zones based on the provided environment, project, and bearer token. It dynamically validates user inputs and prompts for missing values, displaying available options as needed.

## Prerequisites

- **Go (Golang)** must be installed. You can download and install it from the official website: [https://golang.org/dl/](https://golang.org/dl/).

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/karthik78180/gostarterapp.git
   cd gostarterapp
   ```
2. **Build the application**
    Run the following command to compile the code into a binary
   ```bash
   go build -o gostarterapp
   ```
3. **Run the application**
    Use the built binary to execute the program:
    a. **Provide all inputs using flags:**
    You can provide all inputs upfront using the available flags:

    ```bash
    ./gostarterapp --env dev --project OD --token my-bearer-token
    ```
    b. **Prompt for missing inputs:**
    If you don’t provide any of the inputs, the program will prompt you. For example, if you don't pass any flags:
    ```bash
    ./gostarterapp
    ```
## Usage

The CLI accepts three inputs:

1. **Environment (`env`)**: The environment in which the application is running. Valid values are loaded from the `config.json` file.
   - If not provided via the command-line flags, the program will prompt the user to enter a valid environment from the available options.
   
2. **Project (`project`)**: The specific project within the chosen environment.
   - If the environment is `dev`, only `OD` and `ODS` are valid options.
   - For other environments, the valid projects will be loaded from the `config.json` file.
   - If not provided, the program will prompt the user to enter a valid project from the available options.
   
3. **Bearer Token (`token`)**: An authentication token required to fetch the zones.
   - If not provided, the program will prompt the user to enter a token.