# gitlab-webhook-create

## Usage

```
docker run \
    -e GITLAB_TOKEN=mytoken \
    -e WEBHOOKS_FILE_PATH="./webhooks.json" \
    -e WEBHOOK_SECRET_TOKEN="123456" \
    hazim/gitlab-webhook-create 
```