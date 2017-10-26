<?php
	// coded by kosyak <kosyak_ua@yahoo.com>
	
	//// User info add/remove/edit page script.
	
	
function load ($_GET, $_POST)
{
	require_once("init.php");
	require_once("user.class.php");
	require_once(LANG.".language.php");
	
	$user = new User;
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	
	if (!$isLogged)
		return false;
	
	if (isset($_GET['s_ok']))
	{
		$message[] = Localization::$m_data_saved;
	}
	if ($user->getUserInfo($isLogged['user_id']) === false)
		return $user->error;

	$userInfo = $user->data;
	
	if ($user->getUserReplacements($isLogged['user_id']) === false)
		return $user->error;

	$userReplacementList = $user->data;
	if (!empty($userReplacementList))
	{
		foreach ($userReplacementList AS $userReplacement)
		{
			if ($user->getUserInfo($userReplacement['replacement_id']) === false)
				return $user->error;
			$userReplacementIdInfo[$user->data['user_id']] = $user->data;
		}
	}

		
	if ($user->getReplacementingUser($isLogged['user_id']) === false)
		return $user->error;

	$replacementingUserList = $user->data;
	if (!empty($replacementingUserList))
	{
		foreach ($replacementingUserList AS $replacementingUser)
		{
			if ($user->getUserInfo($replacementingUser['user_id']) === false)
				return $user->error;
			$replacementingUserIdInfo[$user->data['user_id']] = $user->data;
		}
	}


	if (isset($_POST['do_save']))
	{
		if (!empty($_POST['old_password']))
		{
			if (md5($_POST['old_password']) != $userInfo['password'])
			{
				$error[] = Localization::$e_bad_old_password;
			}
			if (!$error)
			{
				if (empty($_POST['new_password']) || empty($_POST['new_password_verify']))
				{
					$message[] = Localization::$m_fill_all_password_fields;
					include ("profile_edit.html");
					return true;
				}
				else
				{
					if ($_POST['new_password'] == $_POST['new_password_verify'])
					{
						$data['password'] = md5($_POST['new_password']);
						if ($user->editUserInfo($data, $isLogged['user_id']) === false)
							return $user->error;

						setcookie(C_PASSWORD,$user->genCookies(md5($_POST['new_password'])),NULL,'/');
						$message[] = Localization::$m_data_saved;
						include ("profile_edit.html");
						return true;
					}
					else
					{
						$error[] = Localization::$e_passwords_not_eq;
						include ("profile_edit.html");
						return true;
					}
				}
			}
			else
			{
				include ("profile_edit.html");
				return true;
			}
		}
		else
		{
			if (!empty($_POST['new_password']) || !empty($_POST['new_password_verify']))
			{
				$message[] = Localization::$m_fill_all_password_fields;
			}
			else
			{
				$message[] = Localization::$m_nothing_to_change;
			}
			include ("profile_edit.html");
			return true;
		}
	}
	else if (isset ($_POST['do_add_replacement']))
	{
		if (!empty($_POST['replacement_name']) && !$error)
		{
			$replacement_name = $_POST['replacement_name'];
			if (!is_numeric($replacement_name))
			{
				if ($user->getUserInfo($replacement_name) === false)
					return $user->error;

				if ($user->data)
				{
					$replacementUserInfo = $user->data;
					if ($replacementUserInfo['login'] != $isLogged['login'])
					{
						if (!isset($userReplacementIdInfo[$replacementUserInfo['user_id']]))
						{
							$addData = array(
								'user_id'			=> $isLogged['user_id'],
								'replacement_id'	=> $replacementUserInfo['user_id'],
								'time'				=> time()
							);
							if ($user->addUserReplacement($addData) === false)
								return $user->error;

							$userReplacementIdInfo[$replacementUserInfo['user_id']] = $replacementUserInfo;
							$userReplacementList[] = array
							(
								'user_replacement_id'	=> mysql_insert_id(),
								'user_id'			=> $isLogged['user_id'],
								'replacement_id'	=> $replacementUserInfo['user_id'],
								'time'				=> $addData['gime']
							);
							$message[] = Localization::$m_data_saved;
							include ("replacement_edit.html");
							return true;
						}
						else
						{
							$error[] = Localization::$e_replacement_exist;
							include ("replacement_edit.html");
							return true;
						}
					}
					else
					{
						$error[] = Localization::$e_bad_idea_to_add_yourself;
						include ("replacement_edit.html");
						return true;
					}
				}
				else
				{
					$error[] = Localization::$e_no_user_exist;
					include ("replacement_edit.html");
					return true;
				}
			}
			else
			{
				$error[] = Localization::$e_no_user_exist;
				include ("replacement_edit.html");
				return true;
			}
		}
		else
		{
			$error[] = Localization::$e_fill_required_fields;
			include ("replacement_edit.html");
			return true;
		}

	}
	else if (isset($_GET['delete_replacementing']))
	{
		if ($user->deleteReplecement($_GET['delete_replacementing'], $isLogged['user_id']) === false)
			return $user->error;
		header("Location: profile.php?replacements");
		return true;
	}
	else if (isset($_GET['delete_replacement']))
	{
		if ($user->deleteReplecement($isLogged['user_id'], $_GET['delete_replacement']) === false)
			return $user->error;
		header("Location: profile.php?replacements");
		return true;
	}
	else
	{
		if ($user->getUserInfo($isLogged['user_id']) === false)
			return $user->error;

		$userInfo = $user->data;
		if (isset($_GET['replacements']))
		{
			include ("replacement_edit.html");
			return true;
		}
		else
		{
			include ("profile_edit.html");
			return true;
		}
	}

	return false;
}

include ("request_handler.php");
?>