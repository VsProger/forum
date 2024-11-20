document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('filterButton').addEventListener('click', function() {
        applyFilter();
    });
});

function applyFilter() {
    var selectedGenres = [];
    var checkboxes = document.querySelectorAll('input[name="genre"]:checked');
    checkboxes.forEach(function(checkbox) {
        selectedGenres.push(checkbox.value);
    });
    console.log("Filtered Genres:", selectedGenres);
}

