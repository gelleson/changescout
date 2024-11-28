# ChangeScout: Website Change Detection Service

ChangeScout is currently under active development as a website change detection service developed in Go. Inspired by [ChangeDetection.io](https://changedetection.io/), ChangeScout delivers an efficient and reliable solution for monitoring changes on websites. This project focuses on backend development in Go to leverage the language's capabilities for scalability and efficiency, with the next step aimed at developing the user interface in React.

## Key Features

- **Reliable Change Detection:** Accurately detects changes in website content using robust comparison algorithms and supports various content types.
- **Flexible Scheduling:** Monitor websites at customizable intervals using cron expressions.
- **Multiple Notification Methods:** Receive notifications via Telegram (other methods planned for future releases).
- **User Authentication and Authorization:** Secure management of users and their access to monitored websites.
- **Modern and Clean Design:** Follows best practices for Go development with a clean codebase and design.
- **High Performance:** Designed for performance and scalability to handle numerous websites.
- **Database Persistence:** Persistent storage for website configurations, check results, and user data.

## Architecture

ChangeScout features a microservice-like architecture with distinct components:

- **API:** A GraphQL API offers a flexible interface for interacting with the service.
- **Core Logic:** Handles website checking, change detection, and notification dispatch.
- **Database:** Utilizes an Entity-Relationship model for data persistence.
- **Message Broker (Optional):** Uses a message broker for asynchronous processing of website checks.
- **Scheduler (Optional):** Schedules website checks periodically based on user-defined cron expressions.

## Technology Stack

- **Go:** Programming language.
- **GraphQL:** API specification.
- **GQLGen:** GraphQL server library.
- **Echo:** Web framework.
- **Ent:** ORM for database interactions.
- **Watermill:** Message broker library (optional).
- **goquery:** HTML parsing library.
- **go-difflib & go-diff:** Libraries for efficient diff calculation.
- **Zap:** Logging library.
- **bcrypt:** Password hashing library.
- **golang-jwt:** JWT library for authentication.

## Installation

1. **Clone the Repository:**

    ```bash
    git clone https://github.com/gelleson/changescout.git
    cd changescout
    ```

2. **Install Dependencies:**

    ```bash
    go mod download
    ```

3. **(Optional) Database Setup:** The application supports SQLite3, PostgreSQL, and MySQL. Configure the database connection string as needed. Refer to `config.json` for details.

4. **Run the Application:**

    ```bash
    go run main.go start --help  // View available flags
    ```

## Usage

Access the application primarily through the GraphQL API. Refer to the [GraphQL schema](https://github.com/gelleson/changescout/blob/main/changescout/internal/api/gql/schema/schema.graphqls) for details on available queries and mutations. The GraphQL playground is available at `/playground` once the server is running.

## Contributing

Contributions are welcome! Please open issues for bug reports or feature requests. Pull requests are also welcome, following the established coding style and after necessary testing.

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Future Enhancements

- Support for additional notification methods (email, webhook, etc.).
- Improved UI/UX, with upcoming development focusing on a React-based web interface.
- Advanced features like screenshot comparison.
- Support for more sophisticated selectors.
- Enhanced error handling and reporting.
- Containerization (Docker).

ChangeScout is actively being developed, with new features being added regularly. We are committed to creating a powerful, reliable, and user-friendly website monitoring tool, and the next big step is to build an intuitive user interface in React.