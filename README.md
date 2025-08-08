# Advanced URL Shortening Service

A distributed, sharded, cloud-ready URL shortening service built with Go, supporting scalable database sharding, Redis caching, and Kubernetes deployment.

---

## Features

- **Scalable & Distributed:** Supports sharded databases for horizontal scaling.
- **Cloud Ready:** Easily deployable to AWS or any Kubernetes cluster.
- **High Performance:** Uses Redis for caching and rate limiting.
- **RESTful API:** Built using Gin framework.
- **Easy Local Development:** `.env` based config for local runs.

---

## Local Development

1. **Clone the repo**
2. **Create a `.env` file** in the root:
    ```env
    ENV=local
    DATABASE_URL=postgres://postgres:admin@localhost:5432/URL?sslmode=disable
    REDIS_URL=localhost:6379
    REDIS_PWD=
    ```
3. **Run the app:**
    ```bash
    go run main.go
    ```
4. **API Endpoints:**
    - `POST /shorten` – Shorten a new URL
    - `GET /shorten/:short_url` – Retrieve the original URL
    - `GET /shorten/:short_url/count` – Get redirect count
    - `DELETE /shorten/:short_url` – Delete a short URL

---

## Docker

Build and run the image:

```bash
docker build -t url-shortener .
docker run --env-file .env -p 8080:8080 url-shortener
```

---

## Kubernetes Deployment

### 1. Prepare Secrets

Encode each database URL and Redis password using base64:
```bash
echo -n "postgres://user:pass@host:5432/db" | base64
echo -n "your_redis_password" | base64
```

Create secrets for each shard, replica, index DB, and Redis. Example:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-shard-1
type: Opaque
data:
  url: <base64-encoded-postgres-url>
```
_(Repeat for db-shard-2 to db-shard-5, db-replica-1 to db-replica-5, db-index, redis-credentials)_

### 2. Apply Secrets
```bash
kubectl apply -f secret-db-shard-1.yaml
kubectl apply -f secret-db-shard-2.yaml
...
kubectl apply -f secret-db-index.yaml
kubectl apply -f secret-redis.yaml
```

### 3. Deploy Application

Apply deployment and service manifests:

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
```

### 4. Environment Variables Used in Cloud

| Variable                          | Purpose         | Example Value                 |
|------------------------------------|-----------------|-------------------------------|
| ENV                               | Environment     | cloud                         |
| DATABASE_SHARD_1_URL ... _5       | Write shards    | postgres://...                |
| DATABASE_SHARD_REPLICA_1_URL..._5 | Read replicas   | postgres://...                |
| DATABASE_INDEX_URL                | Index DB        | postgres://...                |
| REDIS_URL                         | Redis host:port | redis-master:6379             |
| REDIS_PWD                         | Redis password  | your_redis_password           |

---

## API Reference

### Shorten a URL
```http
POST /shorten
{
  "original_url": "https://example.com"
}
```
Response:
```json
{
  "id": 123,
  "short_url": "abc123",
  "original_url": "https://example.com"
}
```

### Retrieve Original URL
```http
GET /shorten/abc123
```
Response:
```json
{
  "original_url": "https://example.com"
}
```

### Get Redirect Count
```http
GET /shorten/abc123/count
```
Response:
```json
{
  "url": "abc123",
  "count": 42
}
```

### Delete a Short URL
```http
DELETE /shorten/abc123
```
Response:
```json
{
  "message": "Short URL deleted successfully"
}
```

---

## Contributing

Pull requests and issues are welcome!  
Please open an issue for any bugs or feature requests.

---

## License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0).  
See the [LICENSE](LICENSE) file for details.