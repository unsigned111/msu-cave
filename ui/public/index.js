var firebase= require('firebase');

firebase.initializeApp({
            apiKey: "AIzaSyBMcXd3cG8JoinbmykBtXyWTJZn9P-XYDY",
            databaseURL: "https://msu-cave.firebaseio.com/",
            authDomain: "msu-cave.firebaseapp.com",
         });
firebase.auth().signInWithEmailAndPassword("admin@holter.com", "msu-cave-rocks").catch(function(error) {
            console.log("Add some things in here");
            var errorCode = error.code;
            var errorMessage = error.message;
        });

const pod = 'pod-1/shutdown'; //change later to correct pod id
var ui = firebase.database().ref('installations/holter/ui');
var off_button = ui.child('button01_off');
var power_button = ui.child('button02_power');
var shutdown = ui.child(pod);


//set all pods to false - for testing purposes 
function setFalse() { 
	for (i = 1; i<=5; i++) {
		ui.child('pod-'+i).set( {
			'broadcaster' 	 : false,
			'headset-bridge' : false,
			'lighting' 		 : false,
			'shutdown' 		 : false,
			'sound' 		 : false
		});
	}
}


// call shutdown function when button01_off is true
off_button.on("value", function(snap) {
	if (snap.val() == true) {
		shutdown_procedure();
	}
});

// turn off pi when button02_power is true
power_button.on("value", function(snap) {
	if (snap.val() == true) {
		console.log("power button true - powering off");
		shutdown.set(false);
		off_button.set(false);
		power_button.set(false);
		//execSync('shutdown -h now');
	}
});


function shutdown_procedure() {

	//add shutdown stuffs here
	//

	shutdown.set(true);
}

