docker-compose up --build -d

# # Start recco-dataloader to initialize the recco-db
# # TODO: Move recco-db data files to cloud storage
# if [[ "$@" == *"-d"* ]]; then
#     echo "Starting recco-dataloader container"
#     cd recco-dataloader/ && docker-compose up --build -d
#     cd ..
# fi
