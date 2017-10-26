<?php
require_once (DB_FOR_USE."_core.class.php");

class User
{
	var $data;
	var $error;
	var $db;
	
	function User()
	{
		$className = DB_FOR_USE.'Core';
		$this->db = new $className;
	}
	
	function getUserMarkRules ($rulesData)
	{
		$this->data = "";
		
		if (isset($rulesData[C_SHOW_ORDERED_FOR]))
		{
			if (($rulesData[C_SHOW_ORDERED_FOR] > 0) && ($rulesData[C_SHOW_ORDERED_FOR] < 8))
			{
				$this->data['ordered_for'] = $rulesData[C_SHOW_ORDERED_FOR];
			}
			else
			{
				$this->data['ordered_for'] = date("N", time());
			}
		}
		else
		{
			$this->data['ordered_for'] = date("N", time());
		}
		
		if (isset($rulesData[C_SHOW_PAYED_FOR]))
		{
			if (($rulesData[C_SHOW_PAYED_FOR] > 0) && ($rulesData[C_SHOW_PAYED_FOR] < 8))
			{
				$this->data['payed_for'] = $rulesData[C_SHOW_PAYED_FOR];
			}
			else
			{
				$this->data['payed_for'] = date("N", time());
			}
		}
		else
		{
			$this->data['payed_for'] = date("N", time());
		}
		
		return $this->data;
	}
	
	function setUserMarkRules ($rulesData)
	{
		if (isset($rulesData['ordered_for']))
		{
			setcookie(C_SHOW_ORDERED_FOR,$rulesData['ordered_for'],time()+365*24*3600,'/');
		}
		else
		{
			setcookie(C_SHOW_ORDERED_FOR,SHOW_ORDERED_DATA_FOR,time()+365*24*3600,'/');
		}
		
		if (isset($rulesData['payed_for']))
		{
			setcookie(C_SHOW_PAYED_FOR,$rulesData['payed_for'],time()+365*24*3600,'/');
		}
		else
		{
			setcookie(C_SHOW_PAYED_FOR,SHOW_PAYED_FOR,time()+365*24*3600,'/');
		}
		
		return true;
	}
	
	function isLogged($login, $password)
	{
		$this->data = "";
		if (!empty($password))
		{
			if ($this->getUserInfo($login))
			{
				if ($this->data['activity'] != 0)
				{
					if ($this->genCookies($this->data['password']) == $password)
					{
						return $this->data;
					}
					else if (LDAP_USE)
					{
						if ($this->genCookies($this->data['ldap_password']) == $password)
						{
							return ($this->data);
						}
						else
						{
							return false;
						}
					}
					else
					{
						return false;
					}
				}
				else
				{
					return false;
				}
			}
			else
			{
				return false;
			}
		}
		else
		{
			return false;
		}
	}
	
	function setUserCredential ($userLogin, $userPassword, $savePassword = null)
	{
		setcookie(C_LOGIN,$userLogin,time()+365*24*3600,'/');
		if (!empty($savePassword))
		{
			setcookie(C_PASSWORD,$this->genCookies(md5($userPassword)),time()+365*24*3600,'/');
		}
		else
		{
			setcookie(C_PASSWORD,$this->genCookies(md5($userPassword)),NULL,'/');
		}
		return true;
	}
	
	function getUserInfo ($loginOrId)
	{
		$this->data = "";
		
		$requestClauses['table_name'] = USER_TBL;
		
		if (is_numeric($loginOrId))
		{
			$requestClauses['where_clause'] = sprintf("`user_id` = '%s'", mysql_real_escape_string($loginOrId));
		}
		else
		{
			$requestClauses['where_clause'] = sprintf("`login` = '%s'", mysql_real_escape_string($loginOrId));
		}
		
		if ($this->db->get_data($requestClauses) === true)
		{
			$this->data = $this->db->data[0];
			
			if (empty($this->data) && LDAP_USE && (!is_numeric($loginOrId)))
			{
				// Try to found ldap user.
				$requestClauses['where_clause'] = sprintf("`ldap_login` = '%s'", mysql_real_escape_string($loginOrId));
				
				if ($this->db->get_data($requestClauses) === true)
				{
					$this->data = $this->db->data[0];
					return $this->data;
				}
				else
				{
					$this->error = $this->db->error;
					return false;
				}
			}
			return $this->data;
		}
		else
		{
			
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function getUserInfoByEmailPart ($emailPart)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_TBL,
			'where_clause'	=> sprintf("`email` LIKE '%s"."%%'", mysql_real_escape_string($emailPart)),
		);
		

		if ($this->db->get_data($requestClauses) === true)
		{
			$this->data = $this->db->data[0];
			return $this->data;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function getLdapUserInfo ($login, $password)
	{
		$this->data = "";
		$ldapconn = ldap_connect(LDAP_SERVER, LDAP_PORT);
		if ($ldapconn)
		{
			if ($ldapbind = @ldap_bind($ldapconn, "$login"."@".CORP_EMAIL, "$password"))
			{
				$fields = array("asfname", "assname","samaccountname", "useraccountcontrol");
				if ($search_result = @ldap_search($ldapconn, LDAP_DN, "samaccountname=".$login, $fields))
				{
					$this->data = ldap_get_entries($ldapconn, $search_result);
					ldap_unbind($ldapconn);
					return $this->data;
				}
				else
				{
					return false;
				}
			}
			else
			{
				return false;
			}
		}
		return false;
	}
	
	function getUserReplacements($userId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USERS_REPLACEMENT_TBL,
			'where_clause'	=> sprintf("`user_id` = '%s'", mysql_real_escape_string($userId))
		);

		
		if ($this->db->get_data($requestClauses) === true)
		{
			$this->data = $this->db->data;
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function getReplacementingUser($userId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USERS_REPLACEMENT_TBL,
			'where_clause'	=> sprintf("`replacement_id` = '%s'", mysql_real_escape_string($userId))
		);

		if ($this->db->get_data($requestClauses) === true)
		{
			$this->data = $this->db->data;
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function deleteReplecement($userId, $replacementId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USERS_REPLACEMENT_TBL,
			'where_clause'	=> sprintf("`user_id` = '%s' AND `replacement_id` = '%s'", mysql_real_escape_string($userId),
								mysql_real_escape_string($replacementId))
		);
		
		if ($this->db->delete_data($requestClauses) === true)
		{
			$this->data = $this->db->data;
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function getUserList($orderItem = null, $hashArrayKey = null)
	{
		$this->data = "";
		
		$requestClauses['table_name'] = USER_TBL;
		
		if (!empty($orderItem))
		{
			$requestClauses['order_clause'] = "$orderItem ASC";
		}
		
		if ($this->db->get_data($requestClauses, $hashArrayKey) === true)
		{
			$this->data = $this->db->data;
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}

	}
	
	function addUser ($addData)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_TBL,
			'data'		=> $addData
		);
		
		
		if ($this->db->add_data($requestClauses) === true)
		{
			$this->data = $this->db->data;
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
		
	}
	
	function editUserInfo ($editData, $userId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_TBL,
			'data'		=> $editData,
			'where_clause'	=> sprintf("`user_id` = '%s'", mysql_real_escape_string($userId))
		);
		
		if ($this->db->edit_data($requestClauses) === true)
		{
			$this->data = $this->db->data;
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function setFpKey ($userId, $fpKey)
	{
		$this->data = "";
		
		$data['fp_key'] = $fpKey;
		
		$requestClauses = array (
			'table_name'	=> USER_TBL,
			'data'		=> $data,
			'where_clause'	=> "`user_id` = '$userId'"
		);
		
		if ($this->db->edit_data($requestClauses) === true)
		{
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function deleteUser ($userId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_TBL,
			'where_clause'	=> sprintf("`user_id` = '%s'", mysql_real_escape_string($userId))
		);

		if ($this->db->delete_data($requestClauses) === true)
		{
			$this->data = $this->db->data;
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
		
	}
	
	function addUserReplacement($addData)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USERS_REPLACEMENT_TBL,
			'data'		=> $addData
		);
		
		if ($this->db->add_data($requestClauses) === true)
		{
			$this->data = $this->db->data;
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function genNewPass ($passSize = 6)
	{
		
		// start with a blank password
		$password = "";

		// define possible characters
		$possible = "0123456789abcdefghjklmnopqrstquvwxyz"; 
		
		// set up a counter
		$i = 0; 
		
		// add random characters to $password until $length is reached
		while ($i < $passSize) { 

		// pick a random character from the possible ones
		$char = substr($possible, mt_rand(0, strlen($possible)-1), 1);
		
		// we don't want this character if it's already in the password
		if (!strstr($password, $char)) { 
			$password .= $char;
			$i++;
		}

		}

		// done!
		return $password;

	}
	
	function genCookies ($passString)
	{
		return sha1($_ENV['HTTP_USER_AGENT '] . $passString);
	}
}
?>