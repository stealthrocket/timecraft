package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stealthrocket/timecraft/format"
	"github.com/stealthrocket/timecraft/internal/print/human"
	"github.com/stealthrocket/timecraft/internal/print/jsonprint"
	"github.com/stealthrocket/timecraft/internal/print/textprint"
	"github.com/stealthrocket/timecraft/internal/print/yamlprint"
	"github.com/stealthrocket/timecraft/internal/stream"
	"github.com/stealthrocket/timecraft/internal/timemachine"
)

const getUsage = `
Usage:	timecraft get <resource type> [options]

   The get sub-command gives access to the state of the time machine registry.
   The command must be followed by the name of resources to display, which must
   be one of config, log, module, process, profile, or runtime.
   (the command also accepts plurals and abbreviations of the resource names)

Examples:

   $ timecraft get modules
   MODULE ID     MODULE NAME  SIZE
   9d7b7563baf3  app.wasm     6.82 MiB

   $ timecraft get modules -o json
   {
     "mediaType": "application/vnd.timecraft.module.v1+wasm",
     "digest": "sha256:9d7b7563baf3702cf24ed3688dc9a58faef2d0ac586041cb2dc95df919f5e5f2",
     "size": 7150231,
     "annotations": {
       "timecraft.module.name": "app.wasm"
     }
   }

Options:
   -h, --help           Show this usage information
   -o, --ouptut format  Output format, one of: text, json, yaml
   -r, --registry path  Path to the timecraft registry (default to ~/.timecraft)
`

type resource struct {
	typ       string
	alt       []string
	mediaType format.MediaType
	get       func(context.Context, io.Writer, *timemachine.Registry) stream.WriteCloser[*format.Descriptor]
	describe  func(context.Context, *timemachine.Registry, string) (any, error)
	lookup    func(context.Context, *timemachine.Registry, string) (any, error)
}

var resources = [...]resource{
	{
		typ:       "config",
		alt:       []string{"conf", "configs"},
		mediaType: format.TypeTimecraftConfig,
		get:       getConfigs,
		describe:  describeConfig,
		lookup:    lookupConfig,
	},
	{
		typ:       "log",
		alt:       []string{"logs"},
		mediaType: format.TypeTimecraftManifest,
		describe:  describeLog,
		lookup:    describeLog,
	},
	{
		typ:       "module",
		alt:       []string{"mo", "mod", "mods", "modules"},
		mediaType: format.TypeTimecraftModule,
		get:       getModules,
		describe:  describeModule,
		lookup:    lookupModule,
	},
	{
		typ:       "process",
		alt:       []string{"ps", "proc", "procs", "processes"},
		mediaType: format.TypeTimecraftProcess,
		get:       getProcesses,
		describe:  describeProcess,
		lookup:    lookupProcess,
	},
	{
		typ:       "profile",
		alt:       []string{"prof", "profs", "profiles"},
		mediaType: format.TypeTimecraftProfile,
		get:       getProfiles,
		describe:  describeProfiles,
		lookup:    describeProfiles,
	},
	{
		typ:       "runtime",
		alt:       []string{"rt", "runtimes"},
		mediaType: format.TypeTimecraftRuntime,
		get:       getRuntimes,
		describe:  describeRuntime,
		lookup:    lookupRuntime,
	},
}

func get(ctx context.Context, args []string) error {
	var (
		timeRange    = timemachine.Since(time.Unix(0, 0))
		output       = outputFormat("text")
		registryPath = human.Path("~/.timecraft")
	)

	flagSet := newFlagSet("timecraft get", getUsage)
	customVar(flagSet, &output, "o", "output")
	customVar(flagSet, &registryPath, "r", "registry")
	args = parseFlags(flagSet, args)

	if len(args) != 1 {
		return errors.New(`expected exactly one resource type as argument` + useGet())
	}
	resourceTypeLookup := args[0]
	resource, ok := findResource(resourceTypeLookup, resources[:])
	if !ok {
		matchingResources := findMatchingResources(resourceTypeLookup, resources[:])
		if len(matchingResources) == 0 {
			return fmt.Errorf(`no resources matching '%s'`+useGet(), resourceTypeLookup)
		}
		return fmt.Errorf(`no resources matching '%s'

Did you mean?%s`, resourceTypeLookup, joinResourceTypes(matchingResources, "\n   "))
	}

	registry, err := openRegistry(registryPath)
	if err != nil {
		return err
	}

	// We make a special case for the log segments because they are not
	// immutable objects and therefore don't have descriptors.
	if resource.typ == "log" {
		reader := registry.ListLogManifests(ctx) // TODO: time range
		defer reader.Close()

		var writer stream.WriteCloser[*format.Manifest]
		switch output {
		case "json":
			writer = jsonprint.NewWriter[*format.Manifest](os.Stdout)
		case "yaml":
			writer = yamlprint.NewWriter[*format.Manifest](os.Stdout)
		default:
			writer = getLogs(ctx, os.Stdout, registry)
		}
		defer writer.Close()

		_, err = stream.Copy[*format.Manifest](writer, reader)
		return err
	}

	reader := registry.ListResources(ctx, resource.mediaType, timeRange)
	defer reader.Close()

	var writer stream.WriteCloser[*format.Descriptor]
	switch output {
	case "json":
		writer = jsonprint.NewWriter[*format.Descriptor](os.Stdout)
	case "yaml":
		writer = yamlprint.NewWriter[*format.Descriptor](os.Stdout)
	default:
		writer = resource.get(ctx, os.Stdout, registry)
	}
	defer writer.Close()

	_, err = stream.Copy[*format.Descriptor](writer, reader)
	return err
}

func getConfigs(ctx context.Context, w io.Writer, reg *timemachine.Registry) stream.WriteCloser[*format.Descriptor] {
	type config struct {
		ID      string      `text:"CONFIG ID"`
		Runtime string      `text:"RUNTIME"`
		Modules int         `text:"MODULES"`
		Size    human.Bytes `text:"SIZE"`
	}
	return newTableWriter(w,
		func(c1, c2 config) bool {
			return c1.ID < c2.ID
		},
		func(desc *format.Descriptor) (config, error) {
			c, err := reg.LookupConfig(ctx, desc.Digest)
			if err != nil {
				return config{}, err
			}
			r, err := reg.LookupRuntime(ctx, c.Runtime.Digest)
			if err != nil {
				return config{}, err
			}
			return config{
				ID:      desc.Digest.Short(),
				Runtime: r.Runtime + " (" + r.Version + ")",
				Modules: len(c.Modules),
				Size:    human.Bytes(desc.Size),
			}, nil
		})
}

func getModules(ctx context.Context, w io.Writer, reg *timemachine.Registry) stream.WriteCloser[*format.Descriptor] {
	type module struct {
		ID   string      `text:"MODULE ID"`
		Name string      `text:"MODULE NAME"`
		Size human.Bytes `text:"SIZE"`
	}
	return newTableWriter(w,
		func(m1, m2 module) bool {
			return m1.ID < m2.ID
		},
		func(desc *format.Descriptor) (module, error) {
			name := desc.Annotations["timecraft.module.name"]
			if name == "" {
				name = "(none)"
			}
			return module{
				ID:   desc.Digest.Short(),
				Name: name,
				Size: human.Bytes(desc.Size),
			}, nil
		})
}

func getProcesses(ctx context.Context, w io.Writer, reg *timemachine.Registry) stream.WriteCloser[*format.Descriptor] {
	type process struct {
		ID        format.UUID `text:"PROCESS ID"`
		StartTime human.Time  `text:"START"`
	}
	return newTableWriter(w,
		func(p1, p2 process) bool {
			return time.Time(p1.StartTime).Before(time.Time(p2.StartTime))
		},
		func(desc *format.Descriptor) (process, error) {
			p, err := reg.LookupProcess(ctx, desc.Digest)
			if err != nil {
				return process{}, err
			}
			return process{
				ID:        p.ID,
				StartTime: human.Time(p.StartTime),
			}, nil
		})
}

func getProfiles(ctx context.Context, w io.Writer, reg *timemachine.Registry) stream.WriteCloser[*format.Descriptor] {
	type profile struct {
		ID        string         `text:"PROFILE ID"`
		ProcessID format.UUID    `text:"PROCESS ID"`
		Type      string         `text:"TYPE"`
		StartTime human.Time     `text:"START"`
		Duration  human.Duration `text:"DURATION"`
		Size      human.Bytes    `text:"SIZE"`
	}
	return newTableWriter(w,
		func(p1, p2 profile) bool {
			if p1.ProcessID != p2.ProcessID {
				return bytes.Compare(p1.ProcessID[:], p2.ProcessID[:]) < 0
			}
			if p1.Type != p2.Type {
				return p1.Type < p2.Type
			}
			if !time.Time(p1.StartTime).Equal(time.Time(p2.StartTime)) {
				return time.Time(p1.StartTime).Before(time.Time(p2.StartTime))
			}
			return p1.Duration < p2.Duration
		},
		func(desc *format.Descriptor) (profile, error) {
			processID, _ := uuid.Parse(desc.Annotations["timecraft.process.id"])
			startTime, _ := time.Parse(time.RFC3339Nano, desc.Annotations["timecraft.profile.start"])
			endTime, _ := time.Parse(time.RFC3339Nano, desc.Annotations["timecraft.profile.end"])
			return profile{
				ID:        desc.Digest.Short(),
				ProcessID: processID,
				Type:      desc.Annotations["timecraft.profile.type"],
				StartTime: human.Time(startTime),
				Duration:  human.Duration(endTime.Sub(startTime)),
				Size:      human.Bytes(desc.Size),
			}, nil
		})
}

func getRuntimes(ctx context.Context, w io.Writer, reg *timemachine.Registry) stream.WriteCloser[*format.Descriptor] {
	type runtime struct {
		ID      string `text:"RUNTIME ID"`
		Runtime string `text:"RUNTIME NAME"`
		Version string `text:"VERSION"`
	}
	return newTableWriter(w,
		func(r1, r2 runtime) bool {
			return r1.ID < r2.ID
		},
		func(desc *format.Descriptor) (runtime, error) {
			r, err := reg.LookupRuntime(ctx, desc.Digest)
			if err != nil {
				return runtime{}, err
			}
			return runtime{
				ID:      desc.Digest.Short(),
				Runtime: r.Runtime,
				Version: r.Version,
			}, nil
		})
}

func getLogs(ctx context.Context, w io.Writer, reg *timemachine.Registry) stream.WriteCloser[*format.Manifest] {
	type manifest struct {
		ProcessID format.UUID `text:"PROCESS ID"`
		Segments  human.Count `text:"SEGMENTS"`
		StartTime human.Time  `text:"START"`
		Size      human.Bytes `text:"SIZE"`
	}
	return newTableWriter(w,
		func(m1, m2 manifest) bool {
			return time.Time(m1.StartTime).Before(time.Time(m2.StartTime))
		},
		func(m *format.Manifest) (manifest, error) {
			manifest := manifest{
				ProcessID: m.ProcessID,
				Segments:  human.Count(len(m.Segments)),
				StartTime: human.Time(m.StartTime),
			}
			for _, segment := range m.Segments {
				manifest.Size += human.Bytes(segment.Size)
			}
			return manifest, nil
		})
}

func newTableWriter[T1, T2 any](w io.Writer, orderBy func(T1, T1) bool, conv func(T2) (T1, error)) stream.WriteCloser[T2] {
	tw := textprint.NewTableWriter[T1](w, textprint.OrderBy(orderBy))
	cw := stream.ConvertWriter[T1](tw, conv)
	return stream.NewWriteCloser(cw, tw)
}

func findResource(typ string, options []resource) (resource, bool) {
	for _, option := range options {
		if option.typ == typ {
			return option, true
		}
		for _, alt := range option.alt {
			if alt == typ {
				return option, true
			}
		}
	}
	return resource{}, false
}

func findMatchingResources(typ string, options []resource) (matches []resource) {
	for _, option := range options {
		if prefixLength(option.typ, typ) > 1 || prefixLength(typ, option.typ) > 1 {
			matches = append(matches, option)
		}
	}
	return matches
}

func prefixLength(base, prefix string) int {
	n := 0
	for n < len(base) && n < len(prefix) && base[n] == prefix[n] {
		n++
	}
	return n
}

func joinResourceTypes(resources []resource, prefix string) string {
	s := new(strings.Builder)
	for _, r := range resources {
		s.WriteString(prefix)
		s.WriteString(r.typ)
	}
	return s.String()
}

func useGet() string {
	s := new(strings.Builder)
	s.WriteString("\n\n")
	s.WriteString(`Use 'timecraft get <resource type>' where the supported resource types are:`)
	for _, r := range resources {
		s.WriteString("\n   ")
		s.WriteString(r.typ)
	}
	return s.String()
}
