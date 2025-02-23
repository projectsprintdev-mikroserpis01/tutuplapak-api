# TutupLapak API

![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)

## Description

The application is built on [Go v1.23.4](https://tip.golang.org/doc/go1.22) and [PostgreSQL](https://www.postgresql.org/). It uses [Fiber](https://docs.gofiber.io/) as the HTTP framework and [pgx](https://github.com/jackc/pgx) as the driver and [sqlx](github.com/jmoiron/sqlx) as the query builder.

## Getting Started

1. Ensure you have [Go](https://go.dev/dl/) 1.23 or higher and [Task](https://taskfile.dev/installation/) installed on your machine:

   ```sh
   go version && task --version
   ```

2. Create a copy of the `.env.example` file and rename it to `.env`:

   ```sh
   cp ./config/.env.example ./config/.env
   ```

   Update configuration values as needed.

3. Install all dependencies:

   ```sh
   task
   ```

4. Start the Docker containers:

   ```sh
   task service:up:build
   ```

5. Create the database schema:

   ```sh
   task db:connect
   ```

6. Run database migrations:

   ```sh
   task migrate:up
   ```

---

## Running the Application with Docker Compose

1. Start the Docker containers:

   ```sh
   task service:up:build
   ```

2. To stop the Docker containers:

   ```sh
   task service:down
   ```

3. To build the Docker containers:

   ```sh
   task service:build
   ```

---

## Checking Database Connection

To check the database connection, use the following command:

```sh
task service:db:connect
```

---

## Troubleshooting Database Connection Issues

If you encounter the following error during migration or database connection:
```
error: failed to open database: dial tcp: lookup db on 127.0.0.53:53: server misbehaving
```
It may be due to a DNS resolution issue. You can resolve this by adding `db` to your local `/etc/hosts` file:

```sh
echo "127.0.0.1 db" | sudo tee -a /etc/hosts
```

Then retry:
```sh
task migrate:up
```

---

## Running Load Tests

1. Navigate to the `test` folder and clone the repository:

   ```sh
   git clone https://github.com/ProjectSprint/Batch3Project2TestCase.git
   ```

2. Install `k6` (if you don't have it installed):

   Follow the instructions on the [k6 installation page](https://k6.io/docs/getting-started/installation/) to install `k6` on your machine.

3. Navigate to the folder where this is extracted/cloned in the terminal and run:

   ```sh
   BASE_URL=http://localhost:8080 make pull-test
   ```

4. Ensure that Redis is installed and exposed on port **6379**, then run:

   ```sh
   BASE_URL=http://localhost:8080 k6 run load_test.js
   ```

---

## Installing to Production

Before proceeding, start the VPN:

```sh
sudo openvpn --config /path/to/your/config.ovpn
```
Replace `/path/to/your/config.ovpn` with the correct path to your `.ovpn` configuration file.

1. **Update Go modules**:

   Before building the production binary, update Go modules:

   ```sh
   task
   ```

2. **Build the application for production**:

   ```sh
   task build
   ```

3. **Upload the binary to your EC2 instance using SCP**:

   ```sh
   scp -i /path/to/your-key.pem mybinary ubuntu@<EC2_PUBLIC_IP>:/home/ubuntu/
   ```

   Replace `/path/to/your-key.pem` with the path to your private key, `mybinary` with the binary name, and `<EC2_PUBLIC_IP>` with your EC2 instance’s public IP.

4. **Upload the `.env` configuration file**:

   ```sh
   scp -i /path/to/your-key.pem -r config ubuntu@<EC2_PUBLIC_IP>:/home/ubuntu/
   ```

5. **Login to your EC2 instance**:

   ```sh
   ssh -i /path/to/your-key.pem ubuntu@<EC2_PUBLIC_IP>
   ```

6. **Make the binary executable**:

   ```sh
   chmod +x /home/ubuntu/mybinary
   ```

7. **Run the binary**:

   ```sh
   ./mybinary
   ```

---

## Connecting to Remote Database and Running Migrations

1. **Update Your `.env` File**:  
   Ensure your `.env` file contains the correct production database credentials:
   ```env
   DB_HOST=<PRODUCTION_DB_HOST>
   DB_PORT=<PRODUCTION_DB_PORT>
   DB_USER=<PRODUCTION_DB_USER>
   DB_PASS=<PRODUCTION_DB_PASSWORD>
   DB_NAME=<PRODUCTION_DB_NAME>
   ```

2. **Connect to the Production Database**:  
   ```sh
   task db:connect
   ```

3. **Run Migrations**:  
   ```sh
   task migrate:up
   ```

4. **Rollback Migrations (Optional)**:  
   ```sh
   task migrate:down
   ```

   Or force a specific migration version:
   ```sh
   task migrate:force CLI_ARGS=<VERSION>
   ```

---

## Connecting to the EC2 Instance (Redis Server)

To connect to the EC2 instance where Redis is hosted, use the following command:

```sh
task ec2:connect
```

This will prompt you to select which EC2 instance to connect to. Once connected, you can interact with the Redis server.

> *NOTE: put the key file in the root*

To manually connect using SSH:

```sh
ssh -i /path/to/your-key.pem ubuntu@<EC2_PUBLIC_IP>
```

Replace `/path/to/your-key.pem` with your private key and `<EC2_PUBLIC_IP>` with the public IP of the EC2 instance.

To check if Redis is running:

```sh
redis-cli ping
```

If Redis is running correctly, it should return:

```
PONG
```

---

### Accessing Prometheus and Grafana
**Prometheus UI:** After running docker-compose up, you can access Prometheus at http://localhost:9090.

**Grafana UI:** After running docker-compose up, you can access Grafana at http://localhost:3000 (default credentials: admin / admin).

> *These steps below are automatically done by prometheus.yml and grafana folder in the deploy folder.*

**Add Prometheus as Data Source in Grafana (manually):**

`Go to Configuration → Data Sources → Add Data Source → Select Prometheus.`

Set the URL to http://prometheus:9090 (the name of the Prometheus service in docker-compose.yml).

**Create Grafana Dashboard (manually):**

You can now create dashboards with Prometheus queries like:
- http_requests_total
- rate(http_requests_total[1m])
- http_duration_seconds

