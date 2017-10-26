<?php
require_once (DB_FOR_USE."_core.class.php");

class Cash
{
	var $data;
	var $error;
	var $db;
	
	function Cash()
	{
		$className = DB_FOR_USE.'Core';
		$this->db = new $className;
	}
	
	function addTranzaction ($addData)
	{
		$this->data = "";
		
		$requestClauses = array(
			'table_name'	=> CASH_TRANZACTION_TBL,
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