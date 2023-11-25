const codeEditor = document.getElementById('content');
const lineCounter = document.getElementById('lines');

let lineCountCache = 0;

// Update line counter
function updateLineCounter() {
    const lineCount = codeEditor.value.split('\n').length;

    if (lineCountCache !== lineCount) {
        const outarr = Array.from({length: lineCount}, (_, index) => index + 1);
        lineCounter.value = outarr.join('\n');
    }

    lineCountCache = lineCount;
}

codeEditor.addEventListener('input', updateLineCounter);

codeEditor.addEventListener('keydown', (e) => {
    if (e.key === 'Tab') {
        e.preventDefault();

        const {value, selectionStart, selectionEnd} = codeEditor;
        codeEditor.value = `${value.slice(0, selectionStart)}\t${value.slice(selectionEnd)}`;
        codeEditor.setSelectionRange(selectionStart + 1, selectionStart + 1);
        updateLineCounter();
    }
});

updateLineCounter();
