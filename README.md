# Atrovan - Q1

## Description

### jwt

- role based middleware
- refresh token

### gorm (orm) with postgres

- lessons
- users
- many to many relation - users_lessons

### redis

- cache the lessons when you want to get all
- black list when user logout

## Run

### generate jwt secret files

clone the project then go to `/config/keys`

```console
yeganeh@ubuntu:~$ ssh-keygen -t rsa -b 4096 -m PEM -f private_key.pem
yeganeh@ubuntu:~$ openssl rsa -in private_key.pem -pubout -outform PEM -out public_key.pub
```

notice : Don't add passphrase

### tests

```console
yeganeh@ubuntu:~$ go test -v -coverpkg=./... -coverprofile=profile.out ./tests
```

### Docker

runs app on port 5000 of container and 8080 of host, it depends on postgres and redis.

```console
yeganeh@ubuntu:~$ docker-compose up -d
```

### manual

make sure your redis and postgres are running on the background.

```console
yeganeh@ubuntu:~$ go build
yeganeh@ubuntu:~$ ./Atrovan_Q1
```
