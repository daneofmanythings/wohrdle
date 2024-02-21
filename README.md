# WOHRDLE

![screenshot](/static/bday_image.png)

![code coverage badge](https://github.com/daneofmanythings/wohrdle/actions/workflows/ci.yml/badge.svg)

## Setup
- Install [golang](https://go.dev/doc/install)
- *Optional* Navigate to your desired location to hold the repository. 
```
cd $HOME/Downloads/
```
- Clone the repository.
```
git clone https://github.com/daneofmanythings/wohrdle
```


## Usage
- *Recommended* Install the binary into your \$GOBIN (\$GOPATH/bin) to expose it
as a command in your terminal.
```
go install
```
You can then run it from anywhere with:
```
wohrdle
```
- If you don't want to save the built binary, run this from the project root
to run the program as a one-off.
```
go run main.go
```
- If you want to built the binary but not commit it to your path.
```
go build
```
Then to run the program, run this from the project root, or enter the full path
to where you cloned the repository.
```
./wohrdle
```
