# Playability

Playability is a full stack web app that allows users to search for games, view their details, and submit feedback on the accessibility of games. This is a passion project to help make video games more accessible and inclusive for everyone.

## Table of Contents

- [Features](#features)
- [Technologies](#technologies)
- [Database Schema](#database-schema)
- [Installation](#installation)
  - [Prerequisites](#prerequisites)
  - [Backend Setup](#backend-setup)
  - [Frontend Setup](#frontend-setup)
  - [Database Setup](#database-setup)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Features

- **User Registration:** Secure user sign-up with validation.
- **User Driven Feedback:** Users can submit feedback on the accessibility of games.
- **Game Search:** Search for games by name.
- **Game Details:** View comprehensive details including cover art, platforms, and accessibility features.
- **Accessibility Information:** Display features like closed captions, colorblind modes, controller support, and remapping.
- **Responsive Design:** Mobile-friendly interface using Tailwind CSS.

## Technologies

- **Frontend:**
  - [Nuxt 3](https://nuxt.com/) - Vue.js framework for server-side rendering.
  - [Vue.js](https://vuejs.org/) - JavaScript framework for building user interfaces.
  - [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework.
  - [ShadCN Vue](https://shadcn.com/) - UI component library.
  - [Vee Validate](https://vee-validate.logaretm.com/v4/) - Form validation.
- **Backend:**
  - [Go](https://golang.org/) - Programming language for backend services.
  - [Chi Mux](https://github.com/go-chi/chi) - HTTP router.
  - [PostgreSQL](https://www.postgresql.org/) - Relational database.
  - [IGDB API](https://api-docs.igdb.com/) - Game data source.
  - [PCGamingWiki API](https://pcgamingwiki.com/api.php) - Accessibility information.

## Installation

### Prerequisites

- [Node.js](https://nodejs.org/) v14 or higher
- [pnpm](https://pnpm.io/) package manager
- [Go](https://golang.org/) v1.23.1
- [PostgreSQL](https://www.postgresql.org/) database

### Backend Setup

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/playability.git
   cd playability
   ```

2. **Navigate to the backend directory and install dependencies:**

   ```bash
   cd cmd/api
   go mod download
   ```

3. **Set up environment variables:**

   Create a `.env` file in the root directory with the following:

   ```env
   IGDB_CLIENT_SECRET=your_igdb_client_secret
   ```

4. **Run the backend server:**

   ```bash
   go run main.go
   ```

### Frontend Setup

1. **Navigate to the frontend directory:**

   ```bash
   cd ../frontend
   ```

2. **Install dependencies:**

   ```bash
   pnpm install
   ```

3. **Configure environment variables:**

   Create a `.env` file in the frontend directory if needed.

4. **Run the development server:**

   ```bash
   pnpm run dev
   ```

## Usage

1. **Access the application:**

   Open your browser and navigate to `http://localhost:3000`.

2. **Register a new user:**

   - Go to the registration page.
   - Fill in the required details and submit.

3. **Search for games:**

   - Use the search bar on the homepage to find games.
   - Click on a game from the search results to view details.

## API Endpoints

- **GET `/games?id={game_id}`**
  - Fetch game details by ID.
- **POST `/user/register`**

  - Register a new user.
  - **Body Parameters:**
    - `username` (string, required)
    - `email` (string, required)
    - `password` (string, required)

- **GET `/search?query={search_query}`**
  - Search for games based on query.

## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the repository.**
2. **Create a new branch:**

   ```bash
   git checkout -b feature/YourFeature
   ```

3. **Make your changes and commit them:**

   ```bash
   git commit -m "Add some feature"
   ```

4. **Push to the branch:**

   ```bash
   git push origin feature/YourFeature
   ```

5. **Open a pull request.**

## License

This project is licensed under the [MIT License](LICENSE).
