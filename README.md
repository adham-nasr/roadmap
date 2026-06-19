# NestJS Lambdalith on AWS Serverless

A production-ready NestJS application deployed as a single Lambda function (Lambdalith pattern) using AWS SAM, with full local development support via LocalStack.

## 🏗️ Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────────┐
│           Amazon API Gateway                │
│         (HTTP/ REST API)                    │
└──────────────────┬──────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────┐
│         AWS Lambda (NestJS)                 │
│   ┌─────────────────────────────────┐      │
│   │    NestJS Application            │      │
│   │  ├── Controllers                 │      │
│   │  ├── Services                    │      │
│   │  ├── Modules                     │      │
│   │  └── DTOs                        │      │
│   └─────────────────────────────────┘      │
└─────────────────────────────────────────────┘
```

### How It Works

1. **API Gateway** receives all HTTP requests
2. Routes are proxied to a **single Lambda function**
3. Lambda executes your **NestJS application**
4. NestJS handles routing, business logic, and responses
5. API Gateway returns the response to the client

## 📋 Prerequisites

- [Node.js 22.x](https://nodejs.org/)
- [AWS CLI](https://aws.amazon.com/cli/)
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
- [Docker](https://www.docker.com/get-started)
- [LocalStack](https://localstack.cloud/) (optional, for local AWS emulation)
- [awslocal](https://github.com/localstack/awscli-local) (LocalStack CLI wrapper)

## 🚀 Quick Start

### 1. Clone & Install

```bash
# Install dependencies
npm install

# Install SAM CLI (if not installed)
# macOS
brew install aws-sam-cli

# Linux
pip install aws-sam-cli

# Windows
# Download from: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html
```

### 2. Build Your NestJS Application

```bash
# Build the NestJS app
npm run build

# SAM build
sam build
```

### 3. Local Development

#### Option A: SAM Local (Fastest)
```bash
# Start API locally
sam local start-api

# Your API is now available at:
# http://localhost:3000
```

#### Option B: LocalStack (Full AWS Emulation)
```bash
# Start LocalStack
docker-compose up -d

# Deploy to LocalStack
samlocal deploy --guided \
  --stack-name nest-lambdalith-local \
  --s3-bucket localstack-bucket \
  --no-fail-on-empty-changeset

# Your API is available at:
# http://localhost:4566/_aws/execute-api/{api-id}/Prod/
```

### 4. Deploy to AWS

```bash
# Build
sam build

# Deploy (guided setup)
sam deploy --guided

# Follow the prompts:
# - Stack Name: nest-lambdalith
# - AWS Region: us-east-1 (or your preferred)
# - Confirm changes: y
# - Allow SAM to create IAM roles: y
```

## 📁 Project Structure

```
├── src/
│   ├── main.ts              # Application entry point
│   ├── app.module.ts        # Root module
│   ├── app.controller.ts    # Root controller
│   └── app.service.ts       # Root service
├── dist/                    # Compiled output (generated)
├── roadmap/                 # SAM deployment source
├── template.yaml            # SAM infrastructure
├── package.json
├── tsconfig.json
└── README.md
```

## 🔧 Configuration

### template.yaml Overview

```yaml
AWSTemplateFormatVersion: '2010-09-09'
Transform: AWS::Serverless-2016-10-31

Globals:
  Function:
    Timeout: 30          # Max execution time
    MemorySize: 512      # RAM allocation (MB)

Resources:
  NestLambdalith:
    Type: AWS::Serverless::Function
    Properties:
      Runtime: nodejs22.x
      Handler: dist/main.handler
      CodeUri: roadmap/
      Events:
        ProxyApi:
          Type: Api
          Properties:
            Path: /{proxy+}    # Catch-all route
            Method: ANY        # All HTTP methods
        RootApi:
          Type: Api
          Properties:
            Path: /
            Method: ANY
```


### API Testing with cURL

```bash
# Test locally
curl http://localhost:3000/
curl http://localhost:3000/users

# Test on AWS
curl https://{api-id}.execute-api.{region}.amazonaws.com/Prod/
```

## 📦 Deployment Commands Reference

| Command | Purpose |
|---------|---------|
| `sam build` | Build the SAM application |
| `sam local start-api` | Start local API server |
| `sam local invoke NestLambdalith` | Invoke function locally |
| `sam deploy --guided` | Deploy to AWS |
| `sam logs -n NestLambdalith` | View Lambda logs |
| `sam delete` | Delete the stack |

## 🔍 Debugging

### Local Debugging
```bash
# Run with debug port
sam local start-api --debug-port 5858

# Attach your IDE to port 5858
```

### AWS Debugging
```bash
# View logs in real-time
sam logs -n NestLambdalith --tail

# View specific time range
sam logs -n NestLambdalith --start-time "2024-01-01T00:00:00Z"
```




## 📊 Monitoring

### AWS Console Access
- **Lambda**: AWS Console → Lambda → nest-lambdalith
- **API Gateway**: AWS Console → API Gateway → nest-lambdalith
- **CloudWatch Logs**: AWS Console → CloudWatch → Log Groups → /aws/lambda/nest-lambdalith

### Metrics to Watch
- **Lambda**: Duration, Invocations, Errors, Throttles
- **API Gateway**: 4xx/5xx errors, Latency, Request count

## 🔒 Security Best Practices

1. **Never commit secrets** - Use AWS Secrets Manager or SSM Parameter Store
2. **Enable API Gateway throttling** - Prevents DDoS
3. **Use IAM roles** - Principle of least privilege
4. **Enable logging** - CloudWatch for auditing
5. **Use environment variables** - For configuration

## 📈 Cost Optimization

| Strategy | Savings |
|----------|---------|
| Right-size memory (128MB-512MB) | Up to 75% |
| Use Provisioned Concurrency | For predictable traffic |
| Enable caching in API Gateway | Reduces Lambda invocations |
| Set appropriate timeout | Avoids overpaying |

## 🆚 Lambdalith vs Microservices

| Aspect | Lambdalith (This) | Microservices |
|--------|-------------------|---------------|
| **Cost** | ✅ Lower (one function) | ❌ Higher (many functions) |
| **Cold Start** | ❌ One cold start affects all | ✅ Isolated per service |
| **Complexity** | ✅ Simpler | ❌ Complex orchestration |
| **Deployment** | ✅ One deployment | ❌ Multiple deployments |
| **Scale** | ❌ Scales as one unit | ✅ Granular scaling |

## 📚 Additional Resources

- [NestJS Documentation](https://docs.nestjs.com/)
- [AWS SAM Documentation](https://docs.aws.amazon.com/serverless-application-model/)
- [LocalStack Documentation](https://docs.localstack.cloud/)
- [Serverless-http NPM Package](https://www.npmjs.com/package/serverless-http)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Submit a pull request

## 📝 License

MIT

## 🆘 Support

- Open an issue for bugs
- Discussion for questions
- Pull requests for improvements

---

**Happy Building! 🚀**