# xkcd password generator

This is just a little toy app I built to learn a little about [AWS CDK](https://docs.aws.amazon.com/cdk/) which creates

- a lambda function that returns a password (inspired by the [xkcd](https://xkcd.com/936/) comic)
- an API gateway
- a usage plan with an API key
- a route with query param validation

I built the lambda handler with Go just because I never had before.

![xkcd password strength comic](https://imgs.xkcd.com/comics/password_strength.png)

## Instructions

### Setting up AWS resources

Install, configure, and authenticate the [`aws cli`](https://aws.amazon.com/cli/) tool.

Install the [`CDK toolkit`](https://docs.aws.amazon.com/cdk/v2/guide/cli.html).

[Bootstrap the environment](https://docs.aws.amazon.com/cdk/v2/guide/bootstrapping.html) you want to deploy to. This only needs to be done once per environment.

```bash
cdk bootstrap
```

This provisions some IAM roles to grant permissions needed for deployments and an S3 bucket to store any files that will be deployed. This needs to be run by a user who has a pretty high level of access, maybe your root admin account.

It's generally recommended not to have your root credentials floating around too many places. Once you're bootstrapped you can create an IAM user to run your deployments instead. We just need create and assign a role allowing it to assume the roles CDK bootstrapped. This will give it all the necessary permissions to handle our deployments without giving it the keys to the whole kingdom.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["sts:AssumeRole"],
      "Resource": ["arn:aws:iam::*:role/cdk-*"]
    }
  ]
}
```

### Using this project

You'll need reasonably recent versions of node (nvmrc is set to 20) and pnpm installed, as well as Go to build the lambda handler.

Install dependencies

```bash
pnpm install
```

Compile the handler and deploy the stack with. Once deployed you'll be able to see your stack in CloudFormation and can retrieve your API key from API Gateway.

```bash
pnpm run deploy
```

Once it's up and running you can make `GET` requests to the `/password` endpoint. Optionally provide a `length` param to request how many words are returned - it defaults to 5 if nothing is provided.

```bash
curl -H "Content-Type: application/json" \
     -H "x-api-key: QD178HGkwOxIFkv7Hf7K1Piq5A8zBiJSXYsXjDWO" \
     "https://[appid].execute-api.[region].amazonaws.com/prod/password?length=8"
{"password":"jaguar-janitor-zodiac-enunciate-sniff-feigned-cactus-imitator"}%
```

If you want to tear it all down you can either delete the stack in Cloud Formation or just run

```bash
pnpm run destroy
```

### Should I use this?

Probably not. Your password manager - which we're all definitely using by now, right guys? - probably has a better generator built in. If you were thinking of deploying this or basing your own stack off of it you might want to set rate limiting and throttling on your usage plan.
