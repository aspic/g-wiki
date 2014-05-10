{{ template "header" . }}
<form method="POST">
 <div class="form-group">
  <textarea type="text" class="form-control" rows="25" placeholder="Insert markdown here" name="content">{{ .Content }}</textarea>
 </div>
 <div class="form-inline">
  <div class="form-group col-md-8">
   <input type="text" class="form-control" name="msg" placeholder="Changelog"/>
  </div>
  <div class="form-group col-md-2">
   <input type="text" class="form-control" name="author" placeholder="Author"/>
  </div>
  <div class="form-group col-md-2">
   <button type="submit" class="btn btn-default">
    <span class="glyphicon glyphicon-floppy-disk"></span> Save
   </button>
  </div>
 </div>
</form>
{{ template "footer" . }}
