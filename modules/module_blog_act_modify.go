package modules

import (
	"errors"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_BlogModify() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "blog-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_name := wrap.R.FormValue("name")
		pf_alias := wrap.R.FormValue("alias")
		pf_category := wrap.R.FormValue("category")
		pf_briefly := wrap.R.FormValue("briefly")
		pf_content := wrap.R.FormValue("content")
		pf_active := wrap.R.FormValue("active")

		if pf_active == "" {
			pf_active = "0"
		}

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if !utils.IsNumeric(pf_category) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_name == "" {
			wrap.MsgError(`Please specify post name`)
			return
		}

		if pf_alias == "" {
			pf_alias = utils.GenerateSingleAlias(pf_name)
		}

		if !utils.IsValidSingleAlias(pf_alias) {
			wrap.MsgError(`Please specify correct post alias`)
			return
		}

		// Default is ROOT
		if pf_category == "0" {
			pf_category = "1"
		}

		if pf_id == "0" {
			var lastID int64 = 0
			if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
				// Insert row
				res, err := tx.Exec(
					`INSERT INTO blog_posts SET
						user = ?,
						name = ?,
						alias = ?,
						category = ?,
						briefly = ?,
						content = ?,
						datetime = ?,
						active = ?
					;`,
					wrap.User.A_id,
					pf_name,
					pf_alias,
					utils.StrToInt(pf_category),
					pf_briefly,
					pf_content,
					utils.UnixTimestampToMySqlDateTime(utils.GetCurrentUnixTimestamp()),
					utils.StrToInt(pf_active),
				)
				if err != nil {
					return err
				}

				// Get inserted post id
				lastID, err = res.LastInsertId()
				if err != nil {
					return err
				}

				// Block rows
				if _, err := tx.Exec("SELECT id FROM blog_posts WHERE id = ? FOR UPDATE;", lastID); err != nil {
					return err
				}

				// Insert post and categories relations
				catids := utils.GetPostArrayInt("cats[]", wrap.R)
				if len(catids) > 0 {
					var catsCount int
					err = tx.QueryRow(`
						SELECT
							COUNT(*)
						FROM
							blog_cats
						WHERE
							id IN(` + strings.Join(utils.ArrayOfIntToArrayOfString(catids), ",") + `)
						FOR UPDATE;`,
					).Scan(
						&catsCount,
					)
					if wrap.LogCpError(err) != nil {
						return err
					}
					if len(catids) != catsCount {
						return errors.New("Inner system error")
					}
					var balkInsertArr []string
					for _, el := range catids {
						balkInsertArr = append(balkInsertArr, `(`+utils.Int64ToStr(lastID)+`,`+utils.IntToStr(el)+`)`)
					}
					if _, err = tx.Exec(
						`INSERT INTO blog_cat_post_rel (post_id,category_id) VALUES ` + strings.Join(balkInsertArr, ",") + `;`,
					); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.ResetCacheBlocks()
			wrap.Write(`window.location='/cp/blog/modify/` + utils.Int64ToStr(lastID) + `/';`)
		} else {
			if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
				// Block rows
				if _, err := tx.Exec("SELECT id FROM blog_posts WHERE id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
					return err
				}
				if _, err := tx.Exec("SELECT post_id FROM blog_cat_post_rel WHERE post_id = ? FOR UPDATE;", utils.StrToInt(pf_id)); err != nil {
					return err
				}

				// Update row
				if _, err := tx.Exec(
					`UPDATE blog_posts SET
						name = ?,
						alias = ?,
						category = ?,
						briefly = ?,
						content = ?,
						active = ?
					WHERE
						id = ?
					;`,
					pf_name,
					pf_alias,
					utils.StrToInt(pf_category),
					pf_briefly,
					pf_content,
					utils.StrToInt(pf_active),
					utils.StrToInt(pf_id),
				); err != nil {
					return err
				}

				// Delete post and categories relations
				if _, err := tx.Exec("DELETE FROM blog_cat_post_rel WHERE post_id = ?;", utils.StrToInt(pf_id)); err != nil {
					return err
				}

				// Insert post and categories relations
				catids := utils.GetPostArrayInt("cats[]", wrap.R)
				if len(catids) > 0 {
					var catsCount int
					err := tx.QueryRow(`
						SELECT
							COUNT(*)
						FROM
							blog_cats
						WHERE
							id IN(` + strings.Join(utils.ArrayOfIntToArrayOfString(catids), ",") + `)
						FOR UPDATE;`,
					).Scan(
						&catsCount,
					)
					if wrap.LogCpError(err) != nil {
						return err
					}
					if len(catids) != catsCount {
						return errors.New("Inner system error")
					}
					var balkInsertArr []string
					for _, el := range catids {
						balkInsertArr = append(balkInsertArr, `(`+pf_id+`,`+utils.IntToStr(el)+`)`)
					}
					if _, err := tx.Exec(
						`INSERT INTO blog_cat_post_rel (post_id,category_id) VALUES ` + strings.Join(balkInsertArr, ",") + `;`,
					); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.ResetCacheBlocks()
			wrap.Write(`window.location='/cp/blog/modify/` + pf_id + `/';`)
		}
	})
}
