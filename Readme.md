# Instagram database

## This repository contains realtime example of an instagram database, contains entites like user, post, business, stories, highlights etc.,

### How to run this repo

### Pre requesite
1. docker and docker compose installed in your system.
2. Run below command to start the database along with db migration.
   `docker compose up -d`
3. to connect to the database refer `docker-compose-local.env` file for credentails
4. database schema will be created under `SYS` user, host is `localhost`, port is `5432`, and a database is `postgres` and password presents in env file
5. DB migration will be executed once the db is up.

### How to run the the data-loader script
1. if you have Go installed in your machine run 
   `go run main.go` this command will work on both Windows/MacOS
         or 
   `go build` followed by `./data-loader` incase you are on linux/macOS
2. If you are on windows run `go build` followed by `./data-loader.exe`


## What the script does?
* Extracts data from [instagram_profiles_Github Hashtag_dataset.json](instagram_profiles_Github%20Hashtag_dataset.json) file and loads into 7 different table
* Create a relation between users by making following and followers.
* Create comments for posts based on the comment count mentioned on the data set.