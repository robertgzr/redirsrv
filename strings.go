package main

import (
	"bytes"
	"text/template"
)

const (
	index = `
    WHAT IS THIS?
    a really simple redirection service

    WHAT CAN IT DO?
    if your're authorized, query /adm?bucket=redirs for all available routes
	`
	forbidden = `
    UNAUTHORIZED ACCESS!
    request assistance / access at {{.Email}}
	`
	notfound = `
    ROUTE NOT FOUND!
    seems you've tried to access an invalid route

    BUT THIS SHOULD WORK!?
    if you think this is a mistake, contact me at {{.Email}}
	`
	internalerror = `
    INTERNAL ERROR :(
    oops this was not supposed to happen, tell me about it at {{.Email}}
	`
)

func txt(tmpl string, cfg Config) string {
	var buf bytes.Buffer
	t := template.Must(template.New("txt").Parse(tmpl))
	if err := t.Execute(&buf, &cfg); err != nil {
		panic(err)
	}
	return buf.String()
}
