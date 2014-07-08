SET GOARCH=amd64
SET GOOS=linux
go build -o solver_%GOOS%_%GOARCH% solver.go 