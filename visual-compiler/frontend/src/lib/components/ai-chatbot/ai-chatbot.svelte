<script lang="ts">
    let isOpen = false;
    let messageInput = '';
    let messages: Array<{id: number, text: string, isUser: boolean, timestamp: Date}> = [];

    function toggleChatbot() {
        isOpen = !isOpen;
    }

    function handleSendMessage() {
        if (!messageInput.trim()) return;
        
        // Add user message (placeholder for now)
        messages = [...messages, {
            id: Date.now(),
            text: messageInput.trim(),
            isUser: true,
            timestamp: new Date()
        }];
        
        messageInput = '';
        
        // Placeholder for AI response (will be implemented later)
        setTimeout(() => {
            messages = [...messages, {
                id: Date.now() + 1,
                text: "This is a placeholder AI response. The actual AI functionality will be implemented soon!",
                isUser: false,
                timestamp: new Date()
            }];
        }, 1000);
    }

    function handleKeyPress(event: KeyboardEvent) {
        if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            handleSendMessage();
        }
    }
</script>

<!-- Floating chatbot button -->
<div class="chatbot-container">
    {#if !isOpen}
        <button class="chatbot-toggle" on:click={toggleChatbot} aria-label="Open AI Chatbot">
            <!-- Simple Robot Icon -->
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <!-- Robot head -->
                <rect x="6" y="5" width="12" height="10" rx="2" fill="currentColor"/>
                
                <!-- Robot antenna -->
                <circle cx="12" cy="3" r="1" fill="currentColor"/>
                <line x1="12" y1="4" x2="12" y2="5" stroke="currentColor" stroke-width="1.5"/>
                
                <!-- Robot eyes -->
                <circle cx="9" cy="8" r="1.5" fill="white"/>
                <circle cx="15" cy="8" r="1.5" fill="white"/>
                <circle cx="9" cy="8" r="0.7" fill="currentColor"/>
                <circle cx="15" cy="8" r="0.7" fill="currentColor"/>
                
                <!-- Robot mouth -->
                <rect x="10" y="11" width="4" height="1.5" rx="0.75" fill="white"/>
                
                <!-- Robot body -->
                <rect x="8" y="15" width="8" height="6" rx="1" fill="currentColor" opacity="0.8"/>
                
                <!-- Robot arms -->
                <rect x="4" y="16" width="3" height="1.5" rx="0.75" fill="currentColor"/>
                <rect x="17" y="16" width="3" height="1.5" rx="0.75" fill="currentColor"/>
                
                <!-- Robot chest panel -->
                <rect x="10" y="17" width="4" height="2" rx="0.5" fill="white" opacity="0.7"/>
            </svg>
        </button>
    {/if}

    <!-- Chat window -->
    {#if isOpen}
        <div class="chat-window">
            <!-- Header -->
            <div class="chat-header">
                <div class="chat-title">
                    <div class="ai-icon">
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <!-- Simple Robot Icon -->
                            <rect x="6" y="5" width="12" height="10" rx="2" fill="currentColor"/>
                            <circle cx="12" cy="3" r="1" fill="currentColor"/>
                            <line x1="12" y1="4" x2="12" y2="5" stroke="currentColor" stroke-width="1.5"/>
                            <circle cx="9" cy="8" r="1.5" fill="white"/>
                            <circle cx="15" cy="8" r="1.5" fill="white"/>
                            <circle cx="9" cy="8" r="0.7" fill="currentColor"/>
                            <circle cx="15" cy="8" r="0.7" fill="currentColor"/>
                            <rect x="10" y="11" width="4" height="1.5" rx="0.75" fill="white"/>
                        </svg>
                    </div>
                    <div>
                        <h3>AI Assistant</h3>
                        <span class="status">Online</span>
                    </div>
                </div>
                <button class="close-btn" on:click={toggleChatbot} aria-label="Close chat">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <path d="M18 6L6 18M6 6L18 18" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                    </svg>
                </button>
            </div>

            <!-- Messages area -->
            <div class="messages-container">
                {#if messages.length === 0}
                    <div class="welcome-message">
                        <div class="welcome-icon">
                            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <!-- Simple Robot Icon -->
                                <rect x="6" y="5" width="12" height="10" rx="2" fill="currentColor"/>
                                <circle cx="12" cy="3" r="1" fill="currentColor"/>
                                <line x1="12" y1="4" x2="12" y2="5" stroke="currentColor" stroke-width="1.5"/>
                                <circle cx="9" cy="8" r="1.5" fill="white"/>
                                <circle cx="15" cy="8" r="1.5" fill="white"/>
                                <circle cx="9" cy="8" r="0.7" fill="currentColor"/>
                                <circle cx="15" cy="8" r="0.7" fill="currentColor"/>
                                <rect x="10" y="11" width="4" height="1.5" rx="0.75" fill="white"/>
                                <rect x="8" y="15" width="8" height="6" rx="1" fill="currentColor" opacity="0.8"/>
                                <rect x="10" y="17" width="4" height="2" rx="0.5" fill="white" opacity="0.7"/>
                            </svg>
                        </div>
                        <h4>Hello! I'm your AI Assistant</h4>
                        <p>How can I help you with your Visual Compiler project today?</p>
                    </div>
                {:else}
                    {#each messages as message (message.id)}
                        <div class="message {message.isUser ? 'user-message' : 'ai-message'}">
                            <div class="message-content">
                                {message.text}
                            </div>
                            <div class="message-time">
                                {message.timestamp.toLocaleTimeString([], {hour: '2-digit', minute: '2-digit'})}
                            </div>
                        </div>
                    {/each}
                {/if}
            </div>

            <!-- Input area -->
            <div class="input-area">
                <div class="input-container">
                    <textarea
                        bind:value={messageInput}
                        placeholder="Type your message here..."
                        rows="1"
                        on:keypress={handleKeyPress}
                    ></textarea>
                    <button 
                        class="send-btn" 
                        on:click={handleSendMessage}
                        disabled={!messageInput.trim()}
                        aria-label="Send message"
                    >
                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <path d="M3 12L21 2L18 14L12 12M3 12L10 12M3 12L21 22" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                        </svg>
                    </button>
                </div>
                <div class="input-footer">
                    <span>AI Assistant • Powered by Visual Compiler</span>
                </div>
            </div>
        </div>
    {/if}
</div>

<style>
    .chatbot-container {
        position: fixed;
        bottom: 2rem;
        right: 2rem;
        z-index: 1000;
        font-family: Arial, sans-serif;
    }

    /* Floating button */
    .chatbot-toggle {
        width: 60px;
        height: 60px;
        border-radius: 50%;
        background: linear-gradient(135deg, #8451C7, #AFA2D7);
        border: none;
        color: white;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        box-shadow: 0 4px 20px rgba(132, 81, 199, 0.3);
        transition: all 0.3s ease;
        animation: pulse 2s infinite;
    }

    .chatbot-toggle:hover {
        transform: scale(1.05);
        box-shadow: 0 6px 25px rgba(132, 81, 199, 0.4);
    }

    @keyframes pulse {
        0% { box-shadow: 0 4px 20px rgba(132, 81, 199, 0.3); }
        50% { box-shadow: 0 4px 20px rgba(132, 81, 199, 0.5), 0 0 0 10px rgba(132, 81, 199, 0.1); }
        100% { box-shadow: 0 4px 20px rgba(132, 81, 199, 0.3); }
    }

    /* Chat window */
    .chat-window {
        width: 380px;
        height: 500px;
        background: white;
        border-radius: 12px;
        box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
        display: flex;
        flex-direction: column;
        overflow: hidden;
        animation: slideUp 0.3s ease-out;
    }

    @keyframes slideUp {
        from {
            opacity: 0;
            transform: translateY(20px) scale(0.95);
        }
        to {
            opacity: 1;
            transform: translateY(0) scale(1);
        }
    }

    /* Header */
    .chat-header {
        background: linear-gradient(135deg, #8451C7, #AFA2D7);
        color: white;
        padding: 1rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .chat-title {
        display: flex;
        align-items: center;
        gap: 0.75rem;
    }

    .ai-icon {
        width: 40px;
        height: 40px;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .chat-title h3 {
        margin: 0;
        font-size: 1.1rem;
        font-weight: 600;
    }

    .status {
        font-size: 0.8rem;
        opacity: 0.9;
    }

    .close-btn {
        background: none;
        border: none;
        color: white;
        cursor: pointer;
        padding: 0.5rem;
        border-radius: 4px;
        transition: background-color 0.2s ease;
    }

    .close-btn:hover {
        background: rgba(255, 255, 255, 0.1);
    }

    /* Messages */
    .messages-container {
        flex: 1;
        padding: 1rem;
        overflow-y: auto;
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .welcome-message {
        text-align: center;
        padding: 2rem 1rem;
        color: #6B7280;
    }

    .welcome-icon {
        margin: 0 auto 1rem;
        width: 60px;
        height: 60px;
        background: #F3E8FF;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #8451C7;
    }

    .welcome-message h4 {
        margin: 0 0 0.5rem 0;
        color: #374151;
        font-size: 1.1rem;
    }

    .welcome-message p {
        margin: 0;
        font-size: 0.9rem;
        line-height: 1.4;
    }

    .message {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .user-message {
        align-items: flex-end;
    }

    .ai-message {
        align-items: flex-start;
    }

    .message-content {
        max-width: 70%;
        padding: 0.75rem 1rem;
        border-radius: 18px;
        font-size: 0.9rem;
        line-height: 1.4;
    }

    .user-message .message-content {
        background: linear-gradient(135deg, #8451C7, #AFA2D7);
        color: white;
    }

    .ai-message .message-content {
        background: #F3E8FF;
        color: #4B3A6A;
    }

    .message-time {
        font-size: 0.75rem;
        color: #9CA3AF;
        padding: 0 0.5rem;
    }

    /* Input area */
    .input-area {
        border-top: 1px solid #E5E7EB;
        padding: 1rem;
    }

    .input-container {
        display: flex;
        gap: 0.5rem;
        align-items: flex-end;
    }

    .input-container textarea {
        flex: 1;
        border: 2px solid #AFA2D7;
        border-radius: 20px;
        padding: 0.75rem 1rem;
        font-size: 0.9rem;
        resize: none;
        outline: none;
        font-family: inherit;
        line-height: 1.4;
        max-height: 100px;
    }

    .input-container textarea:focus {
        border-color: #8451C7;
        box-shadow: 0 0 0 3px rgba(132, 81, 199, 0.1);
    }

    .send-btn {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        background: linear-gradient(135deg, #8451C7, #AFA2D7);
        border: none;
        color: white;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all 0.2s ease;
        flex-shrink: 0;
    }

    .send-btn:hover:not(:disabled) {
        transform: scale(1.05);
        box-shadow: 0 2px 10px rgba(132, 81, 199, 0.3);
    }

    .send-btn:disabled {
        opacity: 0.5;
        cursor: not-allowed;
        transform: none;
    }

    .input-footer {
        margin-top: 0.5rem;
        text-align: center;
    }

    .input-footer span {
        font-size: 0.75rem;
        color: #9CA3AF;
    }

    /* Dark mode support */
    :global(html.dark-mode) .chat-window {
        background: #1a202c;
        color: #E5E7EB;
    }

    :global(html.dark-mode) .welcome-message {
        color: #9CA3AF;
    }

    :global(html.dark-mode) .welcome-message h4 {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .welcome-icon {
        background: #2D1B69;
    }

    :global(html.dark-mode) .ai-message .message-content {
        background: #2D1B69;
        color: #C4B5FD;
    }

    :global(html.dark-mode) .input-area {
        border-top-color: #4A5568;
    }

    :global(html.dark-mode) .input-container textarea {
        background: #2d3748;
        color: #E5E7EB;
        border-color: #4a5568;
    }

    :global(html.dark-mode) .input-container textarea:focus {
        border-color: #8451C7;
    }

    /* Responsive design */
    @media (max-width: 480px) {
        .chatbot-container {
            bottom: 1rem;
            right: 1rem;
        }

        .chat-window {
            width: calc(100vw - 2rem);
            height: 80vh;
            max-width: 380px;
        }

        .chatbot-toggle {
            width: 50px;
            height: 50px;
        }
    }
</style>