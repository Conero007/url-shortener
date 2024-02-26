# URL Shortener in Golang

I designed and implemented a production-ready URL Shortener using Golang. This application uses MySQL for robust persistent storage and Redis for efficient caching. To ensure scalability and secure external access, the entire application has been dockerized, with Nginx managing incoming traffic to the Docker container.

## Prerequisites

To run this application, ensure you have the following installed:

- Go 1.21 or later
- MySQL 8.0 or later
- Redis 7 or later
- Nginx 1.23 or later

## Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Conero007/url-shortener.git
   ```

2. **Run the setup script:**

   Execute the bash script present in the `/bash` folder. This script sets up a Docker network on your system, creates a `.env` file with the contents of the `.example.env` file, and also adds an entry for the app URL to your `/etc/hosts` file:

   ```bash
   sh bash/setup.sh
   ```

3. **Configure your environment:**

   Edit the `.env` file with the credentials for your environment, prefilled with values for the configuration present in the GitHub repository:

   ```env
   # App Config
   PORT=3000
   APP_URL=url.shortener.local

   # Database Config
   DB_ADDR=db:3306
   DB_USERNAME=root
   DB_PASSWORD=1234
   DB_NAME=url_shortener

   # Redis Config
   REDIS_ADDR=redis:6379
   REDIS_PASSWORD=
   ```

4. **Configure testing environment:**

   Similarly, configure your `.testing.env` file with the credentials for your testing environment, prefilled with values for the configuration present in the GitHub repository.

5. **Build and start the application:**

   Run the following commands to build and start the application. This will also run the unit tests present in the project. If the tests fail, the application will not start:

   ```bash
   docker-compose up -d --build
   ```

6. **Access the application:**

   If all containers are up and running, the project is successfully up on your system, and its API endpoint can be accessed at the following URL:

   ```
   http://url.shortener.local/<endpoint>
   ```

## Usage

The application provides the following API endpoints:

1. **`/shorten`**: This endpoint is used to create a new short URL. Optionally, you can also provide a custom short key for generating a custom short URL. The request body should contain the following JSON:

   ```json
   {
     "url": "https://www.example.com",
     "custom_short_key": "abc123"
   }
   ```

   - A successful response will contain the following JSON:

   ```json
   {
     "original_url": "https://www.example.com",
     "short_url": "http://url.shortener.local/abc123",
     "expire_time": "2024-02-10"
   }
   ```

   - A failed response will contain the following JSON:

   ```json
   {
     "error": "<error message>"
   }
   ```

2. **`/{key}`**: This endpoint is used to redirect to the original URL. The `key` parameter is the short key of the URL.

   - On a successful response, you will be redirected to the original URL.

   - A failed response will contain the following JSON:

   ```json
   {
     "error": "<error message>"
   }
   ```

Feel free to reach out if you have any questions or need further assistance!