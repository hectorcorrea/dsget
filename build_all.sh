# Builds the executables for Mac, Linux, and Windows
# (requires Go the be installed https://go.dev/)
GOOS=darwin go build -o dsget
GOOS=linux go build -o dsget_linux
GOOS=windows GOARCH=386 go build -o dsget.exe
