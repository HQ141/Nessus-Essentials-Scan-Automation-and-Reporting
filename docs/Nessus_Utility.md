# NessusEssAutomation Utility

This utility enhances **Nessus Essentials** by adding, maintaining, and storing historical data that is not available in the community version of Nessus. The utility is built with **Go**, **PostgreSQL**, and uses **Docker** for deployment.

## Prerequisites

Before proceeding with the installation, ensure you have the following installed on your system:

- **Docker** and **Docker Compose**
- **Go** (for building the utility)
- **PostgreSQL** (database for storing historical data)
- **wkhtmltopdf** (for PDF generation)

### Install Docker

1. **Add Docker's official GPG key**:

   ```bash
   sudo apt-get update
   sudo apt-get install ca-certificates curl
   sudo install -m 0755 -d /etc/apt/keyrings
   sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
   sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add the repository to Apt sources:
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    ```
2. **Install Docker packages:**
    ``` bash
    sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-compose
    ```
### Install Dependencies
1. **Install Go and other dependencies:**
    ```bash
    sudo apt install golang-go wkhtmltopdf
    ```
## Setting Up the Utility
### Clone the Repository
Clone the repository containing the NessusEssAutomation utility:
    
    git clone https://dev.azure.com/vidizmo/VTLS/_git/VIDIZMOSecurity
    cd NessusEssAutomation
### Docker Setup
#### Set Up the Docker Environment File

1. **Move and edit the sampledocker.env file:**

    ```bash
    cd devops
    mv ../docs/sampledocker.env docker.env && nano docker.env
    ```

    Update the .env file with the required values, such as database credentials and other necessary configurations.

2. **Start Docker containers:**

    After updating the `.env` file, start the Docker containers with:
    ```bash
    docker-compose up -d
    ```
### Build and Setup Utility
1. **Run the setup  and install bash script:**
    1. *make the script executable:*
        ```bash
        cd ..
        chmod +x setup_and_install.sh
        ```
    2. *Run the Script with Arguments:*
        * To use `sample.env`
            ```bash
            sudo ./setup_and_install.sh --use-sample-env
            ```
        * To manually enter `.env` values:
            ```bash
            sudo ./setup_and_install.sh --manual-env
            ```
2. Follow the script prompts or amke changes in the `.env` file in `/opt/nessus_ess`.

### Cron Job Setup
1. **Add the utility to cron job for the nessus_ess user:** 
    Edit the cron job for the nessus_ess user to schedule the utility. Run:
    ```bash
    sudo crontab -u nessus_ess -e
    ```
    Then, add the following line to run the utility every Saturday at 12 PM:
    ```bash
    0 12 * * 6 cd /opt/nessus_ess && /opt/nessus_ess/nessus_automation -l -m -c 2>> /var/log/nessus_ess

    ```
3. **Save and exit the crontab editor.**

### Conclusion

The NessusEssAutomation utility is now set up to enhance Nessus Essentials by storing historical data. This setup involves building the utility, configuring Docker and environment files, setting up the user, and scheduling the utility to run at 12 PM every Saturday using cron.

If you encounter any issues or need further assistance, feel free to open an issue or contact support.