# Gorder Server

Gorder Server is a simple Go (Golang) application that allows you to securely call local programs via HTTP requests. It provides a basic API to execute predefined programs on the server.
- [Gorder Server](#gorder-server)
  - [Installation](#installation)
  - [Build](#build)
  - [Usage](#usage)
  - [Configuration](#configuration)
  - [Security](#security)
  - [API Endpoint](#api-endpoint)
    - [POST /call-program](#post-call-program)
  - [License](#license)

<!-- Table of Contents

    Installation
    Usage
    Configuration
    Security
    API Endpoint
    Examples
    License -->
## Installation
Download from [releases](https://github.com/wansyu/gorder/releases/latest).
## Build

Make sure you have Go installed on your machine. Clone this repository and navigate to the project directory:

<!-- bash -->
```
git clone https://github.com/wansyu/gorder.git
cd gorder
```


Build the application:

<!-- bash -->
```
go build
```



## Usage

Once the server is running, you can make HTTP POST requests to the specified route (/call-program by default) to execute local programs.

Run the executable:

<!-- bash -->
```
./gorder -c <config-path> -ip <your-ip-address> -p <port> -r <route-path>
```


Replace `<config-path>`(config.json by default), `<your-ip-address>`(empty by default), `<port>`(8080 by default), and `<route-path>`(/call-program by default) with your desired values.

## Configuration

The server is configured using a JSON file (default: config.json). This file includes details such as the server key, salt for password encryption, and paths to the programs that can be executed.

Example `config.json`:

<!-- json -->
```
{
  "key": "your-server-key",
  "salt": "your-salt-value",
  "program_paths": [
    {"name": "program1", "path": "/path/to/program1", "args": ["arg0"]},
    {"name": "program2", "path": "/path/to/program2", "args": ["arg0", "arg1"]}
  ]
}
```
`"key"` should generate by `authgen.py`, edit the `password` and `salt` as your disired. 

## Security

The server uses scrypt for password-based key derivation to securely validate the incoming requests. Ensure that you set a strong key and salt in the configuration file.
## API Endpoint

###    POST /call-program

Execute a local program. Send a JSON payload with the program name and key for authentication.
    
An example is in `requ.py`.



## License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/wansyu/gorder/blob/main/LICENSE.md) file for details.