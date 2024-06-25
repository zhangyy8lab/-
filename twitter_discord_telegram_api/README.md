# Twitter/Discord/Telegram Api 



## Introduce

In this version, the service needs to invoke third party logins, and here is how each third party uses them



## Env

- Linux/mac
- python3.8



## Install lib/package

- pip3 install -r requirements.txt



## Unit Test

### Twitter

#### setting

1. login `https://developer.x.com/en/portal/dashboard` 
2. Projects & Apps -> application -> setting callbackUrl

3. get Keys and tokens 



#### login

```python
# request - get 
http://localhost:5004/twitter/login
```

```python
# response
{
"oauthUrl": "https://twitter.com/i/oauth2/authorize?response_type=code&client_id=SXVNSjBoU1hjMUNYRjRWdGk2THQ6MTpjaQ&redirect_uri=http://localhost:5004/twitter/callback&scope=tweet.read+users.read+space.read+&state=n0YHbbLhyF1X4Fw4C9S49hLBwvHduz&code_challenge=P1l5VDcX&code_challenge_method=plain"
}
```

#### authorize

```python
# request - get 
https://twitter.com/i/oauth2/authorize?response_type=code&client_id=SXVNSjBoU1hjMUNYRjRWdGk2THQ6MTpjaQ&redirect_uri=http://localhost:5004/twitter/callback&scope=tweet.read+users.read+space.read+&state=n0YHbbLhyF1X4Fw4C9S49hLBwvHduz&code_challenge=P1l5VDcX&code_challenge_method=plain
```

```python
# response
{
"token_type": "bearer",
"expires_in": 7200,
"access_token": "MUNvVFNZOThZS3FzdnRJTXFqUThYYVV6aFpJcXFVVl81UHNXMFRwLWtfWkZsOjE3MTg2MTYxMjU2Mzg6MToxOmF0OjE",
"scope": "users.read space.read tweet.read"
}
```





### Discord

#### setting

1. login `https://discord.com/developers/applications`
2. Application -> My Applications -> Oauth2 
3. Get the information you need eg: `client_id` `client_secret` `bot_token` ...



#### login

```python
# request - get 
http://localhost:5004/discord/login
```

```python
# response 
{
"url": "https://discord.com/oauth2/authorize?client_id=1250016680764244018&redirect_uri=http://localhost:5004/discord/callback&response_type=code&scope=identify%20email%20guilds.members.read%20guilds&state=state"
}
```



#### authorize

```python
# request - get 
https://discord.com/oauth2/authorize?client_id=1250016680764244018&redirect_uri=http://localhost:5004/discord/callback&response_type=code&scope=identify%20email%20guilds.members.read%20guilds&state=state
```

```python
# response 
# You will get a response to whether you are in the specified group, and the following responses are in the response body
[{
		"id": "965918503070728203",
		"name": "TusimaDAO",
		"icon": "5f79db69513db33996d6d0cadddc303d",
		"owner": false,
		"permissions": 0,
		"permissions_new": "140737488355328",
		"features": [
			"THREE_DAY_THREAD_ARCHIVE",
			"COMMUNITY_EXP_LARGE_UNGATED",
			"MEMBER_PROFILES",
			"MEMBER_VERIFICATION_GATE_ENABLED",
			"PRIVATE_THREADS",
			"INVITE_SPLASH",
			"CHANNEL_ICON_EMOJIS_GENERATED",
			"ROLE_ICONS",
			"AUTO_MODERATION",
			"BANNER",
			"ANIMATED_ICON",
			"VANITY_URL",
			"NEWS",
			"SOUNDBOARD",
			"PREVIEW_ENABLED",
			"SEVEN_DAY_THREAD_ARCHIVE",
			"ANIMATED_BANNER",
			"COMMUNITY"
		]
	},
	{
		"id": "973915550965174292",
		"name": "Blockscout",
		"icon": "227c3693a744a597d91f9fffb79e6c6b",
		"owner": false,
		"permissions": 3508289,
		"permissions_new": "1548528987244609",
		"features": [
			"THREE_DAY_THREAD_ARCHIVE",
			"SEVEN_DAY_THREAD_ARCHIVE",
			"MEMBER_PROFILES",
			"MEMBER_VERIFICATION_GATE_ENABLED",
			"AUTO_MODERATION",
			"BANNER",
			"ANIMATED_ICON",
			"VANITY_URL",
			"PRIVATE_THREADS",
			"NEWS",
			"SOUNDBOARD",
			"ROLE_ICONS",
			"PREVIEW_ENABLED",
			"INVITE_SPLASH",
			"CHANNEL_ICON_EMOJIS_GENERATED",
			"ANIMATED_BANNER",
			"COMMUNITY"
		]
	},
	{
		"id": "1250045125208576041",
		"name": "zhangyy27-service",
		"icon": null,
		"owner": true,
		"permissions": 2147483647,
		"permissions_new": "2251799813685247",
		"features": []
	}
]
```





### Telegram

#### setting



#### members

```python
# request - get
http://localhost:5004/telegram/members
```

```python
# response
{
  "ok": true,
  "result": []
}

# You get the information you want on the backend
```





## Config

- Config file path $Project/config/config.py

- Application is configured by category

## Run

- /usr/bin/python3.8 $Project/main.py



## reference

### twitter 

- `https://developer.x.com/en/portal/dashboard`

- `https://developer.x.com/en/docs/authentication/guides/v2-authentication-mapping`
- `https://blog.csdn.net/qq_38935605/article/details/136389612`

- `https://devpress.csdn.net/python/62f99905c6770329307fef15.html#devmenu5`

### discord

`https://discord.com/developers/applications/1250016680764244018/bot`



### telegram

- `https://luoji.men/2023/01/telegram-robot-token-and-chatid-acquisition-tutorial/`

- `https://atjiu.github.io/2023/03/09/telegram-bot-api/`