#!/bin/bash

set -euo pipefail # 严格模式：遇到错误退出，未定义变量报错，管道中任何命令失败都会导致整个管道失败

# 配置变量
readonly REMOTE_HOST="101.132.113.17"
readonly REMOTE_USER="root"
readonly REMOTE_PATH="/data/project/stag01"
readonly SERVICE_NAME="stag01"
readonly BINARY_NAME="mowserver"

# 颜色定义
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# 日志函数
log_info() {
  echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
  echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
  echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
  echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
  if ! command -v "$1" &>/dev/null; then
    log_error "Command '$1' not found. Please install it first."
    exit 1
  fi
}

# 检查远程连接
check_remote_connection() {
  log_info "Checking remote connection..."
  if ! ssh -o ConnectTimeout=10 -o BatchMode=yes "${REMOTE_USER}@${REMOTE_HOST}" exit &>/dev/null; then
    log_error "Cannot connect to ${REMOTE_HOST}. Please check your SSH configuration."
    exit 1
  fi
  log_success "Remote connection OK"
}

# SVN 更新
update_source() {
  log_info "Updating source code from SVN..."
  if svn up; then
    log_success "SVN update completed"
  else
    log_error "SVN update failed"
    exit 1
  fi
}

# 构建项目
build_project() {
  log_info "Building Go project..."
  if go build -o "${BINARY_NAME}" -ldflags "-s -w"; then
    log_success "Build completed successfully"
  else
    log_error "Build failed"
    exit 1
  fi

  # 检查二进制文件是否存在
  if [[ ! -f "${BINARY_NAME}" ]]; then
    log_error "Binary file ${BINARY_NAME} not found after build"
    exit 1
  fi
}

# 停止远程服务
stop_remote_service() {
  log_info "Stopping remote service ${SERVICE_NAME}..."
  if ssh "${REMOTE_USER}@${REMOTE_HOST}" "supervisorctl stop ${SERVICE_NAME}"; then
    log_success "Service stopped successfully"
  else
    log_warning "Failed to stop service (it might not be running)"
  fi
}

# 部署文件
deploy_files() {
  log_info "Deploying configuration files..."

  # 检查本地文件是否存在
  local files_to_check=("configs/gameconfig" "configs/cert" "${BINARY_NAME}")
  for file in "${files_to_check[@]}"; do
    if [[ ! -e "$file" ]]; then
      log_error "File/directory $file not found"
      exit 1
    fi
  done

  # 部署配置文件
  if scp -r configs/gameconfig "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/configs/" &&
    scp -r configs/cert "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/configs/" &&
    scp "${BINARY_NAME}" "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_PATH}/"; then
    log_success "Files deployed successfully"
  else
    log_error "File deployment failed"
    exit 1
  fi
}

# 启动远程服务
start_remote_service() {
  log_info "Starting remote service ${SERVICE_NAME}..."
  if ssh "${REMOTE_USER}@${REMOTE_HOST}" "supervisorctl restart ${SERVICE_NAME}"; then
    log_success "Service started successfully"
  else
    log_error "Failed to start service"
    exit 1
  fi
}

# 主函数
main() {
  log_info "Starting deployment process..."

  # 检查必要的命令
  check_command "svn"
  check_command "go"
  check_command "ssh"
  check_command "scp"

  # 检查远程连接
  check_remote_connection

  # 执行部署步骤
  update_source
  build_project
  stop_remote_service
  deploy_files
  start_remote_service

  log_success "Deployment completed successfully! 🎉"
}

# 信号处理 - 在脚本被中断时清理
cleanup() {
  log_warning "Deployment interrupted. You may need to manually restart the service."
  exit 1
}

trap cleanup SIGINT SIGTERM

# 运行主函数
main "$@"
