{{ template "header" . }}
<form method="POST">
 <div class="form-group">
 <textarea type="text" class="form-control" rows="25" placeholder="Insert markdown here" name="content">{{ .Content }}</textarea><br/>
 <input type="submit" class="btn btn-primary btn-sm" value="Save"/>
 </div>
</form>
{{ template "footer" . }}
