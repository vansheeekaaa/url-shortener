const API_URL = "https://url-shortener-api-gzty.onrender.com";

document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('shorten-form');
    const urlInput = document.getElementById('url-input');
    const expiryInput = document.getElementById('expiry-input');
    const submitBtn = document.getElementById('submit-btn');
    const resultContainer = document.getElementById('result-container');
    const shortUrlLink = document.getElementById('short-url-link');
    const copyBtn = document.getElementById('copy-btn');
    const errorContainer = document.getElementById('error-message');

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        // Reset state
        hideErrors();
        hideResult();
        setLoading(true);

        const url = urlInput.value.trim();
        const expiry = expiryInput.value ? parseInt(expiryInput.value, 10) : 0;

        try {
            const response = await fetch(`${API_URL}/shorten`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    url: url,
                    expiry_seconds: expiry
                })
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Failed to shorten URL');
            }

            showResult(data.short_url);
            
        } catch (err) {
            showError(err.message);
        } finally {
            setLoading(false);
        }
    });

    copyBtn.addEventListener('click', async () => {
        const url = shortUrlLink.href;
        try {
            await navigator.clipboard.writeText(url);
            
            // Visual feedback
            const originalText = copyBtn.textContent;
            copyBtn.textContent = 'Copied!';
            copyBtn.classList.add('success');
            
            setTimeout(() => {
                copyBtn.textContent = originalText;
                copyBtn.classList.remove('success');
            }, 2000);
        } catch (err) {
            console.error('Failed to copy text: ', err);
            prompt("Failed to copy automatically. Copy it here:", url);
        }
    });

    function setLoading(isLoading) {
        if (isLoading) {
            submitBtn.disabled = true;
            submitBtn.textContent = '...';
        } else {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Snip';
        }
    }

    function showResult(shortUrl) {
        shortUrlLink.href = shortUrl;
        shortUrlLink.textContent = shortUrl;
        resultContainer.classList.remove('hidden');
    }

    function hideResult() {
        resultContainer.classList.add('hidden');
    }

    function showError(message) {
        errorContainer.textContent = message;
        errorContainer.classList.remove('hidden');
    }

    function hideErrors() {
        errorContainer.classList.add('hidden');
    }
});
