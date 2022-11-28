# coog-music

## Steps to Run Project

1. [Install PostgreSQL (14.6)](https://www.google.com)
2. After installing, [set up environment paths to use the psql command.](https://www.computerhope.com/issues/ch000549.htm)

a) Add this as a System Variable Path.

`;C:\Program Files\PostgreSQL\14\bin;C:\Program Files\PostgreSQL\14\lib`

3. Open a terminal and run

`psql -U [username]` Then enter password

`CREATE DATABASE [databaseToAddToName];`

Go back to the terminal (Ctrl + C or type `exit`)

`psql -U [username] [databaseToAddToName] < [pathOfDumpFile]`

4. Now the database should be ready to connect to.

5. [Install Go (1.19.3)](https://go.dev/doc/install)

6. Clone this repository

7. Open the repository and update the file `app.env`'s `DB_SOURCE`

`DB_SOURCE=postresql://[username]:[password]@localhost:[5432]/[databaseToAddToName]?sslmode=disable`

8. Finally, in the terminal in this project's root directory, run `go run ./cmd/web/.`

9. Going to `localhost:8080` should bring you to the homepage.
