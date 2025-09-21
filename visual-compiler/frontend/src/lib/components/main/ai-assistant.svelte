<script lang="ts">
    let isOpen = false;
    let activeTab: 'questions' | 'generate' = 'questions';
    let messageInput = '';
    let messages: Array<{id: number, text: string, isUser: boolean, timestamp: Date}> = [];

    function toggleChatbot() {
        isOpen = !isOpen;
    }

    function switchTab(tab: 'questions' | 'generate') {
        activeTab = tab;
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

            <!-- Tab Navigation -->
            <div class="tab-navigation">
                <button 
                    class="tab-btn {activeTab === 'questions' ? 'active' : ''}"
                    on:click={() => switchTab('questions')}
                >
                   
                    Ask Questions
                </button>
                <button 
                    class="tab-btn {activeTab === 'generate' ? 'active' : ''}"
                    on:click={() => switchTab('generate')}
                >
                   
                    Generate Input
                </button>
            </div>

            <!-- Content Area -->
            <div class="content-area">
                {#if activeTab === 'questions'}
                    <!-- Questions Tab Content (Original Chat) -->
                    <div class="messages-container">
                        {#if messages.length === 0}
                            <div class="welcome-message">
                                <div class="welcome-icon">
                                    <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <!-- Phoenix body -->
                                        <path d="M12 3C12 3 8 5 8 9C8 11 9 12 10 13C9 14 8 15 8 17C8 21 12 23 12 23C12 23 16 21 16 17C16 15 15 14 14 13C15 12 16 11 16 9C16 5 12 3 12 3Z" fill="currentColor"/>
                                        <!-- Phoenix wings -->
                                        <path d="M6 7C4 9 3 11 4 13C5 12 6 11 8 10C7 9 6 8 6 7Z" fill="currentColor" opacity="0.8"/>
                                        <path d="M18 7C20 9 21 11 20 13C19 12 18 11 16 10C17 9 18 8 18 7Z" fill="currentColor" opacity="0.8"/>
                                        <!-- Phoenix tail feathers -->
                                        <path d="M12 23C11 21 10 19 11 17C12 18 12 19 12 20C12 19 12 18 13 17C14 19 13 21 12 23Z" fill="currentColor" opacity="0.7"/>
                                        <!-- Phoenix head crest -->
                                        <path d="M12 3C11 4 12 5 12 6C12 5 13 4 12 3Z" fill="currentColor"/>
                                    </svg>
                                </div>
                                <h4>Ask me anything!</h4>
                                <p>I'm here to help you understand compiler concepts, debug issues, or answer questions about your Visual Compiler project.</p>
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

                    <!-- Input area for Questions -->
                    <div class="input-area">
                        <div class="input-container">
                            <textarea
                                bind:value={messageInput}
                                placeholder="Ask me about compilers, syntax, or any coding help..."
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
                {:else}
                    <!-- Generate Input Tab Content -->
                    <div class="generate-container">
                        <div class="generate-welcome">
                            <div class="generate-icon">
                                <svg width="40" height="40" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                    <path d="M12 2L15.09 8.26L22 9L17 14L18.18 21L12 17.77L5.82 21L7 14L2 9L8.91 8.26L12 2Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                </svg>
                            </div>
                            <h4>Generate Phase Input</h4>
                            <p>Let me help you create sample input for your compiler phases.</p>
                        </div>
                        
                        <div class="phase-selection">
                            <h5>Available Phases:</h5>
                            <div class="phase-buttons">
                                <button class="phase-btn" disabled>
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M14 2H6A2 2 0 0 0 4 4V20A2 2 0 0 0 6 22H18A2 2 0 0 0 20 20V8L14 2Z" stroke="currentColor" stroke-width="2"/>
                                        <path d="M14 2V8H20" stroke="currentColor" stroke-width="2"/>
                                    </svg>
                                    Source Code
                                </button>
                                <button class="phase-btn" disabled>
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M4 7V4A2 2 0 0 1 6 2H18A2 2 0 0 1 20 4V7" stroke="currentColor" stroke-width="2"/>
                                        <path d="M20 7H4L2 19H22L20 7Z" stroke="currentColor" stroke-width="2"/>
                                        <path d="M8 12V16" stroke="currentColor" stroke-width="2"/>
                                        <path d="M16 12V16" stroke="currentColor" stroke-width="2"/>
                                    </svg>
                                    Lexer Rules
                                </button>
                                <button class="phase-btn" disabled>
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M22 12H18L15 21L9 3L6 12H2" stroke="currentColor" stroke-width="2"/>
                                    </svg>
                                    Parser Grammar
                                </button>
                                <button class="phase-btn" disabled>
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M9 11H15M9 15H15M17 21H7A2 2 0 0 1 5 19V5A2 2 0 0 1 7 3H14L19 8V19A2 2 0 0 1 17 21Z" stroke="currentColor" stroke-width="2"/>
                                    </svg>
                                    Analyzer Config
                                </button>
                                <button class="phase-btn" disabled>
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M14.5 2H6A2 2 0 0 0 4 4V20A2 2 0 0 0 6 22H18A2 2 0 0 0 20 20V7.5L14.5 2Z" stroke="currentColor" stroke-width="2"/>
                                        <path d="M14 2V8H20" stroke="currentColor" stroke-width="2"/>
                                        <path d="M9 15H15" stroke="currentColor" stroke-width="2"/>
                                        <path d="M9 11H12" stroke="currentColor" stroke-width="2"/>
                                    </svg>
                                    Translator
                                </button>
                                <button class="phase-btn" disabled>
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M13 2L3 14H12L11 22L21 10H12L13 2Z" stroke="currentColor" stroke-width="2"/>
                                    </svg>
                                    Optimizer
                                </button>
                            </div>
                        </div>

                       
                    </div>
                {/if}
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

    /* Tab Navigation */
    .tab-navigation {
        display: flex;
        background: #F8FAFC;
        border-bottom: 1px solid #E5E7EB;
    }

    .tab-btn {
        flex: 1;
        padding: 0.75rem 1rem;
        border: none;
        background: none;
        color: #6B7280;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        font-size: 0.85rem;
        font-weight: 500;
        transition: all 0.2s ease;
        border-bottom: 2px solid transparent;
    }

    .tab-btn:hover {
        color: #8451C7;
        background: rgba(132, 81, 199, 0.05);
    }

    .tab-btn.active {
        color: #8451C7;
        background: white;
        border-bottom-color: #8451C7;
    }

    /* Content Area */
    .content-area {
        flex: 1;
        display: flex;
        flex-direction: column;
        overflow: hidden;
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

    /* Generate Input Tab Styles */
    .generate-container {
        flex: 1;
        padding: 1.5rem;
        display: flex;
        flex-direction: column;
        gap: 2rem;
        overflow-y: auto;
    }

    .generate-welcome {
        text-align: center;
        color: #6B7280;
    }

    .generate-icon {
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

    .generate-welcome h4 {
        margin: 0 0 0.5rem 0;
        color: #374151;
        font-size: 1.1rem;
        font-weight: 600;
    }

    .generate-welcome p {
        margin: 0;
        font-size: 0.9rem;
        line-height: 1.5;
    }

    .phase-selection h5 {
        margin: 0 0 1rem 0;
        color: #374151;
        font-size: 1rem;
        font-weight: 600;
    }

    .phase-buttons {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
        gap: 0.75rem;
    }

    .phase-btn {
        padding: 1rem;
        border: 2px solid #E5E7EB;
        background: white;
        border-radius: 8px;
        cursor: pointer;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.85rem;
        font-weight: 500;
        color: #6B7280;
        transition: all 0.2s ease;
        text-align: center;
    }

    .phase-btn:hover:not(:disabled) {
        border-color: #8451C7;
        color: #8451C7;
        background: rgba(132, 81, 199, 0.02);
        transform: translateY(-1px);
    }

    .phase-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
        background: #F9FAFB;
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

    /* Dark mode - Tab Navigation */
    :global(html.dark-mode) .tab-navigation {
        background: #2d3748;
        border-bottom-color: #4A5568;
    }

    :global(html.dark-mode) .tab-btn {
        color: #A0AEC0;
    }

    :global(html.dark-mode) .tab-btn:hover {
        color: #C4B5FD;
        background: rgba(196, 181, 253, 0.1);
    }

    :global(html.dark-mode) .tab-btn.active {
        color: #C4B5FD;
        background: #1a202c;
        border-bottom-color: #8451C7;
    }

    /* Dark mode - Generate Input */
    :global(html.dark-mode) .generate-container {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .generate-welcome {
        color: #9CA3AF;
    }

    :global(html.dark-mode) .generate-welcome h4 {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .generate-icon {
        background: #2D1B69;
    }

    :global(html.dark-mode) .phase-selection h5 {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .phase-btn {
        background: #2d3748;
        border-color: #4A5568;
        color: #A0AEC0;
    }

    :global(html.dark-mode) .phase-btn:hover:not(:disabled) {
        border-color: #8451C7;
        color: #C4B5FD;
        background: rgba(196, 181, 253, 0.05);
    }

    :global(html.dark-mode) .phase-btn:disabled {
        background: #1a202c;
        border-color: #2d3748;
    }

    :global(html.dark-mode) .coming-soon {
        background: linear-gradient(135deg, rgba(132, 81, 199, 0.1), rgba(175, 162, 215, 0.1));
        border-color: rgba(132, 81, 199, 0.3);
        color: #C4B5FD;
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