const editor = document.getElementById('content');
const filenameSelector = document.getElementById('filename');
const recentPlaksDiv = document.getElementById('recent-plaks');
const lineNumbersDiv = document.getElementById('line-numbers');

function updateLn() {
  const lines = editor.value.split('\n').length;
  const lineNumbers = Array.from({ length: lines }, (_, i) => i + 1).join('\n');

  lineNumbersDiv.textContent = lineNumbers;
}

function updateTitle() {
  document.title = filenameSelector.value == '' ? 'New paste' : ` ${filenameSelector.value} (Unsaved) - Plakken`;
}

function getRecentPlaksFromStorage() {
  return new Set(JSON.parse(localStorage.getItem('recentPlaks')) || []);
}

function updateLocalStorage() {
  localStorage.setItem('recentPlaks', JSON.stringify(Array.from(recentPlaks)));
}

function addRecentPlak(plakId) {
  recentPlaks.add(plakId);
  updateLocalStorage();
}

function deleteRecentPlak(plakId) {
  const plak = document.querySelector(`[href="/${plakId}/"]`).parentElement;

  plak.style.transform = 'translateX(150%)';
  setTimeout(() => plak.remove(), 250);
  
  recentPlaks.delete(plakId);
  updateLocalStorage();
}

function createRecentPlakComponent(plak) {
  const div = document.createElement('div');
  div.classList.add('recent-plak', 'fr');

  const a = document.createElement('a');
  a.href = `/${plak}/`;
  a.textContent = plak;

  const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
  svg.innerHTML = `<line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line>`;
  svg.id = 'cross';
  svg.setAttribute('viewBox', '0 0 24 24');
  svg.onclick = () => deleteRecentPlak(plak);

  div.appendChild(a);
  div.appendChild(svg);

  return div;
}

filenameSelector.addEventListener('input', updateTitle);
editor.addEventListener('input', updateLn);

let recentPlaks = getRecentPlaksFromStorage();

if (recentPlaks === null) {
  recentPlaks = [];
  localStorage.setItem('recentPlaks', JSON.stringify(recentPlaks));
} else {
  for (const plak of recentPlaks)
    recentPlaksDiv.appendChild(createRecentPlakComponent(plak));
}
