<?php
// Developed by kosyak <kosyak_ua@yahoo.com>

//// configurations

//Define DB constant's.
define ("DB_FOR_USE", "mysql");
define("DB_HOST", "localhost");
define ("DB_USER", "user");
define ("DB_PASS", "pass");
define ("DB_NAME", "voracity");

// Define DB table constant's.
define ("TBL_PREFIX", "vi_");											// Tables prefix.
define ("USER_TBL", TBL_PREFIX."users");								// User table name.
define ("WEEK_TBL", TBL_PREFIX."weeks");							// Week's table name.
define ("PROVIDER_TBL", TBL_PREFIX."providers");						// Providers table name.
define ("ORDER_LIST_TBL", TBL_PREFIX."order_list");						// Order list table name.
define ("USER_ORDER_LIST_TBL", TBL_PREFIX."user_order");				// User order list table name.
define ("USERS_REPLACEMENT_TBL", TBL_PREFIX."users_replacement");	// Users replacement list table name.
define ("CASH_TRANZACTION_TBL", TBL_PREFIX."cash_tranzactions");		// Cash tranzactions table...
define ("SYS_CONFIG_TBL", TBL_PREFIX."sys_config"); 					// System configuration
define ("SYS_CATEGORIES_TBL", TBL_PREFIX."sys_categories"); 			// System categories
define ("COMPANIES_TBL", TBL_PREFIX."companies");						// The companies table.


// Define Days constant's
define ("D_MONDAY", 1);
define ("D_TUESDAY", 2);
define ("D_WEDNESDAY", 3);
define ("D_THURSDAY", 4);
define ("D_FRIDAY", 5);
define ("D_SATURDAY", 6);
define ("D_SUNDAY", 7);


$defaultValues["PROJECT_NAME"] = "Voracity";
$defaultValues["SYSTEM_HTTP_ADDRESS"] = "http://172.22.70.89/voracity/";

// Define cookie's constant's.
$defaultValues["C_LOGIN"] = "lk_login";						// Login cookie name.
$defaultValues["C_PASSWORD"] = "lk_password";					// Password cookie name.
$defaultValues["C_SHOW_ORDERED_FOR"] = "lk_show_ordered_for";				// Cookie for SHOW_ORDERED_FOR
$defaultValues["C_SHOW_PAYED_FOR"] = "lk_show_payed_for";		// Cookie for SHOW_PAYED_FOR

$defaultValues["ITEM_PER_PAGE"] = 20;						// How much items displayed per one page.
$defaultValues["RELEASE_VERSION"] = '1.0';					// Version of this release.
$defaultValues["PRICE_CURRENCY"] = "грн.";					// The price currency :)
$defaultValues["PROVIDERS_MULTICHOICE"] = 1;				// 0 - don't allow to choice from different providers to one day; 1 - allow to choice....
$defaultValues["SHOW_PRICE_FOR_USER"] = 1;					// 0 - don't show price's for user (price for order items); 1 - show price's for user...
$defaultValues["SHOW_ORDERED_FOR"] = "week";				// Show users that have some order for a: week; day.
$defaultValues["SHOW_PAYED_FOR"] = "next_day";				// Show users that ordered for lunch for a: next_day;


$defaultValues["DB_DUMP_PATH"] = dirname(__FILE__)."/db_dump";		// Some place where put db dump files. Should be .htaccessed!!! By default it is: ./db_dump

// Define LDAP constant's
$defaultValues["LDAP_USE"] = 0;				// 0 - don't use ldap server.
$defaultValues["LDAP_SERVER"] = "deka";
$defaultValues["LDAP_PORT"] = 389;
$defaultValues["LDAP_DN"] = "OU=People,DC=lv,DC=company,DC=com";		// The base DN for the directory. 


$defaultValues["CORP_EMAIL"] = 'llnw.com';


//SMTP server constant
$defaultValues["SMTP_HOST"] = "ssl://localhost";
$defaultValues["SMTP_PORT"] = "465";
$defaultValues["SMTP_USER"] = "";
$defaultValues["SMTP_PASSWORD"] = "";



// Define Blocked days constant's
$defaultValues["D_MONDAY_BLOCK"] = D_TUESDAY;
$defaultValues["D_TUESDAY_BLOCK"] = D_WEDNESDAY;
$defaultValues["D_WEDNESDAY_BLOCK"] = D_THURSDAY;
$defaultValues["D_THURSDAY_BLOCK"] = D_FRIDAY;
$defaultValues["D_FRIDAY_BLOCK"] = D_MONDAY;
$defaultValues["D_SATURDAY_BLOCK"] = D_MONDAY;
$defaultValues["D_SUNDAY_BLOCK"] = D_MONDAY;
?>
