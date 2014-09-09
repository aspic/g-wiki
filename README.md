# g-wiki

A _KISS_ wiki built on golang with git as the storage back-end. Content
is formatted in [markdown
syntax](http://daringfireball.net/projects/markdown/syntax). The wiki is
rendered with go templates and [bootstrap](http://getbootstrap.com) css.

Current running example: [mehl.no](http://mehl.no:8081/)

## Install

Simply go get it:
	
	go get github.com/tgulacsi/g-wiki

then run it

	g-wiki -http=:8080 -dir=files 

## Develop

Templates are embedded with [go.rice](https://github.com/GeertJohan/go.rice).
If you change a file under templates, either 

	rm templates.rice-box.go

to force g-wiki to load templates from under `templates` directory; or

	go get github.com/GeertJohan/go.rice/rice
	rice embed-go

to regenerate templates.rice-box.go, and be able to just `go build`,
to have a portable binary.
