CREATE TABLE cfs (
   id TEXT PRIMARY KEY,
   location TEXT,
   description TEXT,
   eventtime TIMESTAMP,
   raweventtime TEXT,
   eventdate TEXT
 );

CREATE TABLE location (
   location TEXT PRIMARY KEY,
   latitude REAL,
   longitude REAL,
   ward TEXT,
   neighborhood TEXT,
   zipcode TEXT,
   hasIssue BOOLEAN
 );