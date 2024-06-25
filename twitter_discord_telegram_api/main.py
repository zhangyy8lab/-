import sys

from fastapi import APIRouter
from api.twitter import twitter_oauth, twitter_users
from api.discord import discord
from api.telegram import telegram

sys.path.append('api')

app = APIRouter()

app.include_router(twitter_oauth.twitter_router)
app.include_router(twitter_users.user_router)
app.include_router(discord.api_router)
app.include_router(telegram.members_router)

if __name__ == '__main__':
    import uvicorn
    uvicorn.run(app, host="localhost", port=5004)
