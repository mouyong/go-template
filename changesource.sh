#!/usr/bin/env bash

set -e

# 获取系统发行版及版本号
if [ -r /etc/os-release ]; then
    # shellcheck disable=SC1091
    . /etc/os-release
else
    echo "Unable to detect OS release info, skipping changesource"
    exit 0
fi

DIST_ID=${ID:-}
VERSION_ID=${VERSION_ID:-}

# 用户可选择的镜像源
MIRROR_PROVIDER=${MIRROR_PROVIDER:-"xtom"}  # 默认为 xtom, 可手动改为 tencent, tsinghua

# 定义各版本的 Ubuntu 镜像源
declare -A MIRROR_SOURCE_UBUNTU
MIRROR_SOURCE_UBUNTU["22_xtom"]="mirror-cdn.xtom.com"
MIRROR_SOURCE_UBUNTU["24_xtom"]="mirror-cdn.xtom.com"
MIRROR_SOURCE_UBUNTU["22_tencent"]="mirrors.tencent.com"
MIRROR_SOURCE_UBUNTU["24_tencent"]="mirrors.tencentyun.com"
MIRROR_SOURCE_UBUNTU["22_tsinghua"]="mirrors.tuna.tsinghua.edu.cn"
MIRROR_SOURCE_UBUNTU["24_tsinghua"]="mirrors.tuna.tsinghua.edu.cn"

# 定义各版本的 Debian 镜像源
declare -A MIRROR_SOURCE_DEBIAN
MIRROR_SOURCE_DEBIAN["12_xtom"]="mirrors.xtom.com"
MIRROR_SOURCE_DEBIAN["12_tencent"]="mirrors.tencent.com"
MIRROR_SOURCE_DEBIAN["12_tsinghua"]="mirrors.tuna.tsinghua.edu.cn"

case "$DIST_ID" in
    ubuntu)
        KEY="${VERSION_ID%%.*}_$MIRROR_PROVIDER"
        MIRROR="${MIRROR_SOURCE_UBUNTU[$KEY]}"

        if [ -z "$MIRROR" ]; then
            echo "Unsupported Ubuntu version ($VERSION_ID) or provider ($MIRROR_PROVIDER)"
            exit 1
        fi

        echo "Using mirror: $MIRROR for Ubuntu $VERSION_ID"

        if [ "$VERSION_ID" = "22.04" ]; then
            FILE_PATH="/etc/apt/sources.list"
        elif [ "$VERSION_ID" = "24.04" ]; then
            FILE_PATH="/etc/apt/sources.list.d/ubuntu.sources"
        else
            echo "Unknown Ubuntu version: $VERSION_ID"
            exit 1
        fi

        sed -i "s|ports.ubuntu.com|$MIRROR|g" "$FILE_PATH"
        sed -i "s|archive.ubuntu.com|$MIRROR|g" "$FILE_PATH"
        sed -i "s|security.ubuntu.com|$MIRROR|g" "$FILE_PATH"
        sed -i "s|security-cdn.ubuntu.com|$MIRROR|g" "$FILE_PATH"

        echo "Source changed successfully to $MIRROR"
        ;;
    debian)
        KEY="${VERSION_ID%%.*}_$MIRROR_PROVIDER"
        MIRROR="${MIRROR_SOURCE_DEBIAN[$KEY]}"

        if [ -z "$MIRROR" ]; then
            echo "Unsupported Debian version ($VERSION_ID) or provider ($MIRROR_PROVIDER)"
            exit 1
        fi

        echo "Using mirror: $MIRROR for Debian $VERSION_ID"

        # Debian 12+ 使用新格式的 sources 文件
        if [ "${VERSION_ID%%.*}" -ge 12 ]; then
            FILE_PATH="/etc/apt/sources.list.d/debian.sources"
        else
            FILE_PATH="/etc/apt/sources.list"
        fi

        if [ ! -f "$FILE_PATH" ]; then
            echo "Warning: $FILE_PATH not found, skipping mirror change"
            exit 0
        fi

        sed -i "s|deb.debian.org|$MIRROR|g" "$FILE_PATH"
        sed -i "s|security.debian.org|$MIRROR|g" "$FILE_PATH"
        sed -i "s|security-cdn.debian.org|$MIRROR|g" "$FILE_PATH"

        echo "Source changed successfully to $MIRROR"
        ;;
    *)
        echo "Skipping mirror change for $DIST_ID $VERSION_ID"
        ;;
esac
