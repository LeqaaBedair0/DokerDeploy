# Mini Production Platform

This project simulates a real-world production environment using Microservices, Docker, and Ansible. [cite_start]It demonstrates microservices architecture, Docker multi-stage builds, service orchestration, and secure infrastructure automation [cite: 3-5, 8-15].

## Architecture Diagram

```mermaid
graph TD
  Client[Client/Browser]
  Client -->|HTTP| Auth[Auth Service Go]
  Auth -->|JWT| Task[Task Service Node.js]
  Task -->|Event / HTTP| Notify[Notify Service Go]
  Auth --> DB[(MySQL Database)]
  Task --> DB
  subgraph Docker Network
    Auth
    Task
    Notify
    DB
  end```

## Technologies Used
Auth Service: Go (Gin), JWT, bcrypt

Task Service: Node.js (Express)

Notify Service: Go (Standard Library)

Database: MySQL 8.4

Containerization: Docker & Docker Compose

Automation & Security: Ansible & Ansible Vault

## How to run locally using Docker Compose
Clone the repository.

Copy the example environment file: cp .env.example .env and fill in your local variables.

Build and start the containers:
   docker compose up -d --build

Access the Auth Service on port 8000 and the Task Service on port 3000.

## How to deploy using Ansible (Production)
This project uses Ansible to automate the deployment to a clean VM, including installing Docker, configuring UFW firewalls, and deploying the application using encrypted secrets.

Ensure you have Ansible installed on your control node.

Update the ansible/inventories/dev file with your target server's IP address.

Run the playbook with your Ansible Vault password and Sudo password:
   cd ansible
   ansible-playbook -i inventories/dev site.yml -K --ask-vault-pass
