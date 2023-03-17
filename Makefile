# Copyright © 2022-2023 PaoloB
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

BINARY_NAME=adb-cli

build:
	GOARCH=amd64 GOOS=darwin go build -o executables/${BINARY_NAME}-darwin #main.go
	GOARCH=amd64 GOOS=linux go build -o executables/${BINARY_NAME}-linux #main.go
	GOARCH=amd64 GOOS=windows go build -o executables/${BINARY_NAME}-windows #main.go
	mv executables/${BINARY_NAME}-darwin ${BINARY_NAME}

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm executables/${BINARY_NAME}-darwin
	rm executables/${BINARY_NAME}-linux
	rm executables/${BINARY_NAME}-windows