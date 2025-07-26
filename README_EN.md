# Online Exam System

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.0+-green.svg)](https://vuejs.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)

A comprehensive online examination system with multi-tenant support, AI assistant, and real-time monitoring. Built with modern technology stack for high performance and availability.

English | [中文](./README.md)

## ✨ Features

### 🎯 Core Features
- **Multi-Role Support**: Admin, Teacher, and Student roles with permission separation
- **Exam Management**: Complete exam workflow from creation to statistical analysis
- **Question Bank**: Support for multiple question types (single choice, multiple choice, true/false, essay)
- **Smart Paper Generation**: Manual and AI-powered automatic paper generation
- **Real-time Monitoring**: Real-time exam monitoring with anti-cheating mechanisms
- **Grade Statistics**: Detailed grade analysis and statistical reports

### 🚀 Technical Features
- **Multi-tenant Architecture**: Support for multiple institutions
- **Redis Caching**: High-performance caching for high concurrency
- **AI Integration**: Built-in AI Q&A assistant
- **Responsive Design**: Support for PC and mobile devices
- **Docker Deployment**: One-click deployment, ready to use
- **Performance Optimization**: Database index optimization, load testing verified

## 🛠️ Tech Stack

### Backend
- **Framework**: Go + Gin
- **Database**: PostgreSQL + Redis
- **ORM**: GORM
- **Authentication**: JWT
- **Cache**: Redis
- **Container**: Docker + Docker Compose

### Frontend
- **Framework**: Vue 3 + TypeScript
- **Build Tool**: Vite
- **UI Library**: Element Plus
- **State Management**: Pinia
- **Router**: Vue Router
- **Styling**: Tailwind CSS

## 📦 Quick Start

### Requirements

- Docker & Docker Compose
- Go 1.19+ (for development)
- Node.js 18+ (for development)

### One-Click Deployment (Recommended)

```bash
# Clone the project
git clone https://github.com/your-username/online-exam-system.git
cd online-exam-system

# Start all services
docker-compose up -d

# Access after services are ready
# Frontend: http://localhost:5173
# Backend API: http://localhost:8080
```

### Development Environment

#### 1. Start Database Services

```bash
# Start PostgreSQL and Redis
docker-compose up -d postgres redis
```

#### 2. Backend Service

```bash
cd backend

# Install dependencies
go mod download

# Configure environment variables
cp .env.example .env
# Edit .env file to configure database connection

# Run service
go run main.go
```

#### 3. Frontend Service

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

## 🎮 Default Accounts

The system provides the following test accounts after initialization (password: `admin123`):

| Role | Username | Password | Description |
|------|----------|----------|-------------|
| Admin | admin | admin123 | System Administrator |
| Teacher | teacher1 | admin123 | Teacher Account |
| Student | student1 | admin123 | Student Account |

**⚠️ Please change default passwords in production environment!**

## 📖 API Documentation

### Authentication
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/profile` - Get user profile
- `PUT /api/v1/auth/profile` - Update user profile

### Exam Management
- `GET /api/v1/exams` - Get exam list
- `POST /api/v1/exams` - Create exam
- `POST /api/v1/exams/:id/start` - Start exam
- `POST /api/v1/exams/:id/submit` - Submit exam

### Question Management
- `GET /api/v1/questions` - Get question list
- `POST /api/v1/questions` - Create question
- `PUT /api/v1/questions/:id` - Update question
- `DELETE /api/v1/questions/:id` - Delete question

For more API documentation, see [API Documentation](./docs/API.md)

## 🏗️ Project Structure

```
.
├── backend/                 # Backend service
│   ├── controllers/         # Controllers
│   ├── models/             # Data models
│   ├── middleware/         # Middleware
│   ├── services/           # Business services
│   ├── cache/              # Cache service
│   └── main.go             # Entry point
├── frontend/               # Frontend application
│   ├── src/
│   │   ├── components/     # Components
│   │   ├── views/          # Pages
│   │   ├── stores/         # State management
│   │   └── api/            # API interfaces
│   └── package.json
├── load-test/              # Load testing scripts
├── docker-compose.yml      # Docker orchestration
└── README.md
```

## 🧪 Testing

### Load Testing

The project includes comprehensive k6 load testing scripts supporting 50VU concurrency with 95% response time <100ms:

```bash
# Install k6
brew install k6  # macOS
# or sudo apt install k6  # Ubuntu

# Run load tests
cd load-test
./run-load-tests.sh
```

### Unit Testing

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm run test:unit
```

## 📊 Performance Optimization

- **Database Indexing**: Indexes added for high-frequency query fields
- **Redis Caching**: Cache hot data to reduce database pressure
- **Connection Pooling**: Database connection pool optimization
- **Compression**: Gzip compression to reduce transfer size
- **CDN**: Static resource CDN acceleration

## 🚀 Deployment Guide

### Production Deployment

1. **Environment Configuration**
   ```bash
   # Set production environment variables
   export GIN_MODE=release
   export DB_HOST=your-db-host
   export REDIS_HOST=your-redis-host
   ```

2. **Database Migration**
   ```bash
   # Execute database migration
   psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f backend/migrations/add_indexes.sql
   ```

3. **Start Services**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

### Monitoring and Logging

- Use Prometheus + Grafana for system performance monitoring
- Integrate ELK Stack for log analysis
- Configure alert rules for timely issue detection

## 🤝 Contributing

We welcome all forms of contributions! Please see [Contributing Guide](./CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## 🙏 Acknowledgments

Thanks to all developers who contributed to this project!

## 📞 Contact

- Project Homepage: [GitHub](https://github.com/your-username/online-exam-system)
- Issue Reports: [Issues](https://github.com/your-username/online-exam-system/issues)
- Discussions: [Discussions](https://github.com/your-username/online-exam-system/discussions)

## 🌟 Star History

If this project helps you, please give us a ⭐️!

[![Star History Chart](https://api.star-history.com/svg?repos=your-username/online-exam-system&type=Date)](https://star-history.com/#your-username/online-exam-system&Date)