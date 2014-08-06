#!/bin/sh
#
### BEGIN INIT INFO
# Provides:          ws
# Required-Start:    $network $local_fs $remote_fs $syslog
# Required-Stop:     $network $local_fs $remote_fs $syslog
# Should-Start:      $named
# Should-Stop:
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: A simple protyping web server
# Description:       ws is a simple protyping webserver. It is best suited 
#                    for low traffic, low risk content. 
### END INIT INFO
#

WEBROOT=/Sites/minime.local

PATH=/sbin:/usr/sbin:/bin:/usr/bin:/usr/local/bin:/usr/local/go/bin

do_start () {
	cd $WEBROOT
        if [ -f "$WEBROOT/etc/config.sh" ]; then
            . $WEBROOT/etc/config.sh
        else
            echo "Cannot find etc/config.sh in $(pwd)"
            exit 1
        fi
        mkdir -p $WEBROOT/run
      
        ws &
        WS_PID=$!
        echo $WS_PID > $WEBROOT/run/ws.pid
        echo "ws running with pid "$(cat $WEBROOT/run/ws.pid)
}

do_status () {
	cd $WEBROOT
	if [ -f "$WEBROOT/run/ws.pid" ] ; then
		return 0
	else
		return 4
	fi
}

do_stop () {
	cd $WEBROOT
        if [ -f "$WEBROOT/run/ws.pid" ]; then
        	PID=$(cat $WEBROOT/run/ws.pid) 
                if [ "$PID" != "" ]; then 
                        echo "Shutting down ws with pid: $PID"
                	kill $PID
			rm $WEBROOT/run/ws.pid
                else 
			echo "Cannot find pid of running ws"
                fi
        fi
}

do_reload () {
    do_stop
    do_start
}

case "$1" in
  start|"")
	do_start
	;;
  restart|reload|force-reload)
    do_reload
	#echo "Error: argument '$1' not supported" >&2
	#exit 3
	;;
  stop)
	do_stop
	;;
  status)
	do_status
	exit $?
	;;
  *)
	echo "Usage: motd [start|stop|status]" >&2
	exit 3
	;;
esac

exit 0
