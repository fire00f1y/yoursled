setlocal
REM If you want to build locally, this should point to your project folder
SET GOPATH=C:\Source\GoProjects\yoursled
SET PATH=%PATH%;$GOPATH/bin
SET GOARCH=amd64
SET GOOS=linux
call go install github.com/fire00f1y/yoursled
SET GOOS=darwin
call go install github.com/fire00f1y/yoursled
SET GOOS=windows
call go install github.com/fire00f1y/yoursled
