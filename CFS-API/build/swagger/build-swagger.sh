swagger generate spec -o ../build/swagger/swagger.json

docker run  -d -p 6001:8080 -e BASE_URL=/docs -e SWAGGER_JSON=/swagger/swagger.json -v /home/ubuntu/Projects/CFS/CallsForService/CFS-API/build/swagger:/swagger --name cfs-swagger swaggerapi/swagger-ui