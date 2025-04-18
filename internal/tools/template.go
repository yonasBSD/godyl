package tools

import (
	"fmt"
	"maps"

	"github.com/idelchi/godyl/internal/templates"
	"github.com/idelchi/godyl/internal/tools/sources"
)

func (t *Tool) ToTemplateMap(flatten ...map[string]any) map[string]any {
	templateMap := map[string]any{
		"Name":    t.Name,
		"Env":     t.Env,
		"Values":  t.Values,
		"Version": t.Version.Version,
		"Exe":     t.Exe.Name,
		"Output":  t.Output,
	}

	for _, o := range flatten {
		maps.Copy(templateMap, o)
	}

	return templateMap
}

func TemplateError(err error, name string) error {
	return fmt.Errorf("applying template to %q: %w", name, err)
}

// Template applies templating to various fields of the Tool struct, such as version, path, and checksum.
// It processes these fields using Go templates and updates them with the templated values.
func (t *Tool) TemplateFirst() error {
	values := t.ToTemplateMap(t.Platform.ToMap())

	// Apply templating to Source.Type
	output, err := templates.Apply(t.Source.Type.String(), values)
	if err != nil {
		return TemplateError(err, "source.type")
	}

	t.Source.Type = sources.Type(output)

	// Apply templating to the Skip conditions
	for i := range t.Skip {
		err = templates.ApplyAndSet(&t.Skip[i].Condition, values)
		if err != nil {
			return TemplateError(err, "skip.condition")
		}
	}

	if err := templates.ApplyAndSet(&t.Output, values); err != nil {
		return TemplateError(err, "output")
	}

	values = t.ToTemplateMap(t.Platform.ToMap())

	return nil
}

func (t *Tool) TemplateLast() error {
	values := t.ToTemplateMap(t.Platform.ToMap())

	if err := templates.ApplyAndSet(&t.Exe.Name, values); err != nil {
		return err
	}

	// Apply templating to Exe.Patterns
	for i := range t.Exe.Patterns {
		if err := templates.ApplyAndSet(&t.Exe.Patterns[i], values); err != nil {
			return err
		}
	}

	// Apply templating to Extensions
	for i := range t.Extensions {
		if err := templates.ApplyAndSet(&t.Extensions[i], values); err != nil {
			return err
		}
	}

	// Apply templating to Post commands
	for i, cmd := range t.Commands.Commands {
		output, err := templates.Apply(cmd.String(), values)
		if err != nil {
			return err
		}

		t.Commands.Commands[i].From(output)
	}

	// Apply templating to Hints patterns and weights
	for i := range t.Hints {
		if err := templates.ApplyAndSet(&t.Hints[i].Pattern, values); err != nil {
			return err
		}

		if err := templates.ApplyAndSet(&t.Hints[i].Weight, values); err != nil {
			return err
		}
	}

	if err := templates.ApplyAndSet(&t.Path, values); err != nil {
		return err
	}

	// Template the .Name last
	if err := templates.ApplyAndSet(&t.Name, values); err != nil {
		return err
	}

	return nil
}
