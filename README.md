This project provides a simple Go-based authorization microservice.

**Functionality:**

The service acts as an authorization server, validating incoming requests against a set of rules defined in a `rules.yaml` file. It exposes a single endpoint (`/authorize`) to determine if a user (identified by a role) is allowed to access a specific resource with a given method.

**Getting Started:**

1.  **Install Dependencies:**
    Ensure you have Go installed. Then, in the project directory, run:
