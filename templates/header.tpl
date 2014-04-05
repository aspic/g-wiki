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
 <div class="row">
  <div class="col-md-1">
   <label>{{ .Path }}</label><br />
   <a href="/">Home</a>
  </div>
 <div class="col-md-8">
{{end}}
