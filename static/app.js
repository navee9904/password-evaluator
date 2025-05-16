document.addEventListener('DOMContentLoaded', () => {
    const passwordInput = document.getElementById('password');
    const togglePassword = document.getElementById('togglePassword');
    const realtimeFeedbackDiv = document.getElementById('realtimeFeedback');
    const suggestPasswordBtn = document.getElementById('suggestPassword');
    const checkAnotherBtn = document.getElementById('checkAnother');
    const errorsDiv = document.getElementById('errors');
    const suggestedPasswordDiv = document.getElementById('suggestedPassword');
    const rtStrengthBar = document.getElementById('rtStrengthBar');
    const crackingTimeText = document.getElementById('crackingTimeText');
    const crackingTimeBadge = document.getElementById('crackingTimeBadge');

    // Show/Hide Password Toggle
    togglePassword.addEventListener('click', () => {
        const isPasswordVisible = passwordInput.type === 'text';
        passwordInput.type = isPasswordVisible ? 'password' : 'text';
        togglePassword.querySelector('i').classList.toggle('fa-eye-slash', !isPasswordVisible);
        togglePassword.querySelector('i').classList.toggle('fa-eye', isPasswordVisible);
        togglePassword.setAttribute('aria-label', isPasswordVisible ? 'Show password' : 'Hide password');
    });

    // Format cracking time into a human-readable string
    const formatCrackingTime = (years) => {
        if (years < 0.000001) return "less than a second";
        const seconds = years * 365 * 24 * 60 * 60;
        if (seconds < 60) return `${Math.round(seconds)} seconds`;
        const minutes = seconds / 60;
        if (minutes < 60) return `${Math.round(minutes)} minutes`;
        const hours = minutes / 60;
        if (hours < 24) return `${Math.round(hours)} hours`;
        const days = hours / 24;
        if (days < 365) return `${Math.round(days)} days`;
        return `${Math.round(years)} years`;
    };

    // Real-time Evaluation
    const evaluateRealtime = async (password) => {
        if (!password) {
            realtimeFeedbackDiv.classList.add('hidden');
            suggestPasswordBtn.classList.add('hidden');
            checkAnotherBtn.classList.add('hidden');
            return;
        }

        try {
            const response = await fetch('/api/evaluate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ password })
            });

            if (!response.ok) {
                throw new Error(`HTTP error ${response.status}`);
            }

            const result = await response.json();
            realtimeFeedbackDiv.classList.remove('hidden');
            realtimeFeedbackDiv.classList.add('fade-in');

            document.getElementById('rtLength').textContent = `Length (≥12): ${result.lengthValid ? '✅ Valid' : '❌ Too short'}`;
            document.getElementById('rtUppercase').textContent = `Uppercase: ${result.hasUppercase ? '✅ Present' : '❌ Missing'}`;
            document.getElementById('rtLowercase').textContent = `Lowercase: ${result.hasLowercase ? '✅ Present' : '❌ Missing'}`;
            document.getElementById('rtNumber').textContent = `Numbers: ${result.hasNumber ? '✅ Present' : '❌ Missing'}`;
            document.getElementById('rtSpecial').textContent = `Special Characters (@#$): ${result.hasSpecial ? '✅ Present' : '❌ Missing'}`;
            document.getElementById('rtPattern').textContent = `Common Patterns: ${result.hasCommonPattern ? '❌ Detected' : '✅ None'}`;

            // Format and style cracking time
            const formattedTime = formatCrackingTime(result.crackingTimeYears);
            crackingTimeText.textContent = formattedTime;
            crackingTimeBadge.textContent = result.crackingTimeYears < 1 ? 'Unsafe' : result.crackingTimeYears < 100 ? 'Moderate' : 'Safe';
            crackingTimeBadge.classList.remove('bg-red-500', 'bg-yellow-500', 'bg-green-500', 'text-white');
            if (result.crackingTimeYears < 1) {
                crackingTimeBadge.classList.add('bg-red-500', 'text-white');
            } else if (result.crackingTimeYears < 100) {
                crackingTimeBadge.classList.add('bg-yellow-500', 'text-white');
            } else {
                crackingTimeBadge.classList.add('bg-green-500', 'text-white');
            }

            const rtStrengthEl = document.getElementById('rtStrength');
            rtStrengthEl.textContent = `Strength: ${result.strength}`;
            if (result.strength === 'Strong') {
                rtStrengthEl.classList.remove('text-red-600', 'text-yellow-600');
                rtStrengthEl.classList.add('text-green-600');
                rtStrengthBar.style.width = '100%';
                rtStrengthBar.classList.remove('bg-red-500', 'bg-yellow-500');
                rtStrengthBar.classList.add('bg-green-500');
            } else if (result.strength === 'Medium') {
                rtStrengthEl.classList.remove('text-red-600', 'text-green-600');
                rtStrengthEl.classList.add('text-yellow-600');
                rtStrengthBar.style.width = '66%';
                rtStrengthBar.classList.remove('bg-red-500', 'bg-green-500');
                rtStrengthBar.classList.add('bg-yellow-500');
            } else {
                rtStrengthEl.classList.remove('text-yellow-600', 'text-green-600');
                rtStrengthEl.classList.add('text-red-600');
                rtStrengthBar.style.width = '33%';
                rtStrengthBar.classList.remove('bg-yellow-500', 'bg-green-500');
                rtStrengthBar.classList.add('bg-red-500');
            }

            if (result.suggestedPassword) {
                suggestedPasswordDiv.textContent = `Suggested Strong Password: ${result.suggestedPassword}`;
                suggestedPasswordDiv.classList.remove('hidden');
                suggestedPasswordDiv.classList.add('fade-in');
            } else {
                suggestedPasswordDiv.classList.add('hidden');
            }

            suggestPasswordBtn.classList.remove('hidden');
            checkAnotherBtn.classList.remove('hidden');

            errorsDiv.textContent = '';
        } catch (error) {
            errorsDiv.textContent = `Error: ${error.message}`;
            realtimeFeedbackDiv.classList.add('hidden');
            suggestedPasswordDiv.classList.add('hidden');
            rtStrengthBar.style.width = '0%';
            suggestPasswordBtn.classList.add('hidden');
            checkAnotherBtn.classList.add('hidden');
        }
    };

    passwordInput.addEventListener('input', (e) => {
        evaluateRealtime(e.target.value);
    });

    // Generate Suggested Password
    suggestPasswordBtn.addEventListener('click', async () => {
        try {
            const response = await fetch('/api/suggest', {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' }
            });

            if (!response.ok) {
                throw new Error(`HTTP error ${response.status}`);
            }

            const result = await response.json();
            passwordInput.value = result.suggestedPassword;
            passwordInput.type = 'text';
            togglePassword.querySelector('i').classList.remove('fa-eye');
            togglePassword.querySelector('i').classList.add('fa-eye-slash');
            togglePassword.setAttribute('aria-label', 'Hide password');
            errorsDiv.textContent = '';
            suggestedPasswordDiv.classList.add('hidden');
            rtStrengthBar.style.width = '0%';
            evaluateRealtime(result.suggestedPassword);
        } catch (error) {
            errorsDiv.textContent = `Error: ${error.message}`;
        }
    });

    // Check Another Password
    checkAnotherBtn.addEventListener('click', () => {
        passwordInput.value = '';
        passwordInput.type = 'text';
        togglePassword.querySelector('i').classList.remove('fa-eye');
        togglePassword.querySelector('i').classList.add('fa-eye-slash');
        togglePassword.setAttribute('aria-label', 'Hide password');
        realtimeFeedbackDiv.classList.add('hidden');
        suggestedPasswordDiv.classList.add('hidden');
        errorsDiv.textContent = '';
        rtStrengthBar.style.width = '0%';
        suggestPasswordBtn.classList.add('hidden');
        checkAnotherBtn.classList.add('hidden');
    });

    // Hide buttons initially
    suggestPasswordBtn.classList.add('hidden');
    checkAnotherBtn.classList.add('hidden');
});