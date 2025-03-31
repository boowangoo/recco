# IMAGE="recco-server"
# CONTAINER="recco-server"

# # no image
# if [ -z $(docker images -q $IMAGE) ]; then
#     docker build -t $IMAGE .
# fi

# # no container
# if [ -z $(docker ps -q -f name=$CONTAINER) ]; then
#     docker run --name $CONTAINER -v $(pwd)/recco-server:/recco-server --rm -i -t -p 127.0.0.1:8080:8080 -d $IMAGE
# fi

# # connect to container
# docker attach $CONTAINER
docker-compose up