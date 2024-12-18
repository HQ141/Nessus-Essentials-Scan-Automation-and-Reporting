#!/bin/bash

# Variables
ENV_FILE="/opt/nessus_ess/.env"
TEST_FILE="/opt/nessus_ess/test.txt"
ASSETS_SRC="src/assets"
ASSETS_DEST="/opt/nessus_ess/assets"
SRC_FOLDER="src"
BUILD_OUTPUT="/opt/nessus_ess/nessus_automation"
SAMPLE_ENV="docs/sample.env"

# Usage function
usage() {
    echo "Usage: $0 [--use-sample-env | --manual-env]"
    echo "  --use-sample-env    Use the sample.env file to create the .env file."
    echo "  --manual-env        Manually input values to create the .env file."
    exit 1
}
# Check command-line arguments
if [[ $# -ne 1 ]]; then
    usage
fi

case "$1" in
    --use-sample-env)
        USE_SAMPLE_ENV="true"
        ;;
    --manual-env)
        USE_SAMPLE_ENV="false"
        ;;
    *)
        usage
        ;;
esac

# Create the nessus_ess user if it doesn't already exist
if ! id "nessus_ess" &>/dev/null; then
    echo "Creating user nessus_ess..."
    sudo useradd nessus_ess
else
    echo "User nessus_ess already exists."
fi

# Remove the existing /opt/nessus_ess folder if it exists
if [[ -d "/opt/nessus_ess" ]]; then
    echo "Removing existing /opt/nessus_ess folder..."
    sudo rm -rf /opt/nessus_ess
fi

# Create /opt/nessus_ess folder
echo "Creating /opt/nessus_ess folder..."
sudo mkdir -p /opt/nessus_ess

# Copy the assets folder to /opt/nessus_ess
echo "Copying assets folder to /opt/nessus_ess..."
sudo mkdir -p "$ASSETS_DEST"
sudo cp -r "$ASSETS_SRC/"* "$ASSETS_DEST"

# Handle .env creation based on the USE_SAMPLE_ENV flag
if [[ "$USE_SAMPLE_ENV" == "true" ]]; then
    # Copy sample.env to /opt/nessus_ess/.env
    if [[ -f "$SAMPLE_ENV" ]]; then
        echo "Copying $SAMPLE_ENV to /opt/nessus_ess/.env..."
        sudo cp "$SAMPLE_ENV" "$ENV_FILE"
    else
        echo "Error: $SAMPLE_ENV not found. Exiting."
        exit 1
    fi
else
    # Create the .env file in /opt/nessus_ess and prompt for manual input
    echo "Manual .env setup enabled. Creating .env file in /opt/nessus_ess..."
    sudo touch "$ENV_FILE"

    # Function to prompt the user and write input to the .env file
    write_to_env() {
        local key=$1
        local prompt=$2

        # Read input from the user
        read -p "$prompt: " value

        # Append the key-value pair to the .env file
        echo "$key=$value" | sudo tee -a "$ENV_FILE" > /dev/null
    }

    echo "Setting up .env file..."
    write_to_env "POSTGRES_USER" "Enter the PostgreSQL username"
    write_to_env "POSTGRES_PASSWORD" "Enter the PostgreSQL password"
    write_to_env "POSTGRES_DB" "Enter the PostgreSQL database name"
    write_to_env "POSTGRES_URL" "Enter the PostgreSQL URL"

    write_to_env "SENDGRID_USEREMAIL" "Enter the SendGrid user email"
    write_to_env "SENDGRID_APIKEY" "Enter the SendGrid API key"
    write_to_env "RECEIVERS" "Enter the email receivers (comma-separated)"

    write_to_env "NESSUS_USERNAME" "Enter the Nessus username"
    write_to_env "NESSUS_PASSWORD" "Enter the Nessus password"
    write_to_env "NESSUS_SCANNAME" "Enter the Nessus scan name"
    write_to_env "NESSUS_URL" "Enter the Nessus URL"

    echo "Environment variables have been saved to $ENV_FILE."
fi

# Create a test.txt file in /opt/nessus_ess
echo "Creating test.txt file in /opt/nessus_ess..."
sudo touch "$TEST_FILE"

# Move to the source folder and build the Go application
echo "Building the Go application..."
cd "$SRC_FOLDER" || { echo "Source folder $SRC_FOLDER not found! Exiting."; exit 1; }
sudo go build -o "$BUILD_OUTPUT" -buildvcs=false

if [[ $? -eq 0 ]]; then
    echo "Go application built successfully and saved to $BUILD_OUTPUT."
else
    echo "Error building Go application. Check your source code and dependencies."
    exit 1
fi

# Change ownership of /opt/nessus_ess and its contents to nessus_ess
echo "Changing ownership of /opt/nessus_ess to nessus_ess..."
sudo chown -R nessus_ess:nessus_ess /opt/nessus_ess

echo "Setup completed successfully."
