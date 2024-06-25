# -*- coding: utf-8 -*-
"""
@File        : discord.py
@Author      : Aug
@Time        : 2023/1/6 14:24
@Description :
"""
from fastapi import APIRouter
import requests
from config import config

api_router = APIRouter(prefix="/discord")


def getUserGuilds(headers):
    url = "https://discordapp.com/api/users/@me/guilds"
    response = requests.get(url, headers=headers)
    for item in response.json():
        if item.get("name") == config.Discord_server_name:
            print("exist")
            return item
            # return {"result": "ok"}
            # return response.json()
    # [{'id': '965918503070728203', 'name': 'TusimaDAO', ...}]

    return {"result": "not found"}


@api_router.get("/login")
async def login():
    """
    Returns a discord auth link, please manually redirect the user then it goes to the callback url with the
    query parameter "code" (example: https://callbackurl/?code=isfd78f2UIRFerf) to get the code to use a function
    called getTokens().

    The code can only be used on an active url (callback url) meaning you can only use the code once
    """
    result_url = f"https://discord.com/oauth2/authorize?" \
                 f"client_id={config.Discord_client_id}&" \
                 f"redirect_uri={config.Discord_callback_url}&" \
                 f"response_type=code&" \
                 f"scope={config.Discord_scope}&" \
                 f"state=state"
    return {"url": result_url}


@api_router.get("/callback")
async def callback(code: str = None):
    data = {
        'client_id': config.Discord_client_id,
        'client_secret': config.Discord_client_secret,
        'grant_type': 'authorization_code',
        'code': code,
        'redirect_uri': config.Discord_callback_url
    }

    headers = {
        'Content-Type': 'application/x-www-form-urlencoded'
    }

    tokens = requests.post('https://discord.com/api/oauth2/token', data=data, headers=headers)
    # print("tokens:", tokens.json())
    headers = {
        "Authorization": f'Bearer {tokens.json().get("access_token")}'
    }

    return getUserGuilds(headers=headers)
