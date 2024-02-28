# Reverse Proxy in Golang


This reverse proxy, implemented in Golang, functions as an intermediary between the client and the server. Its primary role is to intercept client requests, perform a REST HTTP call, and return the response to the server. Presently, it exclusively supports HTTP GET requests.


## Video Demo

Check out the [video demo](https://1drv.ms/v/s!AlHEv3VthgvFqLxTUSGzLqHYJpf3_Q?e=cPNYQC) for a comprehensive overview and explanation of the project.

## Purpose and Use Case

A reverse web proxy serves as a valuable tool with diverse applications. One significant use case is within an organization, where it can be implemented to monitor the URLs accessed by clients. This capability provides insight into the communication between clients and servers, allowing for enhanced visibility and control over network traffic.

## Features

- **HTTP GET Requests:** The reverse proxy currently supports HTTP GET requests, making it suitable for retrieving information from servers.
  
### Technical Highlights

1. **Mutex Usage for Deadlock Prevention:**
   - The project employs mutexes to solve deadlock scenarios, ensuring robustness and stability in concurrent operations.

2. **Thread Identification:**
   - Each thread involved in the project is labeled by its unique identifier. This labeling aids in tracing and debugging, enhancing the maintainability of the codebase.
     
     ![image](https://github.com/ferrero-rocher/reverse-web-proxy/assets/60911166/50ca73ae-ca47-4902-96d2-8b97eac5f73a)


## How it Works

1. **Interception:** The proxy intercepts incoming HTTP GET requests from the client.
2. **REST HTTP Call:** It makes a corresponding RESTful HTTP call to the intended server to fetch the requested data.
3. **Response Handling:** The proxy then processes the response obtained from the server.
4. **Return to Server:** The processed response is sent back to the original server that initiated the client request.

## Implementation Details

- **Programming Language:** Golang
- **Versatility:** The proxy is designed to handle various use cases, offering flexibility and extensibility.


## Getting Started

To deploy and test the reverse proxy in your environment, follow these steps:

1. Clone the repository: `git clone https://github.com/your/reverse-proxy.git`
2. Navigate to the project directory: `cd reverse-proxy`
3. Build and run the proxy: `go run proxy.go <port-number>` //if we dont specify port number it will run on 1234 port by default

Feel free to explore and adapt the codebase to suit your specific requirements.

## Contribution

Contributions to enhance and expand the functionality of this reverse proxy are welcome. Feel free to open issues, propose features, or submit pull requests.

## License

This project is licensed under the [MIT License](LICENSE), granting you the freedom to use, modify, and distribute the code.

---

