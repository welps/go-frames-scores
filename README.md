# go-frames-scores

WARNING: Most of this is pretty hacked together


## Requirements

- Get an API key from https://rapidapi.com/tipsters/api/sportscore1

## Getting Started

- Run `SPORTS_API_KEY={API_KEY_FROM_ABOVE} go run cmd/go-frames-scores/main.go` 

## Deployment

- This should work out of the box on railway.app 
- You will need to supply the following environment variables when deployed:
  - `ENVIRONMENT` should be `production`
  - `PUBLIC_URL` should be where it deployed. Example: https://go-frames-scores-production.up.railway.app
  - `SPORTS_API_KEY` should be the API key from above
