package repositories

import (
	"bytes"

	repo "github.com/fsamin/go-repo"

	"github.com/ovh/cds/sdk"
	"github.com/ovh/cds/sdk/log"
)

func (s *Service) processPushFiles(op *sdk.Operation) error {
	r := s.Repo(*op)
	if err := s.checkOrCreateFS(r); err != nil {
		log.Error("Repositories> processCheckout> checkOrCreateFS> [%s] Error %v", op.UUID, err)
		return err
	}

	// Get the git repository
	opts := []repo.Option{repo.WithVerbose()}
	if op.RepositoryStrategy.ConnectionType == "ssh" {
		log.Debug("Repositories> processCheckout> using ssh key %s", op.RepositoryStrategy.SSHKey)
		opts = append(opts, repo.WithSSHAuth([]byte(op.RepositoryStrategy.SSHKeyContent)))
	} else if op.RepositoryStrategy.User != "" && op.RepositoryStrategy.Password != "" {
		opts = append(opts, repo.WithHTTPAuth(op.RepositoryStrategy.User, op.RepositoryStrategy.Password))
	}

	gitRepo, err := repo.New(r.Basedir, opts...)
	if err != nil {
		log.Debug("Repositories> processCheckout> cloning %s into %s", r.URL, r.Basedir)
		if _, err = repo.Clone(r.Basedir, r.URL, opts...); err != nil {
			log.Error("Repositories> processCheckout> Clone> [%s] error %v", op.UUID, err)
			return err
		}
	}

	if err := gitRepo.Pull("origin", op.Setup.PushFiles.ToBranch); err != nil {
		log.Error("Repositories> processCheckout> Pull> [%s] Error: %v", op.UUID, err)
		return err
	}

	if err := gitRepo.CheckoutNewBranch(op.Setup.PushFiles.FromBranch); err != nil {
		log.Error("Repositories> processCheckout> Checkout> [%s] Error: %v", op.UUID, err)
		return err
	}

	defer func(r *repo.Repo) {
		//At the end, reset all the things
		if err := r.ResetHard("HEAD"); err != nil {
			log.Error("Repositories> processCheckout> clean error - reset: %v", err)
		}
		if err := r.Checkout(op.Setup.PushFiles.ToBranch); err != nil {
			log.Error("Repositories> processCheckout> clean error - checkout: %v", err)
		}
		if err := r.DeleteBranch(op.Setup.PushFiles.FromBranch); err != nil {
			log.Error("Repositories> processCheckout> clean error - delete: %v", err)
		}
	}(&gitRepo)

	for path, content := range op.Setup.PushFiles.Files {
		if err := gitRepo.Write(path, bytes.NewBuffer(content)); err != nil {
			log.Error("Repositories> processCheckout> Write> [%s] Error: %v", op.UUID, err)
			return err
		}
		if err := gitRepo.Add(path); err != nil {
			log.Error("Repositories> processCheckout> Add> [%s] Error: %v", op.UUID, err)
			return err
		}
	}

	if err := gitRepo.Commit(op.Setup.Checkout.Commit); err != nil {
		log.Error("Repositories> processCheckout> Commit> [%s] Error: %v", op.UUID, err)
		return err
	}

	if err := gitRepo.Push("origin", op.Setup.PushFiles.FromBranch); err != nil {
		log.Error("Repositories> processCheckout> Push> [%s] Error: %v", op.UUID, err)
		return err
	}

	//Check branch
	currentBranch, err := gitRepo.CurrentBranch()
	if err != nil {
		log.Error("Repositories> processCheckout> CurrentBranch> [%s] error %v", op.UUID, err)
		return err
	}

	if op.Setup.Checkout.Branch == "" {
		op.Setup.Checkout.Branch, err = gitRepo.DefaultBranch()
		if err != nil {
			log.Error("Repositories> processCheckout> DefaultBranch> [%s] error %v", op.UUID, err)
			return err
		}
	}

	if currentBranch != op.Setup.Checkout.Branch {
		log.Debug("Repositories> processCheckout> fetching branch %s from %s", op.Setup.Checkout.Branch, r.URL)
		if err := gitRepo.FetchRemoteBranch("origin", op.Setup.Checkout.Branch); err != nil {
			log.Error("Repositories> processCheckout> FetchRemoteBranch> [%s] error %v", op.UUID, err)
			return err
		}
	}

	//Check commit
	if op.Setup.Checkout.Commit == "" {
		log.Debug("Repositories> processCheckout> pulling branch %s", op.Setup.Checkout.Branch)
		if err := gitRepo.Pull("origin", op.Setup.Checkout.Branch); err != nil {
			log.Error("Repositories> processCheckout> Pull without commit> [%s] error %v", op.UUID, err)
			return err
		}
	} else {
		currentCommit, err := gitRepo.LatestCommit()
		if err != nil {
			log.Error("Repositories> processCheckout> LatestCommit> [%s] error %v", op.UUID, err)
			return err
		}
		if currentCommit.LongHash != op.Setup.Checkout.Commit {
			//Not the same commit
			//Pull and reset HARD the commit
			log.Debug("Repositories> processCheckout> pulling branch %s", op.Setup.Checkout.Branch)
			if err := gitRepo.Pull("origin", op.Setup.Checkout.Branch); err != nil {
				log.Error("Repositories> processCheckout> Pull with commit > [%s] error %v", op.UUID, err)
				return err
			}

			log.Debug("Repositories> processCheckout> reseting commit %s", op.Setup.Checkout.Commit)
			if err := gitRepo.ResetHard(op.Setup.Checkout.Commit); err != nil {
				log.Error("Repositories> processCheckout> ResetHard> [%s] error %v", op.UUID, err)
				return err
			}
		}
	}

	log.Info("Repositories> processCheckout> repository %s ready", r.URL)
	return nil
}
