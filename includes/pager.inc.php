<?php
	// coded by kosyak<kosyak_ua@yahoo.com>
	
	
	// The pager.
	if ($itemCount > $itemPerPage) // Pager should be displayed.
	{
		if ($currentPage <= 0)
		{
			$currentPage = 1;
		}
		$pageCount = ceil($itemCount / $itemPerPage);		// Get total page count.
		$query = "";
		$separator = "?";
		if (!empty($_SERVER['QUERY_STRING']))
		{
			$tokenString = strtok($_SERVER['QUERY_STRING'], '&');
			if (!strstr($tokenString, 'page='))
			{
				$query .= $separator.$tokenString;
				$separator = "&";
			}
			while ($tokenString = strtok('&'))
			{
				if (!strstr($tokenString, 'page='))
				{
					$query .= $separator.$tokenString;
					$separator = "&";
				}
			}
		}

		for ($j = 1, $i = 1; $i<=$pageCount; $i++, $j++)
		{
			if ($i == $currentPage)
			{
				echo "<b>$i</b>&nbsp;";
			}
			else
			{
				echo "<a href=\"${script}${query}${separator}page=${i}\">${i}</a>&nbsp;";
			}
			if ($j == 20)
			{
				echo "<br />";
				$j = 1;
			}
		}
	}
?>