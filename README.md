
go mod init github.com/belieflab/dataview
go install github.com/spf13/cobra-cli@latest
export PATH=$PATH:$(go env GOPATH)/bin
cobra-cli init
go mod tidy