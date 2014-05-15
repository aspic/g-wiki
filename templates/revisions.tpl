{{define "revisions"}}
<div class="row col-md-9">
 <hr />
 <p class="text-muted">Revisions:</p>
 <div class="list-group">
  {{range $log := .Log}}
   {{if $log.Link}}
    <a href="?revision={{$log.Hash}}&revisions=1" class="list-group-item">
   {{else}}
    <a href="?revision={{$log.Hash}}&revisions=1" class="list-group-item active">
   {{end}}
    {{$log.Message}} ({{$log.Time}})
   </a>
   </li>
  {{end}}
 </div>
</div>
{{end}}
