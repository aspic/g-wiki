{{ template "header" .}}
<form method="POST">
 <div class="form-group">
 <textarea type="text" class="row-fluid" placeholder="Insert markdown here" name="content">{{ .Content }}</textarea><br/>
 <input type="submit" />
 </div>
</form>
{{ template "footer" .}}
