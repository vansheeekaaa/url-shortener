const API_URL = "https://snip-url-api.onrender.com";

document.addEventListener('DOMContentLoaded', () => {
    const form           = document.getElementById('shorten-form');
    const urlInput       = document.getElementById('url-input');
    const expiryInput    = document.getElementById('expiry-input');
    const submitBtn      = document.getElementById('submit-btn');
    const resultContainer = document.getElementById('result-container');
    const shortUrlLink   = document.getElementById('short-url-link');
    const copyBtn        = document.getElementById('copy-btn');
    const statsBtn       = document.getElementById('stats-btn');
    const statsContainer = document.getElementById('stats');
    const errorContainer = document.getElementById('error-message');

    let currentShortCode = '';
    let statsInterval    = null;

    //Shorten
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        hideAll();
        setLoading(true);

        const url    = urlInput.value.trim();
        const expiry = parseInt(expiryInput.value, 10) || 0;

        try {
            const res  = await fetch(`${API_URL}/shorten`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ url, expiry_seconds: expiry }),
            });
            const data = await res.json();

            if (!res.ok) throw new Error(data.error || 'Failed to shorten URL');
            currentShortCode = data.short_code;
            showResult(data.short_url);

        } catch (err) {
            showError(err.message);
        } finally {
            setLoading(false);
        }
    });

    //Copy
    copyBtn.addEventListener('click', async () => {
        try {
            await navigator.clipboard.writeText(shortUrlLink.href);
            copyBtn.textContent = 'Copied!';
            copyBtn.classList.add('success');
            setTimeout(() => {
                copyBtn.textContent = 'Copy';
                copyBtn.classList.remove('success');
            }, 2000);
        } catch {
            prompt('Copy it manually:', shortUrlLink.href);
        }
    });

    //View Stats
    statsBtn.addEventListener('click', async () => {
        if (!currentShortCode) return;
        //If stats are already visible, hide and stop polling
        if (!statsContainer.classList.contains('hidden')) {
            stopPolling();
            statsContainer.classList.add('hidden');
            statsBtn.textContent = 'View Stats';
            return;
        }

        statsBtn.textContent = 'Loading…';
        statsBtn.disabled = true;

        try {
            const res  = await fetch(`${API_URL}/stats/${currentShortCode}`);
            const data = await res.json();

            if (!res.ok) throw new Error(data.error || 'Failed to load stats');

            renderStats(data);
            statsContainer.classList.remove('hidden');
            statsBtn.textContent = 'Hide Stats';

            //auto refresh every 10 secs while stats open
            startPolling();
        } catch (err) {
            showError(err.message);
            statsBtn.textContent = 'View Stats';
        } finally {
            statsBtn.disabled = false;
        }
    });

    //polling
    function startPolling() {
        stopPolling();
        statsInterval = setInterval(async () => {
            try {
                const res  = await fetch(`${API_URL}/stats/${currentShortCode}`);
                const data = await res.json();
                if (res.ok) renderStats(data);
            } catch {
                //silently ignore network blips
            }
        }, 10000);
    }

    function stopPolling() {
        clearInterval(statsInterval);
        statsInterval = null;
    }

    //stats card 
    function renderStats(data) {
        const fmt = (iso) => iso ? new Date(iso).toLocaleString() : '—';

        statsContainer.innerHTML = `
            <div class="stats-grid">
                <div class="stat-card stat-clicks">
                    <div class="live-badge"><span class="live-dot"></span>LIVE</div>
                    <span class="stat-value">${data.click_count}</span>
                    <span class="stat-label">Total Clicks</span>
                </div>
                <div class="stat-card">
                    <span class="stat-value">${fmt(data.created_at)}</span>
                    <span class="stat-label">Created</span>
                </div>
                <div class="stat-card">
                    <span class="stat-value">${fmt(data.last_accessed_at)}</span>
                    <span class="stat-label">Last Accessed</span>
                </div>
                <div class="stat-card">
                    <span class="stat-value">${data.expires_at ? fmt(data.expires_at) : 'Never'}</span>
                    <span class="stat-label">Expires</span>
                </div>
            </div>
        `;
    }

    //helpers
    function setLoading(on) {
        submitBtn.disabled    = on;
        submitBtn.textContent = on ? '…' : 'Snip';
    }

    function showResult(shortUrl) {
        shortUrlLink.href        = shortUrl;
        shortUrlLink.textContent = shortUrl;
        resultContainer.classList.remove('hidden');
    }

    function showError(msg) {
        errorContainer.textContent = msg;
        errorContainer.classList.remove('hidden');
    }

    function hideAll() {
        stopPolling();
        resultContainer.classList.add('hidden');
        statsContainer.classList.add('hidden');
        statsContainer.innerHTML = '';
        statsBtn.textContent = 'View Stats';
        errorContainer.classList.add('hidden');
    }
});
