{{ template "header" . }}
<form method="POST">
 <div class="form-group">
  <textarea type="text" class="form-control" rows="25" placeholder="Insert markdown here" name="content">{{ .Content }}</textarea>
 </div>
 <div class="form-group">
  <input type="text" class="form-control" name="msg" placeholder="Changelog"/>
  <button type="submit" class="btn btn-default btn-xs">
   <span class="glyphicon glyphicon-floppy-disk"></span> Save
  </button>
 </div>
</form>
{{ template "footer" . }}
