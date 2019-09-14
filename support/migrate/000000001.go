package migrate

import (
	"golang-fave/engine/sqlw"
)

var Migrations = map[string]func(*sqlw.DB, string) error{
	"000000000": nil,
	"000000001": nil,
	"000000002": Migrate_000000002,
	"000000003": Migrate_000000003,
	"000000004": Migrate_000000004,
	"000000005": Migrate_000000005,
	"000000006": Migrate_000000006,
	"000000007": Migrate_000000007,
	"000000008": Migrate_000000008,
	"000000009": Migrate_000000009,
	"000000010": Migrate_000000010,
}
