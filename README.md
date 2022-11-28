# coog-music

## Steps to Run Project (Windows 10)
1. [Install PostgreSQL (14.6)](https://www.google.com)

To use the dumpfile, you can use DBeaver or PSQL.

### PSQL Method

1. After installing, [set up environment paths to use the psql command.](https://www.computerhope.com/issues/ch000549.htm)

a) Add this as a System Variable Path.

`;C:\Program Files\PostgreSQL\14\bin;C:\Program Files\PostgreSQL\14\lib`

2. Open a terminal and run

`psql -U [username]` Then enter password

`CREATE DATABASE [databaseToAddToName];`

Go back to the terminal (Ctrl + C or type `exit`)

`psql -U [username] [databaseToAddToName] < [pathOfDumpFile]`

### DBeaver Method
 [Install DBeaver](https://dbeaver.io/download/)
 
`To acquire all of the SQL from the dump file: 
Using the database manager DBeaver, you first have to create a new database connection with PostgreSQL. Once you have the connection with PostgreSQL, you must create a new database by right clicking on the Databases tab and clicking Create New Database. Now you right-click on the database you just created and click Tools > Execute script. Select the SQL dump file, then click Start.`

---
1. Now the database should be ready to connect to.

2. [Install Go (1.19.3)](https://go.dev/doc/install)

3. Clone this repository

4. Open the repository and update the file `app.env`'s `DB_SOURCE`

`DB_SOURCE=postresql://[username]:[password]@[hostName]:[5432]/[databaseToAddToName]?sslmode=disable`

5. Finally, in the terminal in this project's root directory, run `go run ./cmd/web/.`

6. Going to `localhost:8080` should bring you to the homepage.


## About the files

1. The front end is located in the `static` and `templates` folders.

2. The `db/migration/000001_init_schema.up.sql` contains just the schema of the database, which isn't needed with the dump file.

3. All the queries are located in `internal/repository/dbrepo/postgres.go`

4. All the requests are handled by `internal/handlers/handlers.go`
