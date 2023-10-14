fetch("https://search.local.meyn.fr:8030/search?q=test").then(function(response) {
    return response.json();
})