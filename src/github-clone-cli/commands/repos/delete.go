package repos

import (
	"flag"
	"fmt"
	cliUtil "github-clone/src/github-clone-cli/util"
	"github-clone/src/repositories"
)

const (
	deleteRepoCommandName = "repo.delete"
)

type DeleteRepoCommand struct {
	singleOperationRepoCommand
}

func NewDeleteRepoCommand() *DeleteRepoCommand {
	command := &DeleteRepoCommand{
		singleOperationRepoCommand{
			flagSet: flag.NewFlagSet(deleteRepoCommandName, flag.ContinueOnError),
		},
	}

	command.addRepoNameFlag()

	return command
}

func (command *DeleteRepoCommand) Run() error {
	fmt.Println(fmt.Sprintf("Deleting Repo %s", command.repo))

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
