package main

import (
	"errors"
	"fmt"
	commands2 "github-clone/src/github-clone-cli/commands"
	repoCommands "github-clone/src/github-clone-cli/commands/repos"
	cliUtil "github-clone/src/github-clone-cli/config"
	"github-clone/src/util"
	"os"
)

func assertCanUseCLI() error {
	if util.StringIsEmpty(cliUtil.GetCLIUser()) {
		return errors.New("you must provide a user for the cli")
	}

	if util.StringIsEmpty(cliUtil.GetCLITable()) {
		return errors.New("you must provide a dynamo db table for the cli")
	}

	return nil
}

func parse(args []string) error {

	if err := assertCanUseCLI(); err != nil {
		return err
	}

	if len(args) < 1 {
		return errors.New("you must pass a sub-command")
	}

	commands := []commands2.Runner{
		repoCommands.NewDeleteRepoCommand(),
		repoCommands.NewFindOneRepoCommand(),
	}

	subcommand := os.Args[1]

	for _, command := range commands {
		if command.Name() == subcommand {
			if err := command.Init(os.Args[2:]); err != nil {
				return fmt.Errorf("failed to initialize subcommand: %s", subcommand)
			}
			fmt.Printf("executing command for user %s\n", cliUtil.GetCLIUser())
			return command.Run()
		}
	}

	return fmt.Errorf("unknown subcommand: %s", subcommand)
}

func main() {
	if err := parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
