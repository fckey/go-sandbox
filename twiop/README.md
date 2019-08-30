# Overview
twiop is client to interact with Twitter API via [go-twitter library](https://github.com/dghubble/go-twitter).
The manager encapsulate auth token so keep the context inside of manager. 

# How to use
Before start to use the module Twitter API access key of user need to be given. 
Follow instruction in [here](https://developer.twitter.com/en/docs/basics/authentication/overview/oauth.html).

Once the tokens are given, set them as environment variables as below.
```
export TWITTER_CONSUMER_KEY=""
export TWITTER_CONSUMER_SECRET=""
export TWITTER_ACCESS_TOKEN=""
export TWITTER_ACCESS_TOKEN_SECRET="" 
```
