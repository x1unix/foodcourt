<?php 
	// Developed by kosyak <kosyak_ua@yahoo.com>
	
	//// Companies add/remove/edit page script.
	
	
	
	function load ($_GET, $_POST)
	{
		require_once("init.php");
		require_once("user.class.php");
		require_once ("company.class.php");
		require_once(LANG.".language.php");
		
		$user = new User;
		$isLogged = $user->isLogged($_COOKIE[C_LOGIN], $_COOKIE[C_PASSWORD]);
		
		if ($isLogged['status'] != 10)
			return false;
		
		if (isset($_GET['add']))
		{
			include ("company_edit.html");
			return true;
		}
		else if (isset($_POST['do_save']))
		{
			$company = new Company;
			$companyInfo = $_POST;
			if ($company->getCompanyInfo($companyInfo['company_name']) === false)
				return $company->error;

			if (!empty($company->data))
			{
				$error[] = Localization::$e_company_exist;
				include ("company_edit.html");
				return true;
			}
			
			$companyData['company_name'] = $companyInfo['company_name'];
			
			if ($company->addCompany($companyData) === false)
				return $company->error;
				
			header ("Location: companies.php");
			return true;
		}
		else
		{
			$company = new Company;
			if ($company->getCompanyList() === false)
				return $company->error;
			
			$companyList = $company->data;
			include ("company_list.html");
			return true;
		}
		
		# We should never be there... otherwise: bug appear.
		return false;
	}
	
	include ("request_handler.php");
?>