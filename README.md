# NessusEssAutomation Utility

This repository contains the **NessusEssAutomation** utility, designed to enhance **Nessus Essentials** by adding, maintaining, and storing historical data. The utility automates the process of retrieving scan data, updating vulnerability records, and generating reports.

## Overview

The utility is built using:

- **Go**: For efficient backend processing.
- **PostgreSQL**: To store and query historical vulnerability data.
- **Docker**: For containerized deployment.
- **wkhtmltopdf**: For generating PDF reports.

---

## Getting Started

The detailed instructions for setting up and configuring the utility are available in the `docs/Nessus_Utility.md` file. This guide covers:

- **Prerequisites**: Required tools and dependencies.
- **Installation Steps**: How to install Docker, Go, and other dependencies.
- **Setup**: Instructions to configure Nessus, PostgreSQL, and environment variables.
- **Cron Job Configuration**: How to automate the utility to run at scheduled intervals.

---

## Understanding the Workflow

A detailed sequence diagram illustrating the utility's workflow is available in `docs/sequence.md`. This diagram provides a visual representation of how the utility interacts with Nessus, PostgreSQL, and other components.

To view the sequence diagram, visit:

[Sequence Diagram: sequence.md](docs/sequence.md)

---

## Documentation

To begin with the utility setup, refer to the main setup guide:

[Setup Guide: Nessus_Utility](docs/Nessus_Utility.md)

---
## Configuring Nessus

Once Nessus Essentials is installed and running, you must configure it to work with this utility. Detailed instructions for setting up Nessus are available in the `docs/Nessus.md` file. This guide includes:

- Setting up scans and schedules.
- Configuring the required settings to ensure seamless integration with the utility.

To access this guide, visit:

[Configuring Nessus: Nessus.md](docs/Nessus.md)


---


## Contributing

If you'd like to contribute to this project, please fork the repository and submit a pull request. For issues or feature requests, please open a ticket in the repository.
