# Nessus Configuration Setup

This text file provides instructions for setting up Nessus Essentials for the first time. It should be used in conjunction with the `nessus_configuration.md` guide.

---

1. **Start Nessus in Docker**:

   - Follow the previous instructions to get Nessus running in Docker (including setting up the `.env` file and starting the services).

2. **Set Up an SSH Tunnel**:

   Since Nessus may not be accessible via a web browser directly due to Docker network settings, we will use an **SSH tunnel** to forward the necessary ports to your local machine. This allows you to access Nessus via `https://localhost:8834` from a browser.

   - Open a terminal on your local machine and execute the following command:
   
     ```bash
     ssh -L 8834:localhost:8834 user@<your-remote-server-ip> -N
     ```

     Where:
     - `8834:localhost:8834` forwards the Nessus port (8834) from the remote server to your local machine.
     - `user` is your username on the remote server running Docker with Nessus.
     - `<your-remote-server-ip>` is the IP address of the server where Docker is running.

   - This command sets up a tunnel that forwards traffic from `localhost:8834` on your machine to the Nessus web interface running on the Docker container.

3. **Access Nessus via Browser**:

   - Now, open your browser and navigate to:
     ```
     https://localhost:8834/
     ```

   - Login using the credentials provided in the `docker.env` file.

4. **Activate Nessus Essentials**:

   - Once logged in, navigate to **Activation Code** in the **Settings** section of the web application.
   - Enter the **valid activation code** for Nessus Essentials.
   - Click **Activate** and wait for Nessus to download and install the required plugins.

5. **Wait for Plugin Installation**:

   - Once the plugins are installed, Nessus will be ready for use. This may take some time, so please be patient during the process.

6. **Create a New Scan**:

   - After the plugins have been installed, navigate to the **Scans** section in the left column.
   - Select **My Scans**.
   - Click the **New Scan** button. This will open a screen for selecting a scan template.

7. **Select a Scan Template**:

   - Choose the appropriate scan template for your needs (e.g., **Basic Network Scan**, **Advanced Scan**).
   - This will open another screen for configuring the scan.

8. **Name the Scan and Configure Settings**:

   - Provide a meaningful name for the scan (e.g., `My Network Scan`).
   - Configure any additional settings as necessary for your environment.

9. **Update the `.env` File**:

   - Navigate to `/opt/nessus_ess/.env` and update the following values:
     - **Scan Name**: Add the scan name you have created in Nessus.

     Example `.env` file:
     ```
     NESSUS_SCAN_NAME=<your-scan-name>
     ```

10. **Provide Credentials for Targets**:

    If your scan requires credentials to access the targets (such as login details for web servers or network devices), you need to configure them in the **Nessus Scan** settings.

    - Navigate to the **Credentials** tab in the scan configuration screen.
    - Add the necessary credentials based on the scan type. Some common options include:
      - **SSH credentials** for Linux/Unix targets.
      - **Windows credentials** for scanning Windows systems.


11. **Schedule the Scan**:

    - In the **Basic Settings** section, navigate to the **Schedule** tab.
    - Set the time according to the time of your cron job (which runs at 12:00 on Saturday).
    - Save the scan by clicking **Save**.

---

Once these steps are completed, the **NessusEssAutomation** utility will be able to interact with your Nessus instance, retrieving scan data and maintaining historical records as per the setup outlined in the `README.md`.

Ensure the scan is running correctly and check the logs for any issues related to the connection or scanning process.

If you face any difficulties, refer to the logs or consult the troubleshooting section in the repository.
