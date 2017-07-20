var firebase= require('firebase');

var ui = firebase.database().ref('installations/holter/ui');
var off_button = ui.child('button01_off');
var power_button = ui.child('button02_power');
var shutdown = ui.child('pod-1/shutdown');

off_button.set(false);
power_button.set(false);

//set all pod systems to false
for (i = 1; i<=5; i++) {
	ui.child('pod-'+i).set( {
		'broadcaster' 	 : false,
		'headset-bridge' : false,
		'lighting' 		 : false,
		'shutdown' 		 : false,
		'sound' 		 : false
	});
}


// call shutdown function when button01_off is true
off_button.on("value", function(snap) {
	if (snap.val() == true) {
		console.log("true");
		shutdown_procedure();
	}
});

// turn off pi when button02_power is true
power_button.on("value", function(snap) {
	if (snap.val() == true) {
		console.log("poweroff true - turning off");

		off_button.set(false);
		shutdown.set(false);
		power_button.set(false);
		//execSync('shutdown -h now');
	}
});


function shutdown_procedure() {

	//TODO add shutdown stuffs

	shutdown.set(true); //set to true if no errors come up
}

