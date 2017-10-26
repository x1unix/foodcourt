<?php

require_once (DB_FOR_USE."_core.class.php");

class Order
{
	var $data;
	var $error;
	var $db;
	
	function Order()
	{
		$className = DB_FOR_USE.'Core';
		$this->db = new $className;
	}
	
	function getWeekOrderList($weekId, $hashArrayKey = null)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> ORDER_LIST_TBL,
			'where_clause'	=> sprintf("`week_id` = '%s'", mysql_real_escape_string($weekId)),
			'order_clause'	=> "`order_list_id` ASC"
			);
		
		
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
	
	function getDayBlockedStatus ($weekId, $dayId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> ORDER_LIST_TBL,
			'select_clause'	=> "blocked",
			'where_clause'	=> sprintf ("`week_id` = '%s' AND `day_id` = '%s'", mysql_real_escape_string ($weekId), mysql_real_escape_string ($dayId)),
			'limit_clause'	=> "1"
		);
		
				
		if ($this->db->get_data($requestClauses) === true)
		{
			$this->data = $this->db->data[0];
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function getOrderedUsers ($weekId, $dayId = null)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_ORDER_LIST_TBL,
			'select_clause'	=> "user_id",
			'where_clause'	=> sprintf("`week_id` = '%s'", mysql_real_escape_string($weekId)),
			'group_clause'	=> "user_id"
		);
		
		if (isset($dayId))
		{
			$requestClauses['where_clause'] .= sprintf(" AND `day_id` = '%s'", mysql_real_escape_string($dayId));
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
	
	function getPayedUsers($weekId = null, $dayId = null, $hashArrayKey = null)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_ORDER_LIST_TBL,
			'where_clause'	=> "`ordered_amount` > 0"
		);
		
		if (isset($weekId))
		{
			$requestClauses['where_clause'] .= sprintf(" AND `week_id` = '%s'", mysql_real_escape_string($weekId));
		}
		
		if (isset($dayId))
		{
			$requestClauses['where_clause'] .= sprintf(" AND `day_id` = '%s'", mysql_real_escape_string($dayId));
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
	
	function getUserOrdersCount($orderId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_ORDER_LIST_TBL,
			'select_clause'	=> "SUM(ordered_item_count) AS COUNT",
			'where_clause'	=> sprintf("`order_list_id` = '%s'", mysql_real_escape_string($orderId))
		);
		
		
		if ($this->db->get_data($requestClauses) === true)
		{
			$this->data = $this->db->data[0]['COUNT'];
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}

	}
	
	function getPayedUserOrdersCount ($orderId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_ORDER_LIST_TBL,
			'select_clause'	=> "SUM(ordered_item_count) AS COUNT",
			'where_clause'	=> sprintf("`order_list_id` = '%s' AND `ordered_amount` > 0", mysql_real_escape_string($orderId))
		);
		
		if ($this->db->get_data($requestClauses) === true)
		{
			$this->data = $this->db->data[0]['COUNT'];
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function addWeekOrder ($addData)
	{
		$this->data = "";
		
		
		$requestClauses = array (
			'table_name'	=> ORDER_LIST_TBL,
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
	
	function addUserOrder($addData)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_ORDER_LIST_TBL,
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
	
	function editUserOrder ($dataArray, $userOrderId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_ORDER_LIST_TBL,
			'data'		=> $dataArray,
			'where_clause'	=> sprintf("`user_order_id` = '%s'", mysql_real_escape_string($userOrderId))
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
	
	function getUsersOrderList ($weekId, $providerId, $dayId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_ORDER_LIST_TBL,
			'where_clause'	=> sprintf("`week_id` = '%s' AND `provider_id` = '%s' AND `day_id` = '%s'",
									mysql_real_escape_string($weekId),
									mysql_real_escape_string($providerId),
									mysql_real_escape_string($dayId))
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
	
	function getUserOrders($userId, $weekId = null, $hashArrayKey = null, $orderedDay = null)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_ORDER_LIST_TBL,
			'where_clause'	=> sprintf("`user_id` = '%s'", mysql_real_escape_string($userId))
		);
		
		if (!empty($weekId))
		{
			$requestClauses['where_clause'] .= sprintf(" AND `week_id`= '%s'", mysql_real_escape_string($weekId));
		}
		
		if ($orderedDay && is_numeric($orderedDay))
		{
			$requestClauses['where_clause'] .= sprintf(" AND `day_id` = '%s'", mysql_real_escape_string($orderedDay));
		}
		
		if ($this->db->get_data($requestClauses, $hashArrayKey))
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
	
	function deleteUserOrders ($userId, $weekId, $blocked = null)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> USER_ORDER_LIST_TBL,
			'where_clause'	=> sprintf("`user_id` = '%s'", mysql_real_escape_string($userId))
		);
		
		if (isset($weekId))
		{
			$requestClauses['where_clause'] .= sprintf(" AND `week_id` = '%s'", mysql_real_escape_string($weekId));
		}
		
		if (isset($blocked))
		{
			$requestClauses['where_clause'] .= sprintf(" AND NOT EXISTS (SELECT * FROM ".ORDER_LIST_TBL." WHERE ".ORDER_LIST_TBL.".order_list_id = ".USER_ORDER_LIST_TBL.".order_list_id AND ".ORDER_LIST_TBL.".blocked = '%s')",
			mysql_real_escape_string($blocked));
		}
		
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
	
	function editWeekOrderInfo ($dataArray, $orderId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> ORDER_LIST_TBL,
			'data'		=> $dataArray,
			'where_clause'	=> sprintf("`order_list_id` = '%s'", mysql_real_escape_string($orderId))
		);
		
		if ($this->db->edit_data($requestClauses))
		{
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function deleteWeekOrderList ($weekId, $dayId = null, $portionNumber = null, $providerId = null)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> ORDER_LIST_TBL,
			'where_clause'	=> sprintf("`week_id` = '%s'", mysql_real_escape_string($weekId))
		);
		
		if ((!empty($dayId)) && (!empty($portionNumber)) && (!empty($providerId)))
		{
			$requestClauses['where_clause'] .= sprintf(" AND `day_id` = '%s' AND `portion_number` >= '%s' AND `provider_id` = '%s'",
													mysql_real_escape_string($dayId),
													mysql_real_escape_string($portionNumber),
													mysql_real_escape_string($providerId));
		}
		
		if ($this->db->delete_data($requestClauses))
		{
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
		
	}
	
	function deleteWeekOrderInfo ($orderId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> ORDER_LIST_TBL,
			'where_clause'	=> sprintf("`order_list_id` = '%s'", mysql_real_escape_string($orderId))
		);

		if ($this->db->delete_data($requestClauses))
		{
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
		
	}
	
	function formatUserOrder($providerList, $orderIdList, $userData)
	{
	
	
		global $s_monday;
		global $s_tuesday;
		global $s_wednesday;
		global $s_thursday; 
		global $s_friday;
		global $s_saturday; 
		global $s_sunday;
		global $e_unknown_chose;
		global $e_day_blocked;
	
		$userOrder = "";
		$dayList[] = D_MONDAY;
		$dayList[] = D_TUESDAY;
		$dayList[] = D_WEDNESDAY;
		$dayList[] = D_THURSDAY;
		$dayList[] = D_FRIDAY;
		$dayList[] = D_SATURDAY;
		$dayList[] = D_SUNDAY;
		// Read user order day by day...
		
		$userPortionCounter = 0;
		foreach ($dayList AS $dayItem)
		{
			if ($userData['order_'.$dayItem] != "no") // If present some choose
			{
				
				// Try to found choosen provider
				foreach ($providerList AS $providerItem)
				{
					if (PROVIDERS_MULTICHOICE ==1 ) {$providerFooter = "_".$providerItem['provider_id']; } else { $providerFooter = ""; } 
					if ($userData['order_'.$dayItem.$providerFooter] == $providerItem['provider_id'])
					{
						// Choosen provider found. Process order
						// Struct will be $userOrder[day_id]['provider_id'] => value
						// $userOrder[day_id][1|2|3|4] => [order_list_id]

						for ($portion = 1; $portion <= 4; $portion++)
						{
						
							if ($providerItem['multichoice'] != 0)
							{
								$portionCounter = 1;
								foreach ($orderIdList AS $orderIdItem)
								{

									#TODO: Optimization required!
									if (($orderIdItem['provider_id'] == $providerItem['provider_id']) && ($orderIdItem['blocked'] == 0))
									{
										
										if (isset($userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_'.$portionCounter]))
										{
											$userOrder[$dayItem][$portion][$userPortionCounter]['order_list_id'] = $userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_'.$portionCounter];
											$userOrder[$dayItem][ $userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_'.$portionCounter]]['provider_id'] = $providerItem['provider_id'];
											
											if ($providerItem['multiitem'] != 0)
											{
												if (is_numeric($userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_'.$portionCounter.'_count']) && ($userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_'.$portionCounter.'_count'] > 0))
												{
													$userOrder[$dayItem][$portion][$userPortionCounter]['ordered_item_count'] = $userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_'.$portionCounter.'_count'];
												}
												else
												{
													$userOrder[$dayItem][$portion][$userPortionCounter]['ordered_item_count'] = 1;
												}
											}
											else
											{
												$userOrder[$dayItem][$portion][$userPortionCounter]['ordered_item_count'] = 1;
											}
											$userPortionCounter++;
										}
										$portionCounter++;
									}
								}
							}
							else
							{

								if (isset($userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion]))
								{
									// User want some portion.
									// Check that the order ID exist in chosen week..
									if ((isset($orderIdList[$userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion]])) && ($orderIdList['blocked'] == 0))
									{
										if (PROVIDERS_MULTICHOICE ==1 )
										{
											$userOrder[$dayItem][$portion][$userPortionCounter]['order_list_id'] = $userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion];
										}
										else
										{
											$userOrder[$dayItem][$portion][0]['order_list_id'] = $userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion];
										}
										$userOrder[$dayItem][ $userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion]]['provider_id'] = $providerItem['provider_id'];
										
										if ($providerItem['multiitem'] != 0)
										{
											if (is_numeric($userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_count']) && ($userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_count'] > 0))
											{
												if (PROVIDERS_MULTICHOICE ==1 )
												{
													$userOrder[$dayItem][$portion][$userPortionCounter]['ordered_item_count'] = $userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_count'];
													$userPortionCounter++;
												}
												else
												{
													$userOrder[$dayItem][$portion][0]['ordered_item_count'] = $userData['provider_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portion.'_count'];
												}
											}
											else
											{
												if (PROVIDERS_MULTICHOICE ==1 )
												{
													$userOrder[$dayItem][$portion][$userPortionCounter]['ordered_item_count'] = 1;
													$userPortionCounter++;
												}
												else
												{
													$userOrder[$dayItem][$portion][0]['ordered_item_count'] = 1;
												}
											}
										}
										else
										{
											if (PROVIDERS_MULTICHOICE ==1 )
											{
												$userOrder[$dayItem][$portion][$userPortionCounter]['ordered_item_count'] = 1;
												$userPortionCounter++;
											}
											else
											{
												$userOrder[$dayItem][$portion][0]['ordered_item_count'] = 1;
											}
										}
										
									}
									else
									{
										// Hm... some hack? Go fuck! Or no... maybe not a hack.. #TODO: find resolution.
										if (!isset($blockedDayId[$dayItem]))
										{
											switch($dayItem) {
												case D_MONDAY: $blockedDay = $s_monday; break;
												case D_TUESDAY: $blockedDay = $s_tuesday; break;
												case D_WEDNESDAY: $blockedDay = $s_wednesday; break;
												case D_THURSDAY: $blockedDay = $s_thursday; break;
												case D_FRIDAY: $blockedDay = $s_friday; break;
												case D_SATURDAY: $blockedDay = $s_saturday; break;
												case D_SUNDAY: $blockedDay = $s_sunday; break;
											}
											$this->error[] = $e_unknown_chose;
											$this->error[] = $e_day_blocked." ".$blockedDay;
											$blockedDayId[$dayItem] = 1;
										}
									}
								}
							}

						}
					}
				}
			}
		}
		
		return $userOrder;
	}
}
?>