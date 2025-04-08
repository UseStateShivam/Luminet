# luminet
``` Plain Text
            Public Client
                 │
                 ▼
        🌍 Public Internet
                 │
                 ▼
        🌐 Tunneling Server
            - Listens for public connections
            - Forwards data between phone & local server
                 │
                 ▼
        🏠 Local Network
                 │
                 ▼
        📞 Local TCP Server
            - Listens on port 8000
            - Receives & responds to messages
```

## What?
Luminet is a developer tool that allows you to expose a local server to the internet securely. It's super handy when you're working on a web app locally and want to test it with others or integrate third-party services like webhooks (e.g., from Stripe, GitHub, or Twilio).
In simple terms:
Luminet creates a public URL (like https://xyz123.luminet.io) that tunnels requests to your local machine, like http://localhost:3000.
🔧 Typical Use Case:
You're running a local server on localhost:3000.
You run luminet http 3000.
Ngrok gives you a public URL: https://your-tunnel-id.luminet.io.
Now, anyone with that URL can access your local app.
🚀 Why it's useful:
Testing webhooks: Services like Stripe or GitHub can hit your local server through the public URL.
Sharing a dev preview: Show your work to teammates or clients without deploying.
Mobile app testing: Access your local backend from a phone or emulator.
🔄 What is a Reverse Proxy?
A reverse proxy is a server that sits in front of one or more backend servers, intercepts requests from clients, and forwards them to the right backend — and then returns the response back to the client.
In simple terms:
It’s like a middleman that hides and protects your actual backend servers.

## Why?
💡Tunneling means creating a secure "path" or "channel" from one network (like the internet) to another (like your local machine) — even if that local machine is behind a firewall, NAT, or not publicly accessible.
Think of it as:
"Hey, internet! Here's a secret shortcut to reach my localhost."
🧱 Why do we need tunneling?
Because your local machine is:
Not assigned a public IP.
Protected by firewalls or NAT (Network Address Translation).
Not reachable from the internet directly.
So if you're running a server on localhost:3000, nobody else can access it unless you set up port forwarding, DNS, security, and a bunch of stuff.
Luminet solves that by:
Running a lightweight agent on your machine.
That agent opens a tunnel to Luminet's servers (which do have public IPs).
Luminet then gives you a public URL that forwards all traffic through the tunnel to your local server.
🔄 Real-world Example:
You're building a payment system using Stripe, and Stripe needs to send a webhook to your server when a payment is successful.
Your server is running on http://localhost:3000/webhook.
Stripe can’t access your localhost directly.
Luminet gives you a public URL like https://abc123.luminet.io/webhook.
You plug that into Stripe’s dashboard, and boom — your local server gets the webhook.

## How?
🧠 How Luminet Works Under the Hood
When you run luminet http 3000, here's what happens:
1. Luminet Agent Connects to Luminet Server
The Luminet client (agent) on your machine opens a persistent, encrypted TLS connection to Luminet's edge servers.
It tells Luminet:
"Expose my local port 3000 to the public internet."
2. Luminet Assigns a Public URL
Luminet’s infrastructure assigns a unique subdomain (e.g., https://abc123.luminet.io)
This subdomain maps directly to your tunnel session.
3. Reverse Proxy (Tunnel) in Action
Someone accesses https://abc123.luminet.io.
Luminet receives the request on its public-facing server.
It routes that request through the open tunnel back to your local server.
Your app processes the request and responds.
The response travels back up the tunnel to the user.
It’s all secure, fast, and seamless.
👥 How Luminet Supports Multiple Users
Luminet is built as a multi-tenant platform, meaning:
Every user session is isolated.
Luminet maintains a registry that maps subdomains to individual connections:
``` json
{
  "abc123.luminet.io": "sessionID_1",
  "user456.luminet.io": "sessionID_2"
}
```
So even if thousands of users run Luminet at the same time, each tunnel is independent and secure.
🔐 Under-the-Hood Features
Feature	            How Luminet Handles It
Traffic             Routing	Requests hit Luminet's edge, then go down your tunnel.
Security	        TLS encryption from the internet → edge → your machine.
Debug Tools	        Luminet provides a local inspector (e.g., localhost:4040) for real-time request logs.
Auth & Limits	    Accounts, rate limits, reserved subdomains, IP whitelisting, and more.
Load Handling	    Luminet scales across global edge servers for performance.
🔍 Bonus: Why the Word "Tunnel"?
Because it’s literally acting like a secure pipe through firewalls and NATs.
Your local server is otherwise “hidden” — Luminet carves out a hidden route from the outside world directly to your dev environment.
