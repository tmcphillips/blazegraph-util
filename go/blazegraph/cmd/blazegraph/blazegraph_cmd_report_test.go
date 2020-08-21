package main

import (
	"strings"
	"testing"

	"github.com/tmcphillips/blazegraph-util/util"
)

func TestBlazegraphCmd_report_static_content(t *testing.T) {

	var outputBuffer strings.Builder
	Main.OutWriter = &outputBuffer
	Main.ErrWriter = &outputBuffer

	t.Run("constant-template", func(t *testing.T) {
		outputBuffer.Reset()
		template := `A constant template.`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			A constant template.
		`)
	})

	t.Run("constant-template-containing-doublequotes", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
			"A constant template"
		`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			"A constant template"
		`)
	})

	t.Run("function-with-quoted-argument", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
			{{up "A constant template"}}
		`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			A CONSTANT TEMPLATE
		`)
	})

	t.Run("function-with-delimited-one-line-argument", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
			{{up '''A constant template'''}}
		`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			A CONSTANT TEMPLATE
		`)
	})

	t.Run("function-with-delimited-one-line-argument-containing-single-quotes", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
			{{up '''A 'constant' template'''}}
		`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			A 'CONSTANT' TEMPLATE
		`)
	})

	t.Run("function-with-delimited-two-line-argument", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
			{{up '''A constant
				template'''}}
		`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			A CONSTANT
			TEMPLATE
		`)
	})

	t.Run("function-with-delimited-two-line-argument-containing-double-quotes", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
			{{up '''A "constant"
				template'''}}
		`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			A "CONSTANT"
			TEMPLATE
		`)
	})

}

func TestBlazegraphCmd_report_two_triples(t *testing.T) {

	var outputBuffer strings.Builder
	Main.OutWriter = &outputBuffer
	Main.ErrWriter = &outputBuffer

	run("blazegraph drop")

	Main.InReader = strings.NewReader(`
		<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
		<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
	`)
	run("blazegraph import --format ttl")

	t.Run("select-piped-to-tabulate", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
			Example select query with tabular output in report\n
			\n
			{{select '''
					prefix ab: <http://tmcphill.net/tags#>
					SELECT ?s ?o
					WHERE
					{ ?s ab:tag ?o }
					ORDER BY ?s
				''' | tabulate}}
		`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			Example select query with tabular output in report

			s                          | o
			==================================
			http://tmcphill.net/data#x | seven
			http://tmcphill.net/data#y | eight
		`)
	})

	t.Run("select-to-variable-to-tabulate", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
				Example select query with tabular output in report\n
				\n
				{{with $tags := (select '''
						prefix ab: <http://tmcphill.net/tags#>
						SELECT ?s ?o
						WHERE
						{ ?s ab:tag ?o }
						ORDER BY ?s
					''')}}{{ tabulate $tags}}{{end}}
			`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
				Example select query with tabular output in report

				s                          | o
				==================================
				http://tmcphill.net/data#x | seven
				http://tmcphill.net/data#y | eight
			`)
	})

	t.Run("select-to-dot-to-tabulate", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
				Example select query with tabular output in report\n
				\n
				{{with (select '''
						prefix ab: <http://tmcphill.net/tags#>
						SELECT ?s ?o
						WHERE
						{ ?s ab:tag ?o }
						ORDER BY ?s
					''')}} {{tabulate .}} {{end}}
			`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
				Example select query with tabular output in report

				s                          | o
				==================================
				http://tmcphill.net/data#x | seven
				http://tmcphill.net/data#y | eight
			`)
	})

	t.Run("select-to-variable-to-range", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
				Example select query with tabular output in report\n
				\n
				{{ with (select '''
						prefix ab: <http://tmcphill.net/tags#>
						SELECT ?s ?o
						WHERE
						{ ?s ab:tag ?o }
						ORDER BY ?s
					''') }}

					Variables:\n
					{{join (.Head.Vars) ", "}}\n
					\n
					Values:\n
					{{range (rows .)}}{{ join . ", " | println}}{{end}}

				{{end}}
			`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
				Example select query with tabular output in report

				Variables:
				s, o

				Values:
				http://tmcphill.net/data#x, seven
				http://tmcphill.net/data#y, eight
			`)
	})

}

func TestBlazegraphCmd_report_multiple_queries(t *testing.T) {

	var outputBuffer strings.Builder
	Main.OutWriter = &outputBuffer
	Main.ErrWriter = &outputBuffer

	run("blazegraph drop")

	Main.InReader = strings.NewReader(`
		<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
		<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
	`)
	run("blazegraph import --format ttl")

	outputBuffer.Reset()
	template := `
		{{prefix "ab" "http://tmcphill.net/tags#"}}

		{{with $subjects := (select '''

				SELECT ?s
				WHERE
				{ ?s ab:tag ?o }
				ORDER BY ?s

			''') | vector }}

			{{range $subject := $subjects }}
				{{with $objects := (select '''

						SELECT ?o
						WHERE
						{ <{{.}}> ab:tag ?o }
						ORDER BY ?o

					''' $subject)}}
					{{tabulate $objects}}
					\n
				{{end}}
			{{end}}

		{{end }}

`
	Main.InReader = strings.NewReader(template)
	run("blazegraph report")
	util.LineContentsEqual(t, outputBuffer.String(), `
		o
		====
		seven

		o
		====
		eight
	`)
}

func TestBlazegraphCmd_report_macros(t *testing.T) {

	var outputBuffer strings.Builder
	Main.OutWriter = &outputBuffer
	Main.ErrWriter = &outputBuffer

	run("blazegraph drop")

	Main.InReader = strings.NewReader(`
		<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
		<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
	`)
	run("blazegraph import --format ttl")

	outputBuffer.Reset()
	template := `

		{{prefix "ab" "http://tmcphill.net/tags#"}}

		{{macro "M1" '''{{select <?
			SELECT ?o 
			WHERE { <{{.}}> ab:tag ?o } 
			ORDER BY ?o 
		?> . | tabulate }}''' }}

		{{with $subjects := (select '''

				SELECT ?s
				WHERE
				{ ?s ab:tag ?o }
				ORDER BY ?s

			''') | vector }}

			{{range $subject := $subjects }}
				{{ expand "M1" $subject }}
			{{end}}

		{{end}}

`
	Main.InReader = strings.NewReader(template)
	run("blazegraph report")
	util.LineContentsEqual(t, outputBuffer.String(), `
		o
		====
		seven
		
		o
		====
		eight
	`)
}

func TestBlazegraphCmd_report_subqueries(t *testing.T) {

	var outputBuffer strings.Builder
	Main.OutWriter = &outputBuffer
	Main.ErrWriter = &outputBuffer

	run("blazegraph drop")

	Main.InReader = strings.NewReader(`
		<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
		<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
	`)
	run("blazegraph import --format ttl")

	outputBuffer.Reset()
	template := `

		{{ prefix "ab" "http://tmcphill.net/tags#" }}

		{{ query "Q1" '''
			SELECT ?s
			WHERE
			{ ?s ab:tag ?o }
			ORDER BY ?s
		''' }}

		{{ query "Q2" '''
			SELECT ?o 
			WHERE { <{{.}}> ab:tag ?o } 
			ORDER BY ?o 
		''' }}

		{{ range (runquery "Q1" | vector) }}
			{{ runquery "Q2" . | tabulate }} \n
		{{ end }}
	`
	Main.InReader = strings.NewReader(template)
	run("blazegraph report")
	util.LineContentsEqual(t, outputBuffer.String(), `
		o
		====
		seven
		
		o
		====
		eight
	`)
}

func TestBlazegraphCmd_report_address_book(t *testing.T) {

	var outputBuffer strings.Builder
	Main.OutWriter = &outputBuffer
	Main.ErrWriter = &outputBuffer

	run("blazegraph drop")
	run("blazegraph import --format jsonld --file testdata/address-book.jsonld")

	t.Run("constant-template", func(t *testing.T) {
		outputBuffer.Reset()
		template := `
			Craig's email addresses\n
			=======================\n
			{{ range (select '''
				PREFIX ab: <http://learningsparql.com/ns/addressbook#>
				SELECT ?email
				WHERE
				{
					?person ab:firstname "Craig"^^<http://www.w3.org/2001/XMLSchema#string> .
					?person ab:email     ?email .
				}
			''' | vector) }}
				{{println .}} 
			{{end}}
		`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			Craig's email addresses
			=======================
			c.ellis@usairwaysgroup.com
			craigellis@yahoo.com
	`)
	})
}

func TestBlazegraphCmd_report_address_book_imports(t *testing.T) {

	var outputBuffer strings.Builder
	Main.OutWriter = &outputBuffer
	Main.ErrWriter = &outputBuffer

	run("blazegraph drop")
	run("blazegraph import --format jsonld --file testdata/address-book.jsonld")

	t.Run("constant-template", func(t *testing.T) {
		outputBuffer.Reset()
		template := `

			{{ include "testdata/address-rules.gst" }}

			{{ prefix "ab" "http://learningsparql.com/ns/addressbook#" }}

			Craig's email addresses\n
			=======================\n
			{{ range (runquery "GetEmailForFirstName" "Craig" | vector) }}
				{{println .}} 
			{{end}}
		`
		Main.InReader = strings.NewReader(template)
		run("blazegraph report")
		util.LineContentsEqual(t, outputBuffer.String(), `
			Craig's email addresses
			=======================
			c.ellis@usairwaysgroup.com
			craigellis@yahoo.com
	`)
	})
}