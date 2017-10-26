<?php
	// Developed by kosyak <kosyak_ua@yahoo.com>
	
	//// Index page script.

	require_once("init.php");
	require_once("user.class.php");
	require_once(LANG.".language.php");
	
	$user = new User;
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	
	if ($isLogged)
	{
		header ("Location: orders.php");
	}
	else
	{
		include ("home.html");
	}
?>