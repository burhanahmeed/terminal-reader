# Terminal Reader 🚀

A powerful command-line RAG (Retrieval-Augmented Generation) tool that lets you chat with any codebase directly from your terminal. Ask questions about functions, understand code structure, and get intelligent answers about your repositories.

## 🚀 Quick Start

### Prerequisites

- **Go 1.22+** - [Download here](https://golang.org/dl/)
- **Git** - [Download here](https://git-scm.com/downloads)
- **Gemini API Key** - [Get one here](https://makersuite.google.com/app/apikey)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/burhanahmeed/terminal-reader.git
   cd terminal-reader
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up your API key:**
   ```bash
   export GEMINI_API_KEY="your_gemini_api_key_here"
   ```
   
   Or create a `.env` file:
   ```bash
   echo "GEMINI_API_KEY=your_gemini_api_key_here" > .env
   ```

4. **Build and run:**
   ```bash
   go run cmd/ragchat/main.go --path /path/to/your/repo
   ```

## 📖 Usage

### Local Repository
```bash
go run cmd/ragchat/main.go --path ~/projects/my-awesome-project
```

### GitHub Repository
```bash
go run cmd/ragchat/main.go --github https://github.com/owner/repo
```

### Example Session
```
📦 Indexing repo: /path/to/repo
✅ Repo indexed. Starting chat (type 'exit' to quit).

> What does the main function do?
Thinking...
The main function initializes the application by setting up the HTTP server, 
configuring routes, and starting the listener on the specified port...

> How does authentication work?
Thinking...
The authentication is handled by the AuthMiddleware function which validates 
JWT tokens and extracts user information...

> exit
```

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Repository    │───▶│   File Loader    │───▶│    Chunker      │
│   (Git/Local)   │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                         │
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Vector Store  │◀───│   Embedder       │◀───│   Chunks        │
│   (SQLite)      │    │   (Gemini)       │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌──────────────────┐
│   Retriever     │    │   LLM Client     │
│                 │    │   (Gemini)       │
└─────────────────┘    └──────────────────┘
         │                       │
         └───────────┬───────────┘
                     ▼
            ┌─────────────────┐
            │   Chat Session   │
            │                 │
            └─────────────────┘
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `GEMINI_API_KEY` | Your Gemini API key | Yes |

### Command Line Options

| Flag | Description | Example |
|------|-------------|---------|
| `--path` | Path to local repository | `--path ~/projects/repo` |
| `--github` | GitHub repository URL | `--github https://github.com/user/repo` |

## 🛠️ Development

### Project Structure
```
terminal-reader/
├── cmd/ragchat/          # Main application
├── internal/
│   ├── embed/           # Embedding functionality
│   ├── llm/             # LLM client
│   ├── repo/            # Repository handling
│   ├── retriever/       # Vector storage & retrieval
│   └── session/         # Chat session management
├── pkg/cache/           # Caching utilities
└── data/                # Local data storage
```

### Building
```bash
# Build the binary
go build -o ragchat cmd/ragchat/main.go

# Run tests
go test ./...

# Run with race detection
go run -race cmd/ragchat/main.go --path .
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Google Gemini](https://ai.google.dev/) for the embedding and LLM capabilities
- [SQLite](https://sqlite.org/) for vector storage
- The Go community for excellent tooling and libraries

---

**Made with ❤️ for developers who want to understand their code better**