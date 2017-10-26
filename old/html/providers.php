<?php
	// coded by kosyak <kosyak_ua@yahoo.com>
	
	//// Provider's info add/remove/edit page script.
	
	function load ($_GET, $_POST)
	{
		require_once("init.php");
		require_once("user.class.php");
		require_once("provider.class.php");
		require_once(LANG.".language.php");
		
		$user = new User;
		$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
		
		if ($isLogged['status'] != 10)
			return false;
		
		if (isset($_GET['s_ok']))
		{
			$message[] = Localization::$m_data_saved;
		}
		
		if (isset($_GET['add']))
		{
			include "provider_edit.html";
			return true;
		}
		else if (isset($_GET['edit']))
		{
			$provider = new Provider;
			if ($provider->getProviderInfo($_GET['edit']) === false)
				return $provider->error;
			
			$providerInfo = $provider->data;
			$edit_id = $_GET['edit'];
			include ("provider_edit.html");
			return true;
		}
		else if (isset($_GET['delete']) && !empty($_GET['delete']))
		{
			$provider = new Provider;
			if ($provider->deleteProvider($_GET['delete']) === false)
				return $provider->error;
			
			header ("Location: ${_SERVER['HTTP_REFERER']}");
			return true;
		}
		else if (isset($_POST['do_save']))
		{
			$error = "";
			if (empty($_POST['name']))
			{
				$error[] = Localization::$e_fill_required_fields;
			}
			
			if (!$error)
			{
				$provider = new Provider;
				if (isset($_POST['edit_id']))
				{
					if ($provider->getProviderInfo($_POST['edit_id']) === false)
						return $provider->error;
					$providerName = $provider->data['name'];
					if ($providerName != $_POST['name'])
					{
						if ($provider->getProviderInfo($_POST['name']))
						{
							$error = $provider->error;
							$error[] = Localization::$e_providername_exist;
						}
						else
						{
							$data['name'] = $_POST['name'];
						}
					}
					if (!$error)
					{
						
						$data['info'] = $_POST['info'];
						$data['first_price'] = $_POST['first_price'];
						$data['second_price'] = $_POST['second_price'];
						$data['third_price'] = $_POST['third_price'];
						if ($_POST['multichoice'] != 0)
						{
							$data['multichoice'] = 1;
						}
						else
						{
							$data['multichoice'] = 0;
						}
						
						if ($_POST['multiitem'] != 0)
						{
							$data['multiitem'] = 1;
						}
						else
						{
							$data['multiitem'] = 0;
						}
	
						if ($provider->editProviderInfo($data, $_POST['edit_id']) === false)
							return $provider->error;
						
						header("Location: providers.php?s_ok&edit=".$_POST['edit_id']);
						return true;
					}
					else
					{
						
						$providerInfo['name'] = $_POST['name'];
						$providerInfo['info'] = $_POST['info'];
						$providerInfo['first_price'] = $_POST['first_price'];
						$providerInfo['second_price'] = $_POST['second_price'];
						$providerInfo['third_price'] = $_POST['third_price'];
						$providerInfo['multichoice'] = $_POST['multichoice'];
						$providerInfo['multiitem'] = $_POST['multiitem'];
						$edit_id = $_POST['edit_id'];
						include ("provider_edit.html");
						return true;
					}
				}
				else
				{
					if (!$provider->getProviderInfo($_POST['name']))
					{
						$data = array(
							"name"			=> $_POST['name'],
							"info"			=> $_POST['info'],
							"first_price"	=> $_POST['first_price'],
							"second_price"	=> $_POST['second_price'],
							"third_price"	=> $_POST['third_price']
						);
						if ($_POST['multichoice'] != 0)
						{
							$data['multichoice'] = 1;
						}
						else
						{
							$data['multichoice'] = 0;
						}
						
						if ($_POST['multiitem'] != 0)
						{
							$data['multiitem'] = 1;
						}
						else
						{
							$data['multiitem'] = 0;
						}
						
						if ($provider->addProvider($data) === false)
							return $provider->error;
						
						header ("Location: providers.php");
						return true;
					}
					else
					{
						$error = $provider->error;
						$error[] = Localization::$e_providername_exist;
						$providerInfo['name'] = $_POST['name'];
						$providerInfo['info'] = $_POST['info'];
						$providerInfo['first_price'] = $_POST['first_price'];
						$providerInfo['second_price'] = $_POST['second_price'];
						$providerInfo['third_price'] = $_POST['third_price'];
						$providerInfo['multichoice'] = $_POST['multichoice'];
						$providerInfo['multiitem'] = $_POST['multiitem'];
						include ("provider_edit.html");
						return true;
					}
				}
			}
			else
			{
				$providerInfo['name'] = $_POST['name'];
				$providerInfo['info'] = $_POST['info'];
				$providerInfo['first_price'] = $_POST['first_price'];
				$providerInfo['second_price'] = $_POST['second_price'];
				$providerInfo['third_price'] = $_POST['third_price'];
				$providerInfo['multichoice'] = $_POST['multichoice'];
				$providerInfo['multiitem'] = $_POST['multiitem'];
				if (isset($_POST['edit_id']))
				{
					$edit_id = $_POST['edit_id'];
				}
				include ("provider_edit.html");
				return true;
			}
		}
		else
		{
			$provider = new Provider;
			if ($provider->getProviderList() === false)
				return $provider->error;
			$providerList = $provider->data;
			include ("providerlist.html");
			return true;
		}
		
		return false;
	}
	
	include ("request_handler.php");
?>