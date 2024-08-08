WORK_DIR=$(dirname "$0")
SERVICE_LIST="producer-worker consumer-worker"
ARGS="--t_endpoint=192.168.49.2:30880 --k_endpoint=192.168.49.2:32059 --k_topic=my-topic --count=200"

usage()
{
	echo " Usage: $0 startall"
	echo "        $0 stopall"
	echo "        $0 restartall"
	echo "        $0 start serviceName"
	echo "        $0 stop serviceName"
	echo "        $0 restart serviceName"
}

has_process()
{
    SELFID=$(id -u)
    PIDS=$(pgrep -u $SELFID -f  "./$1")
    if [ ! -z "$PIDS"  ]
    then
        echo service $1 is running, pid $PIDS
        return 1 
	fi
    return 0
}

start()
{
    has_process $1
    if [ $? -eq 0  ]
    then
        echo "start $1"
		LOG_FILE=$WORK_DIR/logs/$1.log
		mkdir -p $WORK_DIR/logs
		nohup setsid $WORK_DIR/bin/$1 $ARGS &> $LOG_FILE & 
        sleep 1
    fi
}

stop()
{
    SELFID=$(id -u)
    PIDS=$(pgrep -u $SELFID -f  "./$1")
    if [ -n "$PIDS"  ]
    then
        for PID in $PIDS
        do
            kill -2 $PID
            echo service $1 pid $PID is killed
        done
    else
        echo service $1 is not running
    fi
}

startAll()
{
	if [ ! -z $1 ]; then
        start $1
		if [ $? -eq 0 ]; then
			exit 0
		else
			exit 1
		fi
	fi
	start_check
	if [ $? -eq 1 ]; then
		exit 1
	fi
	for S in $SERVICE_LIST
	do
		start $S
	done
}

stopAll()
{
	if [ ! -z $1 ]; then
        stop $1
		if [ $? -eq 0 ]; then
			exit 0
		else
			exit 1
		fi
	fi
	for S in $SERVICE_LIST
	do
		stop $S
	done
	for (( wait_count=0; wait_count<10; wait_count++))
	do
		start_check
		if [ $? -eq 1 ]; then
			wait 5
		else
			break;
		fi
	done
}

restartAll()
{
	if [ ! -z $1 ]; then
        stop $1
        start $1
		if [ $? -eq 0 ]; then
			exit 0
		else
			exit 1
		fi
	fi
	stopAll
	startAll
}

check()
{
    RET=0
	for S in $SERVICE_LIST
	do
		has_process $S
		if [ $? -eq 0 ] 
		then
			echo  "[warning]$S not running"
			RET=1 
		fi
	done
    return $RET
}

start_check()
{
	for S in $SERVICE_LIST
	do
		has_process $S
		if [ $? -eq 1 ] 
		then
			return 1
		fi
	done
    return 0
}

wait()
{
	local secs=$1
	if [ -z $secs ]
	then 
		secs=1
	fi   
	for (( i=0; i<`expr $secs*5`; i++))
	do   
		echo -n "."
		sleep 0.2
	done 
	echo "."
}

case $1 in
start)
	shift 1
	startAll $*
	;;
stop)
	shift 1
	stopAll $*
	;;
restart)
	shift 1
	restartAll $*
	;;
startall)
	startAll
	;;
stopall)
	stopAll
	;;
restartall)
	restartAll
	;;
check)
	check
	;;
status)
	check
	;;
*)
	usage
	;;
esac
