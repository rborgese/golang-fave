package migrate

import (
	"context"

	"golang-fave/engine/sqlw"
)

var Migrations = map[string]func(context.Context, *sqlw.DB, string) error{
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
	"000000011": Migrate_000000011,
	"000000012": Migrate_000000012,
	"000000013": Migrate_000000013,
	"000000014": Migrate_000000014,
	"000000015": Migrate_000000015,
	"000000016": Migrate_000000016,
	"000000017": Migrate_000000017,
	"000000018": Migrate_000000018,
	"000000019": Migrate_000000019,
	"000000020": Migrate_000000020,
}
