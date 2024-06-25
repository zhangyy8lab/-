from fastapi import APIRouter
import requests
from config import config
members_router = APIRouter(prefix="/telegram")

# token 在页中进行配置
# token = "7191113667:AAHN2HyAg8cmDCfzFopunIAcYDRFlcXEI48"


@members_router.get('/members')
async def getUserInfo():

    url = f"https://api.telegram.org/bot{config.TelegramToken}/getUpdates"

    response = requests.get(url=url)
    # print("response.json:", response.json())
    return response.json()
