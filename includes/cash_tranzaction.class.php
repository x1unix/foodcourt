<?php
	
	require_once (DB_FOR_USE."_core.class.php");
	
	class CashTranzaction
	{
		var $data;
		var $error;
		var $db;
		
		function CashTranzaction()
		{
			$className = DB_FOR_USE.'Core';
			$this->db = new $className;
		}
		
		function getTranzactionList($userId = null, $itemPerPage = null, $pageNum = null)
		{
			$this->data = "";
			
			$requestClauses['table_name'] = CASH_TRANZACTION_TBL;
			if (is_numeric($userId))
			{
				$requestClauses['where_clause'] = sprintf("`tranzaction_user_id_from` = '%s' OR `tranzaction_user_id_to` = '%s'",
													mysql_real_escape_string($userId),
													mysql_real_escape_string($userId));
			}
			$requestClauses['order_clause'] .= "`tranzaction_id` ASC";
			
			if (is_numeric($itemPerPage) && is_numeric($pageNum))
			{
				if (($itemPerPage > 0) && ($pageNum > 0))
				{
					$itemFrom = ($pageNum-1)*$itemPerPage;
					$requestClauses['limit_clause'] = "$itemFrom,$itemPerPage";
				}
				else
				{
					$requestClauses['limit_clause'] = "0,".ITEM_PER_PAGE;
				}
			}

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
		
		function getTranzactionsCount ($userId = null)
		{
			$this->data = "";
			
			$requestClauses = array (
				"table_name"		=> CASH_TRANZACTION_TBL,
				"select_clause"	=> "COUNT(*) AS count"
			);
			
			if (!empty($userId))
			{
				$requestClauses['where_clause'] = sprintf("`tranzaction_user_id_from` = '%s' OR `tranzaction_user_id_to` = '%s'",
													mysql_real_escape_string($userId),
													mysql_real_escape_string($userId));
			}
			
			if ($this->db->get_data($requestClauses))
			{
				$this->data = $this->db->data[0]['count'];
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