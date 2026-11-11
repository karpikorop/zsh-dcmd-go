package cli

import (
	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/template"
)

// projectView is the JSON shape of a project. Slices are never nil, so they encode as [].
type projectView struct {
	Name     string        `json:"name"`
	Global   bool          `json:"global"`
	Paths    []string      `json:"paths"`
	Commands []commandView `json:"commands"`
}

type commandView struct {
	ID          int            `json:"id"`
	Description string         `json:"description"`
	Run         string         `json:"run"`
	Variables   []variableView `json:"variables"`
}

type variableView struct {
	Name string `json:"name"`
}

func newProjectView(p *store.Project) projectView {
	view := projectView{
		Name:     p.Name,
		Global:   p.IsGlobal(),
		Paths:    append([]string{}, p.Paths...),
		Commands: make([]commandView, 0, len(p.Commands)),
	}

	for _, c := range p.Commands {
		view.Commands = append(view.Commands, newCommandView(c))
	}
	return view
}

func newCommandView(c store.Command) commandView {
	variables := template.Variables(c.Run)
	view := commandView{
		ID:          c.ID,
		Description: c.Description,
		Run:         c.Run,
		Variables:   make([]variableView, 0, len(variables)),
	}

	for _, v := range variables {
		view.Variables = append(view.Variables, variableView{Name: v.Name})
	}
	return view
}
