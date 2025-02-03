# Postage

## Goals
- Acts as a client forward proxy, allowing backend services to send clients messages in real time.
- Has an admin ui that shows the current state of the proxy. (Clients connected, mock messaging, etc.)
- Handles authentication of clients based on hmac signed tokens.
- Presence hook. Allow backend to register an endpoint that will receive post requests for connected status updates. inform backend services if a client is connected. (answer: is the client online?)
  - Also need a queryable endpoint to check for status as well.
- Scales horizontally to allow for millions of connected devices. (work on this when everthing else is working).
