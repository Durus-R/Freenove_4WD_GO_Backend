#!/bin/bash

# Check if go is in the PATH
if ! command -v go &> /dev/null
then
    echo "Go is not found on the PATH. Exiting."
    exit 1
fi


# Get the go version
go_version=$(go version | awk '{print $3}' | tr -d "go")

# Specify the version to check
desired_version="1.20"

# Check if the installed go version is greater than or equal to the desired version
if echo -e "$desired_version\n$go_version" | sort -V | awk 'END{exit !($0 == ENVIRON["go_version"])}' go_version="$go_version"; then
    echo "Go version is compatible with $desired_version. Installed version is: $go_version"
else
    echo "Go version is not compatible with $desired_version. Installed version is: $go_version"
    exit 1
fi