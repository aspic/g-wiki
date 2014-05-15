{{ template "header" . }}
<div class="row col-md-9">
<form method="POST">
 <div class="form-group col-md-12">
  <textarea type="text" class="form-control" rows="15" placeholder="Insert markdown here" name="content">{{ .Content }}</textarea>
 </div>
 <div class="form-inline col-md-12">
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
</div>
{{ template "footer" . }}
