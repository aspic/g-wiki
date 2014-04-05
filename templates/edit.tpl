{{ template "header" . }}
<form method="POST">
 <div class="form-group">
  <textarea type="text" class="form-control" rows="25" placeholder="Insert markdown here" name="content">{{ .Content }}</textarea>
 </div>
 <div class="form-group">
  <input type="text" class="form-control" name="msg" placeholder="Changelog"/>
  <input type="submit" class="btn btn-primary btn-sm" value="Save"/>
 </div>
</form>
{{ template "footer" . }}
