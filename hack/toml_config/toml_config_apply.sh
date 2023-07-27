#!/bin/bash

# Execute the struct generator
echo "Running the struct generator..."
python3 config_generator.py

# Check if the struct generator executed successfully
if [ $? -eq 0 ]; then
    echo -e "\n\u2714 Struct generator executed successfully."
else
    echo -e "\n\u2718 Struct generator failed."
    exit 1  # Exit with a non-zero status code to indicate failure
fi

# CD into the root directory
echo -e "\nMoving to the root directory..."
cd ../../

# Generate the manifests
echo -e "\nGenerating manifests..."
make manifests
make

# Check if the manifest generation executed successfully
if [ $? -eq 0 ]; then
    echo -e "\n\u2714 Manifests generated successfully."
    echo -e "\nYour new configs are now supported by the operator."
    exit 0  # Exit with a status code of 0 to indicate success
else
    echo -e "\n\u2718 Manifest generation failed."
    exit 1  # Exit with a non-zero status code to indicate failure
fi
