package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/stealthrocket/timecraft/internal/debug"
	"github.com/stealthrocket/timecraft/internal/timemachine"
	"github.com/stealthrocket/timecraft/internal/timemachine/wasicall"
	"github.com/stealthrocket/wasi-go"
	"github.com/stealthrocket/wasi-go/imports"
	"github.com/stealthrocket/wasi-go/imports/wasi_snapshot_preview1"
	"github.com/stealthrocket/wazergo"
)

const replayUsage = `
Usage:	timecraft replay [options] <process id>

Options:
   -c, --config  Path to the timecraft configuration file (overrides TIMECRAFTCONFIG)
   -d, --debug   Start an interactive debugger
   -h, --help    Show this usage information
   -T, --trace   Enable strace-like logging of host function calls
`

func replay(ctx context.Context, args []string) error {
	var (
		debugger = false
		trace    = false
	)

	flagSet := newFlagSet("timecraft replay", replayUsage)
	boolVar(flagSet, &debugger, "d", "debug")
	boolVar(flagSet, &trace, "T", "trace")
	args = parseFlags(flagSet, args)

	if len(args) != 1 {
		return errors.New(`expected exactly one process id as argument`)
	}

	processID, err := uuid.Parse(args[0])
	if err != nil {
		return errors.New(`malformed process id passed as argument (not a UUID)`)
	}
	config, err := loadConfig()
	if err != nil {
		return err
	}
	registry, err := config.openRegistry()
	if err != nil {
		return err
	}

	manifest, err := registry.LookupLogManifest(ctx, processID)
	if err != nil {
		return err
	}
	process, err := registry.LookupProcess(ctx, manifest.Process.Digest)
	if err != nil {
		return err
	}
	processConfig, err := registry.LookupConfig(ctx, process.Config.Digest)
	if err != nil {
		return err
	}
	module, err := registry.LookupModule(ctx, processConfig.Modules[0].Digest)
	if err != nil {
		return err
	}

	logSegment, err := registry.ReadLogSegment(ctx, processID, 0)
	if err != nil {
		return err
	}
	defer logSegment.Close()

	logReader := timemachine.NewLogReader(logSegment, manifest.StartTime)
	defer logReader.Close()

	runtime := config.newRuntime(ctx)
	defer runtime.Close(ctx)

	var debugREPL *debug.REPL
	if debugger {
		debugREPL = debug.NewREPL(os.Stdin, os.Stderr)
		ctx = debug.RegisterFunctionListener(ctx, debugREPL)
	}

	compiledModule, err := runtime.CompileModule(ctx, module.Code)
	if err != nil {
		return err
	}
	defer compiledModule.Close(ctx)

	records := timemachine.NewLogRecordReader(logReader)

	var system wasi.System = wasicall.NewReplay(records)
	defer system.Close(ctx)

	if debugger {
		system = debug.WASIListener(system, debugREPL)
	}

	if trace {
		system = &wasi.Tracer{Writer: os.Stderr, System: system}
	}

	fallback := wasicall.NewObserver(nil, func(ctx context.Context, s wasicall.Syscall) {
		panic(fmt.Sprintf("system call made after log EOF: %s", s.ID()))
	}, nil)
	system = wasicall.NewFallbackSystem(system, fallback)

	hostModule := wasi_snapshot_preview1.NewHostModule(imports.DetectExtensions(compiledModule)...)
	hostModuleInstance := wazergo.MustInstantiate(ctx, runtime, hostModule, wasi_snapshot_preview1.WithWASI(system))
	ctx = wazergo.WithModuleInstance(ctx, hostModuleInstance)

	return instantiate(ctx, runtime, compiledModule, debugREPL)
}