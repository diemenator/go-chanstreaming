#!/bin/bash

# Loop forever until input fails (e.g., EOF)
while true; do
  # Prompt user input
  if ! IFS= read -r line; then
    break  # exit on EOF or input failure
  fi
  # Check if the input is empty
  if [[ -z "$line" ]]; then
    break
  fi
  # Echo back the input
  echo "You said: $line"
done