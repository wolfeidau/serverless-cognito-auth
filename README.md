# serverless-cognito-auth

This serverless application provides a [AWS Cognito](https://aws.amazon.com/cognito/) user pool with supporting [Lambda](https://aws.amazon.com/lambda/) hooks to enable a modern web application to authenticate. It integrates analytics and monitoring out of the box.

# Why?

AWS Cognito has got to the point where providing a simple template with a couple of inline lambdas really doesn't provide enough value, or take advantage of the wide array of features. 

# Features

This application incorporates a range of out of the box features:

* SNS Topic which publishes all sign ups, and sign in events, this can be used to maintain a session table or analytics.
* Optional invitation code which is required to sign up.
* Optional Email domain whitelisting to restrict sign up.
* Optional sign up notifications to an email address.
* Analytics provided by [Amazon Pinpoint](https://aws.amazon.com/pinpoint/) see [Using Amazon Pinpoint Analytics with Amazon Cognito User Pools ](https://docs.aws.amazon.com/cognito/latest/developerguide/cognito-user-pools-pinpoint-integration.html)

# Usage

For usage see [example app.yaml](https://github.com/wolfeidau/serverless-cognito-auth/blob/master/example/app.yaml)

# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au/).