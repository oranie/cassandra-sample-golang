## This project is Amazon DynamoDB training and demo app.


![codebuild](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQUY2QVVVMzh0WE4xZ2hKdVlPbTVrc3hhTlhMc0tNbklpRVIxVm1nUUxONDVGZitKVEtJYjVhZ0lwUVhIRjRoUHh4am0yL213MVV1VWhEVEh1Q1V5dkNZPSIsIml2UGFyYW1ldGVyU3BlYyI6InkzYlJ1dk1uMXczSFhGdzAiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)



Realtime Comment Demo App

![demo](./demo.gif)


## High level architecture
![architecture](./demo_arch.png)

## Envroiment
golang 1.13

If you set API_ENDPOINT environment value 'localhost' is using DynamoDB Local 
 (ex: export API_ENDPOINT=localhost)
 
# Cassandra setting


## Deploy

### local app start

# Production Deploy
##1st step : 1st chalice deploy

## 2nd step : Check your generate APIGateway URL


## 3rd step : set to your api endpoint URL

## 4th step : 2nd chalice deploy


## Python Test


## Data Modeling
Data Modeling:

|name(PK)  |time(SK)  |comment  |chat_room |
|---|---|---|---|
|string  |string(micro sec unixtime)  |string  |string |string|

GSI

|chat_room(PK)  |time(SK)  |comment  |name |
|---|---|---|---|
|string  |string(micro sec unixtime)  |string  |string |



## API

* /chat

return chat client HTML and js.
    
* /chat/comments/add

client sent post request with name,comment txt, get response add comment status

POST value {"name": "oranie", "comment":"hello world"}


* /chat/comments/all

client sent get request, get all comment.
    
* /chat/comments/latest

client sent get request latest 20 comments.

* /chat/comments/latest/{latest_seq_id}

client sent get request with latest chat id, get the difference comments.
    

# License
This library is licensed under the MIT-0 License. See the LICENSE file.