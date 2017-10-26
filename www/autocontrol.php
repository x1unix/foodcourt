<?php

////////// SUSPENDED IN gen_job_queue() FUNCTION... ///////////
	// Developed by kosyak <kosyak_ua@yahoo.com>
	
	//// Automation control for voracity system.
	//// Using config.inc.php for details.
	//// Planned for cron usage.
	
	/// Exit codes:
	///	0 - all is good;
	///	1 - Unknown option in command parameters;
	///	2 - Blocked option found;

################### Main program ##################
	require_once (dirname(__FILE__)."/init.php");

	$globalLogLevel = 4;
	
	$G_LOG_LEVELS = array (
		"CRIT"	=> 1,
		"ERROR"	=> 2,
		"WARNING"	=> 3,
		"INFO"	=> 4,
		"DEBUG"	=> 5
	);
	
	if ($_SERVER['argc'] <= 1)
	{
		print_help();
		exit (0);
	}
	
	$jobList = build_job_list($_SERVER['argc'], $_SERVER['argv'], get_available_options());
	$jobNames = get_job_names($jobList);
	$jobQueue = gen_job_queue($jobList, $jobNames);
	
	execute_jobs($jobQueue);
	

	exit (0);
################# End of Main program ###############







################### Functions ######################

	function print_help ()
	{
		echo "\n";
		echo "Automation control script for voracity system.\n";
		echo "Developed by kosyak <kosyak_ua@yahoo.com>\n";
		echo "\n";
		echo "USAGE: ".$_SERVER['argv'][0]." OPTION1 [OPTION2] [OPTION3] ...\n";
		echo "\n";
		echo "Available options:\n";

		echo "\t--activate_next_week\t\t Activate next week.\n";
		echo "\t--block_next_day\t\t Block next day.\n";
		echo "\t--db_dump\t\t\t Create full system DB dump.\n";
		echo "\t--debug\t\t\t\t Debug output.\n";
		echo "\t--ldap_sync\t\t\t Synchronyze with ldap server.\t[not implemented]\n";
		echo "\t--send_lunch_order\t\t Send users lunch orders.\n";
		echo "\t--send_notify\t\t\t Send notifications to user's.\t[not implemented]\n";
		echo "\t--silent\t\t\t No any output to console.\n";
		echo "\t--install <path_to_tar_file>\t Install new system version.\t[not implemented]\n";
		echo "\t--help\t\t\t\t Print this help message.\n";

	}
	
	function get_available_options()
	{
		$tempOptions["--activate_next_week"] = array(
			"name"			=> "--activate_next_week",
			"implementation"	=> "f_activate_next_week",
			"required"		=> 0,
			"blocked_by"		=> 0,
			"priority"			=> 4,
			"argv_count"		=> 0
		);
		$tempOptions["--block_next_day"] = array(
			"name"			=> "--block_next_day",
			"implementation"	=> "f_block_next_day",
			"required"		=> 0,
			"blocked_by"		=> 0,
			"priority"			=> 3,
			"argv_count"		=> 0
		);
		
		$tempOptions["--db_dump"] = array (
			"name"			=> "--db_dump",
			"implementation"	=> "f_db_dump",
			"required"		=> 0,
			"blocked_by"		=> 0,
			"priority"			=> 7,
			"argv_count"		=> 0
		);
		
		$tempOptions["--debug"] = array (
			"name"			=> "--debug",
			"implementation"		=> "f_debug",
			"required"			=> 0,
			"blocked_by"		=> "--silent",
			"priority"			=> 10,
			"argv_count"		=> 0
		);
		
		$tempOptions["--ldap_sync"] = array (
			"name"			=> "--ldap_sync",
			"implementation"	=> 0,
			"required"			=> 0,
			"blocked_by"		=> 0,
			"priority"			=> 0,
			"argv_count"		=> 0
		);
		
		$tempOptions["--send_lunch_order"] = array (
			"name"			=> "--send_lunch_order",
			"implementation"	=> "f_send_lunch_order",
			"required"		=> 0,
			"blocked_by"		=> 0,
			"priority"			=> 1,
			"argv_count"		=> 0
		);
		
		$tempOptions["--send_notify"] = array (
			"name"			=> "--send_notify",
			"implementation"	=> 0,
			"required"			=> 0,
			"blocked_by"		=> 0,
			"priority"			=> 0,
			"argv_count"		=> 0
		);
		
		$tempOptions["--silent"] = array (
			"name"			=> "--silent",
			"implementation"		=> "f_silent",
			"required"			=> 0,
			"blocked_by"		=> "--debug",
			"priority"			=> 10,
			"argv_count"		=> 0
		);
		
		$tempOptions["--install"] = array (
			"name"			=> "--install",
			"implementation"	=> 0,
			"required"			=> 0,
			"blocked_by"		=> 0,
			"priority"			=> 0,
			"argv_count"		=> 1
		);
		
		$tempOptions["--help"] = array (
			"name"			=> "--help",
			"implementation"	=> "print_help",
			"required"		=> 0,
			"blocked_by"		=> "--silent",
			"priority"			=> 9,
			"argv_count"		=> 0
		);
		
		return $tempOptions;
	}

	function build_job_list($argc, $argv, $availableOptions)
	{
		# Build our job list
		for ($count = 0, $i = 1; $i < $argc; $i++, $count++)
		{
			if (isset($availableOptions[$argv[$i]]))
			{
				# We have info about this option. Process it.
				$jobList[$count] = $availableOptions[$argv[$i]];
				# Check if we need some options/values for our Option :)
				if ($jobList[$count]["argv_count"] > 0)
				{
					# Read required count of options/values
					for ($j = 0; $j < $jobList[$count]["argv_count"]; $j++)
					{
						$i++;
						$jobList[$count]["argv_".$j] = $argv[$i];
					}
				}
			}
			else
			{
				# We don't have any info about option... print help and exit with error code 1.
				echo "Unknown option: ".$argv[$i]."\n";
				print_help();
				exit (1);
			}
		}
		
		return $jobList;
	}
	
	function get_job_names ($jobList)
	{
		if (!empty($jobList))
		{
			foreach ($jobList AS $jobItem)
			{
				$jobNames[$jobItem["name"]] = 1;
			}
			
			return $jobNames;
		}
		else
		{
			return null;
		}
	}
	
	function gen_job_queue($jobList, $jobNames)
	{
		global $G_LOG_LEVELS;
		$jobQueue = "";
		$i = 0;
		if (!empty($jobList))
		{
			foreach ($jobList AS $jobItem)
			{
				$blockedOption = get_blocked_option($jobItem, $jobNames);
				if (!empty($blockedOption))
				{
					$message = "Couldn't run option: \"".$jobItem["name"]."\" [blocked by ".$blockedOption."]";
					print_message($message, $G_LOG_LEVELS['ERROR']);
					exit(2);
				}
				else
				{
					if (empty($jobQueue))
					{
						$jobQueue[$i] = $jobItem;
						$i++;
					}
					else
					{
					
						# Try to find good position for our item
						$j = $i-1;
						$goodChoice = $i;
						while ($j >= 0)
						{
							if ($jobQueue[$j]['priority'] >= $jobItem['priority'])
							{
								break;
							}
							else
							{
								$goodChoice = $j;
								$j--;
							}
						}
						
						# Insert our item to found choice
						if ($goodChoice == $i)
						{
							$jobQueue[$i] = $jobItem;
							$i++;
						}
						else
						{
							for ($j = $i; $j > $goodChoice; $j--)
							{
								$jobQueue[$j] = $jobQueue[$j-1];
							}
							
							$jobQueue[$goodChoice] = $jobItem;
							$i++;
						}
					}
				}
			}
			
			return $jobQueue;
		}
		else
		{
			return null;
		}
	}
	
	function get_blocked_option($jobItem, $jobNames)
	{
		if (!$jobItem["blocked_by"])
		{
			return null;
		}
		else
		{
			$blockedOptions = explode(" ", $jobItem["blocked_by"]);

			if (!empty($blockedOptions))
			{
				
				foreach ($blockedOptions AS $blockedItem)
				{
					if (isset($jobNames[$blockedItem]))
					{
						return $blockedItem;
					}
				}
				
				return null;
			}
			else
			{
				return null;
			}
		}
	}
	
	function execute_jobs($jobQueue)
	{
		global $G_LOG_LEVELS;
		if (!empty($jobQueue))
		{
			foreach ($jobQueue AS $jobItem)
			{
				if ($jobItem['implementation'] !== 0)
				{
					print_message ("Found implementation for \"".$jobItem['name']."\" option", $G_LOG_LEVELS['DEBUG']);
					if (function_exists($jobItem['implementation']))
					{
						print_message ("Function for \"".$jobItem['name']."\" developed. Run it", $G_LOG_LEVELS['DEBUG']);
						$argv = "";
						for ($i =0; $i < $jobItem['argv_count']; $i++)
						{
							$argv[$i] = $jobItem['argv_'.$i];
						}
						$jobItem['result'] = call_user_func_array($jobItem['implementation'], $argv);
						print_message ("Result of \"".$jobItem['name']."\" option processing is: ".$jobItem['result'], $G_LOG_LEVELS['DEBUG']);
					}
					else
					{
						print_message("No any function required to process \"".$jobItem['name']."\" option", $G_LOG_LEVELS['CRIT']);
						$jobItem['result'] = 1;
					}
				}
				else
				{
					print_message("No implementation for option \"".$jobItem['name']."\"", $G_LOG_LEVELS['WARNING']);
				}
			}
		}
		else
		{
			print_message("No job's in queue.", $G_LOG_LEVELS['INFO']);
		}
		
		print_message("All is done. Enjoy\n:)", $G_LOG_LEVELS['INFO']);
	}
	
	function print_message($message, $logLevel)
	{
		global $globalLogLevel;
		
		$LOG_LEVELS_NAMES = array (
			1	=>	"CRIT",
			2	=>	"ERROR",
			3	=>	"WARNING",
			4	=>	"INFO",
			5	=>	"DEBUG"
		);
	

		
		if ($globalLogLevel >= $logLevel)
		{
			echo time()."\t[".$LOG_LEVELS_NAMES[$logLevel]."]\t".$message." \n";
		}
	}
	
	
	function f_debug ()
	{
		global $globalLogLevel;
		global $G_LOG_LEVELS;
		
		$globalLogLevel = $G_LOG_LEVELS['DEBUG'];
		return 0;
	}
	
	function f_silent ()
	{
		global $globalLogLevel;
		global $G_LOG_LEVELS;
		
		$globalLogLevel = $G_LOG_LEVELS['CRIT'];
		return 0;
	}
	
	
	function f_activate_next_week()
	{
		global $G_LOG_LEVELS;
		$returnValue = 0;
		print_message ("Activate next week run", $G_LOG_LEVELS['INFO']);
		require_once ("week.class.php");
		
		$week = new Week();
		
		print_message ("Try to get active week info", $G_LOG_LEVELS['DEBUG']);
		if ($week->getActiveWeekInfo())
		{
			print_message ("Got active week info", $G_LOG_LEVELS['DEBUG']);
			$activeWeekInfo = $week->data;
			if (!empty($activeWeekInfo))
			{
				print_message ("Active week info isn't empty. Nice :)", $G_LOG_LEVELS['DEBUG']);
				print_message ("Try to get next week info", $G_LOG_LEVELS['DEBUG']);
				if ($week->getNextWeekInfo($activeWeekInfo['week_id']))
				{
					$nextWeekInfo = $week->data;
					print_message ("Got next week info", $G_LOG_LEVELS['DEBUG']);
					if (!empty($nextWeekInfo))
					{
						print_message ("Next week info isn't empty, nice :)", $G_LOG_LEVELS['DEBUG']);
						print_message ("Try to deactivate all activated weeks", $G_LOG_LEVELS['DEBUG']);
						if ($week->SetWeekStatus(0))
						{
							print_message ("All weeks deactivated", $G_LOG_LEVELS['DEBUG']);
							print_message ("Activating next week [week_id: ".$nextWeekInfo['week_id']."]", $G_LOG_LEVELS['DEBUG']);
							if ($week->SetWeekStatus(1, $nextWeekInfo['week_id']))
							{
								print_message ("Next week activated", $G_LOG_LEVELS['INFO']);
								$returnValue = 0;
							}
							else
							{
								print_message ("Failed to activate next week: ".$week->error[0], $G_LOG_LEVELS['ERROR']);
								$returnValue = 1;
							}
						}
						else
						{
							print_message ("Failed to deactivate all weeks: ".$week->error[0], $G_LOG_LEVELS['ERROR']);
							$returnValue = 1;
						}
					}
					else
					{
						print_message ("Next week info is empty. Nothing to do", $G_LOG_LEVELS['WARNING']);
						$returnValue = 0;
					}
				}
				else
				{
					print_message ("Failed to get next week info: ".$week->error[0], $G_LOG_LEVELS['ERROR']);
				}
			}
			else
			{
				print_message ("Active week info is empty. Nothing to do", $G_LOG_LEVELS['WARNING']);
				$returnValue = 0;
			}
		}
		else
		{
			print_message ("Failed to get active week info: ".$week->error[0], $G_LOG_LEVELS['ERROR']);
			$returnValue = 1;
		}
		print_message ("Activate next week done", $G_LOG_LEVELS['INFO']);
		return $returnValue;
	}
	
	function f_block_next_day()
	{
		global $G_LOG_LEVELS;
		$returnValue = 0;
		print_message ("Block next day run", $G_LOG_LEVELS['INFO']);
		require_once("week.class.php");
		
		$dayToBlock = array (
			D_MONDAY		=> D_MONDAY_BLOCK,
			D_TUESDAY		=> D_TUESDAY_BLOCK,
			D_WEDNESDAY	=> D_WEDNESDAY_BLOCK,
			D_THURSDAY	=> D_THURSDAY_BLOCK,
			D_FRIDAY		=> D_FRIDAY_BLOCK,
			D_SATURDAY	=> D_SATURDAY_BLOCK,
			D_SUNDAY		=> D_SUNDAY_BLOCK
			
		);

		$week = new Week();
		$currentDayId = date ("N", time());
		print_message ("Current day ID is: $currentDayId", $G_LOG_LEVELS['DEBUG']);
		print_message ("Try to get active week info", $G_LOG_LEVELS['DEBUG']);
		if ($week->getActiveWeekInfo())
		{
			print_message ("Got week info", $G_LOG_LEVELS['DEBUG']);
			$weekInfo = $week->data;
			if (!empty($weekInfo))
			{
				print_message ("Week info isn't empty, nice :)", $G_LOG_LEVELS['DEBUG']);
				$dayIdToBlock = $dayToBlock[$currentDayId];
				print_message ("Day id to block is: $dayIdToBlock", $G_LOG_LEVELS['DEBUG']);
				if ($dayIdToBlock > 0)
				{
					print_message ("Blocking the day", $G_LOG_LEVELS['DEBUG']);
					if ($week->blockWeekDay($weekInfo['week_id'], $dayIdToBlock))
					{
						print_message ("Day blocked", $G_LOG_LEVELS['INFO']);
						$returnValue = 0;
					}
					else
					{
						print_message ("Failed to block day: ".$week->error[0], $G_LOG_LEVELS['ERROR']);
						$returnValue = 1;
					}
				}
				else
				{
					print_message ("Nothing to block, exit", $G_LOG_LEVELS['DEBUG']);
					$returnValue = 0;
				}
			}
			else
			{
				print_message ("Week info is empty... Nothing to do", $G_LOG_LEVELS['DEBUG']);
				$returnValue = 0;
			}
		}
		else
		{
			print_message ("Failed to get active week info: ".$week->error[0], $G_LOG_LEVELS['ERROR']);
			$returnValue = 1;
		}
		print_message ("Block next day done", $G_LOG_LEVELS['INFO']);
		return $returnValue;
	}
	
	function f_db_dump ()
	{
		global $G_LOG_LEVELS;
		print_message("DB dump run", $G_LOG_LEVELS['INFO']);
		$defaultPath = "/var/www/localhost/htdocs/voracity/db_dump/";
		if (DB_HOST && DB_USER && DB_PASS && DB_NAME)
		{
			$host = escapeshellcmd(DB_HOST);
			$user = escapeshellcmd(DB_USER);
			$pass = escapeshellcmd(DB_PASS);
			$name = escapeshellcmd(DB_NAME);
			
			print_message ("Required constants for DB dump defined", $G_LOG_LEVELS['DEBUG']);
			if (exec("/usr/bin/mysqldump -h $host -u $user -p$pass $name", $mySqlDumpOutput))
			{
				if (DB_DUMP_PATH)
				{
					$folderToSave = DB_DUMP_PATH;
				}
				else
				{
					$folderToSave = $defaultPath;
				}
				
				if (file_exists($folderToSave))
				{
					print_message("File \"$folderToSave\" exists", $G_LOG_LEVELS['DEBUG']);
					
				}
				else
				{
					print_message("The folder don't exist. Trying to create it", $G_LOG_LEVELS['DEBUG']);
					if (mkdir($folderToSave))
					{
						print_message("Folder created successfully", $G_LOG_LEVELS['DEBUG']);
					}
					else
					{
						print_message("Couldn't create folder in: $folderToSave", $G_LOG_LEVELS['ERROR']);
						return 1;
					}
				}
				
				
				// We have a folder, trying to write db dump.
				if (is_writeable($folderToSave))
				{
					print_message ("File \"$folderToSave\" is write able", $G_LOG_LEVELS['DEBUG']);
					if (is_dir($folderToSave))
					{
						print_message("\"$folderToSave\" is a folder. Nice...", $G_LOG_LEVELS['DEBUG']);
						$time = time();
						print_message("Directory is write able for us.. ", $G_LOG_LEVELS);
						if (($file = fopen("$folderToSave"."/voracity_".$time.".sql", "w+")) !== false)
						{
							if (!empty($mySqlDumpOutput))
							{
								foreach ($mySqlDumpOutput AS $dump)
								{
									$dump .= "\n";
									if (fwrite($file, $dump) === false)
									{
										print_message("Couldn't write to db dump file", $G_LOG_LEVELS['ERROR']);
										fclose ($file);
										return 1;
									}
								}
								print_message ("DB dump writed to the file", $G_LOG_LEVELS['DEBUG']);
								fclose ($file);
							}
							else
							{
								print_message ("The DB dump is empty. No data to save", $G_LOG_LEVELS['DEBUG']);
							}
						}
						else
						{
							print_message("Couldn't create file: $folderToSave"."/".$time.".sql", $G_LOG_LEVELS['DEBUG']);
							return 1;
						}
					}
					else
					{
						print_message("\"$folderToSave\" isn't a folder. Abort.", $G_LOG_LEVELS['ERROR']);
						return 1;
					}

				}
				else
				{
					print_message("We don't have write permissions to dump folder: $folderToSave", $G_LOG_LEVELS['ERROR']);
					return 1;
				}
			}
			else
			{
				print_message("Couldn't done DB dump. Error uppear: $mySqlDumpOutput", $G_LOG_LEVELS['ERROR']);
			}
			
			print_message("DB dump done", $G_LOG_LEVELS['INFO']);
			return 0;
		}
		else
		{
			print_message("Not all required options set in configuration file.\nRequired: DB_HOST DB_USER DB_PASS DB_NAME", $G_LOG_LEVELS['ERROR']);
			return 1;
		}
	}
	
	function f_send_lunch_order ()
	{
		// Function to send user's lunch orders from all providers
		
		require_once("user.class.php");
		require_once("provider.class.php");
		require_once("week.class.php");
		require_once("order.class.php");
		require_once("messanger.class.php");
		require_once(LANG.".language.php");
		
		global $G_LOG_LEVELS;
		$today = date('N', time());
		$launchInfo['day_id'] = $today;
		print_message ("In send_lunch_order module", $G_LOG_LEVELS['INFO']);
		
		$user = new User;
		
		if ($user->getUserList())
		{
			print_message ("Got user list", $G_LOG_LEVELS['DEBUG']);
			$userList = $user->data;
			$provider = new Provider;
			if ($provider->getProviderList())
			{
				print_message ("Got provider list", $G_LOG_LEVELS['DEBUG']);
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
						$week = new Week;
						if ($week->getActiveWeekInfo())
						{
							print_message("Got active week info", $G_LOG_LEVELS['DEBUG']);
							$weekInfo = $week->data;
							if (!empty($weekInfo))
							{
								$order = new Order;
								if ($order->getWeekOrderList($weekInfo['week_id']))
								{
									print_message ("Got week order list", $G_LOG_LEVELS['DEBUG']);
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
											print_message ("Process provider: ".$providerItem['name'], $G_LOG_LEVELS['INFO']);
											if ($order->getUsersOrderList($weekInfo['week_id'], $providerItem['provider_id'], $launchInfo['day_id']))
											{
												print_message ("Got user order list", $G_LOG_LEVELS['DEBUG']);
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
														print_message ("Processing user: ".$userInfo[$userOrderItem['user_id']]['login']." ID: [".$userOrderItem['user_id']."]", $G_LOG_LEVELS['DEBUG']);
														if (!empty($userInfo[$userOrderItem['user_id']]['email']))
														{
															print_message ("User order is: ".$userOrderItem['portion_1'].", ".$userOrderItem['portion_2'].", ".$userOrderItem['portion_3'].", ".$userOrderItem['portion_4'], $G_LOG_LEVELS['DEBUG']);
															$to = $userInfo[$userOrderItem['user_id']]['email'];
															$body = "\n";
															$body .= "<br><br> Ваше замовлення на сьогодні:<br><br>";
															$body .= "\"".$providerItem['name']."\"<br>";
															$body .= $userOrderItem['portion_1'].", ";
															$body .= $userOrderItem['portion_2'].", ";
															$body .= $userOrderItem['portion_3'].", ";
															$body .= $userOrderItem['portion_4']." ";
															$body .= "<br><br> Смачного!";
															print_message ("Sending message to user", $G_LOG_LEVELS['DEBUG']);
															if ($messanger->sendHtmlEMail($to, $subject, $body))
															{
																print_message ("Send", $G_LOG_LEVELS['DEBUG']);
															}
															else
															{
																print_message ("Couldn't send message to $to", $G_LOG_LEVELS['WARNING']);
															}
														}
														else
														{
															print_message ("We don't have user e-mail address, should be fixed. User ID: [".$userOrderItem['user_id']."]", $G_LOG_LEVELS['WARNING']);
														}
														
														print_message ("Processed user: ".$userInfo[$userOrderItem['user_id']]['login']." ID: [".$userOrderItem['user_id']."]", $G_LOG_LEVELS['DEBUG']);
													}
												}
												else
												{
													print_message ("User order list is empty", $G_LOG_LEVELS['DEBUG']);
												}
											}
											else
											{
												print_message ("Couldn't get user ordre list for provider:".$order->error, $G_LOG_LEVELS['ERROR']);
											}
											print_message ("Provider ".$providerItem['name']." processed", $G_LOG_LEVELS['INFO']);
										}
										
									}
									else
									{
										print_message("No any week order's available", $G_LOG_LEVELS['WARNING']);
										return 0;
									}
								}
								else
								{
									print_message ("Couldn't get week order list: ".$order->error, $G_LOG_LEVELS['ERROR']);
									return 1;
								}
							}
							else
							{
								print_message("We have not active week...", $G_LOG_LEVELS['WARNING']);
								return 0;
							}
						}
						else
						{
							print_message ("Couldn't get info about active week: ".$week->error, $G_LOG_LEVELS['ERROR']);
							return 1;
						}

					}
					else
					{
						print_message ("The provider list is empty", $G_LOG_LEVELS['WARNING']);
						return 0;
					}
				}
				else
				{
					print_message("Strange but we haven't any users in DB", $G_LOG_LEVELS['WARNING']);
					return 0;
				}
			}
			else
			{
				print_message("Couldn't get provider list: ".$provider->error, $G_LOG_LEVELS['ERROR']);
				return 1;
			}
		}
		else
		{
			print_message("Couldn't get user list: ".$user->error, $G_LOG_LEVELS['ERROR']);
			return 1;
		}
		
		return 0;
	}

?>
