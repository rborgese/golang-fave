package fetdata

import (
	"math"
	"strings"

	"golang-fave/engine/wrapper"
	"golang-fave/utils"
)

type ShopPagination struct {
	Num     string
	Link    string
	Current bool
	Dots    bool
}

type Shop struct {
	wrap     *wrapper.Wrapper
	category *ShopCategory
	product  *ShopProduct

	products         []*ShopProduct
	productsCount    int
	productsPerPage  int
	productsMaxPage  int
	productsCurrPage int
	pagination       []*ShopPagination
	paginationPrev   *ShopPagination
	paginationNext   *ShopPagination

	bufferCats map[string][]*ShopCategory
}

func (this *Shop) load() {
	if this == nil {
		return
	}
	sql_nums := `
		SELECT
			COUNT(*)
		FROM
			shop_products
		WHERE
			active = 1
		;
	`
	sql_rows := `
		SELECT
			shop_products.id,
			shop_products.user,
			shop_products.currency,
			shop_products.price,
			shop_products.name,
			shop_products.alias,
			shop_products.briefly,
			shop_products.content,
			UNIX_TIMESTAMP(shop_products.datetime) as datetime,
			shop_products.active,
			users.id,
			users.first_name,
			users.last_name,
			users.email,
			users.admin,
			users.active,
			shop_currencies.id,
			shop_currencies.name,
			shop_currencies.coefficient,
			shop_currencies.code,
			shop_currencies.symbol
		FROM
			shop_products
			LEFT JOIN users ON users.id = shop_products.user
			LEFT JOIN shop_currencies ON shop_currencies.id = shop_products.currency
		WHERE
			shop_products.active = 1
		ORDER BY
			shop_products.id DESC
		LIMIT ?, ?;
	`

	// Category selected
	if this.category != nil {
		var cat_ids []string
		if rows, err := this.wrap.DB.Query(
			`SELECT
				node.id
			FROM
				shop_cats AS node,
				shop_cats AS parent
			WHERE
				node.lft BETWEEN parent.lft AND parent.rgt AND
				node.id > 1 AND
				parent.id = ?
			GROUP BY
				node.id
			ORDER BY
				node.lft ASC
			;`,
			this.category.Id(),
		); err == nil {
			defer rows.Close()
			for rows.Next() {
				var cat_id string
				if err := rows.Scan(&cat_id); err == nil {
					cat_ids = append(cat_ids, cat_id)
				}
			}
		}
		sql_nums = `
			SELECT
				COUNT(*)
			FROM
				(
					SELECT
						COUNT(*)
					FROM
						shop_products
						LEFT JOIN shop_cat_product_rel ON shop_cat_product_rel.product_id = shop_products.id
					WHERE
						shop_products.active = 1 AND
						shop_cat_product_rel.category_id IN (` + strings.Join(cat_ids, ", ") + `)
					GROUP BY
						shop_products.id
				) AS tbl
			;
		`
		sql_rows = `
			SELECT
				shop_products.id,
				shop_products.user,
				shop_products.currency,
				shop_products.price,
				shop_products.name,
				shop_products.alias,
				shop_products.briefly,
				shop_products.content,
				UNIX_TIMESTAMP(shop_products.datetime) AS datetime,
				shop_products.active,
				users.id,
				users.first_name,
				users.last_name,
				users.email,
				users.admin,
				users.active,
				shop_currencies.id,
				shop_currencies.name,
				shop_currencies.coefficient,
				shop_currencies.code,
				shop_currencies.symbol
			FROM
				shop_products
				LEFT JOIN shop_cat_product_rel ON shop_cat_product_rel.product_id = shop_products.id
				LEFT JOIN users ON users.id = shop_products.user
				LEFT JOIN shop_currencies ON shop_currencies.id = shop_products.currency
			WHERE
				shop_products.active = 1 AND
				shop_cat_product_rel.category_id IN (` + strings.Join(cat_ids, ", ") + `)
			GROUP BY
				shop_products.id
			ORDER BY
				shop_products.id DESC
			LIMIT ?, ?;
		`
	}

	if err := this.wrap.DB.QueryRow(sql_nums).Scan(&this.productsCount); err == nil {
		if this.category == nil {
			this.productsPerPage = (*this.wrap.Config).Shop.Pagination.Index
		} else {
			this.productsPerPage = (*this.wrap.Config).Shop.Pagination.Category
		}
		this.productsMaxPage = int(math.Ceil(float64(this.productsCount) / float64(this.productsPerPage)))
		this.productsCurrPage = this.wrap.GetCurrentPage(this.productsMaxPage)
		offset := this.productsCurrPage*this.productsPerPage - this.productsPerPage
		if rows, err := this.wrap.DB.Query(sql_rows, offset, this.productsPerPage); err == nil {
			defer rows.Close()
			for rows.Next() {
				rp := utils.MySql_shop_product{}
				ru := utils.MySql_user{}
				rc := utils.MySql_shop_currency{}
				if err := rows.Scan(
					&rp.A_id,
					&rp.A_user,
					&rp.A_currency,
					&rp.A_price,
					&rp.A_name,
					&rp.A_alias,
					&rp.A_briefly,
					&rp.A_content,
					&rp.A_datetime,
					&rp.A_active,
					&ru.A_id,
					&ru.A_first_name,
					&ru.A_last_name,
					&ru.A_email,
					&ru.A_admin,
					&ru.A_active,
					&rc.A_id,
					&rc.A_name,
					&rc.A_coefficient,
					&rc.A_code,
					&rc.A_symbol,
				); err == nil {
					this.products = append(this.products, &ShopProduct{
						wrap:     this.wrap,
						object:   &rp,
						user:     &User{wrap: this.wrap, object: &ru},
						currency: &Currency{wrap: this.wrap, object: &rc},
					})
				}
			}
		}
	}

	// Build pagination
	if true {
		for i := 1; i < this.productsCurrPage; i++ {
			if this.productsCurrPage >= 5 && i > 1 && i < this.productsCurrPage-1 {
				continue
			}
			if this.productsCurrPage >= 5 && i > 1 && i < this.productsCurrPage {
				this.pagination = append(this.pagination, &ShopPagination{
					Dots: true,
				})
			}
			link := this.wrap.R.URL.Path
			if i > 1 {
				link = link + "?p=" + utils.IntToStr(i)
			}
			this.pagination = append(this.pagination, &ShopPagination{
				Num:     utils.IntToStr(i),
				Link:    link,
				Current: false,
			})
		}

		// Current page
		link := this.wrap.R.URL.Path
		if this.productsCurrPage > 1 {
			link = link + "?p=" + utils.IntToStr(this.productsCurrPage)
		}
		this.pagination = append(this.pagination, &ShopPagination{
			Num:     utils.IntToStr(this.productsCurrPage),
			Link:    link,
			Current: true,
		})

		for i := this.productsCurrPage + 1; i <= this.productsMaxPage; i++ {
			if this.productsCurrPage < this.productsMaxPage-3 && i == this.productsCurrPage+3 {
				this.pagination = append(this.pagination, &ShopPagination{
					Dots: true,
				})
			}
			if this.productsCurrPage < this.productsMaxPage-3 && i > this.productsCurrPage+1 && i <= this.productsMaxPage-1 {
				continue
			}
			link := this.wrap.R.URL.Path
			if i > 1 {
				link = link + "?p=" + utils.IntToStr(i)
			}
			this.pagination = append(this.pagination, &ShopPagination{
				Num:     utils.IntToStr(i),
				Link:    link,
				Current: false,
			})
		}
	} else {
		for i := 1; i <= this.productsMaxPage; i++ {
			link := this.wrap.R.URL.Path
			if i > 1 {
				link = link + "?p=" + utils.IntToStr(i)
			}
			this.pagination = append(this.pagination, &ShopPagination{
				Num:     utils.IntToStr(i),
				Link:    link,
				Current: i == this.productsCurrPage,
			})
		}
	}

	// Pagination prev/next
	if this.productsMaxPage > 1 {
		link := this.wrap.R.URL.Path
		if this.productsCurrPage-1 > 1 {
			link = this.wrap.R.URL.Path + "?p=" + utils.IntToStr(this.productsCurrPage-1)
		}
		this.paginationPrev = &ShopPagination{
			Num:     utils.IntToStr(this.productsCurrPage - 1),
			Link:    link,
			Current: this.productsCurrPage <= 1,
		}
		if this.productsCurrPage >= 1 && this.productsCurrPage < this.productsMaxPage {
			link = this.wrap.R.URL.Path + "?p=" + utils.IntToStr(this.productsCurrPage+1)
		} else {
			link = this.wrap.R.URL.Path + "?p=" + utils.IntToStr(this.productsMaxPage)
		}
		this.paginationNext = &ShopPagination{
			Num:     utils.IntToStr(this.productsCurrPage + 1),
			Link:    link,
			Current: this.productsCurrPage >= this.productsMaxPage,
		}
	}
}

func (this *Shop) Category() *ShopCategory {
	if this == nil {
		return nil
	}
	return this.category
}

func (this *Shop) Product() *ShopProduct {
	if this == nil {
		return nil
	}
	return this.product
}

func (this *Shop) HaveProducts() bool {
	if this == nil {
		return false
	}
	if len(this.products) <= 0 {
		return false
	}
	return true
}

func (this *Shop) Products() []*ShopProduct {
	if this == nil {
		return []*ShopProduct{}
	}
	return this.products
}

func (this *Shop) ProductsCount() int {
	if this == nil {
		return 0
	}
	return this.productsCount
}

func (this *Shop) ProductsPerPage() int {
	if this == nil {
		return 0
	}
	return this.productsPerPage
}

func (this *Shop) ProductsMaxPage() int {
	if this == nil {
		return 0
	}
	return this.productsMaxPage
}

func (this *Shop) ProductsCurrPage() int {
	if this == nil {
		return 0
	}
	return this.productsCurrPage
}

func (this *Shop) Pagination() []*ShopPagination {
	if this == nil {
		return []*ShopPagination{}
	}
	return this.pagination
}

func (this *Shop) PaginationPrev() *ShopPagination {
	if this == nil {
		return nil
	}
	return this.paginationPrev
}

func (this *Shop) PaginationNext() *ShopPagination {
	if this == nil {
		return nil
	}
	return this.paginationNext
}

func (this *Shop) Categories(mlvl int) []*ShopCategory {
	if this == nil {
		return []*ShopCategory{}
	}
	if this.bufferCats == nil {
		this.bufferCats = map[string][]*ShopCategory{}
	}
	key := ""
	where := ``
	if mlvl > 0 {
		where += `AND tbl.depth <= ` + utils.IntToStr(mlvl)
	}
	if _, ok := this.bufferCats[key]; !ok {
		var cats []*ShopCategory
		if rows, err := this.wrap.DB.Query(`
			SELECT
				tbl.*
			FROM
				(
					SELECT
						node.id,
						node.user,
						node.name,
						node.alias,
						node.lft,
						node.rgt,
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
				) AS tbl
			WHERE
				tbl.id > 1
				` + where + `
			;
		`); err == nil {
			defer rows.Close()
			for rows.Next() {
				row := utils.MySql_shop_category{}
				var Depth int
				if err := rows.Scan(&row.A_id, &row.A_user, &row.A_name, &row.A_alias, &row.A_lft, &row.A_rgt, &Depth); err == nil {
					cats = append(cats, &ShopCategory{object: &row, depth: Depth})
				}
			}
		}
		this.bufferCats[key] = cats
	}
	return this.bufferCats[key]
}