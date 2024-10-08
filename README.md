# Personal Blog in Go

A simple personal blog built with Go, featuring a guest section for viewing articles and an admin section for managing articles.

## Features

- **Guest Section**
  - Home Page: Displays a list of all published articles.
  - Article Page: Shows the content of a selected article.

- **Admin Section**
  - Dashboard: Manage articles (add, edit, delete).
  - Basic Authentication: Secure admin access.

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/personal-blog.git
   cd personal-blog
   ```

2. **Build the project**:
   ```bash
   go build -o blog cmd/main.go
   ```

3. **Run the server**:
   ```bash
   ./blog
   ```

4. **Access the blog**:
    - Guest section: `http://localhost:8080`
    - Admin section: `http://localhost:8080/admin` (Username: `admin`, Password: `password`)

## Project Structure

- `cmd/main.go`: The main entry point for the application.
- `templates/`: HTML templates for rendering pages.
- `data/`: Directory to store articles as JSON files.

### roadmap.sh projects
- https://roadmap.sh/projects/personal-blog