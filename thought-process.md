### solution approach   
-  to create a server that can handle thousands of requests we need to 
depend on threading and making the server components independent on each other
as possible so it can scale and easily integrated to a load balancer
- I tried to separate application layers as possible so I decided to follow clean architecture. However,
because I didn't use databases. I depended on service, infrastructure and controller layers.
- to create a storage that traces the unique Ids and won't be affected by the duplicated servers. I used redis cache db set, with expiry of 1 minute for created uniqueIds set.
this insures the uniquness and auto time settings.
- for logging, first I logged the uniqueId and the endpoint response status code to the logs file.
- I added a stream service to send the count of uniqueIds to, thus. I used Kafka message broker and created a kafka producer to send the count parameter along with logging it to the logs file
- at the end, I wrapped everything inside a docker-compose file that contain 4 images: kafka and zookeeper for message broker, redis db and the main go app.
- I would add a consumer service also but I found that it was not required so skipped it.
- application layers are as follows:
  - Infrastructure: definition of Kafka and redis clients
  - services: structs that does business logic for producer and memory logging
  - controllers: contain http endpoints that interact with the client and connects to the services 
  