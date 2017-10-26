<?php
require_once (DB_FOR_USE."_core.class.php");

class Week
{
	var $data;
	var $error;
	var $db;
	
	function Week ()
	{
		$className = DB_FOR_USE.'Core';
		$this->db = new $className;
		
	}
	
	function SetWeekStatus($status, $weekId = null)
	{
		$this->data = "";
		
		$data['active'] = $status;
		
		$requestClauses = array (
			'table_name'	=> WEEK_TBL,
			'data'		=> $data
		);

		if (isset($weekId))
		{
			$requestClauses['where_clause'] = sprintf("`week_id` = '%s'", mysql_real_escape_string($weekId));
		}
		
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
	
	function blockWeekDay ($weekId, $dayId)
	{
		$this->data = "";
		
		$data['blocked'] = 1;
		
		$requestClauses = array (
			'table_name'	=> ORDER_LIST_TBL,
			'data'		=> $data,
			'where_clause'	=> sprintf("`week_id` = '%s' AND `day_id` = '%s'", mysql_real_escape_string($weekId), mysql_real_escape_string($dayId))
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
	
	function getWeekInfo ($nameOrId)
	{
		$this->data = "";
		
		$requestClauses['table_name'] = WEEK_TBL;

		if (is_numeric($nameOrId))
		{
			$requestClauses['where_clause'] = sprintf("`week_id` = '%s'", mysql_real_escape_string($nameOrId));
		}
		else
		{
			$requestClauses['where_clause'] = sprintf("`name` = '%s'", mysql_real_escape_string($nameOrId));
		}
		if ($this->db->get_data($requestClauses) === true)
		{
			$this->data = $this->db->data[0];
			return $this->data;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function getLastWeekInfo ()
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> WEEK_TBL,
			'order_clause'	=> "`week_id` DESC",
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
	
	function getNextWeekInfo($weekId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> WEEK_TBL,
			'where_clause'	=> sprintf("`week_id` > '%s'", mysql_real_escape_string($weekId)),
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
		}
	}
	
	function getActiveWeekInfo()
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> WEEK_TBL,
			'where_clause'	=> "`active` = '1'"
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
	
	function getWeekList()
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> WEEK_TBL,
			'order_clause'	=> "`week_id` ASC"
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
	
	function addWeek ($addData)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> WEEK_TBL,
			'data'		=> $addData
		);

		if ($this->db->add_data($requestClauses) === true)
		{
			return true;
		}
		else
		{
			$this->error = $this->db->error;
			return false;
		}
	}
	
	function editWeekInfo ($dataArray, $weekId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> WEEK_TBL,
			'data'		=> $dataArray,
			'where_clause'	=> sprintf("`week_id` = '%s'", mysql_real_escape_string($weekId))
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
	
	function deleteWeek ($weekId)
	{
		$this->data = "";
		
		$requestClauses = array (
			'table_name'	=> WEEK_TBL,
			'where_clause'	=> sprintf("`week_id` = '%s'", mysql_real_escape_string($weekId))
		);

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
	
	function getJdWeekInfoFromEn ($weekName)
	{
	
		$currentDate = getdate();
		if (empty ($weekName)) return false;
		
		$nameArray = preg_split ("/[ ]+/", $weekName);
		if (empty($nameArray[0])) return false;
		
		if (!is_numeric($nameArray[0])) return false;
		$weekInfo['start']['day'] = $nameArray[0];
		
		if (empty($nameArray[1])) return false;
		$weekInfo['start']['month'] = $nameArray[1];
		
		if (!is_numeric($nameArray[3])) return false;
		$weekInfo['end']['day'] = $nameArray[3];
		
		if (empty($nameArray[4])) return false;
		$weekInfo['end']['month'] = $nameArray[4];
		 #TODO: switch with month
		 
		$weekInfo['start']['year'] = $currentDate['year'];
		 switch ($weekInfo['start']['month'])
		 {
			case "January":
			case "january":
				$weekInfo['start']['month_number'] = 1;
				break;
			case "February":
			case "february":
				$weekInfo['start']['month_number'] = 2;
				break;
			case "March":
			case "march":
				$weekInfo['start']['month_number'] = 3;
				break;
			case "April":
			case "april":
				$weekInfo['start']['month_number'] = 4;
				break;
			case "May":
			case "may":
				$weekInfo['start']['month_number'] = 5;
				break;
			case "June":
			case "june":
				$weekInfo['start']['month_number'] = 6;
				break;
			case "July":
			case "july":
				$weekInfo['start']['month_number'] = 7;
				break;
			case "August":
			case "august":
				$weekInfo['start']['month_number'] = 8;
				break;
			case "September":
			case "september":
				$weekInfo['start']['month_number'] = 9;
				break;
			case "October":
			case "october":
				$weekInfo['start']['month_number'] = 10;
				break;
			case "November":
			case "november":
				$weekInfo['start']['month_number'] = 11;
				break;
			 case "December":
			 case "december":
				$weekInfo['start']['month_number'] = 12;
				break;
			default: 
				return false;
		 }
		 
		 switch ($weekInfo['end']['month'])
		 {
			case "January":
			case "january":
				$weekInfo['end']['month_number'] = 1;
				break;
			case "February":
			case "february":
				$weekInfo['end']['month_number'] = 2;
				break;
			case "March":
			case "march":
				$weekInfo['end']['month_number'] = 3;
				break;
			case "April":
			case "april":
				$weekInfo['end']['month_number'] = 4;
				break;
			case "May":
			case "may":
				$weekInfo['end']['month_number'] = 5;
				break;
			case "June":
			case "june":
				$weekInfo['end']['month_number'] = 6;
				break;
			case "July":
			case "july":
				$weekInfo['end']['month_number'] = 7;
				break;
			case "August":
			case "august":
				$weekInfo['end']['month_number'] = 8;
				break;
			case "September":
			case "september":
				$weekInfo['end']['month_number'] = 9;
				break;
			case "October":
			case "october":
				$weekInfo['end']['month_number'] = 10;
				break;
			case "November":
			case "november":
				$weekInfo['end']['month_number'] = 11;
				break;
			 case "December":
			 case "december":
				$weekInfo['end']['month_number'] = 12;
				break;
			default: 
				return false;
		 }
		 
		 if ($weekInfo['end']['month_number'] < $weekInfo['start']['month_number'])
		 {
			$weekInfo['end']['year'] = $currentDate['year']++;
		 }
		 else
		 {
			$weekInfo['end']['year'] = $currentDate['year'];
		 }
		 
		$weekInfo['start']['jd'] = gregoriantojd ($weekInfo['start']['month_number'], $weekInfo['start']['day'], $weekInfo['start']['year']);
		$weekInfo['end']['jd'] = gregoriantojd ($weekInfo['end']['month_number'], $weekInfo['end']['day'], $weekInfo['end']['year']);
		
		
		$this->data = $weekInfo;
		
		return true;
	}
	
	function getJdWeekInfoFromCurrent()
	{
		$this->data = "";
		
		$currentJd = unixtojd(time());
		
		$currentGregorian = cal_from_jd($currentJd, CAL_GREGORIAN);

		$startJd = $currentJd;
		
		if ($currentGregorian['dow'] <= 0)
		{
			$currentGregorian['dow'] = 7;
		}
		
		while ($currentGregorian['dow'] != 1)
		{
			$startJd--;
			$currentGregorian['dow']--;
		}
		
		$this->data['start']['jd'] = $startJd;
		$this->data['end']['jd'] = $startJd + 6;
		
		return true;
	}
	
	function getJdWeekInfoFromJdName ($weekName)
	{
		if (empty($weekName)) return false;
		
		$nameArray = preg_split("/[ ]+/", $weekName);
		if (empty($nameArray[0])) return false;
		if (!is_numeric($nameArray[0])) return false;
		if (!is_numeric($nameArray[2])) return false;
		
		$this->data['start']['jd'] = $nameArray[0];
		$this->data['end']['jd'] = $nameArray[2];
		return true;

	}
	
	function getJdWeekInfoFromUa ($weekName)
	{
		$currentDate = getdate();
	
		if (empty ($weekName)) return false;
		
		$nameArray = preg_split ("/[ ]+/", $weekName);
		if (empty($nameArray[0])) return false;
		
		if (!is_numeric($nameArray[0])) return false;
		$weekInfo['start']['day'] = $nameArray[0];
		
		if (empty($nameArray[1])) return false;
		$weekInfo['start']['month'] = $nameArray[1];
		
		if (!is_numeric($nameArray[3])) return false;
		$weekInfo['end']['day'] = $nameArray[3];
		
		if (empty($nameArray[4])) return false;
		$weekInfo['end']['month'] = $nameArray[4];
		
		$weekInfo['start']['year'] = $currentDate['year'];
		 switch ($weekInfo['start']['month'])
		 {
			case "Січня":
			case "січня":
				$weekInfo['start']['month_number'] = 1;
				break;
			case "Лютого":
			case "лютого":
				$weekInfo['start']['month_number'] = 2;
				break;
			case "Березня":
			case "березня":
				$weekInfo['start']['month_number'] = 3;
				break;
			case "Квітня":
			case "квітня":
				$weekInfo['start']['month_number'] = 4;
				break;
			case "Травня":
			case "травня":
				$weekInfo['start']['month_number'] = 5;
				break;
			case "Червня":
			case "червня":
				$weekInfo['start']['month_number'] = 6;
				break;
			case "Липня":
			case "липня":
				$weekInfo['start']['month_number'] = 7;
				break;
			case "Серпня":
			case "серпня":
				$weekInfo['start']['month_number'] = 8;
				break;
			case "Вересня":
			case "вересня":
				$weekInfo['start']['month_number'] = 9;
				break;
			case "Жовтня":
			case "жовтня":
				$weekInfo['start']['month_number'] = 10;
				break;
			case "Листопада":
			case "листопада":
				$weekInfo['start']['month_number'] = 11;
				break;
			 case "Грудня":
			 case "грудня":
				$weekInfo['start']['month_number'] = 12;
				break;
			default: 
				return false;
		 }
		 
		 switch ($weekInfo['end']['month'])
		 {
			case "Січня":
			case "січня":
				$weekInfo['end']['month_number'] = 1;
				break;
			case "Лютого":
			case "лютого":
				$weekInfo['end']['month_number'] = 2;
				break;
			case "Березня":
			case "березня":
				$weekInfo['end']['month_number'] = 3;
				break;
			case "Квітня":
			case "квітня":
				$weekInfo['end']['month_number'] = 4;
				break;
			case "Травня":
			case "травня":
				$weekInfo['end']['month_number'] = 5;
				break;
			case "Червня":
			case "червня":
				$weekInfo['end']['month_number'] = 6;
				break;
			case "Липня":
			case "липня":
				$weekInfo['end']['month_number'] = 7;
				break;
			case "Серпня":
			case "серпня":
				$weekInfo['end']['month_number'] = 8;
				break;
			case "Вересня":
			case "вересня":
				$weekInfo['end']['month_number'] = 9;
				break;
			case "Жовтня":
			case "жовтня":
				$weekInfo['end']['month_number'] = 10;
				break;
			case "Листопада":
			case "листопада":
				$weekInfo['end']['month_number'] = 11;
				break;
			 case "Грудня":
			 case "грудня":
				$weekInfo['end']['month_number'] = 12;
				break;
			default: 
				return false;
		 }
		
		 if ($weekInfo['end']['month_number'] < $weekInfo['start']['month_number'])
		 {
			$weekInfo['end']['year'] = $currentDate['year']++;
		 }
		 else
		 {
			$weekInfo['end']['year'] = $currentDate['year'];
		 }
		 
		$weekInfo['start']['jd'] = gregoriantojd ($weekInfo['start']['month_number'], $weekInfo['start']['day'], $weekInfo['start']['year']);
		$weekInfo['end']['jd'] = gregoriantojd ($weekInfo['end']['month_number'], $weekInfo['end']['day'], $weekInfo['end']['year']);
		
		$this->data = $weekInfo;
		
		return true;
	}
	
	function getNextMondayJdId($jdId)
	{
		$this->data = "";
		# TODO: It is very big mistake to do while(true) cycle!!!!! Should be implemented in some any other case.
		while ($calendarInfo = cal_from_jd($jdId , CAL_GREGORIAN))
		{
			if ($calendarInfo['dayname'] == "Monday")
			{
				break;
			}
			else
			{
				$jdId++;
			}
		}
		
		return $jdId;
	}
	
	function getWeekHumanReadableDataFromJd($startJdId)
	{
		$this->data = "";
		
		for ($weekDay = 1; $weekDay <= 7; $weekDay++, $startJdId++)
		{
			$this->data[$weekDay] = cal_from_jd($startJdId, CAL_GREGORIAN);
		}
		
		return true;
	}
	
	function getLocalizatedWeekHumanReadableDataFromJd ($startJdId, $localization)
	{
		$this->data = "";
		
		# Get localization day and month name's.
	/*	global $s_monday;
		global $s_tuesday;
		global $s_wednesday;
		global $s_thursday;
		global $s_friday;
		global $s_saturday;
		global $s_sunday;
		
		global $s_of_january;
		global $s_of_february;
		global $s_of_march;
		global $s_of_april;
		global $s_of_may;
		global $s_of_june;
		global $s_of_july;
		global $s_of_august;
		global $s_of_september;
		global $s_of_october;
		global $s_of_november;
		global $s_of_december;
	*/	
		if ($this->getWeekHumanReadableDataFromJd($startJdId))
		{
			for ($dayCount = 1; $dayCount <= 7; $dayCount++)
			{
				switch ($this->data[$dayCount]['dayname'])
				{
					case "Monday":
						$this->data[$dayCount]['dayname'] = Localization::$s_monday;
						break;
					case "Tuesday":
						$this->data[$dayCount]['dayname'] = Localization::$s_tuesday;
						break;
					case "Wednesday":
						$this->data[$dayCount]['dayname'] = Localization::$s_wednesday;
						break;
					case "Thursday":
						$this->data[$dayCount]['dayname'] = Localization::$s_thursday;
						break;
					case "Friday":
						$this->data[$dayCount]['dayname'] = Localization::$s_friday;
						break;
					case "Saturday":
						$this->data[$dayCount]['dayname'] = Localization::$s_saturday;
						break;
					case "Sunday":
						$this->data[$dayCount]['dayname'] = Localization::$s_sunday;
				}
				
				switch ($this->data[$dayCount]['monthname'])
				{
					case "January":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_january;
						break;
					case "February":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_february;
						break;
					case "March":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_march;
						break;
					case "April":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_april;
						break;
					case "May":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_may;
						break;
					case "June":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_june;
						break;
					case "July":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_july;
						break;
					case "August":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_august;
						break;
					case "September":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_september;
						break;
					case "October":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_october;
						break;
					case "November":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_november;
						break;
					 case "December":
						$this->data[$dayCount]['monthname'] = Localization::$s_of_december;
						break;
				}
			}
		}
		else
		{
			return false;
		}
		
		return true;
	}
	
	
	function getWeekJd($lastWeekInfo = null)
	{
		$this->data = "";
		if (!empty($lastWeekInfo))
		{
			if (($this->getJdWeekInfoFromJdName($lastWeekInfo['name'])) || ($this->getJdWeekInfoFromEn($lastWeekInfo['name'])) || ($this->getJdWeekInfoFromUa($lastWeekInfo['name'])))
			{
				return true;
			}
			else
			{
				return false;
			}
		}
		else
		{
			return false;
		}
	}
	
	function getWeekJdOrCurrent ($lastWeekInfo = null)
	{

		$this->data = "";
		if (!empty($lastWeekInfo))
		{
			if (($this->getWeekJd($lastWeekInfo))) 
			{
				return true;
			}
			else
			{
				$this->getJdWeekInfoFromCurrent();
				return true;
			}
		}
		else
		{
			$this->getJdWeekInfoFromCurrent();
			return true;
		}
	}
	
	function getLocalizationWeekName ($weekInfo, $language)
	{
		$this->data = "";
		if ($this->getWeekJd($weekInfo))
		{
			if ($this->getLocalizatedWeekHumanReadableDataFromJd($this->data['start']['jd'], $language))
			{
				$fullWeekInfo = $this->data;
				$weekInfo['name'] = $fullWeekInfo[1]['day']." ".$fullWeekInfo[1]['monthname']." - ".$fullWeekInfo[7]['day']." ".$fullWeekInfo[7]['monthname'];
			}
		}
		
		return $weekInfo;
	}
}
?>