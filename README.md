# Gator - RSS Feed Aggregator CLI

Gator (gator) is a terminal-based RSS feed aggregator built with Go and PostgreSQL. It allows users to register, follow RSS feeds, scrape new posts, and browse the latest articles—all from the command line.

## 🛠️ Requirements

Before using Gator, ensure the following are installed on your system:

- [Go](https://golang.org/dl/) (version 1.20 or newer recommended)
- [PostgreSQL](https://www.postgresql.org/download/)

## 📦 Installation

You can install Gator globally using the `go install` command:

```bash
go install github.com/JadedPigeon/Gator@latest
```

This compiles the binary and makes it available in your `$GOPATH/bin` or `$HOME/go/bin` directory. Make sure that path is in your system's `PATH`.

## 🚀 Running the App

Before running Gator, you’ll need to:

1. **Set up your PostgreSQL database**
2. **Apply the schema using Goose**
3. **Configure your .gatorconfig.json**

### 1. Create the Config File

You can run `gator register` and `gator login` to generate and update the config automatically. The file typically lives at:

```bash
~/.gatorconfig.json
```

Or you can create it manually:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator_db?sslmode=disable",
  "current_user": ""
}
```

### 2. Apply Migrations

If you’re using Goose, run:

```bash
goose postgres "<your-db-url>" up
```

This will apply all schema changes and prepare the database.

### 3. Run Gator

Once installed, you can run the CLI by typing:

```bash
gator <command> [args]
```

Or if you’re still developing:

```bash
go run . <command> [args]
```

## 🧪 Example Commands

### Register and Login

```bash
gator register alice
gator login alice
```

### Add a Feed

```bash
gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
```

### View Followed Feeds

```bash
gator following
```

### Start Aggregating

```bash
gator agg 60
```

This starts fetching posts from followed feeds every 60 seconds.

### Browse Recent Posts

```bash
gator browse 5
```

This shows the 5 most recent posts from feeds you're following. If no number is given, the default is 2.

## 📚 Development

To run during development:

```bash
go run . <command>
```

To build a production binary:

```bash
go build -o gator
./gator <command>
```

## 🔒 License

This project is for educational purposes and distributed as-is with no warranty.

