# First touch

The first Lab accomplishes the following objectives:

1. Git clone the Dex GitHub repo
1. Run the Dex server in Docker
1. Build the Dex example client
1. Run the example client

## 1. Git clone the Dex GitHub repo

```
mkdir -p ~/gitrepos
cd ~/gitrepos
git clone https://github.com/dexidp/dex.git
cd dex
```

> These labs were performed against the following git branch/SHA:
> ```
> git log --pretty=oneline | head -n1
> ```
> 
> Result:
> ```
> 17a9c1c2429901ce962ce6dda3711e35ab33bdf0 build(deps): bump actions/attest-build-provenance from 1.1.2 to 1.3.1 (#3580)
> ```

## 2. Run the Dex server in Docker

This step uses the dev server configuration in `./examples/config-dev.yaml` in the cloned repository. Since the dev configuration wants to store the SQLite database in `/examples` in the container we mount the host `/tmp` directory to `/examples` in the container because the container does not include that directory:
```
docker run\
  --name dex\
  --detach\
  --publish 5556:5556\
  --volume $PWD/examples/config-dev.yaml:/etc/dex/config-dev.yaml\
  --volume /tmp:/examples\
  docker.io/dexidp/dex:latest-alpine\
  dex serve /etc/dex/config-dev.yaml
```

## Validate that the Dex server started

After running `docker logs dex`, anything other than the following is probably an error:
```
time=2024-07-12T11:25:13.634Z level=INFO msg="Version info" dex_version=master go.version=go1.22.3 go.os=linux go.arch=amd64
time=2024-07-12T11:25:13.634Z level=INFO msg="config issuer" issuer=http://127.0.0.1:5556/dex
time=2024-07-12T11:25:13.724Z level=INFO msg="config storage" storage_type=sqlite3
time=2024-07-12T11:25:13.724Z level=INFO msg="config static client" client_name="Example App"
time=2024-07-12T11:25:13.724Z level=INFO msg="config connector" connector_id=mock
time=2024-07-12T11:25:13.724Z level=INFO msg="config connector: local passwords enabled"
time=2024-07-12T11:25:13.724Z level=INFO msg="config refresh tokens rotation" enabled=true
time=2024-07-12T11:25:13.727Z level=INFO msg="keys expired, rotating"
time=2024-07-12T11:25:13.938Z level=INFO msg="keys rotated" next_rotation=2024-07-12T17:25:13.934Z
time=2024-07-12T11:25:13.938Z level=INFO msg="listening on" server=telemetry address=0.0.0.0:5558
time=2024-07-12T11:25:13.938Z level=INFO msg="listening on" server=http address=0.0.0.0:5556
```

Confirm the Dex server port is exposed by Docker using `docker port dex`. Result should be:
```
5556/tcp -> 0.0.0.0:5556
5556/tcp -> [::]:5556
```

## 3. Build the Dex example client

Use the provided Makefile:
```
make examples
```

## 4. Run the example client

The client is pre-configured to expect Dex on localhost port 5556:
```
./bin/example-app
```

The CLI will wait for CTRL-C in the console.

While the client is running, access the client from the browser on http://127.0.0.1:5555:

![Initial page](resources/example-app-login.jpg)

Just click _Login_, which displays a selection page from the client:

![Log in](resources/login-to-dex.jpg)

In the login page, you see _Log in with Example_. Where does _Example_ come from? Observe `examples/config-dev.yaml` which is what was used to configure the Dex server:
```
...
connectors:
- type: mockCallback
  id: mock
  name: Example
...
```

This will invoke the Dex _Mock_ connecter which is suitable for smoke testing.

After clicking `Log in with Example` the example client will display the mandatory OIDC _Consent_ page:

![Consent](resources/grant-access.jpg)


When you then click _Grant Access_, the OIDC Authorization flow is completed and the example client simply displays a brower page presenting the OIDC information obtained from the Dex Mock connector:
```
ID Token:

eyJhbGciOiJSUzI1NiIsImtpZCI6IjJkYWJiMGEyM2ZjMWU0OWQxMDdiM2FlZjk5MjRhNTFkOGZiYWE4MTQifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2cwd0xUTTROUzB5T0RBNE9TMHdFZ1J0YjJOciIsImF1ZCI6ImV4YW1wbGUtYXBwIiwiZXhwIjoxNzIwODcwNTI0LCJpYXQiOjE3MjA3ODQxMjQsImF0X2hhc2giOiJSVGdGZU1XMm4zZmlRRmYwVzBXSnZRIiwiY19oYXNoIjoibnZYRFZXcFVFX0p1aGFDSVAxMlNvdyIsImVtYWlsIjoia2lsZ29yZUBraWxnb3JlLnRyb3V0IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsIm5hbWUiOiJLaWxnb3JlIFRyb3V0In0.Tf-3_WlvfSANC1cbQqJO1TrqO9PnllGc9uYFSQm9UOjKOW15vUF5Y-ANb2x4wCmMhHg9ATVev8bwJlGbTqs4vblZi7GudQF6SujyqjKvSB0D5CD1hvo-Ncpo6v0I19-neaK5ExpEEr29mVAyxuiX4niIDMfRlxG2-aHW4SDocuB-ASsN81bIyOelN2ZTuo-IkOlZZpptmiyg_dQNG_fVl3wpyrpcDXu6-exVtEhUJyBbEkwTHR8jUo-h7JzCgAXfgnbM73iPq9-dCmI-hYt1a2RnQxvonEGEfeJAi0ZgCuf_MYKiKWJQu5SwULmwLyQqkWU51_u6thOns65GtT7amg

Access Token:

eyJhbGciOiJSUzI1NiIsImtpZCI6IjJkYWJiMGEyM2ZjMWU0OWQxMDdiM2FlZjk5MjRhNTFkOGZiYWE4MTQifQ.eyJpc3MiOiJodHRwOi8vMTI3LjAuMC4xOjU1NTYvZGV4Iiwic3ViIjoiQ2cwd0xUTTROUzB5T0RBNE9TMHdFZ1J0YjJOciIsImF1ZCI6ImV4YW1wbGUtYXBwIiwiZXhwIjoxNzIwODcwNTI0LCJpYXQiOjE3MjA3ODQxMjQsImF0X2hhc2giOiJidE1tUmJZUEkyT2w4eEZHYUpscS13IiwiZW1haWwiOiJraWxnb3JlQGtpbGdvcmUudHJvdXQiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwibmFtZSI6IktpbGdvcmUgVHJvdXQifQ.incACaphaJMdW-wt5-iR7dSsbyzzTM_U1wrPCHwaQxt5PoSA5EMY8GBZsukpEbB7pfwrCnXpYEo1eNbLx3jc3Oa35-lCSRuBVMu6X_RBkIpRB0GzO826ZWulDX_P8YoE6ggdEgKkYlLbmmlqQpqlOEYvvdGMfXsbu7rZ8XM2ti1aDXXLGmoPplA0vdx_U_xqPDtbkW70wfNxErSC99er27u2aVsdKYCOdoyYMa_Kn8-az4xELS6nqI2D_Xd06esi-yYxtlL5y3uDdC25twRpZeMMI447WjhEkh_fPnIMuk_bDVEwxA_SKVXSlmL3InAC77ygOJ869GNNLFwKpw0U8w

Claims:

{
  "iss": "http://127.0.0.1:5556/dex",
  "sub": "Cg0wLTM4NS0yODA4OS0wEgRtb2Nr",
  "aud": "example-app",
  "exp": 1720870524,
  "iat": 1720784124,
  "at_hash": "RTgFeMW2n3fiQFf0W0WJvQ",
  "c_hash": "nvXDVWpUE_JuhaCIP12Sow",
  "email": "kilgore@kilgore.trout",
  "email_verified": true,
  "name": "Kilgore Trout"
}

Refresh Token:

Chlic2Nwa25xeXJocmI2ZDIzcGNibXd1eDRsEhltaTdhem51bnVhcGF2Y3ZxbnQ3NWVjZmFq
```

If we navigate to https://jwt.io and paste the _ID Token_ value into the page we see the _Payload_ from the parsed JWT ID Token that matches the _Claims_ displayed by the example client above:
```
{
  "iss": "http://127.0.0.1:5556/dex",
  "sub": "Cg0wLTM4NS0yODA4OS0wEgRtb2Nr",
  "aud": "example-app",
  "exp": 1720870524,
  "iat": 1720784124,
  "at_hash": "RTgFeMW2n3fiQFf0W0WJvQ",
  "c_hash": "nvXDVWpUE_JuhaCIP12Sow",
  "email": "kilgore@kilgore.trout",
  "email_verified": true,
  "name": "Kilgore Trout"
}
```

## Clean up

1. Press `CTRL-C` to stop the example client.
1. Stop and remove the Dex docker container: `docker rm -f dex`
