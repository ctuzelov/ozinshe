Project Title
A brief tagline describing your project

Overview

    A few sentences explaining the purpose of your application (movie/series database, user management, etc.)
    Key features (filtering, favorites, etc.)

Technologies Used

    Backend: Go, Gin, [database technology]
    Frontend: HTML, CSS, JavaScript (and any frameworks or libraries used)

Getting Started

    Prerequisites:
        Go (version X or later)
        [Database software]
        Node.js (if you have a frontend)
    Clone the repository:
    Bash

    git clone https://github.com/your-username/your-project-name.git

    Используйте код с осторожностью.

Set up the database:

    Create a database.
    Configure database connection settings in [your config file].

Install dependencies:
Bash

go get -u ./...

Используйте код с осторожностью.
Run the application:
Bash

go run main.go

Используйте код с осторожностью.

API Endpoints

    Authentication
        /signup
        /signin
        /signout
    User Management
        /users
        /user/:id
        /change-password
        /change-profile
        /user/:id (Admin only)
    Projects
        ... (Include other project-related endpoints)

Creating a Project

    Endpoint: POST /create-project
    Authentication: Requires Admin authorization
    Request Type: Multipart Form Data
    Required Fields:
        cover: Project cover image file.
        screenshots: Array of project screenshot image files.
        project_type: String ("movie" or "series")
        movie_data (If project_type is "movie"): JSON string containing movie data:
        JSON

         {
             "title": "Movie Title",
             "releaseYear": 2023,
             "description": "Movie description...",
             "popularity": 70,
             "youtubeId": "...",
             "duration": 120,
             "director": "Director Name",
             "producer": "Production Company",
             "genres": ["Action", "Adventure"]
         }

        Используйте код с осторожностью.

        series_data (If project_type is "series"): JSON string containing series data (similar format to movie data)

Example Request (using curl)
Bash

curl -X POST http://localhost:8080/create-project \
 -H "Authorization: Bearer your_admin_token" \
 -F "cover=@path/to/cover.jpg" \
 -F "screenshots[]=@path/to/screenshot1.jpg" \
 -F "screenshots[]=@path/to/screenshot2.jpg" \
 -F "project_type=movie" \
 -F 'movie_data={"title": "The Matrix", ...}'

Используйте код с осторожностью.

Admin Account

    Email: chingizkhan.tuzelov@gmail.com
    Default Password: [Set secure Default password] (Change immediately upon login)
