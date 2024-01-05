# grep-go: Concurrent File Search in Go

A command-line utility written in Go for recursively searching files for a specific word or pattern concurrently.

## Features

- **Concurrent Search:** Utilizes goroutines to concurrently search for the specified word or pattern in multiple files, improving performance.
- **Recursive Search:** Supports recursive search in subdirectories, allowing you to search through the entire directory structure.
- **Case Insensitive:** Option to enable case-insensitive search, making the search more flexible.
- **Informative Output:** Displays lines containing the specified word or pattern.

## Usage

### Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/grep-go.git
cd grep-go
```

Build the executable:

```bash
go build -o grep-go main.go
```

### Usage

Run the executable with the following options:

```bash
./grep-go -path /path/to/search -word searchWord -r -i
```

Options:

- `-path`: Specify the file or directory path to search.
- `-word`: Specify the word or pattern to search for.
- `-r`: Enable recursive search (optional).
- `-i`: Enable case-insensitive search (optional).

Example:

```bash
./grep-go -path ./example -word "example" -r -i
```

## Dependencies

This project uses only the standard libraries provided by Go. No external dependencies need to be installed.

## Contributing

Feel free to contribute by opening issues or submitting pull requests. Contributions are welcome!