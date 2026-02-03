set dotenv-load := true

# Build CLI binary
build-cli:
    go build -o bin/free2kindle ./cmd/cli
    GOOS=linux GOARCH=amd64 go build -o bin/lambda-bootstrap ./cmd/lambda

# Build CLI and convert URL into EPUB
run *ARGS: build-cli
    ./bin/free2kindle convert {{ ARGS }}

# Run linter
lint:
    golangci-lint run

# Run tests
test:
    go test ./...

# Build Lambda binary for Linux
build-lambda:
    GOOS=linux GOARCH=amd64 go build -o cmd/lambda/bootstrap ./cmd/lambda

# Build Lambda zip for deployment
build-lambda-zip: build-lambda
    cd cmd/lambda && zip -j ../../bin/lambda-source.zip bootstrap

# Deploy S3 bucket for Lambda source
deploy-bucket PROJECT_NAME:
    aws cloudformation deploy \
        --template-file infra/bucket.yaml \
        --stack-name {{ PROJECT_NAME }}-bucket \
        --parameter-overrides BucketName={{ PROJECT_NAME }}-lambda-source

# Upload Lambda source zip to S3
upload-zip PROJECT_NAME:
    aws s3 cp bin/lambda-source.zip s3://{{ PROJECT_NAME }}-lambda-source/function-source.zip

# Deploy Lambda infrastructure
deploy PROJECT_NAME:
    aws cloudformation deploy \
        --template-file infra/infra.yaml \
        --stack-name {{ PROJECT_NAME }}-infra \
        --capabilities CAPABILITY_NAMED_IAM \
        --parameter-overrides \
            ProjectName={{ PROJECT_NAME }} \
            SourceBucketName={{ PROJECT_NAME }}-lambda-source \
            MailjetAPIKey="$MAILJET_API_KEY" \
            MailjetAPISecret="$MAILJET_API_SECRET" \
            APIKeySecret="$API_KEY_SECRET" \
            KindleEmail="$F2K_KINDLE_EMAIL" \
            SenderEmail="$F2K_SENDER_EMAIL"

# Full deployment (bucket + upload + infra)
deploy-all PROJECT_NAME: build-lambda-zip
    just deploy-bucket {{ PROJECT_NAME }}
    just upload-zip {{ PROJECT_NAME }}
    just deploy {{ PROJECT_NAME }}

# Destroy Lambda infrastructure
destroy PROJECT_NAME:
    aws cloudformation delete-stack --stack-name {{ PROJECT_NAME }}-infra
    aws cloudformation wait stack-delete-complete --stack-name {{ PROJECT_NAME }}-infra
    aws s3 rm s3://{{ PROJECT_NAME }}-lambda-source --recursive
    aws cloudformation delete-stack --stack-name {{ PROJECT_NAME }}-bucket
    aws cloudformation wait stack-delete-complete --stack-name {{ PROJECT_NAME }}-bucket

# Get Lambda function URL
get-url PROJECT_NAME:
    aws cloudformation describe-stacks \
        --stack-name {{ PROJECT_NAME }}-infra \
        --query "Stacks[0].Outputs[?OutputKey=='FunctionUrl'].OutputValue" \
        --output text

# View CloudFormation stack events
events PROJECT_NAME:
    aws cloudformation describe-stack-events --stack-name {{ PROJECT_NAME }}-infra --query 'StackEvents[*].[Timestamp,ResourceStatus,ResourceType]' --output table

# View Lambda function logs
logs PROJECT_NAME:
    aws logs tail /aws/lambda/{{ PROJECT_NAME }}-api --follow
