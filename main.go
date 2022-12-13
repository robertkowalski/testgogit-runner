package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	url := "https://github.com/robertkowalski/testgogit"
	dir := "test"
	branch := "releases"

	cloneOptions := git.CloneOptions{
		URL: url,
	}

	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	r, err := git.PlainCloneContext(ctx, dir, false, &cloneOptions)
	if err != nil {
		log.Fatalf("Unable to clone %v %s: %v\n", url, branch, err)
	}

	fetchOptions := git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*"},
	}

	pullOptions := git.PullOptions{
		Force: true,
	}

	err = r.Fetch(&fetchOptions)
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		panic("done")
	}

	refName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branch))

	rev := plumbing.Revision(path.Join("refs", "heads", branch))

	var hash *plumbing.Hash
	if hash, err = r.ResolveRevision(rev); err != nil {
		log.Fatalf("Unable to resolve %s: %v\n", branch, err)
	}

	var w *git.Worktree
	if w, err = r.Worktree(); err != nil {
		panic("done")
	}

	if err = w.Pull(&pullOptions); err != nil && err != git.NoErrAlreadyUpToDate {
		fmt.Printf("Unable to pull %s: %v\n", refName, err)
	}
	fmt.Printf("Cloned %s to %s\n", hash, dir)

	if err = w.Checkout(&git.CheckoutOptions{
		Hash:  *hash,
		Force: true,
	}); err != nil {
		log.Fatalf("Unable to checkout %s: %v\n", *hash, err)
	}

	if err = w.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: *hash,
	}); err != nil {
		log.Fatalf("Unable to reset to %s: %v\n", hash.String(), err)
	}

	fmt.Printf("Checked out %s\n", branch)
}
