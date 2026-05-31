# API Versioning Demo - Go

This project demonstrates three common REST API versioning strategies using Go's standard `net/http` package.

It mirrors the Java Spring Boot API versioning demo, but implements the same behaviour in a lightweight Go application with no external web framework.

The three versioning strategies demonstrated are:

1. URI path versioning
2. Custom request header versioning
3. Vendor media type versioning

---

## Technology Stack

* Go
* Standard library only
* `net/http`
* `encoding/json`

---

## What This Project Demonstrates

The application exposes several user endpoints that return different JSON responses depending on the versioning strategy used.

The examples simulate a common API evolution:

Version 1 returns a single `name` field.

Version 2 splits the name into `firstName` and `lastName`, and adds a `status` field.

---

## Running the Application

From the project root:

```bash
go run .
```

The server starts on:

```text
http://localhost:8080
```

You should see:

```text
Server running on http://localhost:8080
```

---

## Endpoint Summary

| Versioning Type | Method | Endpoint                  | Version Selection                                                                |
| --------------- | -----: | ------------------------- | -------------------------------------------------------------------------------- |
| URI path        |    GET | `/api/v1/users/{id}`      | Path contains `v1`                                                               |
| URI path        |    GET | `/api/v2/users/{id}`      | Path contains `v2`                                                               |
| Header          |    GET | `/api/headers/users/{id}` | `API-Version: 1` or `API-Version: 2`                                             |
| Media type      |    GET | `/api/media/users/{id}`   | `Accept: application/vnd.demo.v1+json` or `Accept: application/vnd.demo.v2+json` |

---

## 1. URI Path Versioning

The version is part of the URL.

```text
/api/v1/users/{id}
/api/v2/users/{id}
```

### Version 1

```bash
curl http://localhost:8080/api/v1/users/123
```

Response:

```json
{
  "id": "123",
  "name": "Alice Smith"
}
```

### Version 2

```bash
curl http://localhost:8080/api/v2/users/123
```

Response:

```json
{
  "firstName": "Alice",
  "id": "123",
  "lastName": "Smith",
  "status": "ACTIVE"
}
```

URI versioning is the simplest and most visible strategy. It is easy to test, easy to document, and easy to understand.

---

## 2. Header Versioning

The URL stays the same, but the version is specified using a custom HTTP request header.

```text
/api/headers/users/{id}
```

### Version 1

```bash
curl -H "API-Version: 1" \
  http://localhost:8080/api/headers/users/123
```

Response:

```json
{
  "id": "123",
  "name": "Bob Smith"
}
```

### Version 2

```bash
curl -H "API-Version: 2" \
  http://localhost:8080/api/headers/users/123
```

Response:

```json
{
  "firstName": "Bob",
  "id": "123",
  "lastName": "Smith",
  "status": "ACTIVE"
}
```

### Missing or Unsupported Header

```bash
curl http://localhost:8080/api/headers/users/123
```

Response:

```text
Missing or unsupported API-Version header
```

HTTP status:

```text
400 Bad Request
```

Header versioning keeps URLs clean, but the version is less visible because it is hidden in the request headers.

---

## 3. Vendor Media Type Versioning

The URL stays the same, but the requested version is specified using the HTTP `Accept` header.

```text
/api/media/users/{id}
```

### Version 1

```bash
curl -H "Accept: application/vnd.demo.v1+json" \
  http://localhost:8080/api/media/users/123
```

Response:

```json
{
  "id": "123",
  "name": "Paaa Smith"
}
```

### Version 2

```bash
curl -H "Accept: application/vnd.demo.v2+json" \
  http://localhost:8080/api/media/users/123
```

Response:

```json
{
  "firstName": "Paas",
  "id": "123",
  "lastName": "Smith",
  "status": "ACTIVE"
}
```

### Missing or Unsupported Accept Header

```bash
curl http://localhost:8080/api/media/users/123
```

Response:

```text
Unsupported media type version
```

HTTP status:

```text
406 Not Acceptable
```

Vendor media type versioning uses HTTP content negotiation. The URL identifies the resource, while the `Accept` header identifies the representation requested by the client.

---

## PowerShell Note

On Windows PowerShell, `curl` may be an alias for `Invoke-WebRequest`.

If this causes problems, use `curl.exe` instead:

```powershell
curl.exe -H "API-Version: 1" http://localhost:8080/api/headers/users/123
```

Or use PowerShell syntax:

```powershell
Invoke-WebRequest `
  -Uri "http://localhost:8080/api/headers/users/123" `
  -Headers @{ "API-Version" = "1" }
```

For media type versioning:

```powershell
curl.exe -H "Accept: application/vnd.demo.v1+json" http://localhost:8080/api/media/users/123
```

---

## Code Structure

The application uses Go's standard HTTP request multiplexer.

```go
http.HandleFunc("/api/v1/users/", uriV1)
http.HandleFunc("/api/v2/users/", uriV2)
http.HandleFunc("/api/headers/users/", headerVersioned)
http.HandleFunc("/api/media/users/", mediaVersioned)
```

Each handler extracts the user ID from the URL path:

```go
func lastPathPart(path string) string {
    parts := strings.Split(strings.Trim(path, "/"), "/")
    return parts[len(parts)-1]
}
```

Responses are written as JSON using:

```go
func writeJSON(w http.ResponseWriter, data map[string]any) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
```

---

## Building the Application

Build a binary:

```bash
go build
```

Run the generated binary:

```bash
./go-api-version-demo
```

If you want to specify the output name:

```bash
go build -o api-version-demo
./api-version-demo
```

---

## Creating the Project from Scratch

```bash
mkdir go-api-version-demo
cd go-api-version-demo
go mod init example.com/go-api-version-demo
touch main.go
```

Paste the Go code into `main.go`.

Run:

```bash
go run .
```

---

## Comparison with the Java Spring Boot Version

The Java Spring Boot version uses annotations such as:

```java
@GetMapping("/v1/users/{id}")
@GetMapping(value = "/headers/users/{id}", headers = "API-Version=1")
@GetMapping(value = "/media/users/{id}", produces = "application/vnd.demo.v1+json")
```

Spring automatically routes requests based on the path, headers, and produced media type.

The Go version does the routing manually inside handler functions.

For example, header versioning is handled using:

```go
version := reader.Header.Get("API-Version")
```

and then switching on the value:

```go
switch version {
case "1":
    // return v1 response
case "2":
    // return v2 response
default:
    http.Error(writer, "Missing or unsupported API-Version header", http.StatusBadRequest)
}
```

This makes the Go implementation more explicit and easier to see mechanically, while the Spring Boot implementation is more declarative.

---

## Versioning Strategy Comparison

| Strategy              | Pros                                 | Cons                                       |
| --------------------- | ------------------------------------ | ------------------------------------------ |
| URI path versioning   | Simple, visible, easy to test        | Version is part of the URL                 |
| Header versioning     | Clean URLs                           | Less discoverable, harder to test manually |
| Media type versioning | Good use of HTTP content negotiation | More complex and less familiar             |

For most simple APIs, URI path versioning is the easiest to explain and operate.

Header and media type versioning are useful to understand because they show how API behaviour can be selected using request metadata rather than only URL paths.
