# ngrok
ngrok is a tool that allows you to expose your local development server to the internet by creating a secure tunnel.
It can be used to test how cookies behave over HTTPS if you're running a local HTTP server.

## Initial setup
Install using homebrew
```bash
brew install ngrok/ngrok/ngrok
```
Add auth token. Find the token at ngrok dashboard
```bash
ngrok config add-authtoken <your_auth_token>
```

## Tunneling Local Server with ngrok
Run local server
```bash
flask run --port 5000
```

Expose local server with ngrok
```bash
ngrok http 5000
```

