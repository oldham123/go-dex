# PokeDex API

## Project Overview

This API provides endpoints to retrieve Pokemon information, implementing core features of a PokeDex including:
- Basic Pokemon data retrieval
- Type-based filtering
- Name search functionality
- Comprehensive Pokemon statistics

### Tech Stack
- **Language**: Go
- **Data Storage**: In-memory with CSV data source (Database integration planned)
- **API Style**: RESTful
- **Documentation**: OpenAPI/Swagger (planned)

## Getting Started

### Prerequisites
- Go 1.21 or higher
- Git

### Installation
```bash
# Clone the repository
git clone https://github.com/oldham123/go-dex.git

# Navigate to project directory
cd go-dex

# Install dependencies
go mod download

# Run the server
go run cmd/api/main.go
```

### Development Setup
1. Install Go from [golang.org](https://golang.org)
2. Install recommended VSCode extensions:
   - Go (by Go Team at Google)
   - Git Graph (optional)

## Project Structure
```
pokedex-api/
├── cmd/
│   └── api/            # Application entrypoint
├── internal/           # Private application code
│   ├── models/         # Data models
│   ├── handlers/       # HTTP handlers
│   └── storage/        # Data storage logic
├── data/              # CSV data files
└── docs/              # Documentation
```

## API Endpoints (Planned)

- `GET /pokemon` - List all Pokemon
- `GET /pokemon/{id}` - Get Pokemon by ID
- `GET /pokemon/search` - Search Pokemon by name
- `GET /pokemon/type/{type}` - Filter Pokemon by type

## Contributing

This is a learning project, but suggestions and discussions are welcome! Please feel free to:
1. Open an issue to discuss proposed changes
2. Fork the repository
3. Create a pull request

## Project Status

🚧 **Under Development** 🚧

This project is being actively developed as a learning exercise. Features will be added incrementally, and the API structure may change significantly as development progresses.

### Current Focus
- [ ] Basic project setup
- [ ] CSV data integration
- [ ] Initial API endpoints
- [ ] Data model design

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

**Note**: This is a learning project and is not intended for production use. The Pokemon data and Pokemon name are trademarks of Nintendo/Creatures Inc./GAME FREAK Inc.