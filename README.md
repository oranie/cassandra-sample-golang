## This project is training and demo app using Cassandra.


![codebuild](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiQUY2QVVVMzh0WE4xZ2hKdVlPbTVrc3hhTlhMc0tNbklpRVIxVm1nUUxONDVGZitKVEtJYjVhZ0lwUVhIRjRoUHh4am0yL213MVV1VWhEVEh1Q1V5dkNZPSIsIml2UGFyYW1ldGVyU3BlYyI6InkzYlJ1dk1uMXczSFhGdzAiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)



Realtime Comment Demo App

![demo](./demo.gif)


## High level architecture
![arch](./architecture.png)

## Envroiment

golang 1.13

Docker version 19.03.5, build 633a0ea

## Envroiment Setting
If you need to run app on CI/CD,test and prpduction envroiment.Please overwrite these value.

|EnvroimentName|Description|Default value|
|---|---|---|
|CASSANDRA_ENDPOINT |Cassandra cluster endpoint(ip or domain)|127.0.0.1|
|CASSANDRA_USER|Cassandra cluster user|cassandra|
|CASSANDRA_PASS|Cassandra cluster user pass|cassandra|
|CASSANDRA_PORT|Cassandra cluster network port|9042|
|CASSANDRA_KS|Cassandra cluster key spacename|example|
|APP_ENDPOINT|Chat App use Endpoint|http://127.0.0.1|
|APP_PORT |Chat App use network port|8080|
|APP_ENV|App run envroiment|test|

## Build
```shell script
GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build ./main.go
docker build -t $IMAGE_REPO_NAME .
docker tag $IMAGE_REPO_NAME:$IMAGE_TAG $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/$IMAGE_REPO_NAME:$IMAGE_TAG
docker push $AWS_ACCOUNT_ID.dkr.ecr.$AWS_DEFAULT_REGION.amazonaws.com/$IMAGE_REPO_NAME:$IMAGE_TAG
```

This operation is written in buildspec.yaml for codebuild.

## Deploy
This project build docker container.Any environment where containers and Cassandra can run is OK.
Example using ECR, ECS Fargate, Amazon MCS.

## Cassandra setting(Amazon MCS)

https://docs.aws.amazon.com/mcs/latest/devguide/accessing.html

https://docs.aws.amazon.com/mcs/latest/devguide/getting-started.html

Create Cassandra user
```
aws iam create-service-specific-credential --user-name USERNAME --service-name cassandra.amazonaws.com
```

cqlsh ssl setting
```shell script
curl https://www.amazontrust.com/repository/AmazonRootCA1.pem -O

cat $HOME/.cassandra/cqlshrc

[connection]
port = 9142
factory = cqlshlib.ssl.ssl_transport_factory

[ssl]
validate = true
certfile =  $HOME/AmazonRootCA1.pem
````

```
connect Amazon MCS
cqlsh {endpoint} {port} -u {ServiceUserName} -p {ServicePassword} --ssl

Create KeySpace(This app auto create table.)

CREATE KEYSPACE IF NOT EXISTS example WITH REPLICATION={'class': 'SingleRegionStrategy'};
```

## Create ECR repository
https://docs.aws.amazon.com/ja_jp/AmazonECR/latest/userguide/ECR_AWSCLI.html


### local app start
```shell script
go run ./main.go
```

memo:If you need to connecxt remote cassandra cluster,It might be a good idea to use an SSH tunnel.
```shell script
ssh ${SSH_TUNNEL_HOST} -L 9042:${CASSANDRA.ENDPOINT}:9042
```

# Production Deploy
## 1st step : create service IAM credentials

## 2nd step : 

## 3rd step : 

## 4th step : 


## golang Test


## Data Modeling
Data Modeling:

|name(PK)  |time(clustering column)  |comment  |chat_room |
|---|---|---|---|
|text  |text(micro sec unixtime)  |text  |text |

|chat_room(PK)  |time(clustering column)  |comment  |name |
|---|---|---|---|
|text  |text(micro sec unixtime)  |text  |text |

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