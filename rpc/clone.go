package rpc

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-merkledag"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

// CloneArgs contains the args.
type CloneArgs struct {
	// Cwd is the current working directory.
	Cwd string
	// Dir is the name of the directory to create.
	Dir string
	// ID is the CID of the repo.
	ID cid.Cid
	// Limit is the number of children to fetch.
	Limit int
	// Branch is the name of the branch to clone.
	Branch string
}

// CloneReply contains the reply.
type CloneReply struct {
	// ID is the CID of the commit.
	ID cid.Cid
	// Root is the repo root path.
	Root string
}

// Clone copies a commit tree to the working directory.
func (s *Service) Clone(args *CloneArgs, reply *CloneReply) error {
	ctx := context.Background()

	repo, err := data.GetRepository(ctx, s.client, args.ID)
	if err != nil {
		return err
	}

	if args.Dir == "" {
		args.Dir = repo.Name
	}

	id, ok := repo.Branches[args.Branch]
	if !ok {
		return errors.New("branch does not exist")
	}

	if err := merkledag.FetchGraphWithDepthLimit(ctx, id, args.Limit, s.client); err != nil {
		return err
	}

	commit, err := data.GetCommit(ctx, s.client, id)
	if err != nil {
		return err
	}

	tree, err := s.client.Get(ctx, commit.Tree)
	if err != nil {
		return err
	}

	path := filepath.Join(args.Cwd, args.Dir)
	if err := os.Mkdir(path, 0755); err != nil {
		return err
	}

	reply.ID = id
	reply.Root = path

	return unixfs.Write(ctx, s.client, path, tree)
}