package migrate

import (
	"context"
	"io/ioutil"
	"os"

	ThemeFiles "golang-fave/assets/template"
	"golang-fave/engine/sqlw"
)

func Migrate_000000017(ctx context.Context, db *sqlw.DB, host string) error {
	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/email-new-order-admin.html", ThemeFiles.AllData["email-new-order-admin.html"], 0664); err != nil {
		return err
	}

	if err := ioutil.WriteFile(host+string(os.PathSeparator)+"/template/email-new-order-user.html", ThemeFiles.AllData["email-new-order-user.html"], 0664); err != nil {
		return err
	}

	return nil
}
