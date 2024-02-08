document.addEventListener("DOMContentLoaded", function () {
    const targetNode = document.getElementById('content-wrapper');

    const observer = new MutationObserver(function (mutationsList) {
        for (const mutation of mutationsList) {
            if (mutation.type === 'childList') {
                mutation.addedNodes.forEach(function (addedNode) {
                    if (addedNode.nodeType === Node.ELEMENT_NODE && addedNode.matches('.file')) {
                        // Par exemple, cacher le message d'erreur après 5 secondes
                        const errorDiv = addedNode.querySelector("[id^=error-message]");
                        if (errorDiv) {
                            setTimeout(function () {
                                errorDiv.style.display = "none";
                            }, 5000);
                        }
                    }
                });
            }
        }
    });

    observer.observe(targetNode, {childList: true, subtree: true});
});
