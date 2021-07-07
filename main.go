package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/google/go-github/v36/github"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"golang.org/x/oauth2"
)

func main() {
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")
	bucketName := os.Getenv("S3_BUCKET_NAME")
	haltOnExistingAsset := os.Getenv("HALT_ON_EXISTING_ASSET") != ""

	err := func() error {
		ctx := context.Background()

		s3, err := minio.New(os.Getenv("S3_ENDPOINT"), &minio.Options{
			Creds: credentials.NewEnvAWS(),
			Secure: os.Getenv("S3_DISABLE_SSL") == "",
			Region: os.Getenv("AWS_DEFAULT_REGION"),
		})

		if err != nil {
			return fmt.Errorf("failed to create s3 client: %v", err)
		}

		var httpClient *http.Client

		if accessToken := os.Getenv("GITHUB_ACCESS_TOKEN"); accessToken != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{
					AccessToken: accessToken,
				},
			)
			httpClient = oauth2.NewClient(ctx, ts)
		}

		gh := github.NewClient(httpClient)

		listOptions := &github.ListOptions{
			Page: 1,
			PerPage: 100,
		}

		for listOptions.Page != 0 {
			log.Printf("fetching releases page %d for %s/%s\n", listOptions.Page, owner, repo)

			releases, resp, err := gh.Repositories.ListReleases(ctx, owner, repo, listOptions)
			if err != nil {
				return fmt.Errorf("failed to list releases: %v", err)
			}

			listOptions.Page = resp.NextPage

			for _, release := range releases {
				log.Printf("fetching assets for release %s/%s#%s\n", owner, repo, release.GetTagName())

				assets, _, err := gh.Repositories.ListReleaseAssets(ctx, owner, repo, release.GetID(), &github.ListOptions{
					PerPage: 100,
				})

				if err != nil {
					return fmt.Errorf("failed to list assets: %v", err)
				}

				for _, asset := range assets {
					reader, redirectURL, err := gh.Repositories.DownloadReleaseAsset(ctx, owner, repo, asset.GetID(), http.DefaultClient)
					if err != nil {
						return fmt.Errorf("failed download asset: %v", err)
					}

					if redirectURL != "" {
						req, err := http.NewRequestWithContext(ctx, http.MethodGet, redirectURL, nil)
						if err != nil {
							return fmt.Errorf("failed to create download request: %v", err)
						}

						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return fmt.Errorf("failed download asset: %v", err)
						}

						reader = resp.Body
					}

					assetName := path.Join(repo, release.GetTagName(), asset.GetName())

					obj, _ := s3.GetObject(ctx, bucketName, assetName, minio.GetObjectOptions{})
					_, err = obj.Stat()

					if err != nil {
						log.Printf("storing %s\n", assetName)

						_, err = s3.PutObject(
							ctx,
							bucketName,
							assetName,
							reader,
							int64(asset.GetSize()),
							minio.PutObjectOptions{
								ContentType: asset.GetContentType(),
							},
						)

						if err != nil {
							return fmt.Errorf("failed to upload asset, %v", err)
						}
					} else if haltOnExistingAsset {
						log.Printf("assets already exist, halting")
						return nil
					}
				}
			}
		}

		return nil
	}()

	if err != nil {
		panic(err.Error())
	}
}
