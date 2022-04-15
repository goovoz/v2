#!/usr/bin/env python 
# -*- coding: utf-8 -*-
# @Project : scripts
# @File    : config.py
# @Time    : 2022/4/15 2:59 PM
# @Author  : ClassmateLin
# @Email   : classmatelin.site@gmail.com
# @Site    : classmatelin.top
# @Cron    : * */12 * * * 
# @Desc    :

import os.path
import pathlib
import sys

import yaml

# 根目录
BASE_DIR = os.path.split(os.path.split(pathlib.Path(__file__).absolute())[0])[0]

# 日志目录
LOG_DIR = os.path.join(BASE_DIR, 'logs')

# 配置文件路径
CONF_PATH = os.path.join(BASE_DIR, 'conf', 'config.yaml')

try:
    with open(CONF_PATH, encoding='utf-8') as f:   # demo.yaml内容同上例yaml字符串
        cfg = yaml.safe_load(f)
except FileNotFoundError as e:
    print(f'读取配置文件错误, 找不到文件:{CONF_PATH}')
    sys.exit(0)


JD_CONF = cfg.get('jd', dict())

