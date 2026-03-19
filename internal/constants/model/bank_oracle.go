package model

type BankOracle struct {
	ID              string `json:"id"`
	BankName        string `json:"bank_name"`
	Logo            string `json:"logo"`
	BICCode         string `json:"bic_code"`
	IsEnabled       int    `json:"is_enabled"`
	AccountLength   int    `json:"account_length"`
	HasAlphaNumeric int    `json:"has_alpha_numeric"`
	CreateAt        string `json:"create_at"`
	UpdateAt        string `json:"update_at"`
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
