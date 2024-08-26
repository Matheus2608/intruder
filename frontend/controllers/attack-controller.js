const requestButton = document.getElementById("button-request");
const responseButton = document.getElementById("button-response");
const requestContent = document.getElementById("content-request");
const responseContent = document.getElementById("content-response");

requestButton.addEventListener("click", () => {
	requestButton.classList.add("active");
	responseButton.classList.remove("active");
	requestContent.style.display = "block";
	responseContent.style.display = "none";
});

responseButton.addEventListener("click", () => {
	responseButton.classList.add("active");
	requestButton.classList.remove("active");
	requestContent.style.display = "none";
	responseContent.style.display = "block";
});
