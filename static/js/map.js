let map, infoWindow;

function initMap() {
	map = new google.maps.Map(document.getElementById("map"), {
		center: { lat: 31.0461, lng: 34.8516 },
		zoom: 6,
	});
	infoWindow = new google.maps.InfoWindow();

	const locationButton = document.createElement("button");
	if (!locationButton) {
		console.log('button not found')
		return;
	}

	locationButton.textContent = "Pan to Current Location";
	locationButton.classList.add("bg-white", "text-black", "font-bold", "py-2", "px-4", "rounded", "shadow-md", "hover:shadow-lg", "hover:bg-gray-100", "focus:outline-none", "focus:ring-2", "focus:ring-indigo-500", "focus:ring-offset-2", "focus:ring-offset-gray-100", "transition", "duration-300", "ease-in-out");
	map.controls[google.maps.ControlPosition.TOP_CENTER].push(locationButton);
	locationButton.addEventListener("click", () => {
		if (navigator.geolocation) {
			navigator.geolocation.getCurrentPosition(
				(position) => {
					const pos = {
						lat: position.coords.latitude,
						lng: position.coords.longitude,
					};
					console.log('position', position)

					infoWindow.setPosition(pos);
					infoWindow.setContent("Location found.");
					infoWindow.open(map);
					map.setCenter(pos);
					map.setZoom(15);
				},
				() => {
					handleLocationError(true, infoWindow, map.getCenter());
				},
			);
		} else {
			handleLocationError(false, infoWindow, map.getCenter());
		}
	});
}

function handleLocationError(browserHasGeolocation, infoWindow, pos) {
	infoWindow.setPosition(pos);
	infoWindow.setContent(
		browserHasGeolocation
			? "Error: The Geolocation service failed."
			: "Error: Your browser doesn't support geolocation.",
	);
	infoWindow.open(map);
}

window.initMap = initMap;

