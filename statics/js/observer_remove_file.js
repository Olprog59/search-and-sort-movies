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
                });
            });

        }
    });

    observer.observe(logsNode, {childList: true});
});
