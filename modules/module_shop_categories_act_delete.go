package modules

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopCategoriesDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-categories-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")

		if !utils.IsNumeric(pf_id) || utils.StrToInt(pf_id) <= 1 {
			wrap.MsgError(`Inner system error`)
			return
		}

		err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec("SELECT id FROM shop_cats FOR UPDATE;"); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT category_id FROM shop_cat_product_rel WHERE category_id = ? FOR UPDATE;", pf_id); err != nil {
				return err
			}

			// Process
			if _, err := tx.Exec("SELECT @ml := lft, @mr := rgt FROM shop_cats WHERE id = ?;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("DELETE FROM shop_cats WHERE id = ?;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE shop_cats SET lft = lft - 1, rgt = rgt - 1 WHERE lft > @ml AND rgt < @mr;"); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE shop_cats SET lft = lft - 2 WHERE lft > @mr;"); err != nil {
				return err
			}
			if _, err := tx.Exec("UPDATE shop_cats SET rgt = rgt - 2 WHERE rgt > @mr;"); err != nil {
				return err
			}
			if _, err := tx.Exec("DELETE FROM shop_cat_product_rel WHERE category_id = ?;", pf_id); err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			wrap.MsgError(err.Error())
			return
		}

		// Reload current page
		wrap.Write(`window.location.reload(false);`)
	})
}