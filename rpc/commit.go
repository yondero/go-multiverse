package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/data"
	"github.com/multiverse-vcs/go-multiverse/unixfs"
)

// CommitArgs contains the args.
type CommitArgs struct {
	// Name is the name of the repo.
	Name string
	// Branch is the name of the branch to update.
	Branch string
	// Root is the repo root path.
	Root string
	// Ignore is a list of paths to ignore.
	Ignore []string
	// Message is the description of changes.
	Message string
	// Parent is the CID of the parent commit.
	Parent cid.Cid
}

// CommitReply contains the reply.
type CommitReply struct {
	// Index is the CID of the commit.
	Index cid.Cid
}

// Commit records changes to the repo
func (s *Service) Commit(args *CommitArgs, reply *CommitReply) error {
	ctx := context.Background()

	if args.Branch == "" {
		return errors.New("branch cannot be empty")
	}

	id, err := s.store.GetCid(args.Name)
	if err != nil {
		return err
	}

	repo, err := data.GetRepository(ctx, s.client, id)
	if err != nil {
		return err
	}

	head, ok := repo.Branches[args.Branch]
	if ok && args.Parent != head {
		return errors.New("branch is ahead of parent")
	}

	tree, err := unixfs.Add(ctx, s.client, args.Root, args.Ignore)
	if err != nil {
		return err
	}

	commit := data.NewCommit(tree.Cid(), args.Message)
	if args.Parent.Defined() {
		commit.Parents = append(commit.Parents, args.Parent)
	}

	head, err = data.AddCommit(ctx, s.client, commit)
	if err != nil {
		return err
	}

	repo.Branches[args.Branch] = head

	id, err = data.AddRepository(ctx, s.client, repo)
	if err != nil {
		return err
	}

	reply.Index = head
	return s.store.PutCid(repo.Name, id)
}