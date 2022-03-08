#!/bin/bash

while ! curl --silent -o /dev/null --fail "$1"; do
    echo "Waiting for the application to be ready, retrying in 5s..."
    sleep 5
done

echo "The application is ready"
