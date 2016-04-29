# dockroute

dockroute is a router proxy for internal http service communication between docker containers.
It is currently heavily WORK IN PROGRESS.

## Key Idea

dockroute dynamically discovers the available service over several backends. Than it takes the
http Host header, of a request to determine the target upstream service and port for a call.


