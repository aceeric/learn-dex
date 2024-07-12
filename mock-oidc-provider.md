# Configure an external OIDC Provider

The first Lab used the built-in Dex mock OIDC provider. The second Lab accomplishes the following objectives:

1. Build and run a Mock OIDC Provider
1. Run Dex using the Mock OIDC Provider as a connector
1. Run the example client

***Lab 2 - set the directory context***

This lab is executed in the working directory that you cloned **this** repo into. E.g.:
```
mkdir -p ~/gitrepos
cd ~/gitrepos
git clone https://github.com/aceeric/learn-dex
cd learn-dex
```

## 1. Build and run a Mock OIDC Provider

The Lab uses https://github.com/oauth2-proxy/mockoidc as a Mock OIDC Provider. This provider always returns the same identity, claims, etc. The purpose of this Lab is simply to demonstrate Dex integration with an external OIDC provider. The modified OIDC provider is included in this repo in the `mockoidc/mockoidc` directory. Two minor changes were required to the upstream code:

1. Support a specified port number rather than auto-generating a random port number
2. Add the `name` claim to the ID Token claims list (required, otherwise Dex rejects the claims list...)

Build the Mock OIDC Provider:
```
pushd mockoidc && go build -o bin/mockoidc cmd/main.go && popd
```

Run the Mock OIDC Provider in its own terminal window:
```
mockoidc/bin/mockoidc | jq
```

Result:
```
Press CTRL-C to terminate
{
  "ClientID": "foobar",
  "ClientSecret": "frobozz",
  "Issuer": "http://127.0.0.1:4044/oidc",
  "AccessTTL": 600000000000,
  "RefreshTTL": 3600000000000,
  "CodeChallengeMethodsSupported": [
    "plain",
    "S256"
  ]
}
```

Observe the following:
1. The _Client ID_ is `foobar`
1. The _Client Secret_ is `frobozz`
1. The _Issuer_ is http://127.0.0.1:4044/oidc

The server will run until you press CTRL-C. Leave the Mock OIDC Provider running for now.

## 2. Run Dex with the Mock OIDC Provider as a connector

In order to configure Dex to use this new OIDC Provider, a modified version of the `examples/config-dev.yaml` Dex configuration file from Lab one was used. That modified yaml is in this repo at: `files/config-mockoidc.yaml`. The only change was to specify the new Mock OIDC Provider you created above. Here is the modified portion of the configuration file as compared to the original:

```
connectors:
### Replace mock with mockoidc:
###- type: mockCallback
###  id: mock
###  name: Example
- type: oidc
  id: mockoidc
  name: Mock OIDC
  config:
    issuer: http://127.0.0.1:4044/oidc
    clientID: foobar
    clientSecret: frobozz
    redirectURI: http://127.0.0.1:5556/dex/callback
```

You can see that the connector is specifying the Mock provider on http://127.0.0.1:4044/oidc and the Client ID and Client Secret both match the hard-coded values from the Mock.

Run Dex. Note that we specify `--network=host` since the container will need to communicate with the Mock running as a separate host process. Observe that we supply the new configuration file:
```
docker run\
  --name dex\
  --detach\
  --network=host\
  --publish 5556:5556\
  --volume $PWD/files/config-mockoidc.yaml:/etc/dex/config.yaml\
  --volume /tmp:/examples\
  docker.io/dexidp/dex:latest-alpine\
  dex serve /etc/dex/config.yaml
```

As before, verify Dex is running and was able to connect to the Mock OIDC Provider. `docker logs dex` should produce:
```
time=2024-07-12T18:01:43.822Z level=INFO msg="Version info" dex_version=master go.version=go1.22.3 go.os=linux go.arch=amd64
time=2024-07-12T18:01:43.822Z level=INFO msg="config issuer" issuer=http://127.0.0.1:5556/dex
time=2024-07-12T18:01:43.823Z level=INFO msg="config storage" storage_type=sqlite3
time=2024-07-12T18:01:43.823Z level=INFO msg="config static client" client_name="Example App"
time=2024-07-12T18:01:43.823Z level=INFO msg="config connector" connector_id=mockoidc
time=2024-07-12T18:01:43.823Z level=INFO msg="config connector: local passwords enabled"
time=2024-07-12T18:01:43.823Z level=INFO msg="config refresh tokens rotation" enabled=true
time=2024-07-12T18:01:43.829Z level=INFO msg="listening on" server=telemetry address=0.0.0.0:5558
time=2024-07-12T18:01:43.829Z level=INFO msg="listening on" server=http address=0.0.0.0:5556
```

> Note `config connector" connector_id=mockoidc`

## 3. Run the example client

Just as in the first lab, run the Dex example client. (You may need to modify this depending on where you cloned the Dex repo to in Lab One):
```
~/gitrepos/projects/dex/bin/example-app
```

The CLI will wait for CTRL-C in the console.

While the client is running, access the client from the browser on http://127.0.0.1:5555. You should see the same login screen as in Lab one at this stage:

![Initial page](resources/example-app-login.jpg)

As before, click _Login_, which displays a selection page from the client:

![Log in](resources/login-mock-oidc.jpg)

Notice the new option: _Log in with Mock OIDC_. Click this login option.

After clicking `Log in with Mock OIDC` the example client will display the mandatory OIDC _Consent_ page:

![Consent](resources/grant-access.jpg)


When you then click _Grant Access_, the OIDC Authorization flow is completed **using the external OIDC Provider** and just as before the example client displays a browser page presenting the OIDC information obtained from the OIDC Mock connector:
```


ID Token:

eyJhbGciOiJSUzI1NiIsImtpZCI6IjNhMWUxYmFjZTczNWMwYmJkMzJkNzYxOGRkNjM1ODU0OGZkYzAwMjAifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2dveE1qTTBOVFkzT0Rrd0VnaHRiMk5yYjJsa1l3IiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjA4OTQwOTQsImlhdCI6MTcyMDgwNzY5NCwiYXRfaGFzaCI6IjJXSGtEVTVKLTVBMVNYRU0wRWVaclEiLCJjX2hhc2giOiJyRlpSWi1qNUd5Ukl4X1lfTkpfekRRIiwiZW1haWwiOiJqYW5lLmRvZUBleGFtcGxlLmNvbSIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJuYW1lIjoiamFuZSIsInByZWZlcnJlZF91c2VybmFtZSI6ImphbmUuZG9lIn0.ZZPjHR-ivZ86bW_T7tQMcu61zf2sOQEejpQt_9iaUkSZJMKXNLnGrRjW6bXmr9QX1xKUAAjHcnR_16Zw0wfa5REFK3--4KlxHMK71HyAMyjk7JPO60XKieiqajc4P40gUWxoVLVJCDPexokkXWsPSFXgppEBmkGtFmL7jk8vJOzGBd5Jp7gdcYE-mvToRx78QuXAyawRGKsTZ-LNw4pxHbFw_pcLIKTyowY6Otxwx86qgtP81gOPEkSwsner-7pT24FvZ-6Zv6K242xe1LCqOUldQiL5k7jDRgzkNKqfkFyAJU8u5Yh3S6bgpkFHpdVkeN6dV2yWm4Fe6xfkRSJ_xw

Access Token:

eyJhbGciOiJSUzI1NiIsImtpZCI6IjNhMWUxYmFjZTczNWMwYmJkMzJkNzYxOGRkNjM1ODU0OGZkYzAwMjAifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2dveE1qTTBOVFkzT0Rrd0VnaHRiMk5yYjJsa1l3IiwiYXVkIjoiZXhhbXBsZS1hcHAiLCJleHAiOjE3MjA4OTQwOTQsImlhdCI6MTcyMDgwNzY5NCwiYXRfaGFzaCI6Imw4QXp1aXJTSkdYOHV6MHd2bkVWeUEiLCJlbWFpbCI6ImphbmUuZG9lQGV4YW1wbGUuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJqYW5lIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiamFuZS5kb2UifQ.kLqugnSyiq0LHL_LYA9oHFwXvKV8TGEWBpisPBDYQW-0rVyUSWCAwiksUaOuDCitDftHGWEHO_X2qS8ohlyvJPHjj1abAJMut1QeUeqfYiSQPxHrGy3Z7uQXczi8Pjcp0-TSI_iL-AQjrDJq8H4GFr4X9mj-ijv757cwICD5zHtUq2KmXhtn3nb7ft6x-9ANreBmz-rkCrBbCwsm3nQ3Gn5CXWiVHnL4mtyE_8qcJOoANGLskx16WFM6zlgq7583Ii2ZgZ05IrLpHbvmwwQ2kBTbvgtp2dES6j08OK7IKkAaP3hB2bdOTDDzMEM0Q2-8cbbtWxY7yShOR6rkWDKNEg

Claims:

{
  "iss": "http://127.0.0.1:5556/dex",
  "sub": "CgoxMjM0NTY3ODkwEghtb2Nrb2lkYw",
  "aud": "example-app",
  "exp": 1720894094,
  "iat": 1720807694,
  "at_hash": "2WHkDU5J-5A1SXEM0EeZrQ",
  "c_hash": "rFZRZ-j5GyRIx_Y_NJ_zDQ",
  "email": "jane.doe@example.com",
  "email_verified": true,
  "name": "jane",
  "preferred_username": "jane.doe"
}

Refresh Token:

ChllamZjMmtwNDRyNjRqbzRtNmJubHBuamQyEhliN25zNjRucWF2bnhmeWo2dGdzdW56NWdj
```

Notice the name, email and preferred user name reflect the values hard-coded into the Mock OIDC Provider that was adapted from the upstream OAuth2 Proxy / Mock OIDC Provider.

## Clean up

1. Press `CTRL-C` to stop the example client
1. Press `CTRL-C` to stop the Mock OIDC Provider in its separate terminal window
1. Stop and remove the Dex docker container: `docker rm -f dex`
