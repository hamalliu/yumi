package sql

import (
	"yumi/pkg/external/dbc"
)

func SysmmngCreateTable() {
	//管理账号表
	backPsAccountSql := `CREATE TABLE "back_ps_accounts"
						(
						"id" int primary key,
						"user" Nvarchar(255) not null default('') unique,
						"user_id" Nvarchar(255) not null default('') unique,
						"mobile" Nvarchar(11) default(''),
						"password" Nvarchar(255) default(''),
						"job_number" Nvarchar(255) default('') unique,
						"status" Nvarchar(5) default(''),
						"comment" Nvarchar(512) default(''),
						"register_time" datetime default now(),
						"operator" Nvarchar(255) default(''),
						"operate_time" datetime default now()
						)`
	if _, err := dbc.Get().Exec(backPsAccountSql); err != nil {
		panic(err)
	}

	//管理角色表
	backPsRoleSql := `CREATE TABLE "back_ps_roles"
						(
						"id" int primary key,
						"name" Nvarchar(255) not null default(''),
						"code" Nvarchar(255) not null default('') unique,
						"status" Nvarchar(5) default(''),
						"operator" Nvarchar(255) default(''),
						"operate_time" datetime default now()
						)`
	if _, err := dbc.Get().Exec(backPsRoleSql); err != nil {
		panic(err)
	}

	//管理角色用户关联表
	backPsAccountRoles := `CREATE TABLE "back_ps_account_roles"
						(
						"id" int primary key,
						"ps_account_id" Nvarchar(255) default(''),
						"ps_role_id" Nvarchar(255) default('')
						)`
	if _, err := dbc.Get().Exec(backPsAccountRoles); err != nil {
		panic(err)
	}

	//管理权限表
	backPsPowerSql := `CREATE TABLE "back_ps_powers"
						(
						"id" int primary key,
						"ps_role_id" Nvarchar(255) default(''),
						"ps_account_id" Nvarchar(255) default(''),
						"power_branch" text
						)`
	if _, err := dbc.Get().Exec(backPsPowerSql); err != nil {
		panic(err)
	}

	//管理模块表
	backPsModuleSql := `CREATE TABLE "back_ps_modules"
						(
						"id" int primary key,
						"prt_module_id" int default(0),
						"prt_module_name" Nvarchar(255) default(''),
						"name" Nvarchar(255) default(''),
						"route" Nvarchar(255) default(''),
						"expand" Nvarchar(255) default(''),
						"code" Nvarchar(255) not null default('') unique,
						"cur_func_code" int default(0),
						"cur_sub_code" int default(0),
						"display_order" int default(0),
						"status" Nvarchar(255) default(''),
						"type" Nvarchar(15) default(''),
						"params" Nvarchar(255) default(''),
						"operator" Nvarchar(255) default(''),
						"operate_time" datetime default now()
						)`
	if _, err := dbc.Get().Exec(backPsModuleSql); err != nil {
		panic(err)
	}

	//更新记录表
	updateRecordSql := `CREATE TABLE "back_update_records"
						(
						"id" int primary key,
						"table" Nvarchar(255) default(''),
						"request" text,
						"request_body" text,
						"before_data" text,
						"user_id" Nvarchar(255) default(''),
						"operator" Nvarchar(255) default(''),
						"operate_time" datetime default now()
						)`
	if _, err := dbc.Get().Exec(updateRecordSql); err != nil {
		panic(err)
	}

	//删除记录表
	deleteRecordSql := `CREATE TABLE "back_delete_records"
						(
						"id" int primary key,
						"table" Nvarchar(255) default(''),
						"request" text,
						"request_body" text,
						"before_data" text,
						"user_id" Nvarchar(255) default(''),
						"operator" Nvarchar(255) default(''),
						"operate_time" datetime default now()
						)`
	if _, err := dbc.Get().Exec(deleteRecordSql); err != nil {
		panic(err)
	}

	//支付订单表
	_ = `CREATE TABLE "order_pay" 
		(
		"seq_id" int identity(1,1) primary key,
		"code" varchar(64) default(''), 
		"trade_way" varchar(30) default(''), 
		"seller_key" varchar(255) default(''), 
		"app_id" varchar(255) default(''), 
		"mch_id" varchar(255) default(''), 
		"transaction_id" varchar(255) default(''), 
		"notify_url" varchar(255) default(''), 
		"buyer_logon_id" varchar(255) default(''), 
		"spbill_create_ip" varchar(50) default(''),

		"buyer_account_guid" varchar(255) default(''), 
		"total_fee" int default(0), 
		"body" varchar(1000) default(''), 
		"detail" varchar(1000) default(''), 
		"out_trade_no" varchar(255) default(''), 
		"timeout_express" datetime,
		"pay_expire" datetime, 
		"pay_time" datetime, 
		"cancel_time" datetime, 
		"error_time" datetime, 
		"submit_time" datetime, 
		"status" varchar(30) default(''), 
		"remarks" varchar(500) default('')
		)`

	//支付宝支付商户表
	_ = `CREATE TABLE "ali_pay_merchants" 
		(
		"seq_id" int identity(1,1) primary key, 
		"seller_key" varchar(255) default(''), 
		"app_id" varchar(255) default(''), 
		"private_key" varchar(255) default(''), 
		"public_key" varchar(255) default('')
		)`

	//微信支付商户表
	_ = `CREATE TABLE "wx_pay_merchants" 
		(
		"seq_id" int identity(1,1) primary key, 
		"seller_key" varchar(255) default(''), 
		"app_id" varchar(255) default(''), 
		"mch_id" varchar(255) default(''),
		"private_key" varchar(255) default('')
		)`

}
