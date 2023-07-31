# password-manager
A CLI tool with a server for managing passwords

## Description

This a CLI tool solution for managing passwords with a server that can be deployed on the same computer or the local network. The server **encrypts** and stores the passwords in a MongoDB instance. 
The CLI tool automatically copies the password to the **clipboard**. It is also possible to **silence** the descriptions of the CLI and to **output** the password. It also supports multiple **options** for the generated passwords.

## Getting Started

### Dependencies

* Golang/Go
* Windows 10/11 or Linux
* xclip or xsel if you're using Linux
* Docker

### Installing

* Install the CLI's binary by using the following command:
```
   go install github.com/ChristianSassine/password-manager/pass-cli@main
```
* Clone the repo in the server/storage's node with
```
  git clone https://github.com/ChristianSassine/password-manager.git
```

### Using
#### Server
* Execute `docker compose up` in the cloned repository to run the server and database containers
#### Client
* You will find the CLI's binary in GOBIN if it's defined, otherwise it will be in GOPATH/bin
* You will need to set the server's URL and create a user on the client with:
```
  pass-cli config -u <URL>:<PORT> -c <USER>:<PASSWORD>
  e.g: pass-cli config -u localhost:8080 -c Hello:World
```
* Set the envonmental variables PASS_USERNAME and PASS_PASSWORD to your credentials

## Help

```
pass-cli help
```

## License

This project is licensed under the MIT License - see the LICENSE.md file for details
