package plugins

import (
	"os"
	"path"

	"github.com/kamontat/fthelper/shared/commandline/flags"
	"github.com/kamontat/fthelper/shared/commandline/hooks"
	"github.com/kamontat/fthelper/shared/dotenv"
	"github.com/kamontat/fthelper/shared/fs"
	"github.com/kamontat/fthelper/shared/maps"
)

// Support is example how we can implement plugins support
func SupportEnv(p *PluginParameter) error {
	var wd, err = os.Getwd()
	if err != nil {
		return err
	}

	p.NewFlags(flags.Array{
		Name:    "envs",
		Default: []string{path.Join(wd, ".env")},
		Usage:   "environment files, must follow .env regulation",
		Action: func(data []string) maps.Mapper {
			return maps.New().
				Set("fs.env.type", "file").
				Set("fs.env.mode", "multiple").
				Set("fs.env.fullpath", data)
		},
	})

	p.NewFlags(flags.Bool{
		Name:    "no-env",
		Default: false,
		Usage:   "disabled loading .env files completely",
		Action: func(data bool) maps.Mapper {
			return maps.New().
				Set("internal.flag.noenv", data)
		},
	})

	p.NewHook(hooks.AFTER_FLAG, func(m maps.Mapper) error {
		if disabled := m.Mi("internal").Mi("flag").Bo("noenv", false); disabled {
			return nil
		}

		envs, err := fs.Build("env", m.Mi("fs"))
		if err != nil {
			p.Logger.Warn("cannot found environment file: %v", err)
			return nil
		}

		// write environment value from .env file
		err = dotenv.Overload(envs.Multiple()...)
		if err != nil {
			p.Logger.Warn("dotenv return error: %v", err)
		}
		return nil
	})

	return nil
}
