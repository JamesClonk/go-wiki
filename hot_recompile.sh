#!/bin/bash

GO_PROJECT_MAIN="wiki"
GO_PROJECT_MAIN_EXT="${GO_PROJECT_MAIN}.go"
GO_PROJECT_PID=
CURR_PATH=`pwd`

function go_run {
	kill $GO_PROJECT_PID >/dev/null 2>&1
	for PID in `ps -ef | grep '/tmp/go' | grep -v 'grep' | grep "$GO_PROJECT_MAIN" | awk '{print $2}'`; do
		kill $PID >/dev/null 2>&1
	done

	go run $GO_PROJECT_MAIN_EXT &
	GO_PROJECT_PID=$!
}

go_run

inotifywait -mr --timefmt '%d/%m/%Y %H:%M' --format '%T %w %f' -e close_write $CURR_PATH --excludei ".*\.(txt|sh)" | while read date time dir file; do
    echo "At ${time} on ${date}, ${dir}${file} changed."
    echo "Restarting Go project: [$GO_PROJECT_MAIN_EXT]..."
    go_run
done

