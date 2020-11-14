package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s path/to/products.xml ACCESS_TOKEN\n(got %+v)\n", os.Args[0], os.Args[1:])

		return
	}

	content, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		panic(fmt.Errorf("failed to read file content at path %s: %w", os.Args[1], err))
	}

	pushFile(content, "my-pricing-test", "my-pricing-test.github.io", "public/products.xml", "updated products")
}

func pushFile(content []byte, owner, repo, path, message string) error {
	ctx := context.Background()

	client := github.NewClient(oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Args[2]}),
	))

	file, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)

	if err != nil {
		return fmt.Errorf("failed to push file to github: %w", err)
	}

	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, path, &github.RepositoryContentFileOptions{
		Message: &message, Content: content, SHA: file.SHA,
	})

	if err != nil {
		return fmt.Errorf("failed to push file to github: %w", err)
	}

	return nil
}
