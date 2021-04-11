# Etsy Integration

Developers must authenticate with the Etsy API for their local instance of our
API to work.

The Etsy API is free to use, up to 5000 requests per day.

You can get an API key by signing up on the
[Etsy developer portal](https://www.etsy.com/developers/register).

Once you have your API Key, perform the following steps:

1. Make a copy of the `.env.example` file and name it `.env`
1. Open the `.env` file you just made.
1. Paste your API key between the quotes.
1. Save the file.
1. Run `make stop`.
1. Run `make`.
1. Success.

## Example

```sh
#!/bin/bash

# Developers! Paste your Etsy API Key here for your local API to work.
#
# You can get a key at:
# https://www.etsy.com/developers/register
export ETSY_API_KEY="asdfjklqwertyuiopzxcvbnm"
```

## Did it work?

If the API key was NOT detected, you will see this error:

```log
api              | building...
api              | running...
api              | INFO[0000] Connected to the database.
api              | FATA[0000] must set ETSY_API_KEY
```

If the API key was detected, you will see:

```log
api              | building...
api              | running...
api              | INFO[0000] Connected to the database.
api              | INFO[0000] Starting API server for tier local on port 8080.
```

## But did it really work?

The only way to know is to hit one of the search endpoints and see if there are
any errors in the logs.

Just because an API key was detected on startup does not mean that Etsy will
accept that API key when we make a request.
