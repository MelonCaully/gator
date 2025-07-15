# Gator 

### To run workspace, install and build the following : 

- Install postgreSQl: sudo apt install postgresql postgresql-contrib (Check version with: 'psql --version')

- (Linux only) Change postgres password: sudo passwd postgres

- Start the postgres server in the background: sudo service postgresql start (Connect to the server with: 'sudo -u postgres psql')

- Create a new database: CREATE DATABASE database; (Linux only: ALTER USER postgres PASSWORD 'postgres';)

### Migrations

- Up migration: goose postgres postgres://postgres:Momo34129@localhost:5432/gator up

- Down migration: goose postgres postgres://postgres:Momo34129@localhost:5432/gator down
