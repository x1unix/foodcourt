//		Developed by kosyak <kosyak.ua@gmail.com>

var inputIdArray = [];
var inputIncrementArray = [];


function AddMenuField(id, infoName,priceName, priceCurrency) 
{

	var incrementNumber = 0;
	var isFound = 0;
	if(! document.getElementById && document.createElement) { return; }
	
	var incrementIndex = inputIdArray.indexOf(infoName);
	
	if (incrementIndex >= 0)
	{
		inputIncrementArray[incrementIndex]++;
		incrementNumber = inputIncrementArray[incrementIndex];
	}
	else
	{
		inputIdArray.push(infoName);
		inputIncrementArray.push(0);
		incrementNumber = 0;
	}
	
	var inHere = document.getElementById(id);
	var infoField = document.createElement("input");
	var priceField = document.createElement("input");

	
	infoName = String(infoName + incrementNumber);
	priceName = String(priceName + incrementNumber);
	infoField.name = infoName;
	infoField.type = "text";
	priceField.type = "text";
	priceField.name = priceName;
	priceField.size = 4;
	
	var brTag = document.createElement("br");

	inHere.appendChild(document.createTextNode(incrementNumber + ":"))
	inHere.appendChild(infoField);
	inHere.appendChild(document.createTextNode(" " ))
	inHere.appendChild(priceField);
	inHere.appendChild(document.createTextNode(priceCurrency))
	inHere.appendChild(brTag);

} // function AddFormField()
