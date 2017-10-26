<?php
	// coded by kosyak <kosyak_ua@yahoo.com>
	
	//// Order's info add/remove/edit page script.


function load ($_GET, $_POST)
{
	require_once("init.php");
	require_once(LANG.".language.php");
	require_once("user.class.php");
	require_once("week.class.php");
	require_once("provider.class.php");
	require_once("order.class.php");

	
	$user = new User;
	$scriptName = "orders.php";
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	
	if (isset($_GET['s_ok']))
	{
		$message[] = Localization::$m_data_saved;
	}
	
	if (!$isLogged)
		return false;
		
	if ($user->getReplacementingUser($isLogged['user_id']) === false)
		return $user->error;
		
	$replacementingUserList = $user->data;
	if (!empty($replacementingUserList))
	{
		foreach ($replacementingUserList AS $replacementingUser)
		{
			if ($user->getUserInfo($replacementingUser['user_id']) !== false)
			{
				$replacementingUserIdInfo[$user->data['user_id']] = $user->data;
			}
			else
			{
				$error = $user->error;
				return $error;
			}
		}
	}
	

	if (isset($_POST['do_order_save']))
	{
		if (!empty($_POST['order_id']))
		{
			$orderId = $_POST['order_id'];
			$needOrderForId = $orderId;
			if (empty($replacementingUserIdInfo[$orderId]))
			{
				$error[] = Localization::$e_no_access;
			}
		}
		else
		{
			$needOrderForId = $isLogged['user_id'];
		}
		if ($error)
			return $error;

		$order = new Order;
		$week = new Week;
		$provider = new Provider;
		if (is_numeric($_POST['week_id']))
		{
			$weekId = $_POST['week_id'];
			if ($week->getWeekInfo($weekId) === false)
				return $week->error;
				
			if ($week->data)
			{
				$weekInfo = $week->data;
				if ($provider->getProviderList() === false)
					return $provider->error;
					
				$providerList = $provider->data;
				if (empty($providerList))
					return false;
					
				if ($order->getWeekOrderList($weekId, "order_list_id") === false)
					return $order->error;

				$orderIdList = $order->data;

				if (empty($orderIdList))
					return false;
					
				$userOrder = $order->formatUserOrder($providerList, $orderIdList, $_POST);
				$dayList[] = D_MONDAY;
				$dayList[] = D_TUESDAY;
				$dayList[] = D_WEDNESDAY;
				$dayList[] = D_THURSDAY;
				$dayList[] = D_FRIDAY;
				$dayList[] = D_SATURDAY;
				$dayList[] = D_SUNDAY;

				// User order correct. We should save it...
				if ($order->deleteUserOrders($needOrderForId, $weekId, "1") === false)
					return $order->error;
				
				foreach ($dayList AS $dayItem)
				{
					if (!empty($userOrder[$dayItem]))
					{
						for ($portion = 1; $portion <= 4; $portion++)
						{
							if (!empty($userOrder[$dayItem][$portion]))
							{
								foreach ($userOrder[$dayItem][$portion] AS $portionItem)
								{
									$data = array(
										'day_id'		=> $dayItem,
										'week_id'		=> $weekId,
										'user_id'		=> $needOrderForId,
										'order_list_id'	=> $portionItem['order_list_id'],
										'provider_id'	=> $userOrder[$dayItem][$portionItem['order_list_id']]['provider_id'],
										'ordered_item_count'	=> $portionItem['ordered_item_count']
									);
									if (!$order->addUserOrder($data))
									{
										$error = $order->error;
									}
								}
							}
						}
					}
				}
				if ($error)
					return $error;

				if (isset($_POST['next_week']))
				{
					header ("Location: orders.php?s_ok&next_week&order_id=$orderId");
					return true;
				}
				else
				{
					header ("Location: orders.php?s_ok&order_id=$orderId");
					return true;
				}
			}
			else
			{
				$error[] = Localization::$e_no_week_found;
				include ("messages.html");
				return true;
			}
		}
		else
		{
			echo "Keep smiling :)";
			return true;
		}
	}
	else if (isset($_GET['i_am_lucky']))
	{
		if (!empty($_GET['order_id']))
		{
			$orderId = $_GET['order_id'];
			$needOrderForId = $orderId;
			if (empty($replacementingUserIdInfo[$orderId]) && ($isLogged['status'] != 10))
			{
				$error[] = Localization::$e_no_access;
			}
		}
		else
		{
			$needOrderForId = $isLogged['user_id'];
		}
		if (!$error)
		{
			$order = new Order;
			$week = new Week;
			if ($week->getActiveWeekInfo() === false)
				return $week->error;
				
			$weekInfo = $week->data;
			if (!empty($weekInfo))
			{
				if ($week->getNextWeekInfo($weekInfo['week_id']))
				{
					if (!empty($week->data))
					{
						if (isset($_GET['next_week']))
						{
							$prevWeekInfo = $weekInfo;
							$weekInfo = $week->data;
						}
						else
						{
							$nextWeekInfo = $week->data;
						}
					}
				}
				$provider = new Provider;
				if ($provider->getProviderList() === false)
					return $provider->error;
					
				$providerList = $provider->data;
				if (!empty($providerList))
				{
					if ($order->getWeekOrderList($weekInfo['week_id']) === false)
						return $order->error;
						
					$orderList = $order->data;
					if (!empty($orderList))
					{
						foreach ($orderList AS $orderItem)
						{
							$orderIdToDayProviderPortion[$orderItem['day_id']][$orderItem['provider_id']][$orderItem['portion_number']][] = $orderItem;
						}
						
						// Get providers array
						
						$dayList[] = D_MONDAY;
						$dayList[] = D_TUESDAY;
						$dayList[] = D_WEDNESDAY;
						$dayList[] = D_THURSDAY;
						$dayList[] = D_FRIDAY;
						$dayList[] = D_SATURDAY;
						$dayList[] = D_SUNDAY;
						foreach ($dayList AS $dayItem)
						{
							// Generate provider id
							$providerTmp = $providerList;

							$arrayCount = count($providerTmp);
							while ($arrayCount)
							{
								$providerId = rand(0, count($providerTmp)-1);
								if (!empty($orderIdToDayProviderPortion[$dayItem][$providerTmp[$providerId]['provider_id']]))
								{
									break;
								}
								else
								{
									unset($providerTmp[$providerId]);
									$arrayCount--;
								}
							}

							for ($portion = 1; $portion <= 4; $portion++)
							{
								if (!empty($orderIdToDayProviderPortion[$dayItem][$providerTmp[$providerId]['provider_id']][$portion]))
								{
									$orderChose = rand (0, count($orderIdToDayProviderPortion[$dayItem][$providerTmp[$providerId]['provider_id']][$portion])-1);
									if ($orderIdToDayProviderPortion[$dayItem][$providerTmp[$providerId]['provider_id']][$portion][$orderChose]['blocked'] == 0)
									{
										$userOrder[$dayItem][$portion] = $orderIdToDayProviderPortion[$dayItem][$providerTmp[$providerId]['provider_id']][$portion][$orderChose]['order_list_id'];
										$userOrder[$dayItem]['provider_id'] = $orderIdToDayProviderPortion[$dayItem][$providerTmp[$providerId]['provider_id']][$portion][$orderChose]['provider_id'];
									}
									else
									{
										// break chosen for blocked day.
										break;
									}
								}
							}
						}
						if (!empty($userOrder))
						{
							$weekId = $weekInfo['week_id'];
							if ($order->deleteUserOrders($needOrderForId, $weekId, '1') === false)
								return $order->error;
							
							foreach ($dayList AS $dayItem)
							{
								if (!empty($userOrder[$dayItem]))
								{
									for ($portion = 1; $portion <= 4; $portion++)
									{
										if (!empty($userOrder[$dayItem][$portion]))
										{
											$data = array(
												'day_id'		=> $dayItem,
												'week_id'		=> $weekId,
												'user_id'		=> $needOrderForId,
												'order_list_id'	=> $userOrder[$dayItem][$portion],
												'provider_id'	=> $userOrder[$dayItem]['provider_id']
											);
											if (!$order->addUserOrder($data))
											{
												$error = $order->error;
												break;
											}
										}
									}
									if ($error)
									{
										break;
									}
								}
							}
							if ($error)
								return $error;

							$parameters = explode("?", $_SERVER['HTTP_REFERER']);
							if (!empty($parameters[1]))
							{
								header ("Location: ${_SERVER['HTTP_REFERER']}&s_ok");
								return true;
							}
							else
							{
								header ("Location: ${_SERVER['HTTP_REFERER']}?s_ok");
								return true;
							}
						}
						else
						{
							$parameters = explode("?", $_SERVER['HTTP_REFERER']);
							if (!empty($parameters[1]))
							{
								header ("Location: ${_SERVER['HTTP_REFERER']}&s_ok");
								return true;
							}
							else
							{
								header ("Location: ${_SERVER['HTTP_REFERER']}?s_ok");
								return true;
							}
						}
					}
					else
					{
						include ("orderlist.html");
						return true;
					}
				}
				else
				{
					include ("orderlist.html");
					return true;
				}
			}
			else
			{
				include ("orderlist.html");
				return true;
			}

		}
		else
		{
			include ("messages.html");
			return true;
		}
	}
	else
	{
		if (!empty($_GET['order_id']))
		{
			$orderId = $_GET['order_id'];
			$needOrderForId = $orderId;
			if (empty($replacementingUserIdInfo[$orderId]))
			{
				$error[] = Localization::$e_no_access;
			}
		}
		else
		{
			$needOrderForId = $isLogged['user_id'];
		}
		
		if (!$error)
		{
			$order = new Order;
			$week = new Week;
			if ($week->getActiveWeekInfo() === false)
				return $week->error;

			$weekInfo = $week->data;

			if (!empty($weekInfo))
			{
				if ($week->getNextWeekInfo($weekInfo['week_id']))
				{
					if (!empty($week->data))
					{
						if (isset($_GET['next_week']))
						{
							$prevWeekInfo = $weekInfo;
							$weekInfo = $week->data;
						}
						else
						{
							$nextWeekInfo = $week->data;
						}
					}
				}
				
				# Convert week and next week names. ##############3
				$weekInfo = $week->getLocalizationWeekName($weekInfo, LANG);
				$fullWeekInfo = $week->data;
				$prevWeekInfo = $week->getLocalizationWeekName($prevWeekInfo, LANG);
				$nextWeekInfo = $week->getLocalizationWeekName($nextWeekInfo, LANG);
				########### end of convertions ##############

				$provider = new Provider;
				if ($provider->getProviderList() === false)
					return $provider->error;

				$providerList = $provider->data;
				if (!empty($providerList))
				{
					$dayList[] = D_MONDAY;
					$dayList[] = D_TUESDAY;
					$dayList[] = D_WEDNESDAY;
					$dayList[] = D_THURSDAY;
					$dayList[] = D_FRIDAY;
					$dayList[] = D_SATURDAY;
					$dayList[] = D_SUNDAY;
					
					if (isset($_GET['text_mode']))
					{
						$_GET['user_list'] = "from_order";
						if ($provider->getProviderList("provider_id") === false)
							return $provider->error;

						$providerIdList = $provider->data;
						if (!empty($providerIdList))
						{
							$user = new User;
							if ($user->getUserInfo($needOrderForId) === false)
								return $user->error;

							// We have user list.
							$userListItem = $user->data;
							
							// Try to get week order list
							if ($order->getWeekOrderList($weekInfo['week_id'], "order_list_id") === false)
								return $order->error;

							// We have week order list.
							$weekOrderListByOrderId = $order->data;
							
							if (!empty($weekOrderListByOrderId))
							{
								
								if (!empty($userListItem))
								{
									$userCount = 0;
									{
										if ($order->getUserOrders($userListItem['user_id'], $weekInfo['week_id']) === false)
											return $order->error;

										$userOrder = $order->data;
										$userOrderData[$userCount] = array(
											'user_id'			=> $userListItem['user_id'],
											'login'			=> $userListItem['login'],
											'email'			=> $userListItem['email'],
											'user_name'		=> $userListItem['user_name'],
											);

										if (!empty($userOrder))
										{
											$userOrderData[$userCount]['is_ordered'] = 1;

											foreach ($userOrder AS $userOrderItem)
											{
												$userOrderData[$userCount][$userOrderItem['day_id']][$weekOrderListByOrderId[$userOrderItem['order_list_id']]['portion_number']][] = array (
													'portion_name'	=> $weekOrderListByOrderId[$userOrderItem['order_list_id']]['portion_name'],
													'provider_name'	=> $providerIdList[$userOrderItem['provider_id']]['name'],
													'order_price'		=> $weekOrderListByOrderId[$userOrderItem['order_list_id']]['order_price'],
													'ordered_item_count'	=> $userOrderItem['ordered_item_count']
													);
											}
										}
										else
										{
											$userOrderData[$userCount]['is_ordered'] = 0;
										}
									}
									include ("week_sum.html");
									return true;
								}
								else
								{
									// We should never be there.. but.. we don't know what will be later ;)
									include ("week_sum.html");
									return true;
								}
							}
							else
							{
								include ("week_sum.html");
								return true;
							}
						}
						else
						{
							include ("week_sum.html");
							return true;
						}

					}
					else
					{
					
						if ($order->getWeekOrderList($weekInfo['week_id']) === false)
							return $order->error;

						$weekOrderData = $order->data;
						if (!empty($weekOrderData))
						{
							### Convert data to user data
							$counter = 0;
							foreach ($providerList AS $providerItem)
							{
								$providerToIdList[$providerItem['provider_id']] = $providerItem;
							}
							foreach ($weekOrderData AS $orderItem)
							{
								if (isset($providerToIdList[$orderItem['provider_id']]))
								{
									if (!isset($portionCounter[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]))
									{
										$portionCounter[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']] = 1;
										$isAnyOrderForADay[$orderItem['day_id']] = 1;
									}
									$providerToIdList[$orderItem['provider_id']]['order'][$orderItem['day_id']][$orderItem['portion_number']][$portionCounter[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]] = array(
										'order_list_id'		=> $orderItem['order_list_id'],
										'portion_name'		=> $orderItem['portion_name'],
										'blocked'			=> $orderItem['blocked'],
										'order_price'		=> $orderItem['order_price']
									);
									if ($orderItem['blocked'] == 1)
									{
										$isBlocked[$orderItem['day_id']] = 1;
									}
									
									$portionCounter[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]++;
								}
								else
								{
									# TODO: We need to do something with bad data... why that data available in the DB?
								}
							}
							if ($order->getUserOrders($needOrderForId, $weekInfo['week_id']) === false)
								return $order->error;

							$userOrderList = $order->data;
							if (!empty($userOrderList))
							{
								//
								foreach ($userOrderList AS $userOrderItem)
								{
									$userOrderId['order'][$userOrderItem['order_list_id']] = $userOrderItem;
									$userOrderId['day_'.$userOrderItem['day_id']]['provider_id'] = 1;
									$userOrderId['day_'.$userOrderItem['day_id']][$userOrderItem['provider_id']] = $userOrderItem['provider_id'];
									$userDayOrder[$userOrderItem['day_id']][$userOrderItem['provider_id']] = $userOrderItem['provider_id'];
								}
							}
							include ("orderlist.html");
							return true;

							### End convert data to user data
						}
						else
						{
							include ("orderlist.html");
							return true;
						}
					}
				}
				else
				{
					include ("orderlist.html");
					return true;
				}
			}
			else
			{
				include ("orderlist.html");
				return true;
			}
		}
		else
		{
			include ("messages.html");
			return true;
		}
	}

	return false;
}

include ("request_handler.php");
?>