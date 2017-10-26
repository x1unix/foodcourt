<?php
	// Developed by kosyak <kosyak_ua@yahoo.com>
	
	//// For 'forgot password' functionality.
	
	
function load ($_GET, $_POST)
{
	require_once("init.php");
	require_once("user.class.php");
	require_once ("messanger.class.php");
	require_once(LANG.".language.php");
	
	$scriptName = "forgot_password.php";
	$user = new User;
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	
	$forgotInfo['forgot_password_login'] = $_COOKIE[C_LOGIN];
	
	if ($isLogged)
		return false;
	if (isset($_POST['forgot_password']))
	{
		$forgotInfo['forgot_password_login'] = $_POST['forgot_password_login'];
		if (!empty($forgotInfo['forgot_password_login']))
		{
			if (!is_numeric($forgotInfo['forgot_password_login']))
			{
				if ($user->getUserInfo($forgotInfo['forgot_password_login']) === false)
					return $user->error;
				if (!empty($user->data))
				{
					$userInfo = $user->data;
					$fpKey = md5(time());
					if ($user->setFpKey($userInfo['user_id'], $fpKey) === false)
						return $user->error;
					$messanger = new Messanger();
					if (!empty($userInfo['email']))
					{
						$to = $userInfo['email'];
						$subject = PROJECT_NAME."| Forgot password";
						$body = Localization::$m_asked_for_new_password;
						$body .= "<br>";
						$body .= Localization::$m_if_it_was_you_click_below;
						$body .="<br>";
						$body .= SYSTEM_HTTP_ADDRESS.$scriptName."?id=".$userInfo['user_id']."&fp_key=".$fpKey;
						$body .="<br>";
						$body .= Localization::$m_otherwise_ignore_this_email;
						
						if ($messanger->sendHtmlEMail($to, $subject, $body))
						{
							$message[] = Localization::$m_check_your_email_for_fp_key;
							include ("home.html");
							return true;
						}
						else
						{
							$error[] = Localization::$e_mail_server_error_try_later;
							include ("home.html");
							return true;
						}
					}
					else
					{
						$error[] = Localization::$e_no_email_ask_admin;
						include ("home.html");
						return true;
					}
				}
				else
				{
					$error[] = Localization::$e_username_incorrect_or_not_exist;
					include ("home.html");
					return true;
				}
			}
			else
			{
				$error[] = Localization::$e_username_incorrect_or_not_exist;
				include ("home.html");
				return true;
			}
		}
		else
		{
			$error[] = Localization::$e_fill_required_fields;
			include ("home.html");
			return true;
		}
	}
	else if (isset($_GET['id']))
	{
		if ($user->getUserInfo($_GET['id']) === false)
			return $user->error;

		$userInfo = $user->data;
		if (!empty($userInfo))
		{
			if (!empty($userInfo['fp_key']))
			{
				if ($userInfo['fp_key'] === $_GET['fp_key'])
				{
					$newPass = $user->genNewPass(6);
					$data['password'] = md5($newPass);
					if ($user->editUserInfo($data, $userInfo['user_id']) === false)
						return $user->error;
					if (!empty($userInfo['email']))
					{
						$messanger = new Messanger;
						$to = $userInfo['email'];
						$subject = PROJECT_NAME." | Forgot Password";
						$body = Localization::$m_your_password_changed."<br>";
						$body .= Localization::$m_new_password_is.": $newPass <br>";
						$body .= SYSTEM_HTTP_ADDRESS;
						if ($messanger->sendHtmlEMail($to, $subject, $body))
						{
							$message[] = Localization::$m_check_your_email_for_new_password;
							include ("home.html");
							return true;
						}
						else
						{
							$error[] = Localization::$e_mail_server_error_try_later;
							include ("home.html");
							return true;
						}
					}
					else
					{
						$error[] = Localization::$e_no_email_ask_admin;
						include ("home.html");
						return true;
					}
					
				}
				else
				{
					$error[] = Localization::$e_restore_password_link_is_invalid;
					include ("home.html");
					return true;
				}
			}
			else
			{
				$error[] = Localization::$e_restore_password_link_is_invalid;
				include ("home.html");
				return true;
			}
		}
		else
		{
			$error[] = Localization::$e_restore_password_link_is_invalid;
			include ("home.html");
			return true;
		}
	}
	else
	{
		$forgotInfo['forgot_password_login'] = $_COOKIE[C_LOGIN];
		include ("home.html");
		return true;
	}
	
	return false;
}

include ("request_handler.php");
?>