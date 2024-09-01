let numberOfPayloads = 0;
let payloadVisible = 1;

function addSpecialCharacterBetweenSelectedArea() {
    const textarea = document.getElementById('requestData');
    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const selectedText = textarea.value.substring(start, end);
    const beforeText = textarea.value.substring(0, start);
    const afterText = textarea.value.substring(end);
    textarea.value = beforeText + '§' + selectedText + '§' + afterText;
}

document.getElementById('addSectionButton').addEventListener('click', () => { 
    addSpecialCharacterBetweenSelectedArea()
    numberOfPayloads++
    updateButtons(numberOfPayloads)
    updateTextAreas(numberOfPayloads)
});

document.getElementById('clearButton').addEventListener('click', () => {
    const textarea = document.getElementById('requestData');
    textarea.value = textarea.value.replace(/§/g, '');

    numberOfPayloads = 0;
    payloadVisible = 0;
    updateButtons(numberOfPayloads)
    updateTextAreas(numberOfPayloads)
});

document.getElementById('requestData').addEventListener('input', event => {
    const textarea = event.target;
    const numberOfSpecialCharacters =  (textarea.value.match(/§/g) || []).length
    const expectedNumberOfPayloads = Math.floor(numberOfSpecialCharacters / 2);
    console.log('Expected number of payloads:', expectedNumberOfPayloads);

    if (expectedNumberOfPayloads != numberOfPayloads) {
        numberOfPayloads = expectedNumberOfPayloads;
        updateButtons(expectedNumberOfPayloads)
        updateTextAreas(expectedNumberOfPayloads)
    }    
});

function updateButtons(expectedNumberOfButtons) {
    console.log('Updating buttons to have', expectedNumberOfButtons, 'buttons');
    const container = document.getElementById('buttonsContainer');
    const buttons = container.querySelectorAll('button');
    const numberOfButtons = buttons.length

    if (expectedNumberOfButtons > numberOfButtons) {
        for (let i = numberOfButtons; i < expectedNumberOfButtons; i++) {
            const button = createButton('Payload ' + (i+ 1));
            container.appendChild(button);
        }

        return;
    } 
    
    if (expectedNumberOfButtons < numberOfButtons) {
        for (let i = numberOfButtons; i > expectedNumberOfButtons; i--) {
            container.removeChild(buttons[i - 1]);
        }
    }
}

function createButton(innerText) {
    const button = document.createElement('button');
    button.innerText = innerText;
    button.type = "button"
    button.addEventListener('click', () => {
        const idTextAreaVisible = parseInt(innerText.split(' ')[1], 10);
        makeTextAreaVisible(idTextAreaVisible)
    });
    return button;
} 

function makeTextAreaVisible(idTextAreaVisible) {
    console.log('Making text area', idTextAreaVisible, 'visible');
    const textAreas = document.getElementsByClassName('payload-input')
    if (textAreas[payloadVisible - 1]) {
        textAreas[payloadVisible - 1].style.display = 'none';
    } else {
        console.log('Didnt find text area', payloadVisible - 1);
    }

    payloadVisible = idTextAreaVisible;

    if (textAreas[payloadVisible - 1]) {
        textAreas[payloadVisible - 1].style.display = 'block';
    } else {
        console.log('Didnt find text area', payloadVisible - 1);
    }
}

function updateTextAreas(expectedNumberOfTextAreas) {
    console.log('Updating text areas to have', expectedNumberOfTextAreas, 'text areas');
    content = document.getElementById('content');
    const textAreas = content.querySelectorAll('.payload-input')
    const numberOfTextAreas = textAreas.length

    if (expectedNumberOfTextAreas > numberOfTextAreas) {
        for (let i = numberOfTextAreas; i < expectedNumberOfTextAreas; i++) {
            const textarea = createTextarea(i + 1);
            content.appendChild(textarea);
        }

        return;
    } 
    
    if (expectedNumberOfTextAreas < numberOfTextAreas) {
        for (let i = numberOfTextAreas; i > expectedNumberOfTextAreas; i--) {
            content.removeChild(textAreas[i - 1]);
        }
    }
}

function createTextarea(payloadNumber) {
    const newTextarea = document.createElement('textarea');
    
    newTextarea.setAttribute('data-payload', payloadNumber);
    newTextarea.setAttribute('name', 'payload' + payloadNumber);
    newTextarea.setAttribute('class', 'payload-input');
    newTextarea.setAttribute('placeholder', 'Enter your payload here');
    newTextarea.style.display = 'none';

    return newTextarea;
}