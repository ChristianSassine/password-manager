# password-manager
A CLI tool with a server for managing passwords

## Description

This a CLI tool solution for managing passwords with a server that can be deployed on the same computer or the local network. The server encrypts and stores the passwords in a MongoDB instance. 
The tool automatically copies the password to the clipboard but it is possible to output it and to remove the descriptions of the CLI. The generated passwords can also be customized.

## Getting Started

### Dependencies

* Golang/Go
* Windows 10/11 or Linux
* xclip or xsel if you're using Linux
* Docker

### Installing

* Install the CLI's binary by using the following command ```go install github.com/ChristianSassine/password-manager/pass-cli@main```
* Clone the repo in the server/storage's node with ```git clone https://github.com/ChristianSassine/password-manager.git```

### Running

* Execute `docker compose up` in the repository on the server node to run the database and the server
* You will find the CLI's binary in GOBIN if it's defined, otherwise it will be in GOPATH/bin
* To use the password manager, execute the commands of the cli in a terminal

## Help

```
pass-cli help
```

## License

This project is licensed under the MIT License - see the LICENSE.md file for details
