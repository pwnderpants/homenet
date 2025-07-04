// AI chat interface - using shared utilities
// This file depends on utils.js being loaded first

let isProcessing = false;
let abortController = null;

// Initialize the chat interface
document.addEventListener('DOMContentLoaded', function() {
    Logger.info('AI chat interface initializing...');
    
    const form = document.getElementById('ai-form');
    const input = document.getElementById('ai-input');
    
    // Set up form submission
    form.addEventListener('submit', handleSubmit);
    
    // Set up stop button
    const stopButton = document.getElementById('stop-button');
    stopButton.addEventListener('click', stopRequest);
    
    // Set up keyboard shortcuts
    input.addEventListener('keydown', function(e) {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            if (!isProcessing) {
                form.dispatchEvent(new Event('submit'));
            }
        }
    });
    
    // Auto-resize textarea
    input.addEventListener('input', function() {
        this.style.height = 'auto';
        this.style.height = Math.min(this.scrollHeight, 120) + 'px';
    });
    
    Logger.info('AI chat interface initialized successfully');
});

function handleSubmit(e) {
    e.preventDefault();
    
    if (isProcessing) {
        Logger.warn('Form submission ignored - already processing');
        return;
    }
    
    const input = document.getElementById('ai-input');
    const message = input.value.trim();
    
    if (!message) {
        Logger.warn('Empty message submitted');
        return;
    }
    
    Logger.info('Form submitted with message:', message);
    isProcessing = true;
    handleRequest(message);
}

function handleRequest(message) {
    Logger.info('Starting AI request for message:', message);
    
    const input = document.getElementById('ai-input');
    const output = document.getElementById('ai-output');
    const loadingIndicator = document.getElementById('loading-indicator');
    const sendButton = document.getElementById('send-button');
    const stopButton = document.getElementById('stop-button');
    
    // Add user message
    addUserMessage(message);
    
    // Clear input
    input.value = '';
    input.style.height = 'auto';
    
    // Show loading indicator and stop button, hide send button
    loadingIndicator.classList.remove('hidden');
    sendButton.classList.add('hidden');
    stopButton.classList.remove('hidden');
    
    // Create abort controller for cancellation
    abortController = new AbortController();
    
    // Create form data as URL-encoded string
    const formData = new URLSearchParams();
    formData.append('prompt', message);
    
    // Create AI response container
    const aiResponseId = 'ai-response-' + Date.now();
    addAIResponseContainer(aiResponseId);
    
    Logger.info('Sending request to /ai/query with message:', message);
    
    // Make request to the non-streaming endpoint
    fetch('/ai/query', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData,
        signal: abortController.signal
    })
    .then(response => {
        Logger.debug('Response received, status:', response.status);
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.text();
    })
    .then(responseText => {
        Logger.debug('Raw response received, length:', responseText.length);
        
        const formattedResponse = formatResponseText(responseText);
        Logger.debug('Response formatted successfully');
        
        const aiResponseElement = document.getElementById(aiResponseId);
        aiResponseElement.innerHTML = formattedResponse;
        Logger.info('AI response displayed successfully');
        resetUI();
    })
    .catch(error => {
        Logger.error('Request error:', error);
        const aiResponseElement = document.getElementById(aiResponseId);
        
        // Check if the error is due to user cancellation
        if (error.name === 'AbortError') {
            Logger.info('Request cancelled by user');
            aiResponseElement.innerHTML = '<span class="text-gray-400">Request cancelled by user</span>';
        } else {
            Logger.error('Connection error:', error.message);
            aiResponseElement.innerHTML = '<span class="text-red-400">Error: Failed to connect to AI service</span>';
        }
        
        resetUI();
    });
}

function addUserMessage(message) {
    Logger.debug('Adding user message to chat');
    const output = document.getElementById('ai-output');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'mb-4';
    messageDiv.innerHTML = `<b class="text-blue-400">You:</b> <span class="text-gray-300">${escapeHtml(message)}</span>`;
    output.appendChild(messageDiv);
    scrollToBottom();
}

function addAIResponse(html) {
    Logger.debug('Adding AI response to chat');
    const output = document.getElementById('ai-output');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'mb-4';
    messageDiv.innerHTML = `<b class="text-green-400">AI:</b> <span class="text-gray-300">${html}</span>`;
    output.appendChild(messageDiv);
    scrollToBottom();
}

function addAIResponseContainer(id) {
    Logger.debug('Creating AI response container with ID:', id);
    const output = document.getElementById('ai-output');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'mb-4';
    messageDiv.innerHTML = `<b class="text-green-400">AI:</b> <div id="${id}" class="text-gray-300 mt-2"></div>`;
    output.appendChild(messageDiv);
    scrollToBottom();
}

function addErrorMessage(message) {
    Logger.warn('Adding error message to chat:', message);
    const output = document.getElementById('ai-output');
    const messageDiv = document.createElement('div');
    messageDiv.className = 'mb-4';
    messageDiv.innerHTML = `<b class="text-red-400">Error:</b> <span class="text-red-400">${message}</span>`;
    output.appendChild(messageDiv);
    scrollToBottom();
}

function scrollToBottom() {
    const output = document.getElementById('ai-output');
    const scrollContainer = output.parentElement; // This is the div with overflow-y-auto
    setTimeout(() => {
        scrollContainer.scrollTo({
            top: scrollContainer.scrollHeight,
            behavior: 'smooth'
        });
    }, 10);
}

function clearChat() {
    Logger.info('Clearing chat history');
    const output = document.getElementById('ai-output');
    output.innerHTML = '<div class="text-gray-300">Welcome! Ask me anything and I\'ll help you out.</div>';
}

function stopRequest() {
    if (abortController) {
        Logger.info('Stopping request via abort controller');
        abortController.abort();
    }
    resetUI();
}

function resetUI() {
    Logger.debug('Resetting UI state');
    const loadingIndicator = document.getElementById('loading-indicator');
    const sendButton = document.getElementById('send-button');
    const stopButton = document.getElementById('stop-button');
    const input = document.getElementById('ai-input');
    
    loadingIndicator.classList.add('hidden');
    sendButton.classList.remove('hidden');
    stopButton.classList.add('hidden');
    isProcessing = false;
    input.focus();
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatResponseText(text) {
    Logger.debug('Formatting response text, length:', text.length);
    
    // Simple markdown-like formatting without external library
    let formattedText = text;
    
    // Convert newlines to <br> tags
    formattedText = formattedText.replace(/\n/g, '<br>');
    
    // Basic markdown formatting
    // Bold: **text** or __text__
    formattedText = formattedText.replace(/\*\*(.*?)\*\*/g, '<strong class="font-bold text-white">$1</strong>');
    formattedText = formattedText.replace(/__(.*?)__/g, '<strong class="font-bold text-white">$1</strong>');
    
    // Italic: *text* or _text_
    formattedText = formattedText.replace(/\*(.*?)\*/g, '<em class="italic">$1</em>');
    formattedText = formattedText.replace(/_(.*?)_/g, '<em class="italic">$1</em>');
    
    // Inline code: `code`
    formattedText = formattedText.replace(/`([^`]+)`/g, '<code class="bg-gray-800 px-2 py-1 rounded text-green-400 text-sm">$1</code>');
    
    // Headers: # Header
    formattedText = formattedText.replace(/^### (.*$)/gm, '<h3 class="text-xl font-bold text-white mb-4">$1</h3>');
    formattedText = formattedText.replace(/^## (.*$)/gm, '<h2 class="text-2xl font-bold text-white mb-4">$1</h2>');
    formattedText = formattedText.replace(/^# (.*$)/gm, '<h1 class="text-3xl font-bold text-white mb-4">$1</h1>');
    
    // Handle markdown links: [text](url)
    formattedText = formattedText.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer" class="text-blue-400 hover:text-blue-300 underline">$1</a>');
    
    // Make plain URLs clickable LAST (after all other formatting)
    // More comprehensive URL regex that handles various URL formats
    const urlRegex = /(https?:\/\/(?:www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b(?:[-a-zA-Z0-9()@:%_\+.~#?&//=]*))/g;
    const matches = formattedText.match(urlRegex);
    
    if (matches && matches.length > 0) {
        Logger.debug('Found URLs in response:', matches.length);
        formattedText = formattedText.replace(urlRegex, '<a href="$1" target="_blank" rel="noopener noreferrer" class="text-blue-400 hover:text-blue-300 underline">$1</a>');
    } else {
        // Try a simpler regex as fallback
        const simpleUrlRegex = /(https?:\/\/[^\s<>"']+)/g;
        const simpleMatches = formattedText.match(simpleUrlRegex);
        if (simpleMatches && simpleMatches.length > 0) {
            Logger.debug('Found URLs with simple regex:', simpleMatches.length);
            formattedText = formattedText.replace(simpleUrlRegex, '<a href="$1" target="_blank" rel="noopener noreferrer" class="text-blue-400 hover:text-blue-300 underline">$1</a>');
        }
    }
    
    Logger.debug('Response formatting completed');
    return formattedText;
}

// Initialize log level
initializeLogLevel('ai_log_level', 'INFO'); 