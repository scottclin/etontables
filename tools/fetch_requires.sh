#!/bin/bash

echo ..fetching bencode package 
go get -u github.com/jackpal/bencode-go

echo ..fetching yaml.v2 package
go get -u gopkg.in/yaml.v2