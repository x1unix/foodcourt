<?php
	class mysqlCore
	{
		var $data;
		var $error;
		
		function MysqlCore ()
		{}
		
		function connect($dbHost, $dbUser, $dbPass, $dbName)
		{
			if (empty($dbHost) || empty($dbUser) || empty($dbPass) || empty($dbName))
			{
				return false;
			}
			else
			{
				if (mysql_connect($dbHost,$dbUser,$dbPass))
				{
					if (!mysql_select_db($dbName))
					{
						return false;
					}
					
					return true;
				}
				else
				{
					return false;
				}
			}
		}
		
		function add_data($requestClauses)
		{
			$tableName = $requestClauses['table_name'];
			$addData = $requestClauses['data'];
			$this->data = "";
			if (!empty($addData))
			{
				$query = "INSERT INTO ".$tableName." (";
				$isFirst = true;
				foreach ($addData AS $key=>$value)
				{
					if ($isFirst)
					{
						$query .= sprintf("%s",
						mysql_real_escape_string($key));
						$isFirst = false;
					}
					else
					{

						$query .= sprintf(", %s",
						mysql_real_escape_string($key));
					}
				}
				$query .= ") VALUES (";
				$isFirst = true;
				foreach ($addData AS $key=>$value)
				{
					if ($isFirst)
					{
						$query .= sprintf("'%s'",
						mysql_real_escape_string($value));
						$isFirst = false;
					}
					else
					{
						$query .= sprintf(", '%s'",
						mysql_real_escape_string($value));
					}
				}
				$query .= ")";
				if (mysql_query($query))
				{
					return true;
				}
				else
				{
					$this->error[] = "add_data: sql error: ".mysql_error();
					return false;
				}
			}
			else
			{
				$this->error[] = "add_data: Empty data to add";
				return false;
			}
		}
		
		function edit_data ($requestClauses)
		{
			$this->data = "";
			$tableName = $requestClauses['table_name'];
			$editData = $requestClauses['data'];
			$whereClause = $requestClauses['where_clause'];
			
			$query = "UPDATE ".$tableName." SET ";
			$isFirst = true;
			foreach ($editData AS $key=>$value)
			{
				if ($isFirst)
				{
					$query .= sprintf("%s = '%s'",
					mysql_real_escape_string($key),
					mysql_real_escape_string($value));
					$isFirst = false;
				}
				else
				{
					$query .= sprintf(", %s = '%s'",
					mysql_real_escape_string($key),
					mysql_real_escape_string($value));
				}
			}
			
			if (!empty($whereClause))
			{
				$query .= " WHERE $whereClause";
			}
			
			if (mysql_query($query))
			{
				return true;
			}
			else
			{
				$this->error[] = "edit_data: sql error: ".mysql_error();
				return false;
			}
		}
		
		function delete_data ($requestClauses)
		{
			
			$this->data = "";
			
			$tableName = $requestClauses['table_name'];
			$whereClause = $requestClauses['where_clause'];
			
			$query = "DELETE FROM `".$tableName."`";
			
			if (!empty($whereClause))
			{
				$query .= " WHERE $whereClause";
			}
			if (mysql_query($query))
			{
				return true;
			}
			else
			{
				$this->error[] = "delete_data: sql error: ".mysql_error();
				return false;
			}
		}
		
		function get_data ($requestClauses, $hashArrayKey = null)
		{
			
			$this->data = "";

			$tableName		= $requestClauses['table_name'];
			$whereClause	= $requestClauses['where_clause'];
			$orderClause	= $requestClauses['order_clause'];
			$groupClause	= $requestClauses['group_clause'];
			$limitClause		= $requestClauses['limit_clause'];
			if (!empty($requestClauses['select_clause']))
			{
				$selectClause = $requestClauses['select_clause'];
			}
			else
			{
				$selectClause = "*";
			}
			
			
			$query = "SELECT $selectClause FROM ".$tableName;
			
			if (!empty($whereClause))
			{
				$query .= " WHERE $whereClause";
			}
			
			if (!empty($groupClause))
			{
				$query .= " GROUP BY $groupClause";
			}
			
			if (!empty($orderClause))
			{
				$query .= " ORDER BY $orderClause";
			}

			if (!empty($limitClause))
			{
				$query .= " LIMIT $limitClause";
			}
			
			if ($res = mysql_query($query))
			{
				while ($row = mysql_fetch_assoc($res))
				{
					if (!empty($hashArrayKey))
					{
						$this->data[$row[$hashArrayKey]] = $row;
					}
					else
					{
						$this->data[] = $row;
					}
				}
				return true;
			}
			else
			{
				$this->error[] = "get_data: sql error: ".mysql_error();
				return false;
			}
		}
	}
?>
