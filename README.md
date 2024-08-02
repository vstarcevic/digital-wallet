<h1>Digital Wallet</h1>
<p>Microservices Golang example</p>
<p>Simple golang project, two microservices, each of them using Postgres database as data storage and Kafka and Nats for communication</p>
</p>
    <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue?style=plastic" height="20" alt="License"></a>
</p>

## What can you do with it?

   * Create user
   * Get balance for user
   * Add money to the user
   * Transfer money from one user to another

## Features
 
### Http Requests and responses for user service 
```javascript
POST /user 
{
    "email" : "email@example.com"
}
 
HTTP/1.1. 200
Content-type: application/json
{
    "user_id": 1,
    "email": "email@example.com",
    "created_at": "2000-01-01T12:00:00.00000Z"
}
```
```javascript
POST /balance 
{
    "email" : "email@example.com"
}

HTTP/1.1. 200
Content-type: application/json
{
    "balance": "0",
    "email": "email@example.com"
}
```
### Http Requests and responses for transaction service 
```javascript
POST /add-money 
{
    "user_id" : 1,
    "amount" : 100.00
}

HTTP/1.1. 200
Content-type: application/json
{
    "updated_balance": "100"
}
```

```javascript
POST /transfer-money 
{
    "from_user_id" : 1,
    "to_user_id" : 2,
    "amount_to_transfer" : 100
}

HTTP/1.1. 200
Content-type: application/json
{}
```
### Installation

```bash
# From digital-wallet directory
docker compose up -d
```