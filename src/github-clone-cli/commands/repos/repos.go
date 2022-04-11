package repos

import (
	"errors"
	"flag"
	"fmt"
	"github-clone/src/util"
)

type singleOperationRepoCommand struct {
	flagSet *flag.FlagSet

	repo string
}

func (command *singleOperationRepoCommand) Name() string {
	return command.flagSet.Name()
}

func (command *singleOperationRepoCommand) Init(args []string) error {
	return command.flagSet.Parse(args)
}

func (command *singleOperationRepoCommand) assertCommandIsValid() error {
	if util.StringIsEmpty(command.repo) {
		return errors.New(fmt.Sprintf("invalid repo name provided. repo-name :%s", command.repo))
	}

	return nil
}

func (command *singleOperationRepoCommand) addRepoNameFlag() {
	command.flagSet.StringVar(&command.repo, "repo-name", "", "repository name")
}
