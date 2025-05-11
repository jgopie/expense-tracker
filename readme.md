# Expense Tracker 🚀

*A clean, responsive dashboard for managing finances*

## 🌟 Features

- **User Authentication**  
  Secure signup/login with session management
- **Transaction Management**  
  Add, edit, delete transactions with account assignment
- **Real-Time Balance Tracking**  
  Automatic account balance synchronization
- **Data Export**  
  Download transactions as CSV for external analysis
- **Responsive UI**  
  Works on desktop and mobile devices

## 🛠 Tech Stack

| Category       | Technologies Used                     |
|----------------|---------------------------------------|
| **Backend**    | Go, Fiber, GORM                       |
| **Frontend**   | HTML5, CSS3, Templating               |
| **Database**   | SQLite (production-ready for Postgres)|
| **DevOps**     | Docker, CI/CD-ready                   |

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- SQLite3 (or Docker for containerized setup)

### Installation
```bash
# Clone the repository
git clone https://github.com/yourusername/expense-tracker.git
cd expense-tracker

# Install dependencies
go mod download

# Configure environment (copy and edit example)
cp .env.example .env

# Run migrations and seed data (if needed)
go run main.go migrate

# Start the server
go run main.go