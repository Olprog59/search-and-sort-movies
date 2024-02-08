document.addEventListener("DOMContentLoaded", function () {
    const logsNode = document.getElementById('logs');

    const allFiles = []
    const observer = new MutationObserver(function (mutationsList) {
        const moviesNode = document.getElementById('files');

        if (moviesNode) {
            moviesNode.querySelectorAll('.file').forEach(function (file) {
                const name = file.querySelector("form>input[name='filename']").value;
                allFiles.push(name);
            })
        }
        for (const mutation of mutationsList) {
            mutation.addedNodes.forEach(function (node) {
                // Vérifie si le texte du log contient une des valeurs dans allFiles
                allFiles.forEach(function (file) {
                    if (node.textContent.includes(file) && node.textContent.includes('has been moved to:')) {
                        let fileDiv = document.querySelector(`.file input[value='${file}']`);
                        if (fileDiv){
                            fileDiv.closest('.file').remove();
                        }
                    }
                    if (node.textContent.includes("inconsistency between file name and duration")) {
                        const reload = document.getElementById('reload');
                        reload.classList.add('start');
                        reload.innerHTML = 'La page va être rechargée dans&nbsp;<span class="secondes">10</span>&nbsp;secondes. Terminer ce que vous étiez en train de faire.';
                        setInterval(function () {
                            const secondesElement = reload.querySelector('.secondes');
                            let secondes = parseInt(secondesElement.textContent);
                            if (secondes > 0) {
                                secondesElement.innerHTML = `${secondes - 1}`;
                            }
                        }, 1000);
                        setTimeout(function () {
                            location.reload();
                        }, 10000);
                    }

                });
            });

        }
    });

    observer.observe(logsNode, {childList: true});
});
