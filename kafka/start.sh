WORKDIR=$(dirname "$0")
SERVICE_LIST="producer-starter consumer-starter"
ARGS="--t_endpoint=localhost:7233"	

for S in $SERVICE_LIST
do
    $WORKDIR/bin/$S $ARGS
done
