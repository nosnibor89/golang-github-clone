package commands

import (
	"errors"
	"flag"
	"fmt"
	cliUtil "github-clone/src/github-clone-cli/util"
	"github-clone/src/repositories"
	"github-clone/src/util"
)

const (
	deleteRepoCommandName = "repo.delete"
)

func NewDeleteRepoCommand() *DeleteRepoCommand {
	command := &DeleteRepoCommand{
		flagSet: flag.NewFlagSet(deleteRepoCommandName, flag.ContinueOnError),
	}

	command.flagSet.StringVar(&command.repo, "repo-name", "", "repository name")

	return command
}

type DeleteRepoCommand struct {
	flagSet *flag.FlagSet

	repo string
}

func (command *DeleteRepoCommand) Name() string {
	return command.flagSet.Name()
}

func (command *DeleteRepoCommand) Init(args []string) error {
	return command.flagSet.Parse(args)
}

func (command *DeleteRepoCommand) assertCommandIsValid() error {
	if util.StringIsEmpty(command.repo) {
		return errors.New(fmt.Sprintf("invalid repo name provided. repo-name :%s", command.repo))
	}

	return nil
}

func (command *DeleteRepoCommand) Run() error {
	fmt.Println("(fake) Deleting Repo", command.repo)

	if err := command.assertCommandIsValid(); err != nil {
		return err
	}

	repo := repositories.GithubRepository{}
	if err := repo.Delete(command.repo, cliUtil.GetCLIUser()); err != nil {
		return err
	}

	fmt.Println("repository deleted correctly")
	return nil
}
