#!/bin/bash
set -e;

echo "===================开始更新=================";

if [ -z "$CODE_DIR" ]; then
  echo "设置脚本目录: " + "$CODE_DIR" + "...";
  CODE_DIR="/go/src/scripts";
fi

if [ -z "$SSH_ROOT" ]; then
    SSH_ROOT=~/.ssh/;
fi

if [ ! -d $SSH_ROOT ]; then
  mkdir ~/.ssh/;
  ssh-keyscan github.com > $SSH_ROOT/known_hosts;
fi


if [ ! -d $CODE_DIR ]; then
    mkdir -p $CODE_DIR;
fi


if [ ! -d $CODE_DIR/.git ]; then
  echo "脚本目录为空, 开始初始化...";
  if [ -z "$REPO_URL" ]; then
    echo "设置Git仓库地址: " + "$CODE_DIR" + "...";
    REPO_URL=https://github.com/ClassmateLin/scripts_v2.git;
  fi
  cd $CODE_DIR;
  git init;
  git remote add origin $REPO_URL;
  git pull origin master;
  git branch --set-upstream-to=origin/master master;
else
  echo "正在获取最新脚本..."
  cd $CODE_DIR && git reset --hard HEAD~10 && git pull;
fi

if [ ! -f "$CODE_DIR/config/config.json" ]; then
  echo "初始化脚本配置文件: " + "$CODE_DIR/config/config.json";
  cp $CODE_DIR/.config.json $CODE_DIR/config/config.json;
fi

if [ ! -f "$CODE_DIR/config/crontab.json" ]; then
  echo "初始化用户定时任务配置文件:"  + "$CODE_DIR/config/crontab.json";
  cp $CODE_DIR/.crontab.json $CODE_DIR/config/crontab.json;
fi

if [ ! -f "$CODE_DIR/storage/logs/logger.log" ]; then
  touch $CODE_DIR/storage/logs/logger.log;
fi

echo "正在更新依赖包...";
go get -u;

cd $CODE_DIR;

echo "正在更新配置文件...";
go run $CODE_DIR/tools/update_config.go;

rm -f $CODE_DIR/*.bin;
echo "正在编译脚本...";
go run $CODE_DIR/tools/build_scripts.go;
chmod +x $CODE_DIR/*.bin;

echo "正在更新定时任务...";
go run $CODE_DIR/tools/update_crontab.go;

echo "重启cron..."
/etc/init.d/cron restart;

echo "正在更新docker-entrypoint命令...";
cp $CODE_DIR/shell/docker-entrypoint.sh /bin/docker-entrypoint;
chmod +x /bin/docker-entrypoint;

echo "===================更新完成================="

tail -f /dev/null;