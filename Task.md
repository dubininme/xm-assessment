Technical requirements
Build a microservice to handle companies. It should provide the following operations:
• Create
• Patch
• Delete
• Get (one)
Each company is defined by the following attributes:
• ID (uuid) required
• Name (15 characters) required - unique
• Description (3000 characters) optional
• Amount of Employees (int) required
• Registered (boolean) required
• Type (Corporations | NonProfit | Cooperative | Sole Proprietorship) required
Only authenticated users should have access to create, update and delete companies.
Expectations:
As a deliverable, we expect a GitHub repository (or any other git based repo) with the source
code. We would like the solution to contain clear instructions to set up and execute the project.
We expect the solution to be production ready.
Will be considered a plus:
• On each mutating operation, an event should be produced.
• Dockerize the application to be ready for building the production docker image
• Use docker for setting up the external services such as the database
• REST is suggested, but GRPC is also an option
• JWT for authentication
• Kafka for events
• DB is up to you
• Integration tests are highly appreciated
• Linter
• Configuration file