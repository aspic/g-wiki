{{define "footer"}}
 </div>
<div class="col-md-3">
 <p class="text-muted">Revisions:</p>
 <ol class="list-group" style="margin: 5px;">
  {{range $log := .Log}}
   <li class="list-group-item">
    <a href="?show={{$log.Hash}}" class="btn btn-primary btn-xs">show</a>
    {{$log.Message}} ({{$log.Time}})
   </li>
  {{end}}
 </ol>
</div>
<!-- end row -->
</div>
<div class="row">
 <div class="col-md-9">
  <hr class="text-muted" />
  <a href="https://github.com/aspic/g-wiki"><p class="text-muted text-center">g-wiki on Github</p></a>
</div>
<!-- end container -->
</div>
 </body>
</html>
{{end}}
