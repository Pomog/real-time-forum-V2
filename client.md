## To run the Client server which used HTML to interact with the API server, execute:
```bash
go run ./cmd/client/main.go
```

## To run the API server, open a new terminal and execute:
```bash
go run ./cmd/api/main.go 
```
## Key Components of the Client Server
1. Reading the configurations:
   The Client server loads configuration from the config.json.

2. Parsing HTML (package text/template):
   All the HTML content, and DOM structure generates by JavaScript


1. Configuration Loading:
   The Client server reads its configuration from the config.json file.
   The configuration file specifies important settings such as the backend API server address and the port for the client server.

2. Serving HTML and Static Files:
   The client server serves the HTML content from ./web/public/index.html.
   The static files (such as JavaScript, CSS, and images) are served from the ./web/src directory using a file server (http.FileServer).

3. Template Parsing:
   The HTML content is parsed using Go’s text/template package.
   The template is executed with the backend API server address, which is injected into the HTML at runtime.

4. Request Handling:
   The client server listens for incoming requests on the port specified in the configuration file (conf.Client.Port).
   All requests to /src/ are handled by the file server, serving the static content.
   The root path / serves the main HTML page, which includes JavaScript code that dynamically generates and manipulates the DOM structure.

5. The API server should be running and accessible at the address specified in the client server's configuration.
This allows the client-side JavaScript to interact with the API for features like authentication, posting, and real-time communication.