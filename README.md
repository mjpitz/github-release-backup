# github-release-backup

A minimal tool that backs up GitHub release assets to an S3 compatible backend.
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
export S3_ENDPOINT="sfo2.digitaloceanspaces.com"  # configure the s3 endpoint
export S3_DISABLE_SSL=""                          # set to non-empty value to disable tls
export AWS_DEFAULT_REGION="region"
export AWS_ACCESS_KEY_ID="aws_access_key_id"
export AWS_SECRET_ACCESS_KEY="aws_secret_access_key"
export S3_BUCKET_NAME="bucket_name"

# configure existing asset behavior
export ON_EXISTING_ASSET="halt"      # stop when we encounter an asset that already exists
export ON_EXISTING_ASSET="overwrite" # overwrite assets that already exist
export ON_EXISTING_ASSET=""          # skip assets that already exist

# run the script
github-release-backup
```
