#!/usr/bin/env python 
# -*- coding: utf-8 -*-
# @Project : scripts
# @File    : farm.py
# @Time    : 2022/4/15 2:55 PM
# @Author  : ClassmateLin
# @Email   : classmatelin.site@gmail.com
# @Site    : classmatelin.top
# @Cron    : * */12 * * * 
# @Desc    : 东东农场
import asyncio
import json
from urllib.parse import quote

import aiohttp
from rich.console import Console

from app.jd import JdApp, run_jd


class JdFarm(JdApp):
    """
    东东农场
    """

    def __init__(self, **kwargs):
        """
        :param kwargs:
        """
        super(JdFarm, self).__init__(**kwargs)
        self.headers.update({
            'origin': 'https://carry.m.jd.com',
            'referer': 'https://carry.m.jd.com/babelDiy/Zeus/3KSjXqQabiTuD1cJ28QskrpWoBKT/index.html'
        })

    @property
    def app_name(self):
        """
        :return:
        """
        return "东东农场"

    async def request(self, session, function_id, body=None):
        """
        请求数据
        :param session:
        :param body:
        :param function_id:
        :return:
        """
        try:
            if not body:
                body = dict()
            if 'version' not in body:
                body['version'] = 13
            if 'channel' not in body:
                body['channel'] = 1

            url = 'https://api.m.jd.com/client.action?functionId={}&body={}&appid=wh5'.format(function_id,
                                                                                              quote(json.dumps(body)))
            response = await session.get(url=url)
            data = await response.json()
            await asyncio.sleep(1)
            return data
        except Exception as e:
            self.print('{}, 获取服务器数据错误:{}'.format(self.account, e.args))

    async def init_for_farm(self, session):
        """
        初始化农场数据
        :param session:
        :return:
        """
        data = await self.request(session, 'initForFarm')
        if data['code'] != '0' or 'farmUserPro' not in data:
            return None
        return data['farmUserPro']

    async def run(self):
        """
        功能入口
        :return:
        """
        async with aiohttp.ClientSession(headers=self.headers) as session:
            res = await self.init_for_farm(session)
            self.print(res)


if __name__ == '__main__':
    run_jd(JdFarm)
