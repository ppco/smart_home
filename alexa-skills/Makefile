build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap src/main.go
	zip alexa-skills.zip bootstrap

include .env
create-lambda:
	export AWS_PROFILE=${AWS_PROFILE} && \
	aws lambda create-function --function-name alexa_smart_home \
	--runtime provided.al2023 --handler main \
	--architectures x86_64 \
	--role arn:aws:iam::${ACCOUNT_ID}:role/lambda-role \
	--zip-file fileb://alexa-skills.zip \
	--region ap-northeast-1

include .env
update-lambda:
	export AWS_PROFILE=${AWS_PROFILE} && \
	aws lambda update-function-code --function-name alexa_smart_home \
	--zip-file fileb://alexa-skills.zip \
	--region ap-northeast-1