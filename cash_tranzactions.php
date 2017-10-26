<?php
	// coded by kosyak <kosyak_ua@yahoo.com>
	
	//// Cash tranzactions view script :)
	
function load ($_GET, $_POST)
{
	require_once("init.php");
	require_once ("user.class.php");
	require_once ("cash_tranzaction.class.php");
	require_once (LANG.".language.php");
	
	$user = new User;
	$scriptName = "cash_tranzactions.php";
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);

	if (!$isLogged)
		return false;
	if (isset($_GET['show']) && $isLogged['status'] == 10)
	{	
		if (is_numeric($_GET['show']))
		{
			$showId = $_GET['show'];
		}
		else
		{
			$showId = null;
		}
		
		$cashTranzaction = new CashTranzaction();
		if (!is_numeric($_GET['page']))
		{
			$_GET['page'] = 1;
		}
		
		$currentPage = $_GET['page'];
		$itemPerPage = ITEM_PER_PAGE;
		$cashTranzaction = new CashTranzaction();
		if ($cashTranzaction->getTranzactionsCount($showId) === false)
			return $cashTranzaction->error;
		$itemCount = $cashTranzaction->data;
		if ($cashTranzaction->getTranzactionList($showId, $itemPerPage, $currentPage) === false)
			return $cashTranzaction->error;
		$cashTranzactionList = $cashTranzaction->data;
		if ($user->getUserList(null, 'user_id') === false)
			return $user->error;
		$userListById = $user->data;
		$userListById[0]['user_name'] = "Voracity System";
		if ($user->getUserList ('user_name', null) === false)
			return $user->error;
		$userList = $user->data;
		include ("cash_tranzaction_list.html");
		return true;
	}
	else
	{
		if (!is_numeric($_GET['page']))
		{
			$_GET['page'] = 1;
		}
		
		$currentPage = $_GET['page'];
		$itemPerPage = ITEM_PER_PAGE;
		$cashTranzaction = new CashTranzaction();
		if ($cashTranzaction->getTranzactionsCount($isLogged['user_id']) === false)
			return $cashTranzaction->error;
		$itemCount = $cashTranzaction->data;
		if ($cashTranzaction->getTranzactionList($isLogged['user_id'], $itemPerPage, $currentPage) === false)
			return $cashTranzaction->error;
		$cashTranzactionList = $cashTranzaction->data;
		if ($user->getUserList(null, 'user_id') === false)
			return $user->error;
		$userListById = $user->data;
		$userListById[0]['user_name'] = "Voracity System";
		include ("cash_tranzaction_list.html");
		return true;
	}
	return false;
}

include ("request_handler.php");
?>