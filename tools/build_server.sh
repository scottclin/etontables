#!/bin/bash          
echo Building Server

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
$DIR/fetch_requires.sh

echo ..building
go build server.go

echo done 