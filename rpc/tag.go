package rpc

import (
	"context"
	"errors"

	"github.com/ipfs/go-cid"
	"github.com/multiverse-vcs/go-multiverse/data"
)

// TagArgs contains the args.
type TagArgs struct {
	// Name is the repo name.
	Name string
	// Tag is the name of the tag.
	Tag string
	// Head is the CID of the repo head.
	Head cid.Cid
}

// TagReply contains the reply.
type TagReply struct {
	Tags map[string]cid.Cid
}

// ListTags returns the repo tags.
func (s *Service) ListTags(args *TagArgs, reply *TagReply) error {
	ctx := context.Background()

	id, err := s.store.GetCid(args.Name)
	if err != nil {
		return err
	}

	repo, err := data.GetRepository(ctx, s.node, id)
	if err != nil {
		return err
	}

	reply.Tags = repo.Tags
	return nil
}

// CreateTag creates a new tag.
func (s *Service) CreateTag(args *TagArgs, reply *TagReply) error {
	ctx := context.Background()

	id, err := s.store.GetCid(args.Name)
	if err != nil {
		return err
	}

	repo, err := data.GetRepository(ctx, s.node, id)
	if err != nil {
		return err
	}

	if args.Tag == "" {
		return errors.New("name cannot be empty")
	}

	if _, ok := repo.Tags[args.Tag]; ok {
		return errors.New("tag already exists")
	}

	repo.Tags[args.Tag] = args.Head
	reply.Tags = repo.Tags

	id, err = data.AddRepository(ctx, s.node, repo)
	if err != nil {
		return err
	}

	return s.store.PutCid(repo.Name, id)
}

// DeleteTag deletes an existing tag.
func (s *Service) DeleteTag(args *TagArgs, reply *TagReply) error {
	ctx := context.Background()

	id, err := s.store.GetCid(args.Name)
	if err != nil {
		return err
	}

	repo, err := data.GetRepository(ctx, s.node, id)
	if err != nil {
		return err
	}

	if args.Tag == "" {
		return errors.New("name cannot be empty")
	}

	if _, ok := repo.Tags[args.Tag]; !ok {
		return errors.New("tag does not exists")
	}

	delete(repo.Tags, args.Tag)
	reply.Tags = repo.Tags

	id, err = data.AddRepository(ctx, s.node, repo)
	if err != nil {
		return err
	}

	return s.store.PutCid(repo.Name, id)
}
