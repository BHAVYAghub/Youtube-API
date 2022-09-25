# YouTube API
A Golang application which shows Youtube videos data of any search query. Database will be filled with new data calling youtube API every minute using a cron job.

---
### Project tech Info:
* Language: Project is build in **Golang**.
* DataBase: 
  * **MongoDB** is used as the database for this project. 
  * Indexing: 
    * publishedAt: field is indexed to fetch sorted video data on the basis of the published Date.
    * Title & description: text index is used on video title and description to provide '**Fuzzy Search**' on these fields.
* Containerization: 
  * **Docker** has been used to containerize the application.
  * MongoDB and Application have been containerized in docker-compose file.

---
### Steps to run the application:
###### With Docker:
* Just run this command `docker compose up --build`
###### Without Docker:
* Setup Mongo DB in your local machine.
* Setup Golang in your local machine.
* Open .env file in project's root directory. And Provide all the configs mentioned.
* Run `go run main.go` command in project's root directory.

 Server will be up and running on localhost port 8000 [by default].

---




 



    