WORKDIR=$(dirname "$0")
SERVICE_LIST="producer-starter consumer-starter"
ARGS="--t_endpoint=192.168.49.2:30880"

for S in $SERVICE_LIST
do
    $WORKDIR/bin/$S $ARGS
done
