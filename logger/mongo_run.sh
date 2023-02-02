USER="admin"
PASSWORD="admin"
DATABASE="logs"

echo "Starting MongoDB server at port 27017"

docker run \
    -d --rm --name mongo-local \
    -p 27017:27017 \
    -e MONGO_INITDB_ROOT_USERNAME=$USER \
    -e MONGO_INITDB_ROOT_PASSWORD=$PASSWORD \
    -e MONGO_INITDB_DATABASE=$DATABASE \
    mongo:latest
