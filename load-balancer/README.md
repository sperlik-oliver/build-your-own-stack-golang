# Minimal load balancer

This load balancer is not meant to be used in production. It is a learning project to understand how a minimal load balancer can be programmed using Go. Feel free to use the code for learning purposes.

## The load balancer supports a basic set of features:
- Retry mechanism with configurable interval
- Configurable load balancing modes: ``round robin`` or ``least connections``
- Active async health checking of upstream services with configurable interval

## Build instructions:
Docker: ``docker build . -t load-balancer && docker run load-balancer`` </br>
Native: ``go run .`` or ``go build . -o load-balancer && ./load-balancer``