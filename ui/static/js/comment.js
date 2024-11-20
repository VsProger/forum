function isValidText(text) {
    const pattern = /^[a-zA-Z][a-zA-Z0-9.,!? ]{2,79}$/;
    return pattern.test(text);
}
document.addEventListener('DOMContentLoaded', function() {
    const form = document.querySelector('.formComment');
    const Input = form.querySelector('input[name="text"]');
    const pattern = /^[a-zA-Z][a-zA-Z0-9.,!?_ ]{2,79}$/;
    form.addEventListener('submit', function(event) {
        let errors = [];

        if (!pattern.test(Input.value)) {
            errors.push('Comments can only contain Latin letters, and the following characters: .,!?_');
        } else if (Input.value.length < 4 || Input.value.length > 80) {
            errors.push('Comments must be between 4 and 80 characters long.');
        }
        if (errors.length > 0) {
            event.preventDefault(); // Prevent form submission
            alert(errors.join('\n')); // Display error messages
        }
    });
});