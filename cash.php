<?php 
	// coded by kosyak <kosyak_ua@yahoo.com>
	
	//// Cash add/remove/edit page script.



function load ($_GET, $_POST)
{
	require_once("init.php");
	require_once("user.class.php");
	require_once("week.class.php");
	require_once("provider.class.php");
	require_once("order.class.php");
	require_once ("cash.class.php");
	require_once(LANG.".language.php");
	
	$user = new User;
	$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
	
	if (isset($_GET['s_ok']))
	{
		$message[] = Localization::$m_data_saved;
	}
	
	if ($isLogged['status'] != 10)
		return false;
	if (isset($_POST['do_tranzaction']))
	{
		if (!($_POST['user_id'] > 0) || empty($_POST['tranzaction_type']) || empty($_POST['tranzaction_sum']))
		{
			$error[] = Localization::$e_fill_required_fields;
		}
		
		$tranzactionInfo['user_id'] = $_POST['user_id'];
		$tranzactionInfo['tranzaction_type'] = $_POST['tranzaction_type'];
		$tranzactionInfo['tranzaction_sum'] = $_POST['tranzaction_sum'];
		$tranzactionInfo['tranzaction_comment'] = $_POST['tranzaction_comment'];
		
		if (!$error)
		{
			if ($user->getUserInfo($tranzactionInfo['user_id']) === false)
				return $user->error;
			$userInfo = $user->data;
			if (!empty($userInfo))
			{
				if ( ($tranzactionInfo['tranzaction_type'] != "add") && ($tranzactionInfo['tranzaction_type'] != "remove") && ($tranzactionInfo['tranzaction_type'] != "return") )
				{
					$error[] = Localization::$e_unknown_tranzaction_type;
				}
				
				$tranzactionInfo['tranzaction_sum'] = str_replace (",", ".", $tranzactionInfo['tranzaction_sum']);
				settype ($tranzactionInfo['tranzaction_sum'], "float");
				if (!($tranzactionInfo['tranzaction_sum'] > 0))
				{
					$error[] = Localization::$e_incorrect_sum_value;
				}
				
				if (!$error)
				{
					$userCash = $userInfo['cash_total'];
					if ($tranzactionInfo['tranzaction_type'] == "add")
					{
						$userCash += $tranzactionInfo['tranzaction_sum'];
					}
					else
					{
						$userCash -= $tranzactionInfo['tranzaction_sum'];
					}
					
					if (isset($_POST['confirmation']))
					{
						$data = array (
							'cash_total' => $userCash
						);
						if ($user->editUserInfo($data, $userInfo['user_id']) === false)
							return $user->error;
							
						$message[] = Localization::$m_tranzaction_successfull;
						$tranzactionData = array (
							'tranzaction_user_id_from'	=> $isLogged['user_id'],
							'tranzaction_user_id_to'		=> $userInfo['user_id'],
							'tranzaction_comment'		=> $tranzactionInfo['tranzaction_comment'],
							'tranzaction_amount'		=> $tranzactionInfo['tranzaction_sum'],
							'tranzaction_time'			=> time()
							
						);
						
						switch ($tranzactionInfo['tranzaction_type'])
						{
							case 'add':
								$tranzactionData['tranzaction_type'] = 'a';
								break;
							case 'return':
								$tranzactionData['tranzaction_type'] = 'b';
								break;
							case 'remove':
								$tranzactionData['tranzaction_type'] = 'r';
								break;
						}
						$cash = new Cash();
						if (!$cash->addTranzaction($tranzactionData))
						{
							$error = $cash->error;
						}
						include ("messages.html");
						return true;
					}
					else
					{
						include ("cash_tranzaction_confirmation.html");
						return true;
					}
					
				}
				else
				{
					if ($user->getUserList('user_name') === false)
						return $user->error;
					$userList = $user->data;
					include ("cash_tranzactions.html");
					return true;
				}
			}
			else
			{
				$error[] = $e_unknown_user;
				if ($user->getUserList('user_name') === false)
					return $user->error;
				$userList = $user->data;
				include ("cash_tranzactions.html");
				return true;
			}
		}
		else
		{
			if ($user->getUserList('user_name') === false)
				return $user->error;
			$userList = $user->data;
			include ("cash_tranzactions.html");
			return true;
		}
	}
	else if (!empty($_GET['order_payment']))
	{
		$orderPayment = $_GET['order_payment'];
		$orderPaymentData = preg_split("/[_]+/", $orderPayment);
		
		// Validate user_id
		if (!is_numeric($orderPaymentData[0]))
			$error[0] = Localization::$e_bad_order_payment_format;
		else
			$paymentUserId = $orderPaymentData[0];
		
		// Validate week_id
		if (!is_numeric($orderPaymentData[1]))
			$error[0] = Localization::$e_bad_order_payment_format;
		else
			$paymentWeekId = $orderPaymentData[1];
			
		if (!isset($orderPaymentData[2]))
			$error[0] = Localization::$e_bad_order_payment_format;
		else
		{
			for ($orderPaymentItem = 2; $orderPaymentItem < count($orderPaymentData); $orderPaymentItem++)
			{
				if (!is_numeric($orderPaymentData[$orderPaymentItem]))
				{
					$error[0] = Localization::$e_bad_order_payment_format;
					break;
				}
				else
				{
					$paymentOrderById[$orderPaymentData[$orderPaymentItem]] = 1;
				}
			}
		}
		
		if ($error)
			return $error;
		if ($user->getUserInfo($paymentUserId) === false)
			return $user->error;
		$userInfo = $user->data;
		if (!empty($userInfo))
		{
			$week = new Week();
			if ($week->getWeekInfo($paymentWeekId) === false)
				return $week->error;
			$weekInfo = $week->data;
			if (!empty($weekInfo))
			{
				$order = new Order();
				if ($order->getUserOrders($userInfo['user_id'], $weekInfo['week_id'], 'order_list_id') === false)
					return $order->error;
				$userOrderListById = $order->data;
				if (!empty($userOrderListById))
				{
					if ($order->getWeekOrderList($weekInfo['week_id'], 'order_list_id') === false)
						return $order->error;
					$weekOrderListById = $order->data;
					if (!empty($weekOrderListById))
					{
						$userCash = $userInfo['cash_total'];
						$paymentData = null;
						foreach ($paymentOrderById AS $paymentOrderId => $paymentOrderValue)
						{
							if (!isset($userOrderListById[$paymentOrderId]))
							{
								$error[] = Localization::$e_incorrect_payment_item;
								break;
							}
							else
							{
								if ($userOrderListById[$paymentOrderId]['ordered_amount'] > 0)
								{
									$error[] = Localization::$e_some_item_payed;
									break;
								}
								else
								{
									if ($weekOrderListById[$paymentOrderId]['blocked'])
									{
										if (($userCash - ($weekOrderListById[$paymentOrderId]['order_price']*$userOrderListById[$paymentOrderId]['ordered_item_count'])) >= 0)
										{
											$paymentData[] = array(
												'ordered_amount'	=> ($weekOrderListById[$paymentOrderId]['order_price']*$userOrderListById[$paymentOrderId]['ordered_item_count']),
												'user_order_id'		=> $userOrderListById[$paymentOrderId]['user_order_id'],
												'portion_name'	=> $weekOrderListById[$paymentOrderId]['portion_name'],
												'ordered_item_count' => $userOrderListById[$paymentOrderId]['ordered_item_count']
											);
											$userCash -= ($weekOrderListById[$paymentOrderId]['order_price']*$userOrderListById[$paymentOrderId]['ordered_item_count']);
										}
										else
										{
											$error[] = Localization::$e_not_enough_of_monay;
											break;
										}
									}
									else
									{
										$error[] = Localization::$e_day_should_be_blocked;
										break;
									}
								}
							}
						}
						
						if ($error)
							return $error;
						if (!empty($paymentData))
						{
							$cash = new Cash();
							foreach ($paymentData AS $paymentItem)
							{
								$data = array(
									'ordered_amount'	=> $paymentItem['ordered_amount']
								);
								$userOrderId = $paymentItem['user_order_id'];
								if (!$order->editUserOrder($data, $userOrderId))
								{
									$error = $order->error;
									break;
								}
								else
								{
									
									$data = array (
										'tranzaction_user_id_from'	=> $isLogged['user_id'],
										'tranzaction_user_id_to'		=> $userInfo['user_id'],
										'tranzaction_type'			=> 'r',
										'tranzaction_comment'		=> "Payment for ordered: ".$paymentItem['portion_name']."(x".$paymentItem['ordered_item_count'].")",
										'tranzaction_amount'		=> $paymentItem['ordered_amount'],
										'tranzaction_time'			=> time()
									);
									if (!$cash->addTranzaction($data))
									{
										$error = $cash->error;
										break;
									}
								}
								
							}
							
							if ($error)
								return $error;
							$data = array(
								'cash_total'	=> $userCash
							);
							if ($user->editUserInfo($data, $userInfo['user_id']) === false)
								return $user->error;
							header("Location: ${_SERVER['HTTP_REFERER']}&s_ok");
							return true;
						}
						else
						{
							$message[] = Localization::$m_no_any_order_to_pay;
							include ("messages.html");
							return true;
						}
					}
					else
					{
						$error[] = Localization::$e_no_week_order;
						include ("messages.html");
						return true;
					}
				}
				else
				{
					$error[] = Localization::$e_no_user_order_data;
					include ("messages.html");
					return true;
				}
			}
			else
			{
				$error[] = Localization::$e_no_week_exist;
				include ("messages.html");
				return true;
			}

		}
		else
		{
			$error[] = Localization::$e_no_user_exist;
			include ("messages.html");
			return true;
		}


	}
	else
	{
		if (!empty($_GET['user_id']))
		{
			$tranzactionInfo['user_id'] = $_GET['user_id'];
		}

		if ($user->getUserList('user_name') === false)
			return $user->error;
		$userList = $user->data;
		include ("cash_tranzactions.html");
		return true;
	}
	
	return false;
}

include ("request_handler.php");
?>