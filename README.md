# Worker Queue System

This application is a lightweight, concurrent worker queue system implemented in Go, designed to process incoming JSON POST requests and send transformed data to a webhook endpoint. It uses Goroutines for parallel processing and channels for communication between HTTP handlers and workers.

## Features
- **Concurrent Processing**: Supports multiple workers to handle requests concurrently, improving throughput.
- **Dynamic Key Parsing**: Handles dynamic attributes and traits from request payloads.
- **Webhook Integration**: Sends transformed data to a specified webhook URL.
- **REST API**: Provides an endpoint to submit events.

## Hosted API
The application is hosted at:
```
https://test.shahabazsulthan.cloud/event
```

## How It Works
1. **HTTP Endpoint**: Accepts POST requests with event data.
2. **Dynamic Key Parsing**: Parses `Attributes` and `Traits` dynamically based on specific prefixes (`atrk`, `atrv`, `atrt` for attributes; `uatrk`, `uatrv`, `uatrt` for traits).
3. **Worker Queue**: Enqueues incoming requests to a shared channel (`WorkerQueue`), which is processed by worker Goroutines.
4. **Data Transformation**: Transforms the input data into a standardized format.
5. **Webhook Posting**: Sends the transformed data to a predefined webhook URL.

---

## Endpoints

### 1. **POST `/event`**
Processes an event and enqueues it for transformation and delivery.

#### Request Format
```json
{
  "ev": "event_name",
  "et": "event_type",
  "id": "application_id",
  "uid": "user_id",
  "mid": "message_id",
  "t": "page_title",
  "p": "page_url",
  "l": "browser_language",
  "sc": "screen_size",
  "atrk1": "attribute_key1",
  "atrv1": "attribute_value1",
  "atrt1": "attribute_type1",
  "uatrk1": "trait_key1",
  "uatrv1": "trait_value1",
  "uatrt1": "trait_type1"
}
```

#### Response
- **202 Accepted**: When the request is successfully enqueued.
- **400 Bad Request**: When the request body is invalid.
- **405 Method Not Allowed**: For non-POST requests.

---

## Components

### 1. **InputData**
Represents the raw incoming event data.

### 2. **TransformedData**
The structured and standardized format for the event data.

### 3. **Worker**
Processes data from the `WorkerQueue`, transforms it, and sends it to the webhook.

### 4. **Dynamic Key Parsing**
The `ParseDynamicKeys` function extracts dynamic attributes and traits using specified prefixes.

---

## Running Locally

### Prerequisites
- Go installed on your system.
- A webhook endpoint (or use a tool like [Webhook.site](https://webhook.site)).

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/worker-queue-system.git
   cd worker-queue-system
   ```
2. Run the application:
   ```bash
   go run main.go
   ```
3. Test the endpoint using a tool like `curl` or Postman:
   ```bash
   curl -X POST https://localhost:8080/event -H "Content-Type: application/json" -d @data.json
   ```

---

## Configuration
Modify the `webhookUrl` in `PostToWebhook` to use your desired webhook endpoint:
```go
webhookUrl := "https://your-webhook-url"
```

---

## Deployment
This application can be deployed to any cloud provider. A sample configuration for Docker or Kubernetes can be included for production deployment.

---

## Improvements and Notes
- **Scaling**: Increase the number of workers for higher throughput.
- **Error Handling**: Add retry logic for failed webhook calls.
- **Security**: Secure the endpoint using authentication or IP whitelisting.
- **Metrics**: Integrate monitoring tools for performance tracking.

---

Feel free to contribute or raise issues in the repository.
