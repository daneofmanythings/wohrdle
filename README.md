<div align="center">
# WOHRDLE

![screenshot](/static/bday_image.png)

![code coverage badge](https://github.com/daneofmanythings/wohrdle/actions/workflows/ci.yml/badge.svg)

###### A Harder Way to Play Wordle

[Summary](#summary)
•
[Setup](#setup)
•
[Installation](#installation)

## Summary
</div>
A wordle clone with more challenging options. Freely adjust the word length, guess count, failed
input count, and classic hardmode!

![screenshot](/static/settings_image.png)

## Setup
- Install [golang](https://go.dev/doc/install)
- **[Optional]** Navigate to your desired location to hold the repository. Ex:
```
cd $HOME/Downloads/
```
- Clone the repository.
```
git clone https://github.com/daneofmanythings/wohrdle
```


## Installation
There are three options for how to install and use the app.
#### 1. System Install **\[Recommended\]** 

Install the binary into your \$GOBIN (\$GOPATH/bin) to expose it as a command in your terminal.
```
go install
```
You can then run it from anywhere with: ` wohrdle `

#### 2. Build in and run from the source directory

If you want to build the binary but not commit it to your path.
```
go build
```
Then, to run the program, enter the full path to where you cloned the repository. Ex:
```
$HOME/Downloads/wohrdle
```
or, from the repository root, run: `./wohrdle`

#### 3. Run from the source directory

If you don't want to save the built binary, run this from the repository root
to run the app as a one-off.
```
go run main.go
```

### Note:
The `build` and `install` commands default the binary name to the go package name.
In cases 1 and 2, you may specify a different name for the binary with the `-o` flag. 
Ex:
```
go install -o newname 
```

which can then be ran with `newname`
