from fastapi import APIRouter, Request
from config import config
import requests

twitter_router = APIRouter(prefix="/twitter")

# redirect_uri = "http://localhost:5004/twitter/callback"
# scopes = "tweet.read+users.read+space.read+"
# client_id = "SXVNSjBoU1hjMUNYRjRWdGk2THQ6MTpjaQ"
# client_secret = "5BiHlZcMnW41vdotHhjbjNNAg2PY56akpFRzehCQOau87CaQ4a"
# consumer_key = "aAJzdBu1PMqikhvMZhsmWsCTz"
# consumer_secret = "DesjWwsYTaAM2cKMGEI4QFDkRI9i0gXK8Yu6gr3MSCrBN4SMix"
# access_token = "1770726160583831552-SIHMK9vX8YovG3mZN6TgbQD9AuCe04"
# access_token_secret = "aS1TmHddXfx5jdvAGEqycyyJw4haMDx3AEiUuGymB46qG"
# state = "n0YHbbLhyF1X4Fw4C9S49hLBwvHduz"
# code_challenge = "P1l5VDcX"


@twitter_router.get('/login')
async def authUrlV2():
    authorization_url = f"https://twitter.com/i/oauth2/authorize?" \
                        f"response_type=code&" \
                        f"client_id={config.TwitterClientId}&" \
                        f"redirect_uri={config.TwitterCallbackUrl}&" \
                        f"scope={config.TwitterScopes}&" \
                        f"state={config.TwitterState}&" \
                        f"code_challenge={config.TwitterCodeChallenge}&" \
                        f"code_challenge_method=plain"

    # print("oauthUrl:", authorization_url)
    return {"oauthUrl": authorization_url}


@twitter_router.get("/callback")
async def callbackV2(request: Request):
    # print(request.query_params.items())

    code = request.query_params.get("code")
    payload = {
        "code": code,
        "grant_type": "authorization_code",
        "client_id": config.TwitterClientId,
        "client_secret": config.TwitterClientSecret,
        "redirect_uri": config.TwitterCallbackUrl,
        "code_verifier": "P1l5VDcX",
    }

    headers = {
        "Content-Type": "application/x-www-form-urlencoded",
    }

    response = requests.post(
        "https://api.twitter.com/2/oauth2/token",
        auth=(config.TwitterClientId, config.TwitterClientSecret),
        data=payload,
        headers=headers
    )

    if response.status_code == 200:
        return {"url": "https://x.com/TusimaNetwork"}
        # return response.json()
    else:
        return {"error": response.text}
