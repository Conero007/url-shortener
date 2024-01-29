# URL Shortener in Golang

This project is a straightforward URL shortener application written in Go. It utilizes MySQL to store the original URLs and their corresponding short keys, and Redis for caching the short URLs. The application additionally provides an API for creating new short URLs and redirecting to the original URLs.

## Prerequisites

To run this application, ensure you have the following installed:

- Go 1.21 or later
- MySQL 8.0 or later
- Redis 7 or later

## Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/Conero007/url-shortener.git
   ```

2. Run the bash script present in the `/bash` folder. This script sets up a Docker network on your system and creates a `.env` file with the contents of the `.example.env` file:

   ```bash
   sh bash/setup.sh
   ```

3. Configure your `.env` file with the credentials for your environment:

   ```env
   # App Config
   APP_URL=127.0.0.1
   PORT=3000

   # Database creds
   DB_ADDR=
   DB_USERNAME=
   DB_PASSWORD=
   DB_NAME=

   # Redis creds
   REDIS_ADDR=
   REDIS_PASSWORD=
   ```

4. Configure your `.testing.env` file with the credentials for your testing environment:

   ```env
   # App Config
   APP_URL=127.0.0.1
   PORT=3000

   # Database creds
   DB_ADDR=
   DB_USERNAME=
   DB_PASSWORD=
   DB_NAME=

   # Redis creds
   REDIS_ADDR=
   REDIS_PASSWORD=
   ```

5. Run the following commands to build and start the application. This will also run the unit tests present in the project. If the tests fail, the application will not start:

   ```bash
   docker-compose up -d --build
   ```

6. If all containers are up and running, the project is successfully up on your system, and its API endpoint can be accessed at the following URL:

   ```
   http://127.0.0.1/<endpoint>
   ```

## Usage

The application provides the following API endpoints:

1. `/shorten`: This endpoint is used to create a new short URL. Optionally, you can also provide a custom short key for generating a custom short URL. The request body should contain the following JSON:

   ```json
   {
     "url": "https://www.example.com",
     "custom_short_key": "abc123" // optional
   }
   ```

   - A successful response will contain the following JSON:

   ```json
   {
     "original_url": "https://www.example.com",
     "short_url": "http://localhost:3000/abc123",
     "expire_time": "2024-02-10"
   }
   ```

   - A failed response will contain the following JSON:

   ```json
   {
     "error": "<error message>"
   }
   ```

2. `/{key}`: This endpoint is used to redirect to the original URL. The `key` parameter is the short key of the URL.

   - On a successful response, you will be redirected to the original URL.

   - A failed response will contain the following JSON:

   ```json
   {
     "error": "<error message>"
   }
   ```

## Code Overview

The application is structured as follows:

- `app/`: This directory contains the main application code, including controller and handler logic.
- `database/`: This directory contains the code for handling the migrations for the database.
- `models/`: This directory contains the data models and their logic for the application.
