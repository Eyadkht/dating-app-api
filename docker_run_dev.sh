# Start the database service
docker-compose up -d db

# Wait for some time to make sure the database is up and running before starting the api service
sleep 2

docker-compose up --build api  