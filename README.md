# terminal-reader

## About This Repository

`terminal-reader` is a command-line tool designed to bring Retrieval-Augmented Generation (RAG) capabilities directly to your terminal, allowing you to **chat with the contents of a Git repository**.

This project provides an interactive interface where you can ask questions about a specific codebase hosted on GitHub. It works by:
1.  **Cloning** the specified GitHub repository locally.
2.  **Processing** its files, breaking them down into smaller, manageable chunks.
3.  **Indexing** these chunks for efficient retrieval.
4.  Utilizing a Large Language Model (LLM) to answer your queries based on the retrieved code snippets from the repository.

Think of it as having a knowledgeable assistant who understands your codebase, right in your terminal!

## How to Get Started

If you've opened this repository and want to give `terminal-reader` a try, follow these steps:

### Prerequisites

Before you begin, ensure you have the following installed on your system:

*   **Go (1.22 or newer)**: The project is written in Go. You can download it from [golang.org](https://golang.org/dl/).
*   **Git**: Required for cloning repositories. Download it from [git-scm.com](https://git-scm.com/downloads).
*   **OpenAI API Key (or similar LLM provider)**: While not explicitly shown in the provided snippets, a RAG chatbot typically requires an external LLM. It's highly probable that this application uses an OpenAI-compatible API.
    *   You'll need an API key from a service like OpenAI.
    *   Set it as an environment variable, for example: `export OPENAI_API_KEY="sk-your_api_key_here"`

### Steps to Run

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/YOUR_GITHUB_USER/terminal-reader.git
    cd terminal-reader
    ```
    *(Replace `YOUR_GITHUB_USER` with the actual GitHub username/organization if this project is hosted publicly, otherwise use your local path or an appropriate placeholder).*

2.  **Build the application:**
    Navigate to the project root and build the `ragchat` executable.
    ```bash
    go build -o ragchat ./cmd/ragchat
    ```
    This will create an executable named `ragchat` in your current directory.

3.  **Set your API key (if required):**
    Ensure your LLM API key is set as an environment variable.
    ```bash
    export OPENAI_API_KEY="sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
    ```
    *(Replace with your actual API key)*

4.  **Run the RAG Chat:**
    Now you can run the application, providing the URL of the GitHub repository you want to chat with.

    ```bash
    ./ragchat --repo-url https://github.com/go-chi/chi
    ```
    *(Replace `https://github.com/go-chi/chi` with any public GitHub repository URL you wish to explore).*

    Once started, the application will clone the repository, process its files, and then present you with a prompt where you can type your questions about the codebase.

    **Example Interaction:**
    ```
    Cloning repository... Done.
    Processing files... Done.
    Welcome to RAG Chat! Type your questions about the repo. Type 'exit' or 'quit' to end.

    You: What is the purpose of the 'chi' router?
    AI: The 'chi' router is a lightweight, idiomatic, and composable HTTP router for Go. It's designed for building modular and maintainable web services...
    You: How do I define a new route?
    AI: You can define a new route using methods like `r.Get("/path", handlerFunc)`, `r.Post("/path", handlerFunc)`, etc., on your router instance 'r'...
    You: exit
    ```
    To exit the chat session, simply type `exit` or `quit`.