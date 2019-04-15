# Turnstile

Authentication server

User manager and token provider

## Contents table
- [Turnstile](#turnstile)
  - [Contents table](#contents-table)
  - [Endpoints](#endpoints)
  - [License](#license)

## Endpoints

Register: Creates a new user  
URL: /user/register  
Data multiform [POST]:

  - username:
  - password:

Return [JSON]:

  - Status:
  - Message:  

Login: Log in an existing user  
URL: /user/login  
Data [POST]:

- username:
- password:
  
Return [JSON]:

- Status: 
- Message:
- Token:
- RefreshToken:

Check: Check if the user is already logged in  
URL: /token/check  
Data [POST]:

- token: Normal login token

Return [JSON]:

- Status: 
- Message:

Refresh: Get a new token using the refresh token without credentials  
URL: /token/refresh  
Data [POST]:

- token: Refresh token

Return [JSON]:

- Status: 
- Message:
- Token:
- RefreshToken:

**[Go up](#contents-table)**

## License

This project is licensed under the GPLv3 License. See the [LICENSE](LICENSE) file for the full license text.

**[Go up](#contents-table)**


