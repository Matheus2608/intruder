document.addEventListener('DOMContentLoaded', () => {
    const requestDataTextArea = document.getElementById('requestData');
    const addSectionButton = document.getElementById('addSectionButton');
    const clearButton = document.getElementById('clearButton');
    const contentDiv = document.getElementById('content'); // Container for request/payload textareas

    // --- Event Listeners ---

    // Add payload markers (§) around selected text
    addSectionButton?.addEventListener('click', () => {
        const start = requestDataTextArea.selectionStart;
        const end = requestDataTextArea.selectionEnd;
        const selectedText = requestDataTextArea.value.substring(start, end);

        if (start === end) {
            // If no text selected, maybe insert §§? Or alert.
            alert('Please select the text you want to mark as a payload position.');
            return;
        }
        // Prevent nesting markers like §abc§def§§
        if (selectedText.includes('§')) {
             alert('Cannot place markers within existing markers.');
             return;
        }


        const newValue = requestDataTextArea.value.substring(0, start) +
                         '§' + selectedText + '§' +
                         requestDataTextArea.value.substring(end);

        // Store cursor position before changing value
        const cursorPosition = start + 1; // Position after the opening §

        requestDataTextArea.value = newValue;

        // Restore cursor position (or place it logically)
        requestDataTextArea.focus();
        requestDataTextArea.setSelectionRange(cursorPosition, cursorPosition + selectedText.length);


        updatePayloadInputsUI(); // Update UI based on markers
    });

    // Clear all payload markers (§)
    clearButton?.addEventListener('click', () => {
        // Simple replace might be okay, but safer to track positions if needed later
        requestDataTextArea.value = requestDataTextArea.value.replace(/§/g, '');
        updatePayloadInputsUI(); // Update UI
    });

     // Update payload inputs whenever the request text changes manually (e.g., pasting, typing)
     requestDataTextArea?.addEventListener('input', updatePayloadInputsUI);


    // --- Core UI Update Functions ---

    // Function to dynamically add/remove payload textareas based on § markers
    function updatePayloadInputsUI() {
        if (!requestDataTextArea || !contentDiv) return; // Elements might not exist

        // Count pairs of markers
        const markerCount = Math.floor((requestDataTextArea.value.match(/§/g) || []).length / 2);
        const existingPayloadInputs = contentDiv.querySelectorAll('.payload-input');
        const currentInputCount = existingPayloadInputs.length;

        // Add necessary payload inputs
        for (let i = currentInputCount + 1; i <= markerCount; i++) {
            const newPayloadInput = document.createElement('textarea');
            newPayloadInput.classList.add('payload-input');
            newPayloadInput.name = `payload${i}`;
            newPayloadInput.setAttribute('data-payload', i); // Keep track of index
            newPayloadInput.placeholder = `Payload list ${i} (one payload per line)`;
            contentDiv.appendChild(newPayloadInput);
        }

        // Remove excess payload inputs (from highest index down)
        for (let i = currentInputCount; i > markerCount; i--) {
             if (i > 1) { // Always keep at least payload1
                const inputToRemove = contentDiv.querySelector(`textarea[name="payload${i}"]`);
                if (inputToRemove) {
                    contentDiv.removeChild(inputToRemove);
                }
            } else if (i === 1 && markerCount === 0) {
                 // If only payload1 exists and markers are removed, clear its content? Optional.
                 const payload1Input = contentDiv.querySelector(`textarea[name="payload1"]`);
                 // if (payload1Input) payload1Input.value = ''; // Uncomment to clear
            }
        }
        // Adjust layout/style if needed after adding/removing inputs
        adjustLayout();
    }


    // Adjust layout dynamically based on number of payload inputs
    function adjustLayout() {
         if (!contentDiv) return;
        const requestInput = contentDiv.querySelector('.request-input');
        const payloadInputs = contentDiv.querySelectorAll('.payload-input');
        const totalVisibleInputs = 1 + payloadInputs.length; // 1 for request + number of payload inputs

        if (totalVisibleInputs <= 1) {
             requestInput.style.flexBasis = '100%'; // Request takes full width
        } else if (totalVisibleInputs === 2) {
             requestInput.style.flexBasis = '65%'; // Request takes more space
             payloadInputs.forEach(input => input.style.flexBasis = '35%');
        } else {
            // Distribute width more evenly if many payload inputs
             const percent = Math.floor(100 / totalVisibleInputs);
             requestInput.style.flexBasis = `${percent}%`; // Give request equal share
             payloadInputs.forEach(input => input.style.flexBasis = `${percent}%`);
        }
    }

    // --- Initial Setup ---

    // Initial check in case the textarea already has content with markers on load
    updatePayloadInputsUI();

});