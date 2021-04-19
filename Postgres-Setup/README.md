# Running the Docker Container

docker run -d --name cfs-postgres -e POSTGRES_PASSWORD=<redacted> --net cfs -p 5000:5432 --restart always -v /home/ubuntu/Projects/CFS/cfs-data:/var/lib/postgresql/data postgres

# Connecting to the Postgres Container

docker exec -it <container-name> psql -U <postgresuser>

# Postgres Commands:

- List Databases

\l (slash-L)

- Connect to Database

/c cfs


- List tables

\dt

 - Create Database 

 `CREATE DATABASE dbname;`

- Describe table

`SELECT 
   table_name, 
   column_name, 
   data_type 
FROM 
   information_schema.columns
WHERE 
   table_name = 'city';`


- Create table;

`CREATE TABLE cfs (
   id TEXT PRIMARY KEY,
   location TEXT,
   description TEXT,
   eventtime TIMESTAMP,
   raweventtime TEXT,
   eventdate TEXT
 );`

`CREATE TABLE location (
   location TEXT PRIMARY KEY,
   normalizedlocation TEXT,
   latitude REAL,
   longitude REAL,
   ward TEXT,
   neighborhood TEXT,
   zipcode TEXT,
   hasIssue BOOLEAN
 );`