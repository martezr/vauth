function status() {

    var xmlhttp = new XMLHttpRequest();

    xmlhttp.onreadystatechange = function()
    {
        if (this.readyState == 4 && this.status == 200)
        {
            //Use parse() method to convert JSON string to JSON object
            var responseJsonObj = JSON.parse(this.responseText);

            console.log( responseJsonObj.version );
            vauthSummary = document.getElementById("vauth-summary");
            vauthSummary.innerHTML = responseJsonObj.version
            console.log( responseJsonObj.vault_status );
            vaultStatus = document.getElementById("vault-status");
            vaultStatus.innerHTML = responseJsonObj.vault_status
            console.log( responseJsonObj.vsphere_status );
            vSphereStatus = document.getElementById("vsphere-status");
            vSphereStatus.innerHTML = responseJsonObj.vsphere_status

        }
    };

    xmlhttp.open("GET", "/api/v1/health", true);
    xmlhttp.send();
}