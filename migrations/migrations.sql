-- ##########################################
-- ################ ОСНОВНЫЕ ################
-- ##########################################

-- //////////////////////////
-- //Таблица пользователей://
-- //////////////////////////
create table IF NOT EXISTS users
(
	user_id int auto_increment,
	phone varchar(255) not null,
	email varchar(255) not null,
	first_name tinytext not null,
	last_name tinytext not null,
	unencrypted_pass tinytext not null,
	role tinytext not null,
	created_at timestamp not null,
	confirm_phone tinyint(1) default 0 not null,
	updated_at timestamp default '2000-01-01 00:00:00' not null,
	perek tinyint(1) default 0 not null,
	perek_card varchar(50) default '' not null,
	city tinytext not null,
	age int not null,
	sex tinytext not null,
	confirm_email tinyint(1) default 0 not null,
	unsubscribe_email tinyint(1) default 0 not null,
	constraint users_email_uindex
		unique (email),
	constraint users_phone_uindex
		unique (phone),
	constraint users_user_id_uindex
		unique (user_id)
);

alter table IF NOT EXISTS  users
	add primary key (user_id);

-- ///////////////////////////////////////////
-- //Таблица для хранения ссылок авторизации//
-- ///////////////////////////////////////////
create table IF NOT EXISTS conf_url
(
	user_id int not null,
	email varchar(255) not null,
	token varchar(255) not null,
	created_at timestamp default CURRENT_TIMESTAMP not null,
	constraint conf_url_token_uindex
		unique (token)
);

-- ///////////////////////////////////////////
-- //Таблица для хранения кодов  авторизации//
-- ///////////////////////////////////////////
create table IF NOT EXISTS conf_phone
(
	phone tinytext not null,
	code tinytext not null,
	user_id int not null
);

-- //////////////////////////////
-- //Таблица для хранения чеков//
-- //////////////////////////////
create table IF NOT EXISTS bill_good_response
(
	bill_id int auto_increment,
	user_id int not null,
	code_operation int default 0 not null,
	shop varchar(255) default '' not null,
	fns_url varchar(255) default '' not null,
	seller_address varchar(255) default '' not null,
	kkt_reg_id varchar(255) default '' not null,
	retail_place varchar(255) default '' not null,
	retail_place_address varchar(255) default '' not null,
	date_check varchar(255) default '' not null,
	check_number int default 0 not null,
	total_sum_cop bigint default 0 not null,
	shift_number int default 0 not null,
	operation_type int default 0 not null,
	drive_num varchar(255) default '' not null,
	doc_num bigint default 0 not null,
	fiscal_sign bigint default 0 not null,
	fiscal_doc_format int default 0 not null,
	download_time timestamp default CURRENT_TIMESTAMP not null,
	check_time timestamp default '2000-01-01 00:00:00' null,
	status varchar(255) not null,
	status_for_user varchar(255) not null,
	shop_inn varchar(255) default '' not null,
	perek_by_user tinyint(1) default 0 not null,
	perek_card varchar(255) default '' not null,
	comments varchar(255) default '' not null,
	day_win varchar(255) default '' null,
	url varchar(255) default '' null,
	constraint bill_good_response_bill_id_uindex
		unique (bill_id)
);

alter table bill_good_response
	add primary key (bill_id);

-- //////////////////////////////////////
-- //Таблица для хранения позиций чеков//
-- //////////////////////////////////////
create table IF NOT EXISTS position_in_bill
(
	bill_id int not null,
	user_id int not null,
	name varchar(255) default '' not null,
	price int default 0 not null,
	count double default 0 not null,
	sum bigint default 0 not null,
	marker varchar(255) default '' not null,
	id int auto_increment,
	constraint position_in_bill_id_uindex
		unique (id)
);

alter table position_in_bill
	add primary key (id);

-- ///////////////////////////////////////////////////////////////////
-- //Таблица для хранения возможных наименований акционной продукции//
-- ///////////////////////////////////////////////////////////////////
create table IF NOT EXISTS product_name
(
	name varchar(255) not null,
	id int auto_increment,
	constraint product_name_id_uindex
		unique (id)
);

alter table IF NOT EXISTS product_name
	add primary key (id);

-- ////////////////////////////////////
-- //Таблица для хранения победителей//
-- ////////////////////////////////////
create table IF NOT EXISTS winners
(
	win_id int auto_increment
		primary key,
	user_id int not null,
	bill_id int not null,
	win_time timestamp default CURRENT_TIMESTAMP not null,
	prize_type varchar(255) not null,
	prize_status varchar(255) not null,
	comment varchar(255) default '' not null
);




-- ########################################
-- ################ ЛОГЕРЫ ################
-- ########################################

-- ////////////////////////////////////////
-- //Таблица для хранения логов к гифтери//
-- ////////////////////////////////////////
create table IF NOT EXISTS giftery_log_req
(
	req_id int auto_increment,
	time timestamp default CURRENT_TIMESTAMP not null,
	user_id int not null,
	req text not null,
	bill_id int not null,
	method varchar(255) not null,
	constraint giftery_log_req_req_id_uindex
		unique (req_id)
);

alter table IF NOT EXISTS giftery_log_req
	add primary key (req_id);


-- /////////////////////////////////////////
-- //Таблица для хранения логов от гифтери//
-- /////////////////////////////////////////
create table IF NOT EXISTS giftery_log_response
(
	req_id int not null,
	time timestamp default CURRENT_TIMESTAMP not null,
	response longtext null,
	user_id int not null,
	bill_id int not null,
	method varchar(255) not null,
	constraint giftery_log_response_req_id_uindex
		unique (req_id)
);

-- ////////////////////////////////////
-- //Таблица для хранения логов к смс//
-- ////////////////////////////////////
create table IF NOT EXISTS log_sms
(
	id int auto_increment,
	time timestamp default CURRENT_TIMESTAMP not null,
	phone varchar(255) not null,
	user_id int not null,
	constraint log_sms_id_uindex
		unique (id)
);

alter table log_sms
	add primary key (id);

-- /////////////////////////////////////
-- //Таблица для хранения логов от смс//
-- /////////////////////////////////////
create table IF NOT EXISTS log_sms_resp
(
	resp text not null,
	req_id int not null,
	user_id int not null
);

-- /////////////////////////////////////////////
-- //Таблица для хранения логов к проверкеЧека//
-- /////////////////////////////////////////////
create table IF NOT EXISTS logger_req_checkbill
(
	logger text not null,
	user_id int not null,
	bill_id int not null,
	time timestamp default CURRENT_TIMESTAMP not null
);

-- //////////////////////////////////////////////
-- //Таблица для хранения логов от проверкиЧека//
-- //////////////////////////////////////////////
create table IF NOT EXISTS logger_resp_checkbill
(
	loger longtext null,
	user_id int not null,
	bill_id int not null,
	time timestamp default CURRENT_TIMESTAMP not null
);

-- /////////////////////////////////////
-- //Таблица для логирования обращений//
-- /////////////////////////////////////
create table IF NOT EXISTS os_logger
(
	os_id int auto_increment
		primary key,
	user_id int null,
	contact_name varchar(255) not null,
	contact_email varchar(255) not null,
	contact_phone varchar(255) not null,
	text text not null,
	time_send timestamp default CURRENT_TIMESTAMP not null
);

-- //////////////////////////////////////////////////
-- //Таблица для логирования действий пользователей//
-- //////////////////////////////////////////////////
create table IF NOT EXISTS request_log
(
	id bigint auto_increment
		primary key,
	dump_request varchar(1000) default '' not null,
	ip varchar(100) default '' not null,
	route varchar(100) default '' not null,
	created_at timestamp default CURRENT_TIMESTAMP null
);

-- /////////////////////////////////////
-- //Таблица для логирования розыгрыша//
-- /////////////////////////////////////
create table IF NOT EXISTS winner_logs
(
	num int auto_increment,
	bill_id int not null,
	user_id int not null,
	status varchar(255) not null,
	time timestamp default CURRENT_TIMESTAMP not null,
	dollar float not null,
	time_bill timestamp null,
	time_validate timestamp null,
	constraint winner_logs_bill_id_uindex
		unique (bill_id),
	constraint winner_logs_num_uindex
		unique (num)
);

alter table winner_logs
	add primary key (num);

-- /////////////////////////////////
-- //Таблица для логирования юмани//
-- /////////////////////////////////
create table IF NOT EXISTS yoomoney_log
(
	id int auto_increment,
	user_id int not null,
	type varchar(255) not null,
	error varchar(255) null,
	request_id varchar(255) null,
	status varchar(255) null,
	amount varchar(255) null,
	payment_id varchar(255) default '' null,
	invoice_id varchar(255) default '' null,
	phone varchar(255) not null,
	time_action timestamp default CURRENT_TIMESTAMP not null,
	constraint yoomoney_log_id_uindex
		unique (id)
);

alter table yoomoney_log
	add primary key (id);

-- ########################################################
-- ################ ВРЕМЕННЫЕ/ДЛЯ СКРИПТОВ ################
-- ########################################################

-- /////////////////////////////////////
-- //Таблица для единоразовой рассылки//
-- /////////////////////////////////////
create table IF NOT EXISTS forsend
(
	email varchar(255) null
);

-- ///////////////////////////////////
-- //Таблица для повышенного кэшбека//
-- ///////////////////////////////////
create table IF NOT EXISTS big_cash
(
	phone varchar(255) default '' not null,
	id int auto_increment,
	constraint big_cash_id_uindex
		unique (id)
);

alter table big_cash
	add primary key (id);

