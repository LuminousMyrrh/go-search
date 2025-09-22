POD_NAME=pg-pod
CONTAINER_NAME=searchdb
POSTGRES_PASSWORD=testpassword
POSTGRES_DB=searchdb
DATA_VOLUME=pgdata

# Create a persistent volume if it doesn't exist
if ! podman volume exists $DATA_VOLUME; then
  podman volume create $DATA_VOLUME
  echo "Created Podman volume: $DATA_VOLUME"
fi

# Create pod if it doesn't exist
if ! podman pod exists $POD_NAME; then
  podman pod create --name $POD_NAME -p 5432:5432
  echo "Created pod: $POD_NAME"
fi

# Run or restart PostgreSQL container in the pod
if podman container exists $CONTAINER_NAME; then
  podman rm -f $CONTAINER_NAME
  echo "Removed existing container: $CONTAINER_NAME"
fi

podman run -d \
  --name $CONTAINER_NAME \
  --pod $POD_NAME \
  -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
  -e POSTGRES_DB=$POSTGRES_DB \
  -v ${DATA_VOLUME}:/var/lib/postgresql/data:Z \
  postgres:latest

echo "PostgreSQL container started in pod $POD_NAME"

# Optionally wait for DB initialization here or run migration scripts next
