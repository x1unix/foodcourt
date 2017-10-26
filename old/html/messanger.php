<?php
	// coded by kosyak <kosyak_ua@yahoo.com>
	
	//// User info add/remove/edit page script.
	
	
function load ($_GET, $_POST)
{	
	require_once("init.php");
	require_once("user.class.php");
	require_once("provider.class.php");
	require_once("week.class.php");
	require_once("order.class.php");
	require_once("messanger.class.php");
	require_once(LANG.".language.php");
	
	$user = new User;
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	
	if ($isLogged['status'] != 10)
		return false;
	
	$dayList[] = D_MONDAY;
	$dayList[] = D_TUESDAY;
	$dayList[] = D_WEDNESDAY;
	$dayList[] = D_THURSDAY;
	$dayList[] = D_FRIDAY;
	$dayList[] = D_SATURDAY;
	$dayList[] = D_SUNDAY;
	$today = date('N', time());
	if (isset($_POST['do_send']))
	{
		if ($user->getUserList('user_name', 'user_id') === false)
			return $user->error;

		$provider = new Provider;
		if ($provider->getProviderList() === false)
			return $provider->error;

		$providerList = $provider->data;
		$userList = $user->data;
		if (empty($userList))
			return $error[] = Localization::$e_fill_required_fields;
			
		$messangerInfo['m_body'] = $_POST['m_body'];
		$messangerInfo['send_id'] = $_POST['send_id'];
		if (empty($messangerInfo['m_body']) || empty($messangerInfo['send_id']))
		{
			$error[] = Localization::$e_fill_required_fields;
		}

		if (!$error)
		{
			$messanger = new Messanger;
			$subject = "Voracity info from <".$isLogged['login'].">";
			$body = "\n";
			$body .= htmlspecialchars($messangerInfo['m_body']);
			if ($messangerInfo['send_id'] == 'all')
			{
				foreach ($userList AS $userItem)
				{
					if (!empty($userItem['email']))
					{
						$to = $userItem['email'];
						if ($messanger->sendHtmlEMail($to, $subject, $body))
						{
							if (!$message)
							{
								$message[] = Localization::$m_mess_sent;
							}
						}
						else
						{
							if (!$error)
							{
								$error[] = Localization::$e_mess_server_error;
								$error[] = $to;
							}
							else
							{
								$error[] = $to;
							}
						}
					}
				}
			}
			else if ($messangerInfo['send_id'] == 'all_ordered_for_today')
			{
				$week = new Week();
				$order = new Order();
				if ($week->getActiveWeekInfo())
				{	
					$weekInfo = $week->data;
					if (!empty($weekInfo))
					{
						if ($order->getOrderedUsers($weekInfo['week_id'], date("N", time())))
						{
							$orderedUsers = $order->data;
							if (!empty($orderedUsers))
							{
								foreach ($orderedUsers AS $orderedUserItem)
								{
									if (!empty($userList[$orderedUserItem['user_id']]['email']))
									{
										$to = $userList[$orderedUserItem['user_id']]['email'];
										var_dump ($to);
										if ($messanger->sendHtmlEMail($to, $subject, $body))
										{
											if (!$message)
											{
												$message[] = Localization::$m_mess_sent;
											}
										}
										else
										{
											if (!$error)
											{
												$error[] = Localization::$e_mess_server_error;
												$error[] = $to;
											}
											else
											{
												$error[] = $to;
											}
										}
									}
								}
							}
							else
							{
								$message[] = Localization::$e_no_users_to_send;
							}
						}
						else
						{
							$error = $order->error;
						}
					}
					else
					{
						$message[] = Localization::$e_no_users_to_send;
					}
				}
				else
				{
					$error = $week->error;
				}
			}
			else if ($messangerInfo['send_id'] == 'admins')
			{
				foreach ($userList AS $userItem)
				{
					if (!empty($userItem['email']))
					{
						if ($userItem['status'] == 10)
						{
							$to = $userItem['email'];
							if ($messanger->sendHtmlEMail($to, $subject, $body))
							{
								if (!$message)
								{
									$message[] = Localization::$m_mess_sent;
								}
							}
							else
							{
								if (!$error)
								{
									$error[] = Localization::$e_mess_server_error;
									$error[] = $to;
								}
								else
								{
									$error[] = $to;
								}
							}
						}
					}
				}
			}
			else
			{
				if ($user->getUserInfo($messangerInfo['send_id']))
				{
					if (!empty($user->data['email']))
					{
						$to = $user->data['email'];
						if ($messanger->sendHtmlEMail($to, $subject, $body))
						{
							if (!$message)
							{
								$message[] = Localization::$m_mess_sent;
							}
						}
						else
						{
							if (!$error)
							{
								$error[] = Localization::$e_mess_server_error;
								$error[] = $to;
							}
							else
							{
								$error[] = $to;
							}
						}
					}
					else
					{
						$error[] = Localization::$e_no_user_email;
					}
				}
				else
				{
					$error = $user->error;
					$error[] = Localization::$e_no_users_to_send;
				}
			}
			include ("messanger.html");
			return true;
		}
		else
		{
			include ("messanger.html");
			return true;
		}
	}
	else if (isset($_POST['launch_available']))
	{
		if ($user->getUserList('user_name') === false)
			return $user->error;
			
		$userList = $user->data;
		$provider = new Provider;
		if ($provider->getProviderList() === false)
			return $provider->error;
			
		$providerList = $provider->data;
		if (!empty($userList))
		{
			// Conver userList to userIdList
			foreach ($userList AS $userItem)
			{
				$userInfo[$userItem['user_id']] = $userItem;
			}
				
			if (!empty($providerList))
			{
				$launchInfo['day_id'] = $_POST['day_id'];
				if (!empty($launchInfo['day_id']))
				{
					$week = new Week;
					if ($week->getActiveWeekInfo() === false)
						return $week->error;
					
						$weekInfo = $week->data;
						if (!empty($weekInfo))
						{
							$order = new Order;
							if ($order->getWeekOrderList($weekInfo['week_id']) === false)
								return $order->error;
								
							$weekOrderList = $order->data;
							if (!empty($weekOrderList))
							{
								$messanger = new Messanger;
								$subject = "Your lunch order for today";
								// Convert orderList to OrderIdList
								foreach ($weekOrderList AS $weekOrderItem)
								{
									$orderInfo[$weekOrderItem['order_list_id']] = $weekOrderItem;
								}
								
								foreach($providerList AS $providerItem)
								{
									if (isset($_POST['provider_'.$providerItem['provider_id']]))
									{
										$launchInfo['provider_'.$providerItem['provider_id']] = 'on';
										if ($order->getUsersOrderList($weekInfo['week_id'], $providerItem['provider_id'], $launchInfo['day_id']))
										{
											$usersOrderList = "";
											$usersOrderList = $order->data;

											if (!empty($usersOrderList))
											{
												// Generate users orders by portions.
												if (isset($userOrder))
												{
													unset($userOrder);
												}
												foreach ($usersOrderList AS $userOrderItem)
												{
													$userOrder[$userOrderItem['user_id']]['user_id'] = $userOrderItem['user_id'];
													$userOrder[$userOrderItem['user_id']]['portion_'.$orderInfo[$userOrderItem['order_list_id']]['portion_number']] .= " ".$orderInfo[$userOrderItem['order_list_id']]['portion_name'];
													if ($userOrderItem['ordered_item_count'] > 1)
													{
														$userOrder[$userOrderItem['user_id']]['portion_'.$orderInfo[$userOrderItem['order_list_id']]['portion_number']] .= "(x".$userOrderItem['ordered_item_count'].")";
													}
												}
												
												// Generate mail to send:
												foreach ($userOrder AS $userOrderItem)
												{
													
													if (!empty($userInfo[$userOrderItem['user_id']]['email']))
													{
														$to = $userInfo[$userOrderItem['user_id']]['email'];
														$body = "\n";
														$body .= "<br><br> Ваше замовлення на сьогодні:<br><br>";
														$body .= "\"".$providerItem['name']."\"<br>";
														$body .= $userOrderItem['portion_1'].", ";
														$body .= $userOrderItem['portion_2'].", ";
														$body .= $userOrderItem['portion_3']." ";
														$body .= "<br><br> Смачного!";

														if ($messanger->sendHtmlEMail($to, $subject, $body))
														{
															if (!$message)
															{
																$message[] = Localization::$m_mess_sent;
															}
														}
														else
														{
															if (!$error)
															{
																$error[] = Localization::$e_mess_server_error;
																$error[] = $to;
															}
															else
															{
																$error[] = $to;
															}
														}
													}
												}
											}
										}
										else
										{
											$error = $order->error;
										}
									}
								}
								
								if (!$message && !$error)
								{
									$message[] = Localization::$m_no_users_to_send;
								}
								include ("messanger.html");
								return true;
							}
							else
							{
								$error[] = Localization::$e_no_week_order;
								include ("messanger.html");
								return true; 
							}
						}
						else
						{
							$error[] = Localization::$e_no_active_week;
							include ("messanger.html");
							return true;
						}
					
				}
				else
				{
					$error[] = Localization::$e_no_day_for_send;
					include ("messanger.html");
					return true;
				}
			}
			else
			{
				$error[] = Localization::$e_no_providers;
				include ("messanger.html");
				return true;
			}
		}
		else
		{
			$error[] = Localization::$e_no_users_to_send;
			include ("messanger.html");
			return true;
		}
	}
	else
	{
		$provider = new Provider;
		if ($user->getUserList('user_name') === false)
			return $user->error;

		if ($provider->getProviderList() === false)
			return $provider->error;

		$providerList = $provider->data;
		$userList = $user->data;
		include ("messanger.html");
		return true;
	}
	return false;
}

include ("request_handler.php");
?>