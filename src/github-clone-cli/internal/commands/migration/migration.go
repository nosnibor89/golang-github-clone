package migration

import (
	"errors"
	"flag"
	"fmt"
	"github-clone/src/github-clone-cli/internal/migrations"
	"github-clone/src/util"
	"strings"
)

const (
	executeMigrationCommand = "migration.execute"
)

var knownMigrations = []Migration{
	migrations.UpdateRepoStarsMigration{},
}

type Migration interface {
	Name() string
	Run() error
}

type ExecuteMigrationCommand struct {
	flagSet *flag.FlagSet
	name    string
}

func NewExecuteMigrationCommand() *ExecuteMigrationCommand {
	command := &ExecuteMigrationCommand{
		flagSet: flag.NewFlagSet(executeMigrationCommand, flag.ContinueOnError),
	}

	command.flagSet.StringVar(&command.name, "name", "", "migration name")

	return command
}

func (c *ExecuteMigrationCommand) Name() string {
	return c.flagSet.Name()
}

func (c *ExecuteMigrationCommand) Init(args []string) error {
	return c.flagSet.Parse(args)
}

func (c *ExecuteMigrationCommand) Run() error {

	fmt.Println(fmt.Sprintf("Running migration %s", c.name))

	if err := c.assertCommandIsValid(); err != nil {
		return err
	}

	if migration := *getMigration(c.name); migration != nil {
		return migration.Run()
	}

	return errors.New("no migration was run")
}

func (c ExecuteMigrationCommand) assertCommandIsValid() error {
	if util.StringIsEmpty(c.name) || !migrationExists(c.name) {
		return errors.New(fmt.Sprintf("invalid migration name provided. name :%s", c.name))
	}

	return nil
}

func migrationExists(name string) bool {
	migration := getMigration(name)
	return migration != nil
}

func getMigration(name string) *Migration {
	for _, migration := range knownMigrations {
		if migration.Name() == strings.TrimSpace(name) {
			return &migration
		}
	}

	return nil
}
