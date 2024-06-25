from fastapi import APIRouter, Request
import requests
user_router = APIRouter()


@user_router.get('/me')
async def getUserInfo(request: Request):
    access_token = request.query_params.get("access_token")

    url = "https://api.twitter.com/2/users/me"
    headers = {"Authorization": f"Bearer {access_token}"}

    response = requests.get(url=url, headers=headers)
    print("response.json:", response.json())
    return response.json()
