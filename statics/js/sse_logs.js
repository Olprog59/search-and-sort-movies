if (!!window.EventSource) {
    const source = new EventSource('/logs');

    source.addEventListener('message', function(e) {
        const logDiv = document.getElementById('logs');
        const newMsg = document.createElement('div'); // Créer un nouveau div pour chaque log
        const colorCode = e.data.substring(0, 7);
        newMsg.textContent = e.data.substring(7, e.data.length - 4);
        newMsg.classList.add(color(colorCode));
        logDiv.prepend(newMsg); // Ajouter le nouveau message au conteneur
    }, false);
} else {
    console.log("Votre navigateur ne supporte pas SSE");
}

// ajouter la couleur en récupérer le début de la ligne et la fin
function color(colorCode) {
    const color = colorCode.substring(4, 6);
    switch (color) {
        case "31":
            return "red";
        case "32":
            return "green";
        case "33":
            return "yellow";
        case "34":
            return "purple";
        case "35":
            return "magenta";
        case "36":
            return "teal";
        default:
            return "black";
    }
}
