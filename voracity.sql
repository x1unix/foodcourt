-- phpMyAdmin SQL Dump
-- version 2.11.9.4
-- http://www.phpmyadmin.net
--
-- Хост: localhost
-- Час створення: Бер 09 2009 р., 15:03
-- Версія сервера: 5.0.70
-- Версія PHP: 5.2.8-pl2-gentoo


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- БД: `voracity`
--

-- --------------------------------------------------------

--
-- Структура таблиці `vi_sys_categories`
--

CREATE TABLE IF NOT EXISTS `vi_sys_categories` (
  `id` int(11) NOT NULL auto_increment,
  `name` varchar(128) NOT NULL,
  `localization_var` varchar(128) default NULL,
  PRIMARY KEY  (`id`)
) TYPE=MyISAM  AUTO_INCREMENT=5 ;

--
-- Дамп даних таблиці `vi_sys_categories`
--

INSERT INTO `vi_sys_categories` (`id`, `name`, `localization_var`) VALUES
(1, 'ldap_configuration', 's_ldap_configuration'),
(4, 'db_configuration', 's_db_configuration'),
(3, 'total_system_configuration', 's_total_system_configuration');

-- --------------------------------------------------------

--
-- Структура таблиці `vi_sys_config`
--

CREATE TABLE IF NOT EXISTS `vi_sys_config` (
  `id` int(11) NOT NULL auto_increment,
  `name` varchar(128) NOT NULL,
  `value` text,
  `comment` text,
  `localization_var` varchar(128) default NULL,
  `cat_selected` int(11) NOT NULL default '0',
  `active` int(11) NOT NULL default '1',
  `type` varchar(10) NOT NULL default 'text',
  PRIMARY KEY  (`id`)
) TYPE=MyISAM  AUTO_INCREMENT=11 ;

--
-- Дамп даних таблиці `vi_sys_config`
--

INSERT INTO `vi_sys_config` (`id`, `name`, `value`, `comment`, `localization_var`, `cat_selected`, `active`, `type`) VALUES
(1, 'LDAP_USE', 'true', '1 - use ldap server\r\n0 - don''t use ldap server', 's_use_ldap_server', 1, 1023, 'bool'),
(2, 'LDAP_SERVER', 'deka', 'Ldap server address for use.', 's_ldap_server_address', 1, 1023, 'text'),
(3, 'LDAP_PORT', '389', 'The ldap server port for connect to.', 's_ldap_server_port', 1, 1023, 'int'),
(4, 'LDAP_DN', 'OU=People,DC=lv,DC=lohika,DC=com', 'The base DN for the directory.', 's_ldap_base_dn', 1, 1023, 'text'),
(5, 'PROJECT_NAME', 'Voracity', 'The project name.', 's_project_name', 3, 1023, 'text'),
(6, 'SYSTEM_HTTP_ADDRESS', 'http://localhost/voracity/', 'The system http address.', 's_system_http_address', 3, 1023, 'text'),
(7, 'ITEM_PER_PAGE', '15', 'The count of items to display per 1 page.', 's_item_per_page_count', 3, 1023, 'int'),
(8, 'DB_FOR_USE', 'mysql', 'DB for use.', 's_db_for_use', 4, 1023, 'text'),
(9, 'DB_NAME', 'voracity', 'DB name', 's_db_name', 4, 1023, 'text'),
(10, 'CORP_EMAIL', 'lv.lohika.com', 'Corporative e-mail address.', 's_corp_email', 3, 1023, 'text');
