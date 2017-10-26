<?php
	// Developed by kosyak<kosyak_ua@yahoo.com>
	
	//// Script fro user login.
	
function load ($_GET, $_POST)
{
	require_once("init.php");
	require_once("user.class.php");
	
	$user = new User;
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	
	if (!$isLogged)
	{
		if (!isset($_POST['submit']))
			return false;
		if ($userInfo = $user->getUserInfo($_POST['login']))
		{
			if ($userInfo['password'] == md5($_POST['password']))
			{
				$user->setUserCredential($_POST['login'], $_POST['password'], $_POST['remember_me']);
				header ("Location: ${_SERVER['HTTP_REFERER']}");
				return true;
			}
			else if (LDAP_USE)
			{
				if ($userInfo['ldap_password'] == md5($_POST['password']))
				{
					$user->setUserCredential($_POST['login'], $_POST['password'], $_POST['remember_me']);
					header ("Location: ${_SERVER['HTTP_REFERER']}");
					return true;
				}
				else if (!empty($_POST['password']))
				{
					// We have user but we don't have user ldap password
					if ($user->getLdapUserInfo($_POST['login'], $_POST['password']) === false)
						return $user->error;
					if (!empty($user->data))
					{
						$ldapUserInfo = $user->data;

						if ($user->getUserInfoByEmailPart($_POST['login'].'@') === false)
							return $user->error;
						$userNewInfo = $user->data;
						if ($userNewInfo['user_id'] == $userInfo['user_id'])
						{
							// We have the same user. nice...
							// We have this user.. need to update user ldap login/password
							$data = array(
							'ldap_login'	=> $_POST['login'],
							'ldap_password'	=> md5($_POST['password'])
							);
							if ($user->editUserInfo($data, $userInfo['user_id']) === false)
								return $user->error;
							$user->setUserCredential($_POST['login'], $_POST['password'], $_POST['remember_me']);
							header ("Location: ${_SERVER['HTTP_REFERER']}");
							return true;
						}
						else
						{
							echo "please report the bug";
							return false;
						}
					}
					else
					{
						// Bad user password.
						$isPresentErrMessage = false;
						$parameters = explode("?", $_SERVER['HTTP_REFERER']);
						$parametersItem = explode("&", $parameters[1]);
						foreach ($parametersItem as $item)
						{
							$temp = explode("=", $item);
							if ($temp[0] == "l_err")
							{
								$isPresentErrMessage = true;
							}
						}
						if ($isPresentErrMessage)
						{
							header ("Location: ${_SERVER['HTTP_REFERER']}");
							return true;
						}
						else
						{
							if (!empty($parameters[1]))
							{
								$_SERVER['HTTP_REFERER'] .= "&l_err";
							}
							else
							{
								$_SERVER['HTTP_REFERER'] .= "?l_err";
							}
							header ("Location: ${_SERVER['HTTP_REFERER']}");
							return true;
						}
					}
				}
				else
				{
					$isPresentErrMessage = false;
					$parameters = explode("?", $_SERVER['HTTP_REFERER']);
					$parametersItem = explode("&", $parameters[1]);
					foreach ($parametersItem as $item)
					{
						$temp = explode("=", $item);
						if ($temp[0] == "l_err")
						{
							$isPresentErrMessage = true;
						}
					}
					if ($isPresentErrMessage)
					{
						header ("Location: ${_SERVER['HTTP_REFERER']}");
						return true;
					}
					else
					{
						if (!empty($parameters[1]))
						{
							$_SERVER['HTTP_REFERER'] .= "&l_err";
						}
						else
						{
							$_SERVER['HTTP_REFERER'] .= "?l_err";
						}
						header ("Location: ${_SERVER['HTTP_REFERER']}");
						return true;
					}
				}
			}
			else
			{
				$isPresentErrMessage = false;
				$parameters = explode("?", $_SERVER['HTTP_REFERER']);
				$parametersItem = explode("&", $parameters[1]);
				foreach ($parametersItem as $item)
				{
					$temp = explode("=", $item);
					if ($temp[0] == "l_err")
					{
						$isPresentErrMessage = true;
					}
				}
				if ($isPresentErrMessage)
				{
					header ("Location: ${_SERVER['HTTP_REFERER']}");
					return true;
				}
				else
				{
					if (!empty($parameters[1]))
					{
						$_SERVER['HTTP_REFERER'] .= "&l_err";
					}
					else
					{
						$_SERVER['HTTP_REFERER'] .= "?l_err";
					}
					header ("Location: ${_SERVER['HTTP_REFERER']}");
					return true;
				}
			}
		}
		else
		{
			if ($user->error)
				return $user->error;
			$loginError = 1;
			
			if (LDAP_USE && !empty($_POST['password']))
			{
				// Try to get user from ldap server
				
				if ($user->getLdapUserInfo($_POST['login'], $_POST['password']))
				{
					$ldapUserInfo = $user->data;
	
					$loginError = 0;
					if ($user->getUserInfoByEmailPart($_POST['login'].'@') === false)
						return $user->error;
					$userInfo = $user->data;
	
					if (!empty($userInfo))
					{
						// We have this user.. need to update user ldap login/password
						$data = array(
						'ldap_login'	=> $_POST['login'],
						'ldap_password'	=> md5($_POST['password'])
						);
						if ($user->editUserInfo($data, $userInfo['user_id']) === false)
							return $user->error;
						$user->setUserCredential($_POST['login'], $_POST['password'], $_POST['remember_me']);
						header ("Location: ${_SERVER['HTTP_REFERER']}");
						return true;
	
					}
					else
					{
						// We have a new user.. Should be add
						$data = array(
							'login'				=> $_POST['login'],
							'ldap_login'			=> $_POST['login'],
							'password'			=> md5($_POST['password']),
							'ldap_password'		=> md5($_POST['password']),
							'email'				=> $_POST['login'] . "@" . CORP_EMAIL,
							'status'				=> 1,
							'user_name'			=> $_POST['login'],
							'time'				=> time()
						);
						if ($user->addUser($data) === false)
							return $user->error;
						$user->setUserCredential($_POST['login'], $_POST['password'], $_POST['remember_me']);
						header ("Location: ${_SERVER['HTTP_REFERER']}");
						return true;
					}
				}
			}

			if ($loginError)
			{
				$isPresentErrMessage = false;
				$parameters = explode("?", $_SERVER['HTTP_REFERER']);
				$parametersItem = explode("&", $parameters[1]);

				foreach ($parametersItem as $item)
				{
					$temp = explode("=", $item);
					if ($temp[0] == "l_err")
					{
						$isPresentErrMessage = true;
					}
				}
				if ($isPresentErrMessage)
				{
					header ("Location: ${_SERVER['HTTP_REFERER']}");
					return true;
				}
				else
				{
					if (!empty($parameters[1]))
					{
						$_SERVER['HTTP_REFERER'] .= "&l_err";
					}
					else
					{
						$_SERVER['HTTP_REFERER'] .= "?l_err";
					}
					header ("Location: ${_SERVER['HTTP_REFERER']}");
					return true;
				}
			}
		}
	}
	elseif (isset($_GET['logout']))
	{
		setcookie(C_PASSWORD,NULL,NULL,'/');
		header ("Location: index.php");
		return true;
	}
	else
	{
		header ("Location: index.php");
		return true;
	}
	return false;
}

include ("request_handler.php");
?>
