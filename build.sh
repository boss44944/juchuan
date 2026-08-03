#!/usr/bin/env bash
# ==============================================
#  菊传 (Juchuan) 构建助手
#  交互式菜单：编译前端 / 同步静态资源 / 编译后端 / 打包 macOS 应用
# ==============================================
set -uo pipefail

# ---------- 路径配置（脚本放在项目根目录即可） ----------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
FRONTEND_DIR="${SCRIPT_DIR}/frontend"
BACKEND_DIR="${SCRIPT_DIR}/backend"
DIST_DIR="${FRONTEND_DIR}/dist"
STATIC_DIR="${BACKEND_DIR}/static"   # 后端 go:embed static/* 需要的位置

# ---------- 颜色输出 ----------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${CYAN}[信息]${NC} $*"; }
ok()    { echo -e "${GREEN}[完成]${NC} $*"; }
warn()  { echo -e "${YELLOW}[提示]${NC} $*"; }
err()   { echo -e "${RED}[错误]${NC} $*"; }

# ---------- 前置检查 ----------
check_paths() {
  if [ ! -d "${FRONTEND_DIR}" ]; then
    err "前端目录不存在: ${FRONTEND_DIR}"; exit 1
  fi
  if [ ! -d "${BACKEND_DIR}" ]; then
    err "后端目录不存在: ${BACKEND_DIR}"; exit 1
  fi
  if [ ! -d "${STATIC_DIR}" ]; then
    warn "后端 static 目录不存在，将自动创建: ${STATIC_DIR}"
    mkdir -p "${STATIC_DIR}"
  fi
}

# ---------- 功能函数 ----------
install_frontend() {
  info "安装前端依赖 (npm install) ..."
  ( cd "${FRONTEND_DIR}" && npm install --no-fund --no-audit ) || { err "依赖安装失败"; return 1; }
  ok "前端依赖安装完成"
}

build_frontend() {
  if [ ! -d "${FRONTEND_DIR}/node_modules" ]; then
    warn "未检测到 node_modules，请先选择「安装前端依赖」。"
    return 1
  fi
  info "编译前端 (npm run build) ..."
  ( cd "${FRONTEND_DIR}" && npm run build ) || { err "前端编译失败"; return 1; }
  ok "前端编译完成 → ${DIST_DIR}"
}

sync_to_static() {
  info "同步产物到 backend/static ..."
  rsync -a --delete "${DIST_DIR}/" "${STATIC_DIR}/" || { err "同步失败"; return 1; }
  ok "静态资源已同步 → ${STATIC_DIR}"
}

build_backend() {
  info "编译后端 (go build ./...) ..."
  ( cd "${BACKEND_DIR}" && go build ./... ) || { err "后端编译失败"; return 1; }
  ok "后端编译完成"
}

package_macos() {
  info "打包 macOS 应用 (build-mac-app.sh) ..."
  ( cd "${BACKEND_DIR}" && bash build-mac-app.sh ) || { err "macOS 打包失败"; return 1; }
  ok "macOS 应用打包完成"
}

dev_frontend() {
  warn "启动前端开发服务器 (Ctrl+C 停止) ..."
  ( cd "${FRONTEND_DIR}" && npm run dev )
}

run_backend() {
  warn "启动后端服务 (Ctrl+C 停止) ..."
  ( cd "${BACKEND_DIR}" && go run . )
}

# ---------- 菜单 ----------
show_menu() {
  echo
  echo "=============================================="
  echo "         菊传 (Juchuan) 构建助手"
  echo "=============================================="
  echo "  1) 安装前端依赖 (npm install)"
  echo "  2) 编译前端 → 同步到 backend/static"
  echo "  3) 编译前端 + 后端 (go build)"
  echo "  4) 编译前端 + 打包 macOS 应用"
  echo "  5) 只编译后端 (go build)"
  echo "  6) 启动前端开发服务器 (vite dev)"
  echo "  7) 启动后端服务 (go run .)"
  echo "  0) 退出"
  echo "----------------------------------------------"
  printf "请选择 [0-7]: "
}

# ---------- 主循环 ----------
main() {
  check_paths

  while true; do
    show_menu
    if ! read -r choice; then
      echo
      ok "再见！"
      break
    fi
    case "${choice}" in
      1) install_frontend ;;
      2) build_frontend && sync_to_static ;;
      3) build_frontend && sync_to_static && build_backend ;;
      4) build_frontend && sync_to_static && package_macos ;;
      5) build_backend ;;
      6) dev_frontend ;;
      7) run_backend ;;
      0|q|Q) ok "再见！"; break ;;
      *) warn "无效选项，请重新输入。" ;;
    esac
  done
}

main "$@"
