<?php
require_once (DB_FOR_USE."_core.class.php");

class Provider
{
	var $data;
	var $error;
	var $db;
	
	function Provider ()
	{
		$className = DB_FOR_USE.'Core';
		$this->db = new $className;
	}
	
	function getProviderInfo ($nameOrId)
	{
		$this->data = "";
		
		$requestClauses['table_name'] = PROVIDER_TBL;

		if (is_numeric($nameOrId))
		{
			$requestClauses['where_clause'] = sprintf("`provider_id` = '%s'", mysql_real_escape_string($nameOrId));
		}
		else
		{
			$requestClauses['where_clause'] =  sprintf("`name` = '%s'", mysql_real_escape_string($nameOrId));
		}
		
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
	
	function getProviderList($hashArrayKey = null)
	{
		$this->data = "";
		
		$requestClauses['table_name'] = PROVIDER_TBL;
		
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
	
	function addProvider ($addData)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> PROVIDER_TBL,
			'data'		=> $addData
		);
		
		if ($this->db->add_data($requestClauses) === true)
		{
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function editProviderInfo ($dataArray, $providerId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> PROVIDER_TBL,
			'data'		=> $dataArray,
			'where_clause'	=> sprintf("`provider_id` = '%s'", mysql_real_escape_string($providerId))
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
	
	function deleteProvider ($providerId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> PROVIDER_TBL,
			'where_clause'	=> sprintf("`provider_id` = '%s'", mysql_real_escape_string($providerId))
		);
		
		if ($this->db->delete_data($requestClauses) === true)
		{
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
		
	}
}
?>