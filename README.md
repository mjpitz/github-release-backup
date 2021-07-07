# github-release-backup

A simple tool that backs up GitHub release assets to an S3 compatible backend.
By default, GitHub actions will only host artifacts for a maximum of 90 days.
This tool can be used to perform a one-time backup of all release assets.

## Installation

```shell script
go get github.com/mjpitz/github-release-backup
```

## Usage

```shell script
# github configuration
export GITHUB_OWNER="depscloud"
export GITHUB_REPO="depscloud"
export GITHUB_ACCESS_TOKEN="github_access_token" 

# s3 configuration
export S3_ENDPOINT="sfo2.digitaloceanspaces.com"
export S3_DISABLE_SSL="1"
export AWS_DEFAULT_REGION="region"
export AWS_ACCESS_KEY_ID="aws_access_key_id"
export AWS_SECRET_ACCESS_KEY="aws_secret_access_key"
export S3_BUCKET_NAME="bucket_name"

# halt on existing assets (useful for successive runs)
export HALT_ON_EXISTING_ASSET="1"

# run the script
github-release-backup
```
