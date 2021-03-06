package template

var VarBlogPostHtmlFile = []byte(`{{template "header.html" .}}
<div class="card mb-4">
	<div class="card-body">
		{{if $.Data.IsUserLoggedIn}}
			{{if $.Data.CurrentUser.IsAdmin}}
				<a href="/cp/blog/modify/{{$.Data.Blog.Post.Id}}/" target="_blank" style="float:right;">Edit</a>
			{{end}}
		{{end}}
		<h2 class="card-title">{{$.Data.Blog.Post.Name}}</h2>
		<div class="page-content">
			{{$.Data.Blog.Post.Briefly}}
			{{$.Data.Blog.Post.Content}}
		</div>
	</div>
	<div class="card-footer text-muted">
		<div>Published on {{$.Data.Blog.Post.DateTimeFormat "02/01/2006, 15:04:05"}}</div>
		<div>Author: {{$.Data.Blog.Post.User.FirstName}} {{$.Data.Blog.Post.User.LastName}}</div>
	</div>
</div>
{{template "footer.html" .}}`)
