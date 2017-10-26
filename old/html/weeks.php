<?php
	// Developed by kosyak <kosyak_ua@yahoo.com>
	
	//// Week's info add/remove/edit page script.
	
	
function load()
{	

	require_once("init.php");
	require_once("user.class.php");
	require_once("week.class.php");
	require_once("provider.class.php");
	require_once("order.class.php");
	require_once(LANG.".language.php");
	
	$user = new User;
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	
	
	
	$dayLocalizedName[D_MONDAY] = Localization::$s_monday;
	$dayLocalizedName[D_TUESDAY] = Localization::$s_tuesday;
	$dayLocalizedName[D_WEDNESDAY] = Localization::$s_wednesday;
	$dayLocalizedName[D_THURSDAY] = Localization::$s_thursday;
	$dayLocalizedName[D_FRIDAY] = Localization::$s_friday;
	$dayLocalizedName[D_SATURDAY] = Localization::$s_saturday;
	$dayLocalizedName[D_SUNDAY] = Localization::$s_sunday;
	
	$scriptName = "weeks.php";
	
	if (isset($_GET['s_ok']))
	{
		$message[] = Localization::$m_data_saved;
	}
	
	if ($isLogged['status'] != 10)
		return false;

	if (isset($_GET['add']))
	{
		// Get last available week info
		$week = new Week;
		if ($week->getLastWeekInfo() === false)
			return $week->error;

		$lastWeekInfo = $week->data;
		$week->getWeekJdOrCurrent($lastWeekInfo);
		$weekDateInfo = $week->data;
		$newWeekJdStartDay = $week->getNextMondayJdId($weekDateInfo['end']['jd']);
		if ($week->getLocalizatedWeekHumanReadableDataFromJd($newWeekJdStartDay, LANG) === false)
			return $week->error;

		$fullWeekInfo = $week->data;
		$weekInfo['name'] = $fullWeekInfo[1]['day']." ".$fullWeekInfo[1]['monthname']." - ".$fullWeekInfo[7]['day']." ".$fullWeekInfo[7]['monthname'];
		// The old variant...
		$provider = new Provider;
		if ($provider->getProviderList() === false)
			return $provider->error;

		$dayList[] = D_MONDAY;
		$dayList[] = D_TUESDAY;
		$dayList[] = D_WEDNESDAY;
		$dayList[] = D_THURSDAY;
		$dayList[] = D_FRIDAY;
		$dayList[] = D_SATURDAY;
		$dayList[] = D_SUNDAY;
		$providerList = $provider->data;
		include "week_edit.html";
		return true;
	}
	else if (isset($_GET['do_active']))
	{
		$week = new Week;
		if ($week->SetWeekStatus('0') === false)
			return $week->error;

		if ($week->SetWeekStatus('1', $_GET['do_active']) === false)
			return $week->error;

		header ("Location: ${_SERVER['HTTP_REFERER']}");
		return true;
	}
	else if (isset($_GET['sum']))
	{
		$provider = new Provider;
		$week = new Week;
		$order = new Order;
		$ordered_for = "";
		if (isset($_GET['user_list']))
		{
			if (isset($_GET['ordered_for']))
			{
				if (($_GET['ordered_for'] >= 1) && ($_GET['orderd_for'] <= 7))
				{
					$ordered_for = $_GET['ordered_for'];
				}
			}
			
			if ($week->getWeekInfo($_GET['sum']) === false)
				return $week->error;

			$weekInfo = $week->data;
			if ($week->getWeekJd($weekInfo))
			{
				if ($week->getLocalizatedWeekHumanReadableDataFromJd($week->data['start']['jd'], LANG))
				{
					$fullWeekInfo = $week->data;
					$weekInfo['name'] = $fullWeekInfo[1]['day']." ".$fullWeekInfo[1]['monthname']." - ".$fullWeekInfo[7]['day']." ".$fullWeekInfo[7]['monthname'];
				}
			}
			
			if ($provider->getProviderList("provider_id") === false)
				return $provider->error;

			$providerIdList = $provider->data;
			if (!empty($providerIdList))
			{
				$user = new User;
				if ($user->getUserList('user_name') === false)
					return $user->error;

				// We have a user list.
				$userList = $user->data;
				
				// Try to get week order list
				if ($order->getWeekOrderList($_GET['sum'], "order_list_id") === false)
					return $order->error;

				// We have week order list.
				$weekOrderListByOrderId = $order->data;
				
				if (!empty($weekOrderListByOrderId))
				{
					
					if (!empty($userList))
					{
						$userCount = 0;
						foreach($userList AS $userListItem)
						{
							if ($order->getUserOrders($userListItem['user_id'], $_GET['sum'], null, $ordered_for))
							{
								$userOrder = $order->data;
								if (!empty($userOrder))
								{
									$userOrderData[$userCount] = array(
										'user_id'			=> $userListItem['user_id'],
										'login'			=> $userListItem['login'],
										'email'			=> $userListItem['email'],
										'user_name'		=> $userListItem['user_name'],
										);
								
									$userOrderData[$userCount]['is_ordered'] = 1;

									foreach ($userOrder AS $userOrderItem)
									{
										$userOrderData[$userCount][$userOrderItem['day_id']][$weekOrderListByOrderId[$userOrderItem['order_list_id']]['portion_number']][] = array (
											'portion_name'		=> $weekOrderListByOrderId[$userOrderItem['order_list_id']]['portion_name'],
											'provider_name'		=> $providerIdList[$userOrderItem['provider_id']]['name'],
											'order_price'			=> $weekOrderListByOrderId[$userOrderItem['order_list_id']]['order_price'],
											'ordered_item_count'	=> $userOrderItem['ordered_item_count']
											);
									}
									$userCount++;
								}
							/*	else
								{
									
									#$userOrderData[$userCount]['is_ordered'] = 0;
									# ToDo: do something one...
									unset($userOrderData[$userCount]);
								}
							*/
							
							}
							else
							{
								$error = $order->error;
								break;
							}
						}
						
						if ($error)
							return $error;

						$dayList[] = D_MONDAY;
						$dayList[] = D_TUESDAY;
						$dayList[] = D_WEDNESDAY;
						$dayList[] = D_THURSDAY;
						$dayList[] = D_FRIDAY;
						$dayList[] = D_SATURDAY;
						$dayList[] = D_SUNDAY;
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
			if ($week->getWeekInfo($_GET['sum']))
			{
				$weekInfo = $week->getLocalizationWeekName($week->data, LANG);
				$fullWeekInfo = $week->data;
				
				if ($provider->getProviderList() === false)
					return $provider->error;
					
				$providerList = $provider->data;
				if (!empty($providerList))
				{
					foreach ($providerList AS $providerItem)
					{
						$providerIdInfo[$providerItem['provider_id']] = $providerItem;
					}
					
					if ($order->getWeekOrderList($weekInfo['week_id']) === false)
						return $order->error;
						
					$orderList = $order->data;
					if (!empty($orderList))
					{	
						// Get count of all orders...
						foreach($orderList AS $orderItem)
						{
							if ($order->getUserOrdersCount($orderItem['order_list_id']))
							{
								if (!$order->data)
								{
									$order->data = 0;
								}
								$orderIdToCount[$orderItem['order_list_id']] = array( 
											'portion_count'			=> $order->data,
											'portion_name'	=> $orderItem['portion_name'],
											'sum_price'		=> ($order->data * $orderItem['order_price'])
								);
								
								$orderIdToDayAndProvider[$orderItem['day_id']][$orderItem['provider_id']][] = $orderItem['order_list_id'];
								if ($order->getPayedUserOrdersCount($orderItem['order_list_id']))
								{
									$orderIdToCount[$orderItem['order_list_id']]['payed_portion_count'] = $order->data;
									$orderIdToCount[$orderItem['order_list_id']]['payed_sum_price'] = ($order->data*$orderItem['order_price']);
									
								}
								else
								{
									$error = $order->error;
									break;
								}
							}
							else
							{
								$error = $order->error;
								break;
							}
						}
						if ($error)
							return $error;

						$dayList[] = D_MONDAY;
						$dayList[] = D_TUESDAY;
						$dayList[] = D_WEDNESDAY;
						$dayList[] = D_THURSDAY;
						$dayList[] = D_FRIDAY;
						$dayList[] = D_SATURDAY;
						$dayList[] = D_SUNDAY;
						include ("week_sum.html");
						return true;
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
				$error = $week->error;
				$error[] = Localization::$e_unknown_chose;
				include ("messages.html");
				return true;
			}
		}
	}
	else if (isset($_GET['delete']))
	{
		$week = new Week;
		$order = new Order;
		if ($order->deleteWeekOrderList($_GET['delete']) === false)
			return $order->error;

		if ($week->deleteWeek($_GET['delete']) === false)
			return $week->error;

		header ("Location: ${_SERVER['HTTP_REFERER']}");
		return true;
	}
	else if (isset($_GET['edit']))
	{
		$week = new Week;
		$order = new Order;
		$provider = new Provider;
		$edit_id = $_GET['edit'];
		if ($week->getWeekInfo($edit_id) === false)
			return $week->error;

		$weekInfo = $week->getLocalizationWeekName($week->data, LANG);
		$fullWeekInfo = $week->data;
		
		// The old variant...
		if ($provider->getProviderList() === false)
			return $provider->error;

		$providerList = $provider->data;
		$dayList[] = D_MONDAY;
		$dayList[] = D_TUESDAY;
		$dayList[] = D_WEDNESDAY;
		$dayList[] = D_THURSDAY;
		$dayList[] = D_FRIDAY;
		$dayList[] = D_SATURDAY;
		$dayList[] = D_SUNDAY;
		if ($order->getWeekOrderList($weekInfo['week_id']) === false)
			return $order->error;

		$weekOrderData = $order->data;
		#### Conver code data to user data
		if (!empty($weekOrderData))
		{
			foreach ($weekOrderData AS $orderItem)
			{
				if (!isset($providerForDay[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]))
				{
					$providerForDay[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']] = 1;
				}
				$weekInfo['order_'.$orderItem['provider_id'].'_'.$orderItem['day_id'].'_'.$orderItem['portion_number'].'_'.$providerForDay[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]] = $orderItem['portion_name'];
				$weekInfo['price_'.$orderItem['provider_id'].'_'.$orderItem['day_id'].'_'.$orderItem['portion_number'].'_'.$providerForDay[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]] = $orderItem['order_price'];
				$providerForDay[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]++;
				if ($orderItem['blocked'] == 1)
				{
					$weekInfo['day_'.$orderItem['day_id']] = "off";
				}
				else
				{
					$weekInfo['day_'.$orderItem['day_id']] = "on";
				}
			}
		}
		#### End code data converting
		include ("week_edit.html");
		return true;
	}
	else if (isset($_POST['do_save']))
	{
		$provider = new Provider;
		$week = new Week;
		$weekInfo['name'] = $_POST['name'];
		if (is_numeric($_POST['edit_id']))
		{
			$edit_id = $_POST['edit_id'];
		}
		
		if ($provider->getProviderList() === false)
			return $provider->error;
	
		$providerList = $provider->data;
		######## Save orders for provider's #######
		if (!empty($providerList))
		{
			$dayList[] = D_MONDAY;
			$dayList[] = D_TUESDAY;
			$dayList[] = D_WEDNESDAY;
			$dayList[] = D_THURSDAY;
			$dayList[] = D_FRIDAY;
			$dayList[] = D_SATURDAY;
			$dayList[] = D_SUNDAY;
			
			foreach ($providerList AS $providerItem)
			{
				### Read all orders for all days
				foreach ($dayList AS $dayItem)
				{
					for ($portionCase = 1; $portionCase <=4; $portionCase++)
					{
						$counter = 1; // Portion number counter
						while (!empty($_POST['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]))
						{	
							$weekInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter] =  array(
								'provider_id'		=> $providerItem['provider_id'],
								'day_id'			=> $dayItem,
								'portion_number'	=> $portionCase,
								'portion_name'		=> $_POST['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter],
								'portion_counter'	=> $counter,
								'order_price'		=> str_replace (',', '.', $_POST['price_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter])
							);
							$counter++;
						}
					}
					if ($_POST['day_'.$dayItem] == 'on' || $_POST['day_'.$dayItem] == 'off')
					{
						$dayActivity = $_POST['day_'.$dayItem];
					}
					else
					{
						$dayActivity = 'off';
					}
					$weekInfo['day_'.$dayItem] = $dayActivity;
					$weekDayActivity[$dayItem] = $dayActivity;
				}
			}
			### End read all orders for all days
			if (empty($weekInfo['name']))
			{
				$error[] = Localization::$e_fill_required_fields;
			}
			
			if (is_numeric($weekInfo['name']))
			{
				$error[] = Localization::$e_week_shouldnt_be_numeric;
			}
			
			if (!$error)
			{
				$week->getWeekJd($weekInfo);
				$weekDateInfo = $week->data;
				if (!empty($weekDateInfo))
				{
					
					$weekInfo['name'] = $weekDateInfo['start']['jd']." - ".$weekDateInfo['end']['jd'];
				}
				$order = new Order;
				if (isset($edit_id))
				{
					if ($week->getWeekInfo($edit_id))
					{
						$weekName = $week->data['name'];
						if ($weekName != $weekInfo['name'])
						{
							if ($week->getWeekInfo($weekInfo['name']))
							{
								$error[] = Localization::$e_week_exist;
							}
							else
							{
								$data['name'] = $weekInfo['name'];
							}
						}
						else
						{
							$data['name'] = $weekName;
						}
					}
					else
					{
						$error = $week->error;
						$error[] = Localization::$e_no_week_found;
					}
				}
				else
				{
					if (!$week->getWeekInfo($weekInfo['name']))
					{
						$data = array(
							"name"		=> $weekInfo['name'],
						);
						if ($week->addWeek($data))
						{
							$edit_id = mysql_insert_id();
							$_POST['edit_id'] = $edit_id;
						}
						else
						{
							$error = $week->error;
						}
					}
					else
					{
						$error = $week->error;
						$error[] = Localization::$e_week_exist;
					}
				}
				if (!$error)
				{
					if ($week->editWeekInfo($data, $edit_id) === false)
						return $week->error;

					## Start to edit the orders
					if ($order->getWeekOrderList($edit_id) === false)
						return $order->error;

					$orderList = $order->data;
					## Convert order list
					if (!empty($orderList))
					{
						foreach ($orderList AS $orderItem)
						{
							if (!isset($providerOldForDay[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]))
							{
								$providerOldForDay[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']] = 1;
							}
							$weekOldInfo['order_'.$orderItem['provider_id'].'_'.$orderItem['day_id'].'_'.$orderItem['portion_number'].'_'.$providerOldForDay[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]] = array(
								'portion_name'				=> $orderItem['portion_name'],
								'blocked'			=> $orderItem['blocked'],
								'order_list_id'		=> $orderItem['order_list_id'],
								'order_price'		=> $orderItem['order_price']
							); 
							$providerOldForDay[$orderItem['provider_id']][$orderItem['day_id']][$orderItem['portion_number']]++;
							if ($orderItem['blocked'] == 1)
							{
								$weekOldInfo['day_'.$orderItem['day_id']] = "off";
							}
							else
							{
								$weekOldInfo['day_'.$orderItem['day_id']] = "on";
							}
						}
					}
					## End convert order list
					if (!empty($weekInfo))
					{
						foreach ($providerList AS $providerItem)
						{
							### Read all orders for all days
							foreach ($dayList AS $dayItem)
							{
								for ($portionCase = 1; $portionCase <=4; $portionCase++)
								{
									$counter = 1;
									while (!empty($weekInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]))
									{
										if (($weekInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['portion_name'] != $weekOldInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['portion_name']) || ($weekInfo['day_'.$dayItem] != $weekOldInfo['day_'.$dayItem])
											|| ($weekInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['order_price'] != $weekOldInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['order_price']))
										{
											if ($weekInfo['day_'.$dayItem] == 'off')
											{
												$isBlocked = 1;
											}
											else
											{
												$isBlocked = 0;
											}
											if (isset($weekOldInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['order_list_id']))
											{
												$data = array(
												'portion_name'	=> $weekInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['portion_name'],
												'order_price'		=> $weekInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['order_price'],
												'blocked'	=> $isBlocked
												);

												if (!$order->editWeekOrderInfo($data, $weekOldInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['order_list_id']))
												{
													$error = $order->error;
													break;
												}
											}
											else
											{
												// We have a new data.
												$data = array(
												'week_id'			=> $edit_id,
												'day_id'			=> $dayItem,
												'portion_number'	=> $portionCase,
												'portion_name'		=> $weekInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['portion_name'],
												'provider_id'		=> $providerItem['provider_id'],
												'blocked'			=> $isBlocked,
												'order_price'		=> $weekInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['order_price']
												);

												if (!$order->addWeekOrder($data))
												{
													$error = $order->error;
													break;
												}
											}
										}
										$counter++;
									}
									if (!$error)
									{
										// Delete all old data...
										while (!empty($weekOldInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]))
										{
											if (!$order->deleteWeekOrderInfo($weekOldInfo['order_'.$providerItem['provider_id'].'_'.$dayItem.'_'.$portionCase.'_'.$counter]['order_list_id']))
											{
												$error = $order->error;
												break;
											}
											$counter++;
										}
									}
									else
									{
										break;
									}
								}
							}
						}
					}
					else
					{
						if ($order->deleteWeekOrderList($edit_id) === false)
							return $order->error;

						header ("Location: weeks.php?edit=${edit_id}&s_ok");
						return true;
					}
					
					if ($error)
						return $error;
					
					header ("Location: weeks.php?edit=${edit_id}&s_ok");
					return true;
					## End of order's editing
				}
				else
				{
					$weekInfo = $_POST;
					include ("week_edit.html");
					return true;
				}
			}
			else
			{
				$weekInfo = $_POST;
				include ("week_edit.html");
				return true;
			}
		}
		######### End save orders for provider's #### 
	}
	else
	{
		$week = new Week;
		if ($week->getWeekList() === false)
			return $week->error;

		$tmpWeekList = $week->data;
		
		foreach ($tmpWeekList AS $weekItem)
		{
			# Convert week data to localizated human readable view...
			$weekList[] = $week->getLocalizationWeekName($weekItem, LANG);
		}
		
		include ("weeklist.html");
		return true;
	}
	
	return false;
}

include ("request_handler.php");
?>