<?php

	// Developed by kosyak <kosyak_ua@yahoo.com>
	
	//// Class for message's sending
	
require_once("Mail.php");
	
class Messanger
{
	var $data;
	var $error;
	
	function Messanger ()
	{}
	
	function sendHtmlEMail($to, $subject, $body)
	{
		
	if (!function_exists("quoted_printable_encode")) {
	/**
	* Process a string to fit the requirements of RFC2045 section 6.7. Note that
	* this works, but replaces more characters than the minimum set. For readability
	* the spaces and CRLF pairs aren't encoded though.
	*/
		function quoted_printable_encode($string)
		{
			return preg_replace('/[^\r\n]{73}[^=\r\n]{2}/', "$0=\r\n",
			str_replace("%", "=", str_replace("%0D%0A", "\r\n",
			str_replace("%20"," ",rawurlencode($string)))));
		}
	}
		
		$head = '<html><head><meta http-equiv=3D"Content-Type" content=3D"text/html; charset=
=3Dutf-8"></head><body>';
		$foot = '</body></html>';
		$from = "Voracity System <noreply@lohika.com>";
		#$body = $head . str_replace ("%","=",rawurlencode($body)) . $foot;
		
		$body = $head.quoted_printable_encode(nl2br($body)).$foot;
		
		$headers = array ('MIME-Version' => '1.0',
		'Content-Type' => 'text/html; charset="utf-8"',
		'Content-Transfer-Encoding'	=> 'quoted-printable',
		'From' => $from,
		'To' => $to,
		'Subject' => $subject);
		$smtp = Mail::factory('smtp',
			array ('host' => SMTP_HOST,
				'port' => SMTP_PORT,
				'auth' => true,
				'username' => SMTP_USER,
				'password' => SMTP_PASSWORD
			)
		);
		$mail = $smtp->send($to, $headers, $body);

		if (PEAR::isError($mail))
		{
			return false;
		}
		else
		{
 			return true;
		}
	}
	
	
}
?>