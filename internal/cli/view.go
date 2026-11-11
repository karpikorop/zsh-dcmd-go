package cli

import (
	"maps"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/template"
)

// projectView is the JSON shape of a project. Slices and maps are never nil, so they encode as [] and {}.
type projectView struct {
	Name     string            `json:"name"`
	Global   bool              `json:"global"`
	Paths    []string          `json:"paths"`
	Vars     map[string]string `json:"vars"`
	Commands []commandView     `json:"commands"`
}

type commandView struct {
	ID          int            `json:"id"`
	Description string         `json:"description"`
	Run         string         `json:"run"`
	Confirm     bool           `json:"confirm"`
	Variables   []variableView `json:"variables"`
}

// variableView tells an agent which values it must pass. Required means there is no project value and no default.
type variableView struct {
	Name     string  `json:"name"`
	Default  *string `json:"default,omitempty"`
	Required bool    `json:"required"`
}

func newProjectView(p *store.Project) projectView {
	view := projectView{
		Name:     p.Name,
		Global:   p.IsGlobal(),
		Paths:    append([]string{}, p.Paths...),
		Vars:     make(map[string]string, len(p.Vars)),
		Commands: make([]commandView, 0, len(p.Commands)),
	}
	maps.Copy(view.Vars, p.Vars)

	for _, c := range p.Commands {
		view.Commands = append(view.Commands, newCommandView(c, p.Vars))
	}
	return view
}

func newCommandView(c store.Command, projectVars map[string]string) commandView {
	variables := template.Variables(c.Run)
	view := commandView{
		ID:          c.ID,
		Description: c.Description,
		Run:         c.Run,
		Confirm:     c.Confirm,
		Variables:   make([]variableView, 0, len(variables)),
	}

	for _, v := range variables {
		view.Variables = append(view.Variables, newVariableView(v, projectVars))
	}
	return view
}

func newVariableView(v template.Variable, projectVars map[string]string) variableView {
	_, hasProjectValue := projectVars[v.Name]
	view := variableView{Name: v.Name, Required: !hasProjectValue && !v.HasDefault}

	if v.HasDefault {
		view.Default = &v.Default
	}
	return view
}
