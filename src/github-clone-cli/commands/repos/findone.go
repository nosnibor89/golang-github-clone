package repos

import (
	"flag"
	"fmt"
	cliUtil "github-clone/src/github-clone-cli/util"
	"github-clone/src/repositories"
	util2 "github-clone/src/util"
	"github.com/aws/aws-sdk-go/aws/awsutil"
)

const (
	findOneCommandName = "repo.findone"
)

type FindOneRepoCommand struct {
	singleOperationRepoCommand
}

func NewFindOneRepoCommand() *FindOneRepoCommand {
	command := &FindOneRepoCommand{
		singleOperationRepoCommand{
			flagSet: flag.NewFlagSet(findOneCommandName, flag.ContinueOnError),
		},
	}

	command.addRepoNameFlag()

	return command
}

func (command *FindOneRepoCommand) Run() error {
	fmt.Printf("Finding Repo %s\n", command.repo)

	if err := command.assertCommandIsValid(); err != nil {
		return err
	}

	repository := repositories.GithubRepository{}
	found := repository.FindOne(command.repo, cliUtil.GetCLIUser())

	if util2.IsModelEmpty(found) {
		util2.PrintRed(fmt.Sprintf("No repo found with name %v", command.repo))
	} else {
		util2.PrintCyan(awsutil.Prettify(found))
	}

	return nil
}
