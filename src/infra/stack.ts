import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import * as lambda from "aws-cdk-lib/aws-lambda";
import { RestApi, LambdaIntegration } from "aws-cdk-lib/aws-apigateway";

declare global {
  namespace NodeJS {
    interface ProcessEnv {
      API_KEY: string;
      API_CLIENT: string;
    }
  }
}
export class PasswordGeneratorStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    const myFunction = new lambda.Function(this, "XkcdPasswordGenerator", {
      code: lambda.Code.fromAsset("dist"),
      handler: "bootstrap",
      runtime: lambda.Runtime.PROVIDED_AL2,
      architecture: lambda.Architecture.ARM_64,
    });

    const gateway = new RestApi(this, "XkcdPasswordApi", {
      defaultCorsPreflightOptions: {
        allowOrigins: ["*"],
        allowMethods: ["GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"],
      },
    });

    const plan = gateway.addUsagePlan("UsagePlan", {
      apiStages: [{ api: gateway, stage: gateway.deploymentStage }],
    });

    const apiKey = gateway.addApiKey(process.env.API_CLIENT, {
      apiKeyName: process.env.API_CLIENT,
      value: process.env.API_KEY,
    });

    plan.addApiKey(apiKey);

    const integration = new LambdaIntegration(myFunction);
    const passwordRoute = gateway.root.addResource("password");
    passwordRoute.addMethod("GET", integration, {
      apiKeyRequired: true,
      requestParameters: { "method.request.querystring.path": false },
    });
  }
}
