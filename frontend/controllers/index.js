document.getElementById("add-section-sign-button").addEventListener("click", function () {
	const textarea = document.getElementById("request-data-textarea");
	const start = textarea.selectionStart;
	const end = textarea.selectionEnd;
	const selectedText = textarea.value.substring(start, end);
	const beforeText = textarea.value.substring(0, start);
	const afterText = textarea.value.substring(end);
	textarea.value = beforeText + '§' + selectedText + '§' + afterText;
});

document.getElementById("clear-section-sign-button").addEventListener("click", function () {
	const textarea = document.getElementById("request-data-textarea");
	textarea.value = textarea.value.replace(/§/g, '');
});

// document.getElementById("request-data-textarea").addEventListener("input", function() {
// 	const textarea = this;
// 	const formattedText = document.getElementById("formatted-text");
// 	const text = textarea.value;

// 	const highlightedText = text.replace(/§(.*?)§/g, '<span class="highlight">$1</span>');

// 	formattedText.innerHTML = highlightedText;
// });
