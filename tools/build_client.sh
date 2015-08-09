#!/bin/bash          
echo Building client

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
$DIR/fetch_requires.sh

echo ..building
go build client.go

echo done 