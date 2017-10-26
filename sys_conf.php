<?php


function load ()
{

	require_once("init.php");
	require_once(LANG.".language.php");
	require_once("user.class.php");
	require_once("sysconfig.class.php");
	require_once("validator.class.php");

	$user = new User;
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	$sysconfig = new SysConfig;
	$validator = new Validator;
	$error = array();

	$type_list = array(
		'text' => Localization::$s_type_text,
		'int' => Localization::$s_type_int,
		'bool' => Localization::$s_type_bool,
		'select' => Localization::$s_type_select,
		'password' => Localization::$s_type_password,
	);

	if (!$isLogged['status']==10)
		return false;

	if (isset($_GET['sys_cat']))
	{
		if ($sysconfig->getCategories() === false)
			return $sysconfig->error;

		$catList = $sysconfig->data;
		include "sys_conf_cat.html";
		return true;
	}
	else if (isset($_GET['add_cat']))
	{
		include "sys_conf_cat_edit.html";
	}
	else if (isset($_GET['edit_cat']))
	{
		if ($sysconfig->getCategories() === false)
			return $sysconfig->error;
		
		$catList = $sysconfig->data;
		$edit_id = mysql_real_escape_string($_GET['edit_cat']);
		$catName = $catList[$_GET['edit_cat']]['name'];
		$catLocalVar = $catList[$_GET['edit_cat']]['localization_var'];

		if ($catName){
			include "sys_conf_cat_edit.html";
			return true;
		}else{
			$error[] = Localization::$e_no_item_exists;
			include "sys_conf_cat.html";
			return true;
		}
	}
	else if (!empty($_GET['delete_cat']))
	{
		if ($validator->idCatexists($_GET['delete_cat']))
		{
			if ($sysconfig->deleteCategory($_GET['delete_cat']))
			{
				$message[] = Localization::$m_data_saved;
			}
			else
			{
				$error = array_merge($error, $sysconfig->error);
			}
		}
		else
		{
			$error[] = Localization::$e_no_item_exists;
		}

		if ($sysconfig->getCategories() === false)
			return $sysconfig->error;

		$catList = $sysconfig->data;
		include "sys_conf_cat.html";
		return true;
	}
	else if (isset($_GET['add']))
	{
		if ($sysconfig->getCategories() === false)
			return $sysconfig->error;

		$catList = $sysconfig->data;
		include "sys_conf_edit.html";
		return true;
	}
	else if (!empty($_GET['edit']))
	{
		if ($sysconfig->getCategories() === false)
			return $sysconfig->error;

		$catList = $sysconfig->data;
		if ($sysconfig->loadConfig($_GET['edit']) === false)
			return $sysconfig->error;
		
		$sysConfig = array_values($sysconfig->data);
		$sysConfig = $sysConfig[0];

		if ($sysConfig)
		{
			$edit_id = $_GET['edit'];
			include "sys_conf_edit.html";
			return true;
		}
		else
		{
			$error[] = Localization::$e_no_item_exists;
			if ($sysconfig->loadConfig() === false)
				return $sysconfig->error;
			
			$configList = $sysconfig->data;
			include "sys_conf_const.html";
			return true;
		}

	}
	else if (!empty($_GET['delete']))
	{
		if ($validator->idConfigexists($_GET['delete']))
		{
			if ($sysconfig->deleteConfig($_GET['delete']) === false)
				return $sysconfig->error;

			$message[] = $m_data_saved;
		}
		else
		{
			$error[] = Localization::$e_no_item_exists;
		}

		if ($sysconfig->loadConfig() === false)
			return $sysconfig->error;

		$configList = $sysconfig->data;
		include "sys_conf_const.html";
		return true;
	}
	else if (isset($_POST['save_cat']))
	{
		if (!$_POST['name'])
		{
			$error[] = Localization::$e_fill_required_fields;
			$edit_id = $_POST['edit_id'];
			$catName = $_POST['name'];
			$catLocalVar = $_POST['localization_var'];
			include "sys_conf_cat_edit.html";
			return true;
		}
		else
		{
			if ($id = $_POST['edit_id']){
				if ($validator->idCatexists($id))
				{
					if ($validator->validateCategoryName($_POST['name'], $id))
					{
						if ($sysconfig->editCategory($id, $_POST['name'], $_POST['localization_var']))
						{
							$message[] = Localization::$m_data_saved;
							if ($sysconfig->getCategories() === false)
								return $sysconfig->error;
							
							$catList = $sysconfig->data;
							include "sys_conf_cat.html";
							return true;
						}
						else
						{
							$error = array_merge($error, $sysconfig->error);
							$catName = $_POST['name'];
							$catLocalVar = $_POST['localization_var'];
							$edit_id = $id;
							include "sys_conf_cat_edit.html";
							return true;
						}
					}
					else
					{
						$error[] = Localization::$e_variable_exists;
						$catName = $_POST['name'];
						$catLocalVar = $_POST['localization_var'];
						$edit_id = $id;
						include "sys_conf_cat_edit.html";
						return true;
					}
				}
				else
				{
					$error[] = Localization::$e_no_item_exists;
					if ($sysconfig->getCategories() === false)
						return $sysconfig->error;
					
					$catList = $sysconfig->data;
					include "sys_conf_cat.html";
					return true;
				}
			}
			else
			{
				if ($validator->validateCategoryName($_POST['name']))
				{
					if($sysconfig->addCategory($_POST['name'],$_POST['localization_var']))
					{
						$message[] = Localization::$m_data_saved;
						if ($sysconfig->getCategories() === false)
							return $sysconfig->error;
						
						$catList = $sysconfig->data;
						include "sys_conf_cat.html";
						return true;
					}
					else
					{
						$error = array_merge($error, $sysconfig->error);
						$catName = $_POST['name'];
						$catLocalVar = $_POST['localization_var'];
						include "sys_conf_cat_edit.html";
						return true;
					}
				}
				else
				{
					$error[] = Localization::$e_item_exists;
					$catName = $_POST['name'];
					$catLocalVar = $_POST['localization_var'];
					include "sys_conf_cat_edit.html";
					return true;
				}
			}
		}
	}
	else if (isset($_POST['save_cfg']))
	{
		if ($sysconfig->getCategories() === false)
			return $sysconfig->error;

		$catList = $sysconfig->data;
		if (!$_POST['name'])
		{
			$error[] = Localization::$e_fill_required_fields;
			$sysConfig = $_POST;
			$edit_id = $_POST['edit_id'];
			include "sys_conf_edit.html";
			return true;
		}
		elseif (!$catList[$_POST['cat_selected']] && $_POST['cat_selected']!=0)
		{
			$error[] = Localization::$e_invalid_category;
			$sysConfig = $_POST;
			$edit_id = $_POST['edit_id'];
			include "sys_conf_edit.html";
			return true;
		}
		elseif (!in_array($_POST['type'], array_keys($type_list)))
		{
			$error[] = Localization::$e_unknown_type;
			$sysConfig = $_POST;
			$edit_id = $_POST['edit_id'];
			include "sys_conf_edit.html";
			return true;
		}
		elseif (!$validator->validateType($_POST['value'], $_POST['type']))
		{
			$error[] = Localization::$e_invalid_value;
			$sysConfig = $_POST;
			$edit_id = $_POST['edit_id'];
			include "sys_conf_edit.html";
			return true;
		}
		else
		{
			$config = array (
				'name' => $_POST['name'],
				'value' => $_POST['value'],
				'comment' => $_POST['comment'],
				'localization_var' => $_POST['localization_var'],
				'cat_selected' => $_POST['cat_selected'],
				'active' => $_POST['active'],
				'type' => $_POST['type']
			);
			if ($id = $_POST['edit_id'])
			{
				if ($validator->idConfigexists($id))
				{
					if ($validator->validateConfigName($config['name'], $id))
					{
						if ($sysconfig->editConfig($id, $config) === false)
							return $sysconfig->error;

						$message[] = Localization::$m_data_saved;
						if ($sysconfig->loadConfig() === false)
							return $sysconfig->error;
							
						$configList = $sysconfig->data;
						include "sys_conf_const.html";
						return true;
					}
					else
					{
						$error[] = Localization::$e_variable_exists;
						$edit_id = $id;
						$sysConfig = $_POST;
						include "sys_conf_edit.html";
						return true;
					}
				}
				else
				{
					$error[] = Localization::$e_no_item_exists;
					$edit_id = $id;
					$sysConfig = $_POST;
					include "sys_conf_edit.html";
					return true;
				}
			}
			else
			{
				if ($validator->validateConfigName($config['name']))
				{
					if ($sysconfig->saveConfig($config) === false)
						return $sysconfig->error;

						$message[] = Localization::$m_data_saved;
						if ($sysconfig->loadConfig() === false)
							return $sysconfig->error;
						
						$configList = $sysconfig->data;
						include "sys_conf_const.html";
						return true;
				}
				else
				{
					$error[] = Localization::$e_variable_exists;
					$sysConfig = $_POST;
					include "sys_conf_edit.html";
					return true;
				}
			}
		}

	}
	else if (isset($_POST['save_values']))
	{
		if ($sysconfig->getCategories() === false)
			return $sysconfig->error;
		
		$catList = $sysconfig->data;
		foreach ($_POST['varArray'] as $id=>$array)
		{
			if (!$validator->validateType($array['value'], $array['type']))
			{
				$wrong_type = $array['value'];
			}
		}
		if ($wrong_type)
		{
			$error[] = $e_invalid_value . ": $wrong_type";
			$wrong_config = $_POST['varArray'];
		}
		else
		{
			if ($sysconfig->updateVariables($_POST['varArray']) === false)
				return $sysconfig->error;
				
			$message[] = Localization::$m_data_saved;
		}
		
		if ($wrong_config)
		{	
			if ($sysconfig->loadConfig() === false)
				return $sysconfig->error;
			
			$configList = $sysconfig->data;
			# array_merge_recursive does not work properly here
			foreach ($configList as $id=>$config)
			{
				if ($wrong_config[$id])
				{
					foreach ($config as $key=>$value)
					{
						if ($wrong_config[$id][$key])
						{
							$configList[$id][$key] = $wrong_config[$id][$key];
						}
					}
				}
			}
			# $configList = array_merge($configList, $wrong_config);
			if ($sysconfig->getCategories() === false)
				return $sysconfig->error;
			
			$catList = $sysconfig->data;
			include "sys_conf_const.html";
			return true;
		}
		else
		{
			header ("Location: ${_SERVER['HTTP_REFERER']}");
			return true;
		}
	}
	else
	{
		if ($sysconfig->getCategories() === false)
			return $sysconfig->error;
		
		$catList = $sysconfig->data;
		if ($sysconfig->loadConfig() === false)
			return $sysconfig->error;
		
		$configList = $sysconfig->data;
		include "sys_conf_const.html";
		return true;

	}
		
	return false;
}

	include ("request_handler.php");
?>
