package modules

import (
	"html"
	"io/ioutil"
	"net/http"
	"strings"

	"golang-fave/engine/assets"
	"golang-fave/engine/builder"
	"golang-fave/engine/consts"
	"golang-fave/engine/fetdata"
	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) index_TemplateNameToValue(filename string) string {
	if i := strings.LastIndex(filename, "."); i > -1 {
		return filename[:i]
	}
	return filename
}

func (this *Modules) index_GetTemplateSelectOptions(wrap *wrapper.Wrapper, template string) string {
	result := ``

	// index.html
	result += `<option title="index.html" value="index"`
	if template == "index" {
		result += ` selected`
	}
	result += `>index.html</option>`

	// page.html
	result += `<option title="page.html" value="page"`
	if template == "" || template == "page" {
		result += ` selected`
	}
	result += `>page.html</option>`

	// User templates
	if files, err := ioutil.ReadDir(wrap.DTemplate); err == nil {
		for _, file := range files {
			if len(file.Name()) > 0 && file.Name()[0] == '.' {
				continue
			}
			if len(file.Name()) > 0 && strings.ToLower(file.Name()) == "robots.txt" {
				continue
			}
			if !wrap.IsSystemMountedTemplateFile(file.Name()) {
				value := this.index_TemplateNameToValue(file.Name())
				result += `<option title="` + file.Name() + `" value="` + value + `"`
				if template == value {
					result += ` selected`
				}
				result += `>` + file.Name() + `</option>`
			}
		}
	}

	return result
}

func (this *Modules) RegisterModule_Index() *Module {
	return this.newModule(MInfo{
		Mount: "index",
		Name:  "Pages",
		Order: 0,
		Icon:  assets.SysSvgIconPage,
		Sub: &[]MISub{
			{Mount: "default", Name: "List of pages", Show: true, Icon: assets.SysSvgIconList},
			{Mount: "add", Name: "Add new page", Show: true, Icon: assets.SysSvgIconPlus},
			{Mount: "modify", Name: "Modify page", Show: false},
		},
	}, func(wrap *wrapper.Wrapper) {
		// Front-end
		row := &utils.MySql_page{}
		rou := &utils.MySql_user{}
		err := wrap.DB.QueryRow(
			wrap.R.Context(),
			`SELECT
				fave_pages.id,
				fave_pages.user,
				fave_pages.template,
				fave_pages.name,
				fave_pages.alias,
				fave_pages.content,
				fave_pages.meta_title,
				fave_pages.meta_keywords,
				fave_pages.meta_description,
				UNIX_TIMESTAMP(fave_pages.datetime) as datetime,
				fave_pages.active,
				fave_users.id,
				fave_users.first_name,
				fave_users.last_name,
				fave_users.email,
				fave_users.admin,
				fave_users.active
			FROM
				fave_pages
				LEFT JOIN fave_users ON fave_users.id = fave_pages.user
			WHERE
				fave_pages.active = 1 and
				fave_pages.alias = ?
			LIMIT 1;`,
			wrap.R.URL.Path,
		).Scan(
			&row.A_id,
			&row.A_user,
			&row.A_template,
			&row.A_name,
			&row.A_alias,
			&row.A_content,
			&row.A_meta_title,
			&row.A_meta_keywords,
			&row.A_meta_description,
			&row.A_datetime,
			&row.A_active,
			&rou.A_id,
			&rou.A_first_name,
			&rou.A_last_name,
			&rou.A_email,
			&rou.A_admin,
			&rou.A_active,
		)
		if err != nil && err != wrapper.ErrNoRows {
			// System error 500
			wrap.LogCpError(&err)
			utils.SystemErrorPageEngine(wrap.W, err)
			return
		} else if err == wrapper.ErrNoRows {
			// User error 404 page
			wrap.RenderFrontEnd("404", fetdata.New(wrap, true, nil, nil), http.StatusNotFound)
			return
		}

		// Render template
		wrap.RenderFrontEnd(row.A_template, fetdata.New(wrap, false, row, rou), http.StatusOK)
	}, func(wrap *wrapper.Wrapper) (string, string, string) {
		content := ""
		sidebar := ""
		if wrap.CurrSubModule == "" || wrap.CurrSubModule == "default" {
			content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
				{Name: "List of pages"},
			})
			content += builder.DataTable(
				wrap,
				"fave_pages",
				"id",
				"DESC",
				&[]builder.DataTableRow{
					{
						DBField: "id",
					},
					{
						DBField: "template",
					},
					{
						DBField:     "name",
						NameInTable: "Page / URL",
						CallBack: func(values *[]string) string {
							name := `<a href="/cp/` + wrap.CurrModule + `/modify/` + (*values)[0] + `/">` + html.EscapeString((*values)[2]) + `</a>`
							alias := html.EscapeString((*values)[3])
							template := html.EscapeString((*values)[1]) + ".html"
							return `<div>` + name + `</div><div class="template"><small>` + template + `</small></div><div><small>` + alias + `</small></div>`
						},
					},
					{
						DBField: "alias",
					},
					{
						DBField:     "datetime",
						DBExp:       "UNIX_TIMESTAMP(`datetime`)",
						NameInTable: "Date / Time",
						Classes:     "d-none d-md-table-cell",
						CallBack: func(values *[]string) string {
							t := int64(utils.StrToInt((*values)[4]))
							return `<div>` + utils.UnixTimestampToFormat(t, "02.01.2006") + `</div>` +
								`<div><small>` + utils.UnixTimestampToFormat(t, "15:04:05") + `</small></div>`
						},
					},
					{
						DBField:     "active",
						NameInTable: "Active",
						Classes:     "d-none d-sm-table-cell",
						CallBack: func(values *[]string) string {
							return builder.CheckBox(utils.StrToInt((*values)[5]))
						},
					},
				},
				func(values *[]string) string {
					return builder.DataTableAction(&[]builder.DataTableActionRow{
						{
							Icon:   assets.SysSvgIconView,
							Href:   (*values)[3],
							Hint:   "View",
							Target: "_blank",
						},
						{
							Icon: assets.SysSvgIconEdit,
							Href: "/cp/" + wrap.CurrModule + "/modify/" + (*values)[0] + "/",
							Hint: "Edit",
						},
						{
							Icon: assets.SysSvgIconRemove,
							Href: "javascript:fave.ActionDataTableDelete(this,'index-delete','" +
								(*values)[0] + "','Are you sure want to delete page?');",
							Hint:    "Delete",
							Classes: "delete",
						},
					})
				},
				"/cp/"+wrap.CurrModule+"/",
				nil,
				nil,
				true,
			)
		} else if wrap.CurrSubModule == "add" || wrap.CurrSubModule == "modify" {
			if wrap.CurrSubModule == "add" {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Add new page"},
				})
			} else {
				content += this.getBreadCrumbs(wrap, &[]consts.BreadCrumb{
					{Name: "Modify page"},
				})
			}

			data := utils.MySql_page{
				A_id:               0,
				A_user:             0,
				A_template:         "",
				A_name:             "",
				A_alias:            "",
				A_content:          "",
				A_meta_title:       "",
				A_meta_keywords:    "",
				A_meta_description: "",
				A_datetime:         0,
				A_active:           0,
			}

			if wrap.CurrSubModule == "modify" {
				if len(wrap.UrlArgs) != 3 {
					return "", "", ""
				}
				if !utils.IsNumeric(wrap.UrlArgs[2]) {
					return "", "", ""
				}
				err := wrap.DB.QueryRow(
					wrap.R.Context(),
					`SELECT
						id,
						user,
						template,
						name,
						alias,
						content,
						meta_title,
						meta_keywords,
						meta_description,
						active
					FROM
						fave_pages
					WHERE
						id = ?
					LIMIT 1;`,
					utils.StrToInt(wrap.UrlArgs[2]),
				).Scan(
					&data.A_id,
					&data.A_user,
					&data.A_template,
					&data.A_name,
					&data.A_alias,
					&data.A_content,
					&data.A_meta_title,
					&data.A_meta_keywords,
					&data.A_meta_description,
					&data.A_active,
				)
				if *wrap.LogCpError(&err) != nil {
					return "", "", ""
				}
			}

			btn_caption := "Add"
			if wrap.CurrSubModule == "modify" {
				btn_caption = "Save"
			}

			content += builder.DataForm(wrap, []builder.DataFormField{
				{
					Kind:  builder.DFKHidden,
					Name:  "action",
					Value: "index-modify",
				},
				{
					Kind:  builder.DFKHidden,
					Name:  "id",
					Value: utils.IntToStr(data.A_id),
				},
				{
					Kind:     builder.DFKText,
					Caption:  "Page name",
					Name:     "name",
					Value:    data.A_name,
					Required: true,
					Min:      "1",
					Max:      "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Page alias",
					Name:    "alias",
					Value:   data.A_alias,
					Hint:    "Example: /about-us/ or /about-us.html",
					Max:     "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Page template",
					Name:    "template",
					Value:   "0",
					CallBack: func(field *builder.DataFormField) string {
						return `<div class="form-group n2">` +
							`<div class="row">` +
							`<div class="col-md-3">` +
							`<label for="lbl_template">Page template</label>` +
							`</div>` +
							`<div class="col-md-9">` +
							`<div>` +
							`<select class="selectpicker form-control" id="lbl_template" name="template" data-live-search="true">` +
							this.index_GetTemplateSelectOptions(wrap, data.A_template) +
							`</select>` +
							`</div>` +
							`</div>` +
							`</div>` +
							`</div>`
					},
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Page content",
					Name:    "content",
					Value:   data.A_content,
					Classes: "wysiwyg",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Meta title",
					Name:    "meta_title",
					Value:   data.A_meta_title,
					Max:     "255",
				},
				{
					Kind:    builder.DFKText,
					Caption: "Meta keywords",
					Name:    "meta_keywords",
					Value:   data.A_meta_keywords,
					Max:     "255",
				},
				{
					Kind:    builder.DFKTextArea,
					Caption: "Meta description",
					Name:    "meta_description",
					Value:   data.A_meta_description,
					Max:     "510",
				},
				{
					Kind:    builder.DFKCheckBox,
					Caption: "Active",
					Name:    "active",
					Value:   utils.IntToStr(data.A_active),
				},
				{
					Kind:   builder.DFKSubmit,
					Value:  btn_caption,
					Target: "add-edit-button",
				},
			})

			if wrap.CurrSubModule == "add" {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Add</button>`
			} else {
				sidebar += `<button class="btn btn-primary btn-sidebar" id="add-edit-button">Save</button>`
			}
		}
		return this.getSidebarModules(wrap), content, sidebar
	})
}
