<?php
	// coded by kosyak <kosyak_ua@yahoo.com>
	
	//// User info add/remove/edit page script.

function load()
{
	require_once("init.php");
	require_once("user.class.php");
	require_once("week.class.php");
	require_once ("cash.class.php");
	require_once("provider.class.php");
	require_once("order.class.php");
	require_once(LANG.".language.php");
	
	$user = new User;
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	$scriptName = "users.php";
	
	$dayLocalizedName[D_MONDAY] = Localization::$s_monday;
	$dayLocalizedName[D_TUESDAY] = Localization::$s_tuesday;
	$dayLocalizedName[D_WEDNESDAY] = Localization::$s_wednesday;
	$dayLocalizedName[D_THURSDAY] = Localization::$s_thursday;
	$dayLocalizedName[D_FRIDAY] = Localization::$s_friday;
	$dayLocalizedName[D_SATURDAY] = Localization::$s_saturday;
	$dayLocalizedName[D_SUNDAY] = Localization::$s_sunday;
	
	if ($isLogged['status'] != 10)
		return false;

	if (isset($_GET['s_ok']))
	{
		$message[] = Localization::$m_data_saved;
	}
	
	if (isset($_GET['add']))
	{
		include ("user_edit.html");
		return true;
	}
	else if (isset($_GET['change_activity']))
	{
		$user_id = $_GET['user_id'];
		$data['activity'] = $_GET['change_activity'];
		if ($user->editUserInfo($data, $user_id) === false)
			return $user->error;

		header ("Location: ${_SERVER['HTTP_REFERER']}");
		return true;
	}
	else if (isset($_GET['edit']))
	{
		if ($user->getUserInfo($_GET['edit']) === false)
			return $user->error;
		

		$userInfo = $user->data;
		$edit_id = $_GET['edit'];
		$orderId = $_GET['edit'];
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
					if ($provider->getProviderList("provider_id") === false)
						return $provider->error;
					
					$providerIdList = $provider->data;
					if (!empty($providerIdList))
					{
						$_GET['user_list'] = "from_order";
						if ($order->getWeekOrderList($weekInfo['week_id'], "order_list_id") === false)
							return $order->error;

						// We have week order list.
						$weekOrderListByOrderId = $order->data;
							
						if (!empty($weekOrderListByOrderId))
						{
							$userListItem = $userInfo;
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
								
								include ("user_edit.html");
								return true;

							}
							else
							{
								// We should never be there.. but.. we don't know what will be later ;)
								include ("user_edit.html");
								return true;
							}
						}
						else
						{
							include ("user_edit.html");
							return true;
						}
					}
					else
					{
						include ("user_edit.html");
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
						if ($order->getUserOrders($userInfo['user_id'], $weekInfo['week_id']) === false)
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
						include ("user_edit.html");
						return true;
						### End convert data to user data
					}
					else
					{
						include ("user_edit.html");
						return true;
					}
				}
			}
			else
			{
				include ("user_edit.html");
				return true;
			}
		}
		else
		{
			include ("user_edit.html");
			return true;
		}

	}
	else if (isset($_GET['delete']) && !empty($_GET['delete']))
	{
		if ($user->deleteUser($_GET['delete']) === false)
			return $user->error;

		header ("Location: ${_SERVER['HTTP_REFERER']}");
		return true;
	}
	else if (isset($_POST['do_save']))
	{
		$error = "";
		if (empty($_POST['login']))
		{
			$error[] = Localization::$e_fill_required_fields;
		}
		if (empty($_POST['password'])  && !isset($_POST['edit_id']))
		{
			$error[] = Localization::$e_fill_required_fields;
		}
		if (!$error)
		{
			if (isset($_POST['edit_id']))
			{
				if ($user->getUserInfo($_POST['edit_id']) === false)
					return $user->error;
				

				$userLogin = $user->data['login'];
				if ($userLogin != $_POST['login'])
				{
					if ($user->getUserInfo($_POST['login']))
					{
						$error = $user->error;
						$error[] = Localization::$e_username_exist;
					}
					else
					{
						$data['login'] = $_POST['login'];
					}
				}
				if (!$error)
				{
					if (!empty($_POST['password']))
					{
						$data['password'] = md5($_POST['password']);
					}
					$data['status'] = $_POST['status'];
					$data['email'] = $_POST['email'];
					$data['user_name'] = $_POST['user_name'];
					$data['activity'] = $_POST['activity'];
					if ($user->editUserInfo($data, $_POST['edit_id']) === false)
						return $user->error;
					
					if ($isLogged['user_id'] == $_POST['edit_id'])
					{
						if (!empty($_POST['password']))
						{
							setcookie(C_PASSWORD,md5($_POST['password']),NULL,'/');
						}
						setcookie(C_LOGIN,$_POST['login'],time()+365*24*3600,'/');
					}
					header ("Location: users.php?s_ok&edit=".$_POST['edit_id']);
					return true;
				}
				else
				{
					
					$userInfo['login'] = $_POST['login'];
					$userInfo['status'] = $_POST['status'];
					$userInfo['activity'] = $_POST['activity'];
					$edit_id = $_POST['edit_id'];
					include ("user_edit.html");
					return true;
				}
			}
			else
			{
				if (!$user->getUserInfo($_POST['login']))
				{
					$data = array(
						"login"		=> $_POST['login'],
						"status"	=> $_POST['status'],
						"password"	=> md5($_POST['password']),
						"email"		=> $_POST['email'],
						"user_name"	=> $_POST['user_name'],
						"activity"		=> $_POST['activity'],
						"time"		=> time()
					);
					if ($user->addUser($data) === false)
						return $user->error;

					header ("Location: users.php");
					return true;
				}
				else
				{
					if (!empty($user->error))
						return $user->error;
					
					$error[] = $e_username_exist;
					$userInfo['login'] = $_POST['login'];
					$userInfo['status'] = $_POST['status'];
					$userInfo['email'] = $_POST['email'];
					$userInfo['activity'] = $_POST['activity'];
					include ("user_edit.html");
					return true;
				}
			}
		}
		else
		{
			$userInfo['login'] = $_POST['login'];
			$userInfo['status'] = $_POST['status'];
			$userInfo['activity'] = $_POST['activity'];
			if (isset($_POST['edit_id']))
			{
				$edit_id = $_POST['edit_id'];
			}
			include ("user_edit.html");
			return true;
		}
	}
	else if (isset($_POST['do_order_save']))
	{
		$order = new Order;
		$week = new Week;
		$provider = new Provider;
		if (is_numeric($_POST['week_id']) && is_numeric($_POST['user_id']))
		{
			$weekId = $_POST['week_id'];
			$userId = $_POST['user_id'];
			$needOrderForId = $userId;
			if ($week->getWeekInfo($weekId))
			{
				$weekInfo = $week->data;
				if ($provider->getProviderList() === false)
					return $provider->error;
				
				$providerList = $provider->data;
				if (!empty($providerList))
				{
					if ($order->getWeekOrderList($weekId, "order_list_id") === false)
						return $order->error;
					
					$orderIdList = $order->data;
					
					
					if (!empty($orderIdList))
					{
						
						$dayList[] = D_MONDAY;
						$dayList[] = D_TUESDAY;
						$dayList[] = D_WEDNESDAY;
						$dayList[] = D_THURSDAY;
						$dayList[] = D_FRIDAY;
						$dayList[] = D_SATURDAY;
						$dayList[] = D_SUNDAY;
						
						$userOrder = $order->formatUserOrder($providerList, $orderIdList, $_POST);
						// User order correct. We should save it...
						if ($order->deleteUserOrders($userId, $weekId, "1") === false)
							return $order->error;
						
						foreach ($dayList AS $dayItem)
						{
							if (!empty($userOrder[$dayItem]))
							{
								for ($portion = 1; $portion <= 3; $portion++)
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
							header ("Location: users.php?s_ok&next_week&edit=$userId");
							return true;
						}
						else
						{
							header ("Location: users.php?s_ok&edit=$userId");
							return true;
						}
					}
					else
					{
						header ("Location: index.php");
						return true;
					}
				}
				else
				{
					header ("Location: index.php");
					return true;
				}
			}
			else
			{
				$error = $week->error;
				$error[] = Localization::$e_no_week_found;
				include ("messages.html");
				return true;
			}
		}
		else
		{
			# Some hack :-?
			echo "Keep smiling :)";
			return true;
		}
	}
	else if (isset($_POST['mark_rules']))
	{
		$user->setUserMarkRules($_POST);
		header ("Location: ${_SERVER['HTTP_REFERER']}");
		return true;
	}
	else if (isset ($_GET['do_block']))
	{
		$tmp = preg_split ('/[_]/', $_GET['do_block']);

		if (count($tmp) < 2 || count ($tmp) > 2)
		{
			$error[] = Localization::$e_incorrect_values_received;
			return $error;
		}
		else
		{
			$blockInfo['week_id'] = $tmp[0];
			$blockInfo['day_id'] = $tmp[1];
			
			if (!is_numeric ($blockInfo['week_id']) || !is_numeric ($blockInfo['day_id']))
			{
				$error[] = Localization::$e_incorrect_values_received;
				return $error;
			}
			else
			{
				$week = new Week();
				if ($week->getActiveWeekInfo() === false)
					return $week->error;
				
				$weekInfo = $week->data;
				if ($weekInfo['week_id'] == $blockInfo['week_id'])
				{
					if ($week->blockWeekDay($blockInfo['week_id'], $blockInfo['day_id']) === false)
						return $week->error;
					
					header ("Location: users.php");
					return true;
				}
				else
				{
					$error[] = Localization::$e_incorrect_week_id;
					return $error;
				}
			}
		}
		
	}
	else if (isset ($_GET['pay_for_all']))
	{
		$markRules = $user->getUserMarkRules ($_COOKIE);
		$tmp = preg_split ('/[_]/', $_GET['pay_for_all']);

		if (count($tmp) < 2 || count ($tmp) > 2)
		{
			$error[] = Localization::$e_incorrect_values_received;
			return $error;
		}
		else
		{
			$payInfo['week_id'] = $tmp[0];
			$payInfo['day_id'] = $tmp[1];
			
			if (!is_numeric ($payInfo['week_id']) || !is_numeric ($payInfo['day_id']))
			{
				$error[] = Localization::$e_incorrect_values_received;
				return $error;
			}
			else
			{
				if ($markRules['ordered_for'] != $payInfo['day_id'])
				{
					header ("Location: users.php");
					return true;
				}
				else
				{
					if ($user->getUserList('user_name', 'user_id') === false)
						return $user->error;

					$userList = $user->data;
					$week = new Week();
					if ($week->getActiveWeekInfo() === false)
						return $week->error;
					
					$weekInfo = $week->data;
					if ($weekInfo['week_id'] == $payInfo['week_id'])
					{
						#TODO: process payment...
						$order = new Order();
						if ($order->getOrderedUsers($payInfo['week_id'], $payInfo['day_id'], 'order_list_id') === false)
							return $order->error;
						
						$orderedUserList = $order->data;
						if ($order->getDayBlockedStatus($payInfo['week_id'], $payInfo['day_id']) === false)
							return $order->error;
						
						$dayInfo['blocked'] = $order->data['blocked'];
						$dayInfo['day_id'] = $payInfo['day_id'];
						if ($dayInfo['blocked'])
						{
							if (!empty($orderedUserList))
							{
								if ($order->getWeekOrderList($payInfo['week_id'], 'order_list_id'))
								{
									$weekOrderListById = $order->data;
									if (!empty($weekOrderListById))
									{
										$cash = new Cash();
										foreach ($orderedUserList AS $orderedUserItem)
										{
											#TODO: process users... should be method to generate user payment info...
											$userCash = $userList[$orderedUserItem['user_id']]['cash_total'];
											$paymentData = "";
											$userError = 0;
											if ($order->getUserOrders($orderedUserItem['user_id'], $payInfo['week_id'], 'order_list_id'))
											{
												$userOrderListById = $order->data;
												if (!empty($userOrderListById))
												{
													// process user.
													foreach ($userOrderListById AS $paymentOrderId => $paymentOrderValue)
													{
														if ($weekOrderListById[$paymentOrderId]['day_id'] == $payInfo['day_id'])
														{
															if ($userOrderListById[$paymentOrderId]['ordered_amount'] > 0)
															{
																break;
															}
															
															if ($weekOrderListById[$paymentOrderId]['blocked'])
															{
																if (($userCash - ($weekOrderListById[$paymentOrderId]['order_price'] * $userOrderListById[$paymentOrderId]['ordered_item_count'])) >= 0)
																{
																	$paymentData[] = array(
																		'ordered_amount'	=> ($weekOrderListById[$paymentOrderId]['order_price'] * $userOrderListById[$paymentOrderId]['ordered_item_count']),
																		'user_order_id'		=> $userOrderListById[$paymentOrderId]['user_order_id'],
																		'portion_name'	=> $weekOrderListById[$paymentOrderId]['portion_name'],
																		'ordered_item_count' => $userOrderListById[$paymentOrderId]['ordered_item_count']
																	);
																	$userCash -= ($weekOrderListById[$paymentOrderId]['order_price']*$userOrderListById[$paymentOrderId]['ordered_item_count']);
																}
																else
																{
																	$error[] = $userList[$orderedUserItem['user_id']]['user_name']." <".$userList[$orderedUserItem['user_id']]['login'].">: ".Localization::$e_not_enough_of_monay;
																	$userError = 1;
																	break;
																}
															}
															else
															{
																$error[] = $userList[$orderedUserItem['user_id']]['user_name']." <".$userList[$orderedUserItem['user_id']]['login'].">: ".Localization::$e_day_should_be_blocked;
																$userError = 1;
																break;
															}
														}
													}
													
					
													if (!$userError && !empty($paymentData))
													{
														foreach ($paymentData AS $paymentItem)
														{
															$data = array(
																'ordered_amount'	=> $paymentItem['ordered_amount']
															);
															$userOrderId = $paymentItem['user_order_id'];
															if (!$order->editUserOrder($data, $userOrderId))
															{
																$error[] = $userList[$orderedUserItem['user_id']]['user_name']." <".$userList[$orderedUserItem['user_id']]['login'].">: ".$order->error[0];
																$userError = 1;
																break;
															}
															else
															{
																
																$data = array (
																	'tranzaction_user_id_from'	=> $isLogged['user_id'],
																	'tranzaction_user_id_to'		=> $orderedUserItem['user_id'],
																	'tranzaction_type'			=> 'r',
																	'tranzaction_comment'		=> "Payment for ordered: ".$paymentItem['portion_name']."(x".$paymentItem['ordered_item_count'].")",
																	'tranzaction_amount'		=> $paymentItem['ordered_amount'],
																	'tranzaction_time'			=> time()
																);
																if (!$cash->addTranzaction($data))
																{
																	$error[] = $userList[$orderedUserItem['user_id']]['user_name']." <".$userList[$orderedUserItem['user_id']]['login'].">: ".$cash->error[0];
																	$userError = 1;
																	break;
																}
															}
															
														}
														
														if (!$userError)
														{
															$data = array(
																'cash_total'	=> $userCash
															);
															if (!$user->editUserInfo($data, $orderedUserItem['user_id']))
															{
																$error[] = $userList[$orderedUserItem['user_id']]['user_name']." <".$userList[$orderedUserItem['user_id']]['login'].">: ".$user->error[0];
															}
														}
													}
													
												}
											}
											else
											{
												$error[] = $userList[$orderedUserItem['user_id']]['user_name']." <".$userList[$orderedUserItem['user_id']]['login'].">: ".$order->error[0];
											}
											
										}
									}
								}
								else
								{
									$error[] = $week->error[0];
								}
							}
						}
						else
						{
							$error[] = Localization::$e_day_should_be_blocked;
						}
							
						
						if (!empty($orderedUserList))
						{
							foreach ($orderedUserList AS $orderedUserItem)
							{
								$orderedUser[$orderedUserItem['user_id']] = 1;
							}
							
							if ($markRules['payed_for'] > 7)
							{
								$payedDay = 1;
								if ($week->getNextWeekInfo($weekInfo['week_id']))
								{
									if (!empty($week->data))
									{
										$requiredWeekId = $week->data['week_id'];
									}
									else
									{
										$requiredWeekId = 0;
									}
								}
								else
								{
									$requiredWeekId = 0;
									$error[] = $week->error[0];
								}
							}
							else
							{
								$requiredWeekId = $weekInfo['week_id'];
							}
							
							if ($order->getPayedUsers($requiredWeekId, $markRules['payed_for'], 'user_id'))
							{
								$payedUsersById = $order->data;
							
								#TODO: Should be some universale method to sort arrays....
								$freePlace = 0;

									foreach ($userList AS $userItem)
									{
										if (!empty($orderedUser[$userItem['user_id']]))
										{
											$lastFreePlace = $freePlace;
											
											while (($lastFreePlace > 0))
											{
												if ( (strcmp($userItem['user_name'], $sortedUserList[$lastFreePlace-1]['user_name']) > 0) && (!empty($orderedUser[$sortedUserList[$lastFreePlace-1]['user_id']])))
												{
													break;
												}
												$sortedUserList[$lastFreePlace] = $sortedUserList[$lastFreePlace-1];
												$lastFreePlace--;
											}
											
											
											$sortedUserList[$lastFreePlace] = $userItem;
										}
										else
										{
											$sortedUserList[$freePlace] = $userItem;
										}
										$freePlace++;
									}
								
									$userList = $sortedUserList;

								
							}
							else
							{
								$error[] = $order->error[0];
							}
							
						}
						
						include ("userlist.html");
						return true;
						

						#ToDo: gen responce for user.
						
					}
					else
					{
						$error[] = Localization::$e_incorrect_week_id;
						return $error;
					}
				}
			}
		}
	}
	else
	{
		if ($user->getUserList('user_name') === false)
			return $user->error;

		$userList = $user->data;
		$order = new Order;
		$week = new Week;
		
		$markRules = $user->getUserMarkRules ($_COOKIE);
		if ($week->getActiveWeekInfo() === false)
			return $week->error;

		$weekInfo = $week->data;
		if (!empty($weekInfo))
		{

			if ($order->getOrderedUsers($weekInfo['week_id'], $markRules['ordered_for']) === false)
				return $order->error;

			if (!empty($order->data))
			{
				foreach ($order->data AS $orderedUserItem)
				{
					$orderedUser[$orderedUserItem['user_id']] = 1;
				}
				
				if ($markRules['payed_for'] > 7)
				{
					$payedDay = 1;
					if ($week->getNextWeekInfo($weekInfo['week_id']))
					{
						if (!empty($week->data))
						{
							$requiredWeekId = $week->data['week_id'];
						}
						else
						{
							$requiredWeekId = 0;
						}
					}
					else
					{
						$requiredWeekId = 0;
						$error = $week->error;
					}
				}
				else
				{
					$requiredWeekId = $weekInfo['week_id'];
				}
				
				if ($order->getPayedUsers($requiredWeekId, $markRules['payed_for'], 'user_id'))
				{
					$payedUsersById = $order->data;
				
					#TODO: Should be some universale method to sort arrays....
					$freePlace = 0;
					foreach ($userList AS $userItem)
					{
						if (!empty($orderedUser[$userItem['user_id']]))
						{
							$lastFreePlace = $freePlace;
							
							while (($lastFreePlace > 0))
							{
								if ( (strcmp($userItem['user_name'], $sortedUserList[$lastFreePlace-1]['user_name']) > 0) && (!empty($orderedUser[$sortedUserList[$lastFreePlace-1]['user_id']])))
								{
									break;
								}
								$sortedUserList[$lastFreePlace] = $sortedUserList[$lastFreePlace-1];
								$lastFreePlace--;
							}
							
							
							$sortedUserList[$lastFreePlace] = $userItem;
						}
						else
						{
							$sortedUserList[$freePlace] = $userItem;
						}
						$freePlace++;
					}
				
					$userList = $sortedUserList;
					
					
				}
				else
				{
					$error = $order->error;
				}
				
			}
			
			if ($order->getDayBlockedStatus($weekInfo['week_id'], $markRules['ordered_for']))
			{
				$dayInfo['day_id'] = $markRules['ordered_for'];
				$dayInfo['blocked'] = $order->data['blocked'];
			}
			else
			{
				$error = $order->error;
			}
			
			if ($error)
				return $error;

			include ("userlist.html");
			return true;
		}
		else
		{
			include ("userlist.html");
			return true;
		}
	}
	
	return false;
}

include ("request_handler.php");
?>