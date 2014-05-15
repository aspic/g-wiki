{{define "header"}}
<!doctype html>
<head>
 <meta charset="UTF-8">
 <title>go-wiki!</title>
 <link href="//netdna.bootstrapcdn.com/bootstrap/3.0.0/css/bootstrap.min.css" rel="stylesheet">
 <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
<div class="container">
 <div class="row col-md-9">
   <ol class="breadcrumb">
    {{range $dir := .Dirs }}
     {{if $dir.Active }}
      <li class="active">{{$dir.Name}}</li>
     {{ else }}
      <li><a href="{{ $dir.Path }}">{{$dir.Name}}</a></li>
     {{ end }}
    {{ end }}
   </ol>
   {{ if .Revision}}<p class="text-muted">Revision: {{.Revision}}</p>{{end}}
 </div>
{{end}}
