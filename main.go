package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
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
		panic(fmt.Errorf("failed to read file content at path %s: %s\n%w", os.Args[1], err.Error(), err))
	}

	err = pushFile(content, "kivi-pu", "products", "products.xml", "updated products")

	if err != nil {
		panic(err)
	}
}

func pushFile(content []byte, owner, repo, path, message string) error {
	ctx := context.Background()

	httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	client := github.NewClient(oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Args[2]}),
	))

	file, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, nil)

	if err != nil {
		return fmt.Errorf("failed to get file from github: %s\n%w", err.Error(), err)
	}

	_, _, err = client.Repositories.UpdateFile(ctx, owner, repo, path, &github.RepositoryContentFileOptions{
		Message: &message, Content: content, SHA: file.SHA,
	})

	if err != nil {
		return fmt.Errorf("failed to push file to github: %s\n%w", err.Error(), err)
	}

	return nil
}
