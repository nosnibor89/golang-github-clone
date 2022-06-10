package repos

import (
	"flag"
	"fmt"
	"github-clone/src/database/repository"
	cliUtil "github-clone/src/github-clone-cli/internal/config"
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

	if err := repository.Delete(command.repo, cliUtil.GetCLIUser()); err != nil {
		return err
	}

	fmt.Println("repository deleted correctly")
	return nil
}
