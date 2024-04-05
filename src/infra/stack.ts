import { Size, Stack, RemovalPolicy } from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import {
  RestApi,
  LambdaIntegration,
  Cors,
  MethodLoggingLevel,
  AccessLogFormat,
  LogGroupLogDestination,
} from "aws-cdk-lib/aws-apigateway";
import { LogGroup, RetentionDays } from "aws-cdk-lib/aws-logs";

export class PasswordGeneratorStack extends Stack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    const logs = new LogGroup(this, `${id}-public-api-logs`, {
      logGroupName: `${id}-public-api-logs`,
      removalPolicy: RemovalPolicy.DESTROY,
      retention: RetentionDays.ONE_WEEK,
    });

    const restApi = new RestApi(this, `${id}-public-api`, {
      restApiName: `${id}-public-api`,
      minCompressionSize: Size.bytes(0),
      deployOptions: {
        loggingLevel: MethodLoggingLevel.INFO,
        accessLogFormat: AccessLogFormat.jsonWithStandardFields(),
        accessLogDestination: new LogGroupLogDestination(logs),
      },
      defaultCorsPreflightOptions: {
        allowOrigins: Cors.ALL_ORIGINS,
        allowMethods: Cors.ALL_METHODS,
        allowHeaders: Cors.DEFAULT_HEADERS,
      },
    });

    const handler = new lambda.Function(this, `${id}-handler`, {
      functionName: "xkcd-password-generator-handler",
      code: lambda.Code.fromAsset("dist"),
      handler: "bootstrap",
      runtime: lambda.Runtime.PROVIDED_AL2,
      architecture: lambda.Architecture.ARM_64,
    });

    const plan = restApi.addUsagePlan(`${id}-public-api-usage-plan`, {
      apiStages: [{ api: restApi, stage: restApi.deploymentStage }],
    });

    const apiKey = restApi.addApiKey(`${id}-public-api-key`, {
      apiKeyName: `${id}-public-api-usage-plan`,
    });

    plan.addApiKey(apiKey);

    const integration = new LambdaIntegration(handler);

    restApi.root.resourceForPath("/password").addMethod("GET", integration, {
      apiKeyRequired: true,
      requestParameters: {
        "method.request.querystring.length": false,
      },
      requestValidatorOptions: {
        validateRequestParameters: true,
      },
    });
  }
}
