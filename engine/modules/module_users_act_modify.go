package modules

import (
	"context"

	"golang-fave/engine/utils"
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterAction_UsersModify() *Action {
	return this.newAction(AInfo{
		Mount:     "users-modify",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := utils.Trim(wrap.R.FormValue("id"))
		pf_first_name := utils.Trim(wrap.R.FormValue("first_name"))
		pf_last_name := utils.Trim(wrap.R.FormValue("last_name"))
		pf_email := utils.Trim(wrap.R.FormValue("email"))
		pf_password := utils.Trim(wrap.R.FormValue("password"))
		pf_admin := utils.Trim(wrap.R.FormValue("admin"))
		pf_active := utils.Trim(wrap.R.FormValue("active"))

		if pf_admin == "" {
			pf_admin = "0"
		}

		if pf_active == "" {
			pf_active = "0"
		}

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_email == "" {
			wrap.MsgError(`Please specify user email`)
			return
		}

		if !utils.IsValidEmail(pf_email) {
			wrap.MsgError(`Please specify correct user email`)
			return
		}

		// First user always super admin
		// Rewrite active and admin status
		if pf_id == "1" {
			pf_admin = "1"
			pf_active = "1"
		}

		if pf_id == "0" {
			// Add new user
			if pf_password == "" {
				wrap.MsgError(`Please specify user password`)
				return
			}

			var lastID int64 = 0
			if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
				res, err := tx.Exec(
					ctx,
					`INSERT INTO fave_users SET
						first_name = ?,
						last_name = ?,
						email = ?,
						password = MD5(?),
						admin = ?,
						active = ?
					;`,
					pf_first_name,
					pf_last_name,
					pf_email,
					pf_password,
					pf_admin,
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
				return nil
			}); err != nil {
				wrap.MsgError(err.Error())
				return
			}
			wrap.ResetCacheBlocks()
			wrap.Write(`window.location='/cp/users/modify/` + utils.Int64ToStr(lastID) + `/';`)
		} else {
			// Update user
			if pf_password == "" {
				if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
					_, err := tx.Exec(
						ctx,
						`UPDATE fave_users SET
							first_name = ?,
							last_name = ?,
							email = ?,
							admin = ?,
							active = ?
						WHERE
							id = ?
						;`,
						pf_first_name,
						pf_last_name,
						pf_email,
						pf_admin,
						utils.StrToInt(pf_active),
						utils.StrToInt(pf_id),
					)
					if err != nil {
						return err
					}
					return nil
				}); err != nil {
					wrap.MsgError(err.Error())
					return
				}
			} else {
				if err := wrap.DB.Transaction(wrap.R.Context(), func(ctx context.Context, tx *wrapper.Tx) error {
					_, err := tx.Exec(
						ctx,
						`UPDATE fave_users SET
							first_name = ?,
							last_name = ?,
							email = ?,
							password = MD5(?)
						WHERE
							id = ?
						;`,
						pf_first_name,
						pf_last_name,
						pf_email,
						pf_password,
						utils.StrToInt(pf_id),
					)
					if err != nil {
						return err
					}
					return nil
				}); err != nil {
					wrap.MsgError(err.Error())
					return
				}
			}
			wrap.ResetCacheBlocks()
			wrap.Write(`window.location='/cp/users/modify/` + pf_id + `/';`)
		}
	})
}
