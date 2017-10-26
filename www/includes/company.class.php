<?php
require_once (DB_FOR_USE."_core.class.php");

class Company
{
	var $data;
	var $error;
	var $db;
	
	function Company ()
	{
		$className = DB_FOR_USE.'Core';
		$this->db = new $className;
	}
	
	function getCompanyList ()
	{
		$this->data = "";
		$requestClauses['table_name'] = COMPANIES_TBL;
		
		if ($this->db->get_data($requestClauses))
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
	
	function getCompanyInfo ($companyIdOrName)
	{
		$this->data = "";
		$requestClauses['table_name'] = COMPANIES_TBL;
		
		if (is_numeric($companyIdOrName))
			$requestClauses['where_clause'] = sprintf("comnapy_id = '%d'", mysql_real_escape_string($companyIdOrName));
		else 
			$requestClauses['where_clause'] = sprintf("company_name = '%s'", mysql_real_escape_string($companyIdOrName));
			
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
	
	function addCompany ($addData)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> COMPANIES_TBL,
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
	
}
?>