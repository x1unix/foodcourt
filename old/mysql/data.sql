-- phpMyAdmin SQL Dump
-- version 2.11.9.4
-- http://www.phpmyadmin.net
--
-- Хост: localhost
-- Час створення: Чрв 09 2009 р., 17:05
-- Версія сервера: 5.0.70
-- Версія PHP: 5.2.9-pl2-gentoo


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- БД: `voracity`
--

-- --------------------------------------------------------

--
-- Структура таблиці `vi_cash_tranzactions`
--

use `voracity`;

CREATE TABLE IF NOT EXISTS `vi_cash_tranzactions` (
  `tranzaction_id` bigint(20) NOT NULL auto_increment,
    `tranzaction_user_id_from` int(11) NOT NULL,
      `tranzaction_user_id_to` int(11) NOT NULL,
        `tranzaction_type` varchar(1) NOT NULL COMMENT 'a - add; b - return (back); r - remove',
	  `tranzaction_comment` text,
	    `tranzaction_amount` varchar(10) NOT NULL,
	      `tranzaction_time` int(11) NOT NULL,
	        PRIMARY KEY  (`tranzaction_id`)
		) TYPE=MyISAM  AUTO_INCREMENT=3579 ;

		-- --------------------------------------------------------

		--
		-- Структура таблиці `vi_companies`
		--

		CREATE TABLE IF NOT EXISTS `vi_companies` (
		  `company_id` int(11) NOT NULL auto_increment,
		    `company_name` varchar(128) NOT NULL,
		      PRIMARY KEY  (`company_id`)
		      ) TYPE=MyISAM AUTO_INCREMENT=1 ;

		      -- --------------------------------------------------------

		      --
		      -- Структура таблиці `vi_order_list`
		      --

		      CREATE TABLE IF NOT EXISTS `vi_order_list` (
		        `order_list_id` int(11) NOT NULL auto_increment,
			  `week_id` int(11) NOT NULL,
			    `day_id` int(11) NOT NULL,
			      `portion_number` int(11) NOT NULL,
			        `portion_name` varchar(128) NOT NULL,
				  `provider_id` int(11) NOT NULL,
				    `blocked` int(11) NOT NULL default '0',
				      `order_price` varchar(6) NOT NULL default '0',
				        PRIMARY KEY  (`order_list_id`)
					) TYPE=MyISAM  AUTO_INCREMENT=1956 ;

					-- --------------------------------------------------------

					--
					-- Структура таблиці `vi_providers`
					--

					CREATE TABLE IF NOT EXISTS `vi_providers` (
					  `provider_id` int(11) NOT NULL auto_increment,
					    `name` varchar(128) NOT NULL,
					      `info` text,
					        `first_price` varchar(10) default NULL,
						  `second_price` varchar(10) default NULL,
						    `third_price` varchar(10) default NULL,
						      `multiitem` int(11) NOT NULL default '0' COMMENT '0 - don''t allow users to set count of portions; 1 - allow...',
						        `multichoice` int(11) NOT NULL default '0' COMMENT '0 - no multi choice; 1 (or other) - multichoice',
							  PRIMARY KEY  (`provider_id`)
							  ) TYPE=MyISAM  AUTO_INCREMENT=13 ;

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

										    -- --------------------------------------------------------

										    --
										    -- Структура таблиці `vi_users`
										    --

										    CREATE TABLE IF NOT EXISTS `vi_users` (
										      `user_id` int(11) NOT NULL auto_increment,
										        `login` varchar(100) NOT NULL,
											  `ldap_login` varchar(128) default NULL,
											    `password` varchar(32) NOT NULL,
											      `ldap_password` varchar(32) NOT NULL,
											        `email` varchar(100) NOT NULL,
												  `status` int(11) NOT NULL default '1',
												    `user_name` varchar(128) default NULL,
												      `activity` int(11) NOT NULL default '1' COMMENT '0 - disabled; 1 - enabled',
												        `cash_total` varchar(10) NOT NULL default '0',
													  `fp_key` varchar(32) default NULL COMMENT 'Forgot password key',
													    `time` int(11) NOT NULL,
													      PRIMARY KEY  (`user_id`)
													      ) TYPE=MyISAM  AUTO_INCREMENT=188 ;

													      -- --------------------------------------------------------

													      --
													      -- Структура таблиці `vi_users_replacement`
													      --

													      CREATE TABLE IF NOT EXISTS `vi_users_replacement` (
													        `user_replacement_id` bigint(11) NOT NULL auto_increment,
														  `user_id` int(11) NOT NULL,
														    `replacement_id` int(11) NOT NULL,
														      `time` int(11) NOT NULL,
														        PRIMARY KEY  (`user_replacement_id`)
															) TYPE=MyISAM  AUTO_INCREMENT=35 ;

															-- --------------------------------------------------------

															--
															-- Структура таблиці `vi_user_order`
															--

															CREATE TABLE IF NOT EXISTS `vi_user_order` (
															  `user_order_id` int(11) NOT NULL auto_increment,
															    `week_id` int(11) NOT NULL,
															      `order_list_id` int(11) NOT NULL,
															        `user_id` int(11) NOT NULL,
																  `day_id` int(11) NOT NULL,
																    `provider_id` int(11) NOT NULL,
																      `ordered_item_count` int(11) NOT NULL default '1',
																        `ordered_amount` varchar(10) NOT NULL default '0' COMMENT 'How much user pay for that item',
																	  PRIMARY KEY  (`user_order_id`)
																	  ) TYPE=MyISAM  AUTO_INCREMENT=32382 ;

																	  -- --------------------------------------------------------

																	  --
																	  -- Структура таблиці `vi_user_to_company`
																	  --

																	  CREATE TABLE IF NOT EXISTS `vi_user_to_company` (
																	    `user_id` int(11) NOT NULL,
																	      `company_id` int(11) NOT NULL,
																	        KEY `user_id` (`user_id`,`company_id`)
																		) TYPE=MyISAM;

																		-- --------------------------------------------------------

																		--
																		-- Структура таблиці `vi_weeks`
																		--

																		CREATE TABLE IF NOT EXISTS `vi_weeks` (
																		  `week_id` int(11) NOT NULL auto_increment,
																		    `name` varchar(128) NOT NULL,
																		      `active` int(11) NOT NULL default '0',
																		        PRIMARY KEY  (`week_id`)
																			) TYPE=MyISAM  AUTO_INCREMENT=47 ;

