{{define "footer"}}
 </div>
<div class="col-md-3">
 <p class="text-muted">Revisions:</p>
 <div class="list-group">
  {{range $log := .Log}}
   {{if $log.Link}}
    <a href="?revision={{$log.Hash}}" class="list-group-item">
   {{else}}
    <a href="?revision={{$log.Hash}}" class="list-group-item active">
   {{end}}
    {{$log.Message}} ({{$log.Time}})
   </a>
   </li>
  {{end}}
 </div>
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
