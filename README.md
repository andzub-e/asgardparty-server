## Adgard Party Game Engine

REST API for Asgard Party slot game

------

#### Requirements

- Golang 1.15+
- Docker (optional)

------

#### Environments

| Name                   | Default           | Description                         |
|------------------------|-------------------| ----------------------------------- |
| HTTP_HOST              | 0.0.0.0           | server host                         |
| HTTP_PORT              | 8086              | server port                         |
| ASGARD_PARTY_SSL       | false             | If you use https, it must be 'true' |
| ASGARD_PARTY_CERT_PATH | ./ssl/ejawsdk.crt | Path to ssl certificate             |
| SGARD_PARTY_CERT_KEY   | ./ssl/ejawsdk.key | Path to ssl certificate key         |

------

#### How to run

Setup all required environments from table and run commands for build and run:

```bash
go build -o asgardparty-server -v ./cmd/slot && ./asgardparty-server
```

or

```bash
go run ./cmd/slot
```

For Docker run command (at first you need to run overlord in your docker):

```bash
docker compose up -d
```

------

#### Game flow

Every game begins from request to **/core/state** endpoint. In response you get **session_token** for next requests.

Then user can place a bet by calling **/core/wager** endpoint.

Endpoints **/core/free_spins** and **/core/spins_history** needs for getting client info about his free spins and spins history respectively.

------

### Endpoints descriptions

more info in **/swagger/index.html**

**State**

Returns user state from bet-overlord service.

Endpoint: */core/state*

Method: *GET*

*Request query parameters:*

```json
{
  "integrator": "mock", 
  "game": "test", 
  "params": {
    "currency": "RUB", 
    "integrator": "mock", 
    "game": "test", 
    "user_id": "NQaoqnlZAZ7FnMPuIrac6kzdaduISZd4VIGoDsWZ", 
    "userlocale": "da-DK"}
}
```

*Response example:*

status 200

```json
{
  "THE_EJAW_SLOT": "string",
  "balance": 0,
  "currency": 0,
  "default_wager": 0,
  "error": "string",
  "freespinid": 0,
  "last_wager": 0,
  "reels": {
    "amount": 0,
    "is_autospin": true,
    "is_cheat_stops": true,
    "is_turbospin": true,
    "spins": [
      {}
    ]
  },
  "session_token": "string",
  "spins_indexes": {
    "base_stage_index": 0,
    "bonus_spin_index": 0
  },
  "total_wins": 0,
  "username": "string",
  "wager_levels": [
    0
  ],
  "wallet_play_id": "string"
}
```

status 400

```json
{
  "code": 400,
  "message": "INTERNAL_SERVER_ERROR"
}
```

status 500

```json
{
  "code": 500,
  "message": "INTERNAL_SERVER_ERROR"
}
```

------

**Wager**

Make a bet (spin)

Endpoint: */core/wager*

Method: *POST*

*Request body parameters:*

| Parameter name | Type    | Description                   | Required |
| -------------- | ------- | ----------------------------- | -------- |
| session_token  | string  | token for current session     | required |
| wager          | integer | amount of money for bet       | required |
| currency       | string  | current currency              | required |
| freespin_id    | string  | id of a free spin from casino | optional |

*Response example:*

status 200

```json
{
  "THE_EJAW_SLOT": "string",
  "balance": 0,
  "currency": 0,
  "default_wager": 0,
  "error": "string",
  "freespinid": 0,
  "last_wager": 0,
  "reels": {
    "amount": 0,
    "is_autospin": true,
    "is_cheat_stops": true,
    "is_turbospin": true,
    "spins": [
      {}
    ]
  },
  "session_token": "string",
  "spins_indexes": {
    "base_stage_index": 0,
    "bonus_spin_index": 0
  },
  "total_wins": 0,
  "username": "string",
  "wager_levels": [
    0
  ],
  "wallet_play_id": "string"
}
```

status 400

```json
{
  "code": 400,
  "message": "INTERNAL_SERVER_ERROR"
}
```

status 500

```json
{
  "code": 500,
  "message": "INTERNAL_SERVER_ERROR"
}
```