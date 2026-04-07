package model

type BankOracle struct {
	ID              string `sqlx:"id" json:"id"`
	BankName        string `sqlx:"bank_name" json:"bank_name"`
	Logo            string `sqlx:"logo" json:"logo"`
	BICCode         string `sqlx:"bic_code" json:"bic_code"`
	IsEnabled       int    `sqlx:"is_enabled" json:"is_enabled"`
	AccountLength   int    `sqlx:"account_length" json:"account_length"`
	HasAlphaNumeric int    `sqlx:"has_alpha_numeric" json:"has_alpha_numeric"`
	CreateAt        string `sqlx:"create_at" json:"create_at"`
	UpdateAt        string `sqlx:"update_at" json:"update_at"`
}

// create table banks(
//   id uuid primary key default gen_random_uuid(),
//   bank_name varchar(32) not null,
//   logo varchar(255) not null,
//   bic_code varchar(16) not null,
//   is_enabled int default 1,
//   account_length int not null,
//   has_alpha_numeric int default 0,
//   create_at timestamp default now(),
//   update_at timestamp default now()
// );
