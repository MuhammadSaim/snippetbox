var navLinks = document.querySelectorAll('nav a');
for (var i = 0; i < navLinks.length; i++) {
    var link = navLinks[i];
    if (link.getAttribute('href') == window.location.pathname) {
        link.classList.add('live');
        break;
    }
}

// highlightjs
document.addEventListener('DOMContentLoaded', (event) => {
    // Ensure hljs is defined
    if (typeof hljs !== 'undefined') {
        document.querySelectorAll('pre code').forEach((block) => {
            hljs.highlightElement(block);
        });
    } else {
        console.error('Highlight.js failed to load.');
    }
});
