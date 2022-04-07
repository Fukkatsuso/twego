# twego

subscribing to real-time stream of public tweets

## Installation

```bash
$ go install github.com/Fukkatsuso/twego
```

## Usage

### Twitter Auth/Unauth

```bash
$ twego auth --bearer $TWITTER_BEARER_TOKEN
or
$ twego auth --key $TWITTER_API_KEY --secret $TWITTER_API_SECRET
```

Then `Twitter Bearer Token` is saved in `$HOME/.twego/config.toml`.

To delete the token, please execute `unauth` command.

```bash
$ twego unauth
```

### Set Rules

```bash
$ twego rules list
ID                    VALUE                 TAG
1234567890123456789   twitter -is:retweet   twitter
$ twego rules add "golang -is:retweet" --tag "golang"
ID                    VALUE                TAG
0123456789012345678   golang -is:retweet   golang
$ twego rules delete 1234567890123456789 0123456789012345678
1234567890123456789
0123456789012345678
```

### Start Streaming

```bash
$ twego stream
```

## Use Docker

```bash
$ git clone https://github.com/Fukkatsuso/twego.git
$ cd twego
$ docker build ./ -t twego --target runner
$ docker volume create twego-vol
$ docker run --rm -v twego-vol:/root/.twego twego auth --bearer $TWITTER_BEARER_TOKEN
$ docker run --rm -v twego-vol:/root/.twego twego rules list
$ docker run --rm -v twego-vol:/root/.twego twego stream
```

For development,

```bash
$ docker-compose up -d
$ docker-compose exec twego go run main.go auth --bearer $TWITTER_BEARER_TOKEN
$ docker-compose exec twego go run main.go rules list
$ docker-compose exec twego go run main.go stream
```

## References

- [Filtered Stream - Twitter](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction)
- [cobra](https://cobra.dev/)
