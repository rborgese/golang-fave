package fetdata

import (
	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type ShopCategory struct {
	wrap   *wrapper.Wrapper
	object *utils.MySql_shop_category

	user *User

	bufferCats map[int]*utils.MySql_shop_category
}

func (this *ShopCategory) load(cache *map[int]*utils.MySql_shop_category) *ShopCategory {
	if this == nil {
		return this
	}
	if cache != nil {
		this.bufferCats = (*cache)
		return this
	}
	if this.bufferCats == nil {
		this.bufferCats = map[int]*utils.MySql_shop_category{}
	}
	if rows, err := this.wrap.DB.Query(`
		SELECT
			main.id,
			main.user,
			main.name,
			main.alias,
			main.lft,
			main.rgt,
			depth.depth,
			MAX(main.parent_id) AS parent_id
		FROM
			(
				SELECT
					node.id,
					node.user,
					node.name,
					node.alias,
					node.lft,
					node.rgt,
					parent.id AS parent_id
				FROM
					shop_cats AS node,
					shop_cats AS parent
				WHERE
					node.lft BETWEEN parent.lft AND parent.rgt AND
					node.id > 1
				ORDER BY
					node.lft ASC
			) AS main
			LEFT JOIN (
				SELECT
					node.id,
					(COUNT(parent.id) - 1) AS depth
				FROM
					shop_cats AS node,
					shop_cats AS parent
				WHERE
					node.lft BETWEEN parent.lft AND parent.rgt
				GROUP BY
					node.id
				ORDER BY
					node.lft ASC
			) AS depth ON depth.id = main.id
		WHERE
			main.id > 1 AND
			main.id <> main.parent_id
		GROUP BY
			main.id
		ORDER BY
			main.lft ASC
		;
	`); err == nil {
		defer rows.Close()
		for rows.Next() {
			row := utils.MySql_shop_category{}
			if err := rows.Scan(
				&row.A_id,
				&row.A_user,
				&row.A_name,
				&row.A_alias,
				&row.A_lft,
				&row.A_rgt,
				&row.A_depth,
				&row.A_parent,
			); err == nil {
				this.bufferCats[row.A_id] = &row
			}
		}
	}
	return this
}

func (this *ShopCategory) loadById(id int) {
	if this == nil {
		return
	}
	if this.object != nil {
		return
	}
	this.object = &utils.MySql_shop_category{}
	if err := this.wrap.DB.QueryRow(`
		SELECT
			main.id,
			main.user,
			main.name,
			main.alias,
			main.lft,
			main.rgt,
			depth.depth,
			MAX(main.parent_id) AS parent_id
		FROM
			(
				SELECT
					node.id,
					node.user,
					node.name,
					node.alias,
					node.lft,
					node.rgt,
					parent.id AS parent_id
				FROM
					shop_cats AS node,
					shop_cats AS parent
				WHERE
					node.lft BETWEEN parent.lft AND parent.rgt AND
					node.id > 1
				ORDER BY
					node.lft ASC
			) AS main
			LEFT JOIN (
				SELECT
					node.id,
					(COUNT(parent.id) - 1) AS depth
				FROM
					shop_cats AS node,
					shop_cats AS parent
				WHERE
					node.lft BETWEEN parent.lft AND parent.rgt
				GROUP BY
					node.id
				ORDER BY
					node.lft ASC
			) AS depth ON depth.id = main.id
		WHERE
			main.id > 1 AND
			main.id <> main.parent_id AND
			main.id = ?
		GROUP BY
			main.id
		;`,
		id,
	).Scan(
		&this.object.A_id,
		&this.object.A_user,
		&this.object.A_name,
		&this.object.A_alias,
		&this.object.A_lft,
		&this.object.A_rgt,
		&this.object.A_depth,
		&this.object.A_parent,
	); err != nil {
		return
	}
}

func (this *ShopCategory) Id() int {
	if this == nil {
		return 0
	}
	return this.object.A_id
}

func (this *ShopCategory) User() *User {
	if this == nil {
		return nil
	}
	if this.user != nil {
		return this.user
	}
	this.user = (&User{wrap: this.wrap}).load()
	this.user.loadById(this.object.A_user)
	return this.user
}

func (this *ShopCategory) Name() string {
	if this == nil {
		return ""
	}
	return this.object.A_name
}

func (this *ShopCategory) Alias() string {
	if this == nil {
		return ""
	}
	return this.object.A_alias
}

func (this *ShopCategory) Left() int {
	if this == nil {
		return 0
	}
	return this.object.A_lft
}

func (this *ShopCategory) Right() int {
	if this == nil {
		return 0
	}
	return this.object.A_rgt
}

func (this *ShopCategory) Permalink() string {
	if this == nil {
		return ""
	}
	return "/shop/category/" + this.object.A_alias + "/"
}

func (this *ShopCategory) Level() int {
	if this == nil {
		return 0
	}
	return this.object.A_depth
}

func (this *ShopCategory) Parent() *ShopCategory {
	if this == nil {
		return nil
	}
	if this.bufferCats == nil {
		return nil
	}
	if _, ok := this.bufferCats[this.object.A_parent]; !ok {
		return nil
	}
	cat := &ShopCategory{wrap: this.wrap, object: this.bufferCats[this.object.A_parent]}
	return cat.load(&this.bufferCats)
}

func (this *ShopCategory) HaveChilds() bool {
	// TODO: Add ability
	if this == nil {
		return false
	}
	return false
}
