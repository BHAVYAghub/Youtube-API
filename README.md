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
### Approach used:
* 

---
### Steps to run the application:
* Set configs in .env file.
###### With Docker:
* Just run this command `docker compose up --build`
###### Without Docker:
* Setup Mongo DB in your local machine.
* Setup Golang in your local machine.
* Open .env file in project's root directory. And Provide all the configs mentioned.
* Run `go run main.go` command in project's root directory.

Note: 
* Server will be up and running on localhost port 8000 [by default].
* while using docker approach, make sure `localhost` is set to `db` in `MONGO_URI` config present in .env file due to container network. 

---

### Configs Info:

The following configs are to be set in .env file.<br/>
See .env file for example values.

| Name                    | Description                                                                                   |
|-------------------------|-----------------------------------------------------------------------------------------------|
| MONGODB_URI             | Describes the mongoDB connection URI.                                                         |  
| MONGODB_COLLECTION_NAME | Describes the mongoDB connection name.                                                        |
| MONGODB_DATABASE_NAME   | Describes the mongoDB database name.                                                          |
| YT_API_KEY              | Describes the API-key to be used for calling YT endpoints.                                    |
| YT_QUERY_STRING         | Describes the query string to be used to fetch YT records and save in DB.                     |
| YT_API_BASE_URL         | Describes the YT Api Base URL.                                                                |
| YT_FETCH_RECORDS_AFTER  | Describes the time after which the YT records are to be fetched and saved.                    |
| YT_API_FETCH_INTERVAL   | Describes the interval to call the scheduler to fetch youtube records in minutes. (default 1) |
| PORT                    | Describes the port on which application will run.                                             |
| DISABLE_CRON            | kill switch to disable the youTube data fetching scheduler. (default false)                   |


---

### Rest Endpoints details:

* The GET endpoints mentioned below supports pagination via query params.
  * `page` default as 1.
  * `limit` default as 5.


| METHOD | ENDPOINT         | Description                                                                                                                                                       |
|--------|------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| GET    | /youtube/findAll | Fetches the youtube records saved in DB, sorted in reverse chronological order of their publishing date-time.                                                     |  
| GET    | /youtube/find    | Fetches the youtube records saved in DB, sorted in reverse chronological order of their publishing date-time on the basis of `search` query param (fuzzy serach). |

---
Note: Carefully set the `YT_FETCH_RECORDS_AFTER` config as this might exhaust your quota.


 



    