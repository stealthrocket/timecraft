package timecraft

import (
	"context"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
)

const defaultFunction = "_start"

func runModule(ctx context.Context, runtime wazero.Runtime, compiledModule wazero.CompiledModule, function string) error {
	if function == "" {
		function = defaultFunction
	}

	module, err := runtime.InstantiateModule(ctx, compiledModule, wazero.NewModuleConfig().
		WithStartFunctions())
	if err != nil {
		return err
	}
	defer module.Close(ctx)

	fn := module.ExportedFunction(function)
	if fn == nil {
		return fmt.Errorf("function %q not found in guest", function)
	}
	_, err = fn.Call(ctx)
	switch err {
	case context.Canceled, context.DeadlineExceeded:
		err = nil
	}

	switch e := err.(type) {
	case *sys.ExitError:
		switch exitCode := e.ExitCode(); exitCode {
		case 0:
			err = nil
		default:
			err = ExitError(exitCode)
		}
	}

	return err
}

// ExitError indicates a WebAssembly module exited with a non-zero exit code.
type ExitError uint32

func (e ExitError) Error() string {
	return fmt.Sprintf("module exited with code %d", uint32(e))
}
