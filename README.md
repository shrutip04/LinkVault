```text     _       _     __     __            _ _
| |   (_)_ __ | | __ \ \   / /_ _ _   _| | |_
| |   | | '_ \| |/ /  \ \ / / _` | | | | | __|
| |___| | | | |   <    \ V / (_| | |_| | | |_
|_____|_|_| |_|_|\_\    \_/ \__,_|\__,_|_|\__|

```

# 🔗 LinkVault

<p align="center">
  <strong>A Full-Stack URL Shortening & Link Management Platform</strong>
</p>

<p align="center">
  Create • Protect • Organize • Track
</p>

<p align="center">
  <img src="https://skillicons.dev/icons?i=go,react,js,sqlite,css,git,github" alt="Tech Stack" />
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Backend-Go%20%7C%20Gin-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Frontend-React-61DAFB?style=for-the-badge&logo=react&logoColor=black" alt="React" />
  <img src="https://img.shields.io/badge/Database-SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite" />
  <img src="https://img.shields.io/badge/API-REST-6E6E6E?style=for-the-badge" alt="REST API" />
</p>

---

## 📌 Overview

**LinkVault** is a full-stack web application that allows users to create, manage, protect, and track shortened URLs.

Unlike a basic URL shortener, LinkVault treats each shortened URL as a manageable resource with features such as **authentication, password protection, expiration, categorization, and click tracking**.

The application uses a **React frontend**, a **Go/Gin REST API**, and **SQLite** for persistent data storage.

---

## ✨ Features

<table>
<tr>
<td width="33%" align="center">

### 🔐 Authentication

User registration and login with securely hashed passwords.

</td>
<td width="33%" align="center">

### 🔗 URL Shortening

Generate unique short URLs and redirect users to the original destination.

</td>
<td width="33%" align="center">

### 🛡️ Protected Links

Add an optional password to restrict access to a shortened URL.

</td>
</tr>

<tr>
<td width="33%" align="center">

### ⏳ Expiration

Set an expiration time for links and prevent access after expiry.

</td>
<td width="33%" align="center">

### 🗂️ Categories

Organize links into categories for easier management.

</td>
<td width="33%" align="center">

### 📊 Click Tracking

Track total clicks and the last time a link was accessed.

</td>
</tr>
</table>

---

## 🏗️ Architecture

```text
                         ┌──────────────────────┐
                         │    React Frontend    │
                         │     linkvault-ui     │
                         └──────────┬───────────┘
                                    │
                              REST / JSON
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │     Go + Gin API     │
                         │                      │
                         │  Routes → Middleware │
                         │      → Handlers      │
                         │      → Models        │
                         └──────────┬───────────┘
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │        SQLite        │
                         │     Persistent DB    │
                         └──────────────────────┘
```

### Request Flow

```text
User
 │
 ▼
React UI
 │
 │ HTTP Request / JSON
 ▼
Gin Router
 │
 ▼
Authentication Middleware
 │
 ▼
Handler
 │
 ▼
SQLite
 │
 ▼
JSON Response
 │
 ▼
React UI
```

---

## 🛠️ Tech Stack

<table>
<tr>
<th>Layer</th>
<th>Technology</th>
<th>Purpose</th>
</tr>

<tr>
<td><strong>Frontend</strong></td>
<td>⚛️ React · JavaScript · CSS</td>
<td>User interface and client-side interaction</td>
</tr>

<tr>
<td><strong>Backend</strong></td>
<td>🐹 Go · Gin</td>
<td>REST API and application logic</td>
</tr>

<tr>
<td><strong>Database</strong></td>
<td>🗄️ SQLite</td>
<td>Persistent storage</td>
</tr>

<tr>
<td><strong>Security</strong></td>
<td>🔐 bcrypt · Middleware</td>
<td>Password hashing and authentication</td>
</tr>

<tr>
<td><strong>Communication</strong></td>
<td>🌐 REST · JSON · CORS</td>
<td>Frontend-backend communication</td>
</tr>

<tr>
<td><strong>Tools</strong></td>
<td>Git · GitHub · npm</td>
<td>Version control and development</td>
</tr>
</table>

---

## 📁 Project Structure

<table>
<tr>
<th>Backend — Go</th>
<th>Frontend — React</th>
</tr>

<tr>
<td valign="top">

```text
config/
database/
handlers/
middleware/
models/
routes/
static/
templates/
utils/

main.go
go.mod
go.sum
```

</td>

<td valign="top">

```text
linkvault-ui/
│
├── public/
└── src/
    ├── components/
    ├── pages/
    ├── api.js
    ├── App.js
    └── index.js
```

</td>
</tr>
</table>

### Backend Components

| Component     | Responsibility                            |
| ------------- | ----------------------------------------- |
| `routes/`     | Defines API endpoints                     |
| `handlers/`   | Processes HTTP requests                   |
| `middleware/` | Handles authentication and request checks |
| `models/`     | Represents application data               |
| `database/`   | Initializes and manages database access   |
| `utils/`      | Shared utility functionality              |
| `config/`     | Application configuration                 |
| `main.go`     | Backend entry point                       |

### Frontend Components

| Component     | Responsibility                                           |
| ------------- | -------------------------------------------------------- |
| `pages/`      | Login, registration, dashboard and link creation screens |
| `components/` | Reusable UI components                                   |
| `api.js`      | API communication                                        |
| `App.js`      | Application routing and structure                        |
| `public/`     | Static frontend assets                                   |

---

## 🗄️ Database Design

LinkVault uses SQLite with two primary entities:

<table>
<tr>
<th>👤 users</th>
<th>🔗 links</th>
</tr>

<tr>
<td valign="top">

```text
id
username
email
password
created_at
```

</td>

<td valign="top">

```text
id
user_id
original
short
clicks
created_at
last_accessed
expires_at
password
category
```

</td>
</tr>
</table>

### Relationship

```text
┌──────────────┐             ┌──────────────┐
│     User     │  1      N  │     Link     │
│              │────────────▶│              │
│     id       │             │   user_id    │
└──────────────┘             └──────────────┘
```

Each shortened link belongs to a user through `user_id`.

---

## 🔐 Security

<table>
<tr>
<td width="33%" align="center">

### 🔑 Password Hashing

User passwords are hashed using **bcrypt** before being stored.

</td>

<td width="33%" align="center">

### 🛡️ Authentication

Protected operations are handled through authentication middleware.

</td>

<td width="33%" align="center">

### 🔒 Link Protection

Individual links can have an additional password.

</td>
</tr>
</table>

### Password Flow

```text
User Password
      │
      ▼
    bcrypt
      │
      ▼
 Password Hash
      │
      ▼
   SQLite
```

Passwords are not stored as plaintext.

---

## 🔄 Application Flow

### Creating a Link

```text
User
 │
 ▼
Create Link Form
 │
 ▼
React API Request
 │
 ▼
Go / Gin Handler
 │
 ▼
Validate Request
 │
 ▼
Generate Short URL
 │
 ▼
Store Link in SQLite
 │
 ▼
Return Link Details
 │
 ▼
Display in Dashboard
```

### Accessing a Short URL

```text
Short URL
    │
    ▼
Go Backend
    │
    ▼
Find Link
    │
    ├── Expired? ──▶ Reject
    │
    ├── Protected? ──▶ Verify Password
    │
    ▼
Increment Click Count
    │
    ▼
Update Last Accessed
    │
    ▼
Redirect to Original URL
```

---

## ⚙️ Getting Started

### Prerequisites

| Requirement                    | Purpose                         |
| ------------------------------ | ------------------------------- |
| [Go](https://go.dev/)          | Run the backend                 |
| [Node.js](https://nodejs.org/) | Run the React frontend          |
| npm                            | Install frontend dependencies   |
| Git                            | Clone and manage the repository |

---

### 1. Clone the Repository

```bash
git clone https://github.com/shrutip04/LinkVault.git
cd LinkVault
```

---

### 2. Start the Backend

```bash
go mod download
go run main.go
```

Backend:

```text
http://localhost:8080
```

---

### 3. Configure the Frontend

Open another terminal:

```bash
cd linkvault-ui
```

Create a `.env` file:

```env
REACT_APP_API_URL=http://localhost:8080
```

---

### 4. Install Dependencies

```bash
npm install
```

---

### 5. Start the Frontend

```bash
npm start
```

Frontend:

```text
http://localhost:3000
```

---

## 🚀 Running LinkVault

Both applications need to run simultaneously.

<table>
<tr>
<th>🐹 Backend</th>
<th>⚛️ Frontend</th>
</tr>

<tr>
<td>

```bash
cd LinkVault
go run main.go
```

<strong>Port:</strong> 8080

</td>

<td>

```bash
cd LinkVault/linkvault-ui
npm start
```

<strong>Port:</strong> 3000

</td>
</tr>
</table>

Open **http://localhost:3000** in your browser.

---

## 🌐 Frontend ↔ Backend Communication

The frontend communicates with the Go backend through REST API requests.

```text
┌─────────────────┐
│  React :3000    │
└────────┬────────┘
         │
         │ HTTP / JSON
         ▼
┌─────────────────┐
│   Go + Gin      │
│     :8080       │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│     SQLite      │
└─────────────────┘
```

CORS is configured on the backend to allow communication between the frontend development server and API.

---

## 📡 API Overview

The application exposes REST endpoints for authentication and link management.

| Area               | Operations                   |
| ------------------ | ---------------------------- |
| **Authentication** | Register · Login             |
| **Links**          | Create · Retrieve · Redirect |
| **Management**     | View and manage user links   |
| **Tracking**       | Click count · Last accessed  |
| **Protection**     | Link passwords · Expiration  |
| **Organization**   | Categories                   |

---

## 📈 Project Status

<table>
<tr>
<th>Area</th>
<th>Status</th>
</tr>

<tr>
<td>Go Backend</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>Gin REST API</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>SQLite Database</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>React Frontend</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>Authentication</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>URL Shortening</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>Password-Protected Links</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>Link Expiration</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>Categories</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>Click Tracking</td>
<td>✅ Complete</td>
</tr>

<tr>
<td>Frontend ↔ Backend Integration</td>
<td>✅ Complete</td>
</tr>

</table>

---

## 🗺️ Future Improvements

<table>
<tr>
<td>

* 🔎 Link search & filtering
* ✏️ Custom short aliases
* 📊 Advanced analytics
* 📱 QR code generation
* 🚦 API rate limiting

</td>

<td>

* 🧪 Unit & integration tests
* 📖 API documentation
* 🐳 Docker support
* ⚙️ CI/CD pipeline
* ☁️ Production deployment

</td>
</tr>
</table>

---

## 🎯 Key Technical Concepts

LinkVault demonstrates practical implementation of:

`REST APIs` · `Go` · `Gin` · `React` · `SQLite` · `Authentication` · `bcrypt` · `Middleware` · `CORS` · `JSON` · `Database Relationships` · `Git`

---

## 👩‍💻 Author

<p align="center">
  <strong>Shruti Pawar</strong><br>
  Computer Engineering Student
</p>

<p align="center">
  <a href="https://github.com/shrutip04">GitHub</a>
  &nbsp;•&nbsp;
  <a href="https://github.com/shrutip04/LinkVault">LinkVault Repository</a>
</p>

---

<p align="center">
  <strong>🔗 LinkVault — Create. Protect. Organize. Track.</strong>
</p>

