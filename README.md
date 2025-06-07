# Authorization Microservice

## Project Overview

This project provides a simple Go-based authorization microservice. The service acts as an authorization server, validating incoming requests against a set of rules defined in a `rules.yaml` file. It exposes a single endpoint (`/authorize`) to determine if a user (identified by a role) is allowed to access a specific resource with a given method.

The authorization decision is based solely on the `role` field provided in the JSON payload of the `/authorize` request and the rules defined in `rules.yaml`. It does not perform any prior authentication of the HTTP request itself (e.g., token validation).

## API Endpoint

### `/authorize`

*   **Method:** `POST`
*   **Description:** Checks if a given role is authorized to perform an action on a resource.
*   **Request Body (JSON):**
    ```json
    {
        "resource": "/some/resource",
        "method": "GET",
        "role": "user_role"
    }
    ```
*   **Response Body (JSON):**
    *   On successful authorization check (authorized or not):
        ```json
        {
            "authorized": true
        }
        ```
        or
        ```json
        {
            "authorized": false
        }
        ```
    *   On invalid request (e.g., malformed JSON):
        HTTP Status Code `400 Bad Request` with an error message.

## Rules Configuration (`rules.yaml`)

The service uses a `rules.yaml` file to define authorization rules. This file must be present in the same directory where the service is run.

**Structure:**

The YAML file should contain a list of rules under the `rules` key. Each rule specifies:
*   `resource`: The path or identifier of the resource.
*   `allowed_methods`: A list of HTTP methods permitted for this resource and roles.
*   `roles`: A list of roles that are granted access according to the `allowed_methods`.

**Example `rules.yaml`:**

```yaml
rules:
  - resource: /example/resource
    allowed_methods:
      - GET
      - POST
    roles:
      - admin
      - editor
  - resource: /another/resource
    allowed_methods:
      - GET
    roles:
      - viewer
```

## Getting Started

1.  **Prerequisites:**
    *   Ensure you have Go installed (version 1.18 or newer recommended).

2.  **Clone the Repository (if applicable):**
    ```bash
    # If you haven't cloned yet
    # git clone <repository-url>
    # cd <repository-directory>
    ```

3.  **Build the Service:**
    Navigate to the project root directory (where `main.go` and `go.mod` are located) and run:
    ```bash
    go build -o authservice_app main.go
    ```
    This will create an executable file named `authservice_app`.

4.  **Prepare `rules.yaml`:**
    Create a `rules.yaml` file in the same directory as the `authservice_app` executable, following the structure described above.

5.  **Run the Service:**
    ```bash
    ./authservice_app
    ```
    The service will start, typically listening on port 8080 (this can be configured in `main.go`).

## Running Tests

To run the automated tests for the `authservice` package, navigate to the project root and execute:
```bash
go test ./...
```
This command will discover and run all tests in the current directory and its subdirectories.

## Contributing

Please read `CONTRIBUTING.md` for details on how to contribute to this project.

## License

This project is licensed under the MIT License - see the `LICENSE` file for details.
