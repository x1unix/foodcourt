<?php
//	Developed by kosyak <kosyak_ua@yahoo.com>
//	Initialyze function for constant's, language, db... etc for the system

//	Init output:
//	0: 	no errors
//	1:	System required files error;
//	2:	Required class isn't appear.
//	3:	Required constant isn't appear.
//	4:	DB connect error;

function init()
{
	// Install path's for the system.

	setIncludesPath();
	
	if (include ('config.inc.php'))
	{
		// Config file is appear, nice... 
		if (!defined('DB_FOR_USE'))
			return 3;
			
		if (include (DB_FOR_USE.'_core.class.php'))
		{
			$res = connectToDb();
			if ($res !== true)
				return $res;
			
			
			
			// Init client ENV.
			$res = initClientEnv();
			if ($res !== true)
				return $res;
			
			// Init default values if needed
			$res = initDefEnv($defaultValues);
			if ($res !== true)
				return $res;
			
			// Init language
			$res = initSystemLanguage();
			if ($res !== true)
				return $res;
		}
		else
		{
			return 2;
		}
	}
	else
	{
		return 1;
	}
	
	return 0;
}

function connectToDb()
{
	$className = DB_FOR_USE.'Core';
	$db = new $className;
	if ($db->connect(DB_HOST, DB_USER, DB_PASS, DB_NAME))
	{
		return true;
	}
	else
	{
		return 4;
	}
	
}

function initSystemLanguage()
{
	if (isset($_COOKIE['lang']))
	{
		$lang = $_COOKIE['lang'];
		if (file_exists("./includes/".$lang.".language.php"))
		{
			define ("LANG", $lang);
		}
		else
		{
			if (file_exists(dirname(__FILE__)."/includes/ua.language.php"))
				define ("LANG", "ua");
			else
				return 1;
		}
	}
	else
	{
		if (file_exists(dirname(__FILE__)."/includes/ua.language.php"))
			define ("LANG", "ua");
		else
			return 1;
	}
	
	return true;
}

function setIncludesPath()
{
	
	set_include_path(dirname(__FILE__) . '/includes' . PATH_SEPARATOR .
	dirname(__FILE__) . '/html-includes' . PATH_SEPARATOR .
	dirname(__FILE__) . '/Net' . PATH_SEPARATOR .
	dirname(__FILE__) . '/Mail' . PATH_SEPARATOR .
	dirname(__FILE__) . PATH_SEPARATOR .
	get_include_path());
}

function initClientEnv()
{
	require_once("sysconfig.class.php");
	
	$sysConfig = new SysConfig();
	if ($sysConfig->assignVariables())
		return true;
		
	return 100;
}

function initDefEnv ($defaultValues = null)
{
	if (!empty($defaultValues))
	{
		foreach ($defaultValues AS $key => $value)
		{
			if (!defined($key))
			{
				define ($key, $value);
			}
		}
	}
	
	return true;
}

$res = init ();

if ($res !== 0)
{
	echo "init::Error::$res";
	exit (1);
}

?>