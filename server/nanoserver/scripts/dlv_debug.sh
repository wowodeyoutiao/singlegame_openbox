#!/bin/bash

# 使用方式: ./dlv_debug.sh <binary_path> [start|stop|status|restart]

BINARY_PATH="$1"
ACTION="${2:-start}"

# 检查是否传入了二进制路径
if [[ -z "$BINARY_PATH" ]]; then
  echo "Usage: $0 <binary-path> [start|stop|status|restart]"
  exit 1
fi

# 定义日志和 pid 文件路径
DLV_LOG="./dlv.log"
APP_LOG="./myapp.log"
PID_FILE="./dlv.pid"
PORT="12345"

function start_debug() {
  if [[ -f "$PID_FILE" ]]; then
    echo "PID file exists. Is another instance running?"
    exit 1
  fi

  echo "Starting dlv debug server for $BINARY_PATH..."

  nohup dlv exec "$BINARY_PATH" \
    --headless --listen=0.0.0.0:$PORT \
    --api-version=2 \
    --accept-multiclient \
    --continue \
    --log --log-dest="$DLV_LOG" \
    -- "$@" >"$APP_LOG" 2>&1 &

  echo $! >"$PID_FILE"

  echo "Started with PID $(cat $PID_FILE)"
  echo "Logs: $APP_LOG and $DLV_LOG"
}

function stop_debug() {
  if [[ ! -f "$PID_FILE" ]]; then
    echo "PID file not found. Is the process running?"
    return 1
  fi

  PID=$(cat "$PID_FILE")

  if ps -p $PID >/dev/null; then
    echo "Stopping dlv debug server (PID: $PID)..."
    kill $PID && rm -f "$PID_FILE"
    echo "Stopped."
  else
    echo "Process with PID $PID not found. Removing stale PID file."
    rm -f "$PID_FILE"
  fi
}

function status_debug() {
  if [[ -f "$PID_FILE" ]]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID >/dev/null; then
      echo "Debug server is running with PID: $PID"
    else
      echo "Stale PID file found. Process not running."
    fi
  else
    echo "Debug server is not running."
  fi
}

# 主逻辑
case "$ACTION" in
start)
  shift 2
  start_debug "$@"
  ;;
stop)
  stop_debug
  ;;
restart)
  stop_debug
  sleep 1
  start_debug
  ;;
status)
  status_debug
  ;;
*)
  echo "Unknown action: $ACTION"
  echo "Usage: $0 <binary-path> [start|stop|status|restart]"
  exit 1
  ;;
esac
