<?php
		
	if (function_exists('load'))
	{
		if (($res = load ($_GET, $_POST)) !== true)
		{
			if ($res === false)
			{
				header ("Location: index.php");
			}
			else
			{
				$error = $res;
				include ("messages.html");
			}
		}
	}
	else
	{
		header ("Location: index.php");
	}

?>