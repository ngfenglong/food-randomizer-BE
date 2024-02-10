# Time-To-Makan (TTM) API 🍲

This repository contains the API for the Time-To-Makan (TTM) project, a dedicated platform for easing the decision-making process of where to eat in the Novena area, Singapore. It powers both a web application and a Telegram bot to generate random dining places based on user preferences. The backend is designed for extensibility to support additional locations and filters in the future.

> 🚨 This is an active project and is continuously evolving. Detailed documentation will be provided as new features are developed.

## Table of Contents
- [Technology Stack](#technology-stack)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [Project Structure](#project-structure)
- [Contribution](#contribution)
- [Contact](#contact)

## Technology Stack 💻
- **Language:** Go
- **Database:** MySQL

## Getting Started
1. Clone the repository.
2. Navigate to the project directory and update the "app.env" file with your database details for local development.
3. Run the project using the Makefile:
   
   ```
   make start

## Usage 🛠️
The API serves as the backend for the TTM web application and the Telegram bot, handling place, category, and location management. It's capable of operating independently as a standalone server or in conjunction with the front-end services.

## Project Structure 🌳
```
cmd/
└── main.go            # Application entry point
pkg/
├── auth/              # Authentication logic
├── category/          # Category management
├── config/            # Configuration handling
├── database/          # Database operations
├── http/              # HTTP server and routing
├── location/          # Location management
├── middleware/        # Middleware
├── models/            # Data models
├── place/             # Place management
└── utils/             # Utility functions
data-seeding/          # Database seeding scripts
dist/                  # Distribution binaries
```
> **Note:** The "dist" directory is excluded from this repository as it is generated during the build process and is not tracked in version control.
## Contribution 🤝
Contributions, feature ideas, and bug reports are highly appreciated. To contribute, please fork the repository, make your changes, and submit a pull request.

## Contact 📬
For inquiries or further information about this project, reach out to [zell_dev@hotmail.com](mailto:zell_dev@hotmail.com).
