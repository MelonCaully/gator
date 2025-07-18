# Gator 
### To run workspace, install and build the following : 

- (Using Ubuntu) Install go: 'sudo apt update' & 'sudo apt install golang-go'

- Install postgreSQL: 'sudo apt install postgresql postgresql-contrib' (Check version with: 'psql --version')

- (Linux only) Change postgres password: 'sudo passwd postgres'

- Start the postgres server in the background: 'sudo service postgresql start' (Connect to the server with: 'sudo -u postgres psql')

- Create a new database: 'CREATE DATABASE gator;' 

- Connect to the new database: '\c gator' (Linux only: 'ALTER USER postgres PASSWORD 'password';')

### How to run project:

- Build the project in the root of the project: 'go build'

- Run the project: 'go run . (commands)'

### Commands

- login: sets the current user to a user registered in the database

- register: registers a given username into the database

- reset: resets the users table

- users: grabs all instances of a user in the users table

- agg: grabs all metadata/data from website

- addFeed: stores feed in database

- feeds: prints out all feeds in the database

- follow: follows feed with current user

- following: shows you all feeds followed by current user

- unfollow: unfollows given feed from the current user

- browse: allows the user to see all posts of a given feed they follow

### Migrations

- Up migration: 'goose postgres postgres://postgres:Momo34129@localhost:5432/gator up'

- Down migration: 'goose postgres postgres://postgres:Momo34129@localhost:5432/gator down'

### Built With 

- WSL and Ubuntu

- go version 1.24

- PostgreSQL
