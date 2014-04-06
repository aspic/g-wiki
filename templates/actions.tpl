{{ define "actions" }}
 <form method="POST">
  <div class="form-group">
   <button type="submit" class="btn btn-default btn-xs">
    <span class="glyphicon glyphicon-edit"></span> Edit
   </button>
   <input type="hidden" name="edit" value="true" />
  </div>
 </form>
{{end}}
