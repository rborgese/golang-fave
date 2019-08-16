package modules

import (
	"os"
	"path/filepath"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

func (this *Modules) RegisterAction_ShopUploadDelete() *Action {
	return this.newAction(AInfo{
		WantDB:    true,
		Mount:     "shop-upload-delete",
		WantAdmin: true,
	}, func(wrap *wrapper.Wrapper) {
		pf_id := wrap.R.FormValue("id")
		pf_file := wrap.R.FormValue("file")

		if !utils.IsNumeric(pf_id) {
			wrap.MsgError(`Inner system error`)
			return
		}

		if pf_file == "" {
			wrap.MsgError(`Inner system error`)
			return
		}

		if err := wrap.DB.Transaction(func(tx *wrapper.Tx) error {
			// Block rows
			if _, err := tx.Exec("SELECT id FROM shop_products WHERE id = ? FOR UPDATE;", pf_id); err != nil {
				return err
			}
			if _, err := tx.Exec("SELECT product_id FROM shop_product_images WHERE product_id = ? FOR UPDATE;", pf_id); err != nil {
				return err
			}

			// Delete row
			if _, err := tx.Exec("DELETE FROM shop_product_images WHERE product_id = ? AND filename = ?;", pf_id, pf_file); err != nil {
				return err
			}

			// Delete file
			target_file_full := wrap.DHtdocs + string(os.PathSeparator) + "products" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + pf_id + string(os.PathSeparator) + pf_file
			os.Remove(target_file_full)

			// Delete thumbnails
			pattern := wrap.DHtdocs + string(os.PathSeparator) + "products" + string(os.PathSeparator) + "images" + string(os.PathSeparator) + pf_id + string(os.PathSeparator) + "thumb-*-" + pf_file
			if files, err := filepath.Glob(pattern); err == nil {
				for _, file := range files {
					if err := os.Remove(file); err != nil {
						wrap.LogError("[upload delete] Thumbnail file (%s) delete error: %s", file, err.Error())
					}
				}
			}

			return nil
		}); err != nil {
			wrap.MsgError(err.Error())
			return
		}

		wrap.Write(`$('#list-images a').each(function(i, e) { if($(e).attr('title') == '` + pf_file + `') { $(e).parent().remove(); return; } });`)
	})
}