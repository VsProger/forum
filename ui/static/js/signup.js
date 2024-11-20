document.addEventListener('DOMContentLoaded', function() {
    const form = document.querySelector('.form-signup');
    const usernameInput = form.querySelector('input[name="username"]');
    const emailInput = form.querySelector('input[name="email"]');
    const passwordInput = form.querySelector('input[name="password"]');
    const usernameHint = form.querySelector('.username-hint');
    const emailHint = form.querySelector('.email-hint');
    const passwordHint = form.querySelector('.password-hint');
    const usernamePattern = /^[a-zA-Z0-9]+$/;
    const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    const passwordPattern = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]{8,}$/;

    usernameInput.addEventListener('input', function() {
        if (!usernamePattern.test(usernameInput.value)) {
            usernameHint.textContent = 'Username can only contain letters and numbers.';
        } else {
            usernameHint.textContent = '';
        }
    });

    emailInput.addEventListener('input', function() {
        if (!emailPattern.test(emailInput.value)) {
            emailHint.textContent = 'Email is invalid.';
        } else {
            emailHint.textContent = '';
        }
    });

    passwordInput.addEventListener('input', function() {
        if (!passwordPattern.test(passwordInput.value)) {
            passwordHint.textContent = 'Password must be at least 8 characters long, include uppercase and lowercase letters, contain a number, and have a special character (e.g., !@#$%^&*).';
        } else {
            passwordHint.textContent = '';
        }
    });

    form.addEventListener('submit', function(event) {
        let errors = [];

        if (!usernamePattern.test(usernameInput.value)) {
            errors.push('Username can only contain letters and numbers.');
        }

        if (!emailPattern.test(emailInput.value)) {
            errors.push('Email is invalid.');
        }

        if (!passwordPattern.test(passwordInput.value)) {
            errors.push('Password must be at least 6 characters long, include uppercase and lowercase letters, contain a number, and have a special character (e.g., !@#$%^&*).');
        }

        if (errors.length > 0) {
            event.preventDefault(); // Prevent form submission
            alert(errors.join('\n')); // Display error messages
        }
    });
});
