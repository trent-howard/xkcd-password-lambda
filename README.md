# xkcd password generator

This is just a little toy app I built to learn a little about [AWS CDK](https://docs.aws.amazon.com/cdk/) which creates

- a lambda function that returns a password (inspired by the [xkcd](https://xkcd.com/936/) comic)
- an API gateway
- a usage plan with an API key
- a route with query param validation

I built the lambda handler with Go just because I never had before.

![xkcd password strength comic](https://imgs.xkcd.com/comics/password_strength.png)

## Instructions

You'll need reasonably recent versions of node and npm installed, as well as Go to build the lambda handler.

Be sure to have installed and configured/authenticated [`aws cli`](https://aws.amazon.com/cli/) and `cdk` too.

Install dependencies

```bash
npm install
```

Copy the `.env.template` to `.env` and update the variables inside. `API_KEY` needs to be **at least 20 characters** long.

Compile the handler and deploy the stack with

```bash
npm run deploy
```

If you want to tear it all down run

```bash
npm run destroy
```

Once it's up and running you can make `GET` requests to the `/password` endpoint. Optionally provide a `length` param to request how many words are returned - it defaults to 5 if nothing is provided.
```bash
curl -H "Content-Type: application/json" \
     -H "x-api-key: wifeless-opossum-jiffy-liverwurst-hygienist" \
     "https://[appid].execute-api.[region].amazonaws.com/prod/password?length=8"
{"password":"jaguar-janitor-zodiac-enunciate-sniff-feigned-cactus-imitator"}%
```

### Should I use this?

Probably not - your password manager - which we're all definitely using by now, right guys? - probably has a better generator built in. If you were thinking of deploying this or basing your own stack off of it you might want to set rate limiting and throttling on your usage plan.

