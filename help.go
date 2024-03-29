package main

import (
	"context"
	"fmt"
	"strings"
)

const helpUsage = `
Usage:	timecraft <command> [options]

Registry Commands:
   describe  Show detailed information about specific resources
   export    Export resources to local files
   get       Display resources from the time machine registry

Runtime Commands:
   run       Run a WebAssembly module, and optionally trace execution
   replay    Replay a recorded trace of execution

Debugging Commands:
   logs      Print the logs for a module execution
   profile   Generate performance profile from execution records
   trace     Generate traces from execution records

Other Commands:
   config    View or edit the timecraft configuration
   help      Show usage information about timecraft commands
   version   Show the timecraft version information

Global Options:
   -c, --config path  Path to the timecraft configuration file (overrides TIMECRAFTCONFIG)
   -h, --help         Show usage information

For a description of each command, run 'timecraft help <command>'.`

func help(ctx context.Context, args []string) error {
	flagSet := newFlagSet("timecraft help", helpUsage)

	args, err := parseFlags(flagSet, args)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		args = []string{"help"}
	}

	for i, cmd := range args {
		var msg string

		if i != 0 {
			fmt.Println("---")
		}

		switch cmd {
		case "config":
			msg = configUsage
		case "describe":
			msg = describeUsage
		case "export":
			msg = exportUsage
		case "get":
			msg = getUsage
		case "help":
			msg = helpUsage
		case "logs":
			msg = logsUsage
		case "profile":
			msg = profileUsage
		case "run":
			msg = runUsage
		case "replay":
			msg = replayUsage
		case "trace":
			msg = traceUsage
		case "version":
			msg = versionUsage
		default:
			perrorf("timecraft help %s: unknown command", cmd)
			return exitCode(2)
		}

		fmt.Println(strings.TrimSpace(msg))
	}
	return nil
}
