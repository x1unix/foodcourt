function uncheckIfChecked (element)
{
	if (element.newValue == "checked")
	{
		element.newValue = "unchecked"
		element.checked = false
	}
	else
	{
		element.newValue = "checked"
		element.checked = true
	}
}

function checkRadioButton (elementName)
{
	elementName.checked = true
	elementName.newValue = "checked"
}

function uncheckRadioButton (elementName)
{
	elementName.checked = false
	elementName.newValue = "unchekced"
}
