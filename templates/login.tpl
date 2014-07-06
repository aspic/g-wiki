{{ template "header" . }}
<div class="row col-md-9">
This page is password protected.
<form method="POST">
  {{if .Config}}
  <div class="form-group">
   <input type="text" class="form-control" name="password" placeholder="password"/>
  </div>
  {{end}}
  <div class="form-group">
   <button type="submit" class="btn btn-default">
    <span class="glyphicon glyphicon-ok-circle"></span> Enter
   </button>
  </div>
 </div>
</form>
</div>
{{ template "footer" . }}
