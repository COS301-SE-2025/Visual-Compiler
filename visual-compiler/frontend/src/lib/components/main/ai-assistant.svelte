<script lang="ts">
    import { activePhase, setActivePhase, type PhaseType } from '../../stores/pipeline';
    import { AddToast } from '$lib/stores/toast';
    import { get } from 'svelte/store';
    import { projectName } from '$lib/stores/project';
    
    let isOpen = false;
    let activeTab: 'questions' | 'generate' = 'questions';
    let messageInput = '';
    let messages: Array<{id: number, text: string, isUser: boolean, timestamp: Date}> = [];
    
    // Subscribe to the active phase
    $: currentPhase = $activePhase;
    
    // Phase configuration for generate input
    const phaseConfig = {
        source: {
            name: 'Source Code',
            icon: 'M14 2H6A2 2 0 0 0 4 4V20A2 2 0 0 0 6 22H18A2 2 0 0 0 20 20V8L14 2Z',
            iconExtra: 'M14 2V8H20',
            description: 'Generate sample source code for compilation'
        },
        lexer: {
            name: 'Lexer Rules',
            icon: 'M4 7V4A2 2 0 0 1 6 2H18A2 2 0 0 1 20 4V7',
            iconExtra: 'M20 7H4L2 19H22L20 7ZM8 12V16M16 12V16',
            description: 'Generate lexical analysis rules and tokens'
        },
        parser: {
            name: 'Parser Grammar',
            icon: 'M22 12H18L15 21L9 3L6 12H2',
            iconExtra: '',
            description: 'Generate parsing grammar and syntax rules'
        },
        analyser: {
            name: 'Analyzer Config',
            icon: 'M9 11H15M9 15H15M17 21H7A2 2 0 0 1 5 19V5A2 2 0 0 1 7 3H14L19 8V19A2 2 0 0 1 17 21Z',
            iconExtra: '',
            description: 'Generate semantic analysis configuration'
        },
        translator: {
            name: 'Translator',
            icon: 'M14.5 2H6A2 2 0 0 0 4 4V20A2 2 0 0 0 6 22H18A2 2 0 0 0 20 20V7.5L14.5 2Z',
            iconExtra: 'M14 2V8H20M9 15H15M9 11H12',
            description: 'Generate code translation templates'
        },
        optimizer: {
            name: 'Optimizer',
            icon: 'M12 2L2 7V10C2 16 6 20.5 12 22C18 20.5 22 16 22 10V7L12 2Z',
            iconExtra: 'M9 12L11 14L15 10',
            description: 'Generate optimization rules and configurations'
        }
    };

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
                text: "This is a placeholder AI response",
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
    
    // Replace the generatePhaseInput function with this implementation:
    async function generatePhaseInput(phase: PhaseType) {
        if (!phase) return;
        
        const accessToken = sessionStorage.getItem('access_token') || 
                           sessionStorage.getItem('authToken') || 
                           localStorage.getItem('access_token') || 
                           localStorage.getItem('authToken') || 
                           localStorage.getItem('token');
        
        console.log('=== AI Generation Debug Info ===');
        console.log('Phase:', phase);
        console.log('Access Token exists:', !!accessToken);
        console.log('Access Token (first 20 chars):', accessToken ? accessToken.substring(0, 20) + '...' : 'null');
        
        if (!accessToken) {
            AddToast('Authentication required: Please log in to use AI generation', 'error');
            return;
        }

        // Add user message to show generation started
        messages = [...messages, {
            id: Date.now(),
            text: `Generate ${phaseConfig[phase].name} input`,
            isUser: true,
            timestamp: new Date()
        }];
        
        // Switch to questions tab to show the interaction
        activeTab = 'questions';

        try {
            // Show loading message
            const loadingMessageId = Date.now() + 1;
            messages = [...messages, {
                id: loadingMessageId,
                text: `Generating ${phaseConfig[phase].name} input...`,
                isUser: false,
                timestamp: new Date()
            }];

            // Prepare request body
            const requestBody = {
                phase: phase,
                artefact: " " 
            };
            
            console.log('Request body:', JSON.stringify(requestBody, null, 2));
            console.log('Request URL:', 'http://localhost:8080/api/ai/generate');

            const response = await fetch('http://localhost:8080/api/ai/generate', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${accessToken}`
                },
                body: JSON.stringify(requestBody)
            });

            console.log('Response status:', response.status);
            console.log('Response ok:', response.ok);
            console.log('Response headers:', Object.fromEntries(response.headers.entries()));

            const responseText = await response.text();
            console.log('Raw response text:', responseText);

            let data;
            try {
                data = JSON.parse(responseText);
                console.log('Parsed response data:', data);
            } catch (parseError) {
                console.error('JSON parse error:', parseError);
                throw new Error(`Invalid JSON response: ${responseText}`);
            }

            if (!response.ok) {
                console.error('Error response data:', data);
                throw new Error(data.error || data.details || `HTTP ${response.status} error`);
            }
            
            // Remove loading message and add success message
            messages = messages.filter(msg => msg.id !== loadingMessageId);
            
            if (data.response) {
                console.log('AI generated response:', data.response);
                
                // For source code phase, trigger the code input update
                if (phase === 'source') {
                    console.log('Dispatching ai-source-generated event');
                    // Dispatch a custom event that the code-input component can listen to
                    window.dispatchEvent(new CustomEvent('ai-source-generated', {
                        detail: { code: data.response }
                    }));
                    
                    messages = [...messages, {
                        id: Date.now() + 2,
                        text: `✅ Source code generated successfully! The code has been automatically inserted into your source code input area. You can review and modify it as needed.`,
                        isUser: false,
                        timestamp: new Date()
                    }];
                    
                    AddToast('AI source code generated and inserted successfully!', 'success');
                } else {
                    // For other phases (to be implemented later)
                    messages = [...messages, {
                        id: Date.now() + 2,
                        text: `Here's the generated ${phaseConfig[phase].name} input:\n\n${data.response}`,
                        isUser: false,
                        timestamp: new Date()
                    }];
                    
                    AddToast(`${phaseConfig[phase].name} input generated successfully!`, 'success');
                }
            } else {
                throw new Error('No content generated in response');
            }
            
        } catch (error) {
            console.error('=== AI Generation Error ===');
            console.error('Error object:', error);
            console.error('Error message:', error.message);
            console.error('Error stack:', error.stack);
            
            // Remove loading message if it exists
            messages = messages.filter(msg => !msg.text.includes('Generating'));
            
            messages = [...messages, {
                id: Date.now() + 3,
                text: `❌ Failed to generate ${phaseConfig[phase].name} input: ${error.message}. Please try again.`,
                isUser: false,
                timestamp: new Date()
            }];
            
            AddToast(`AI generation failed: ${error.message}`, 'error');
        }
    }

    function generatePhaseInputOld(phase: PhaseType) {
        if (!phase) return;
        
        // Add placeholder functionality for generating input
        messages = [...messages, {
            id: Date.now(),
            text: `Generate input for ${phaseConfig[phase].name}`,
            isUser: true,
            timestamp: new Date()
        }];
        
        // Switch to questions tab to show the generated content
        activeTab = 'questions';
        
        // Placeholder AI response for input generation
        setTimeout(() => {
            messages = [...messages, {
                id: Date.now() + 1,
                text: `Here's a sample input for the ${phaseConfig[phase].name} phase. (This will be replaced with actual AI-generated content)`,
                isUser: false,
                timestamp: new Date()
            }];
        }, 1000);
    }
</script>

<!-- Floating chatbot button -->
<div class="chatbot-container">
    {#if !isOpen}
        <button class="chatbot-toggle" on:click={toggleChatbot} aria-label="Open AI Chatbot">
            <!-- Robot Icon based on your image -->
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <!-- Robot antenna -->
                <circle cx="12" cy="2.5" r="1.5" stroke="currentColor" stroke-width="2" fill="none"/>
                <circle cx="12" cy="2.5" r="0.5" fill="currentColor"/>
                <line x1="12" y1="4" x2="12" y2="6" stroke="currentColor" stroke-width="2"/>
                
                <!-- Robot head -->
                <rect x="5" y="6" width="14" height="10" rx="3" stroke="currentColor" stroke-width="2" fill="none"/>
                
                <!-- Robot eyes -->
                <circle cx="9" cy="10" r="1.5" fill="currentColor"/>
                <circle cx="9" cy="10" r="0.5" fill="white"/>
                <circle cx="15" cy="10" r="1.5" fill="currentColor"/>
                <circle cx="15" cy="10" r="0.5" fill="white"/>
                
                <!-- Robot mouth -->
                <path d="M10 13 Q12 15 14 13" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round"/>
                
                <!-- Robot arms/handles -->
                <rect x="2" y="11" width="3" height="4" rx="1.5" stroke="currentColor" stroke-width="2" fill="none"/>
                <rect x="19" y="11" width="3" height="4" rx="1.5" stroke="currentColor" stroke-width="2" fill="none"/>
                
                <!-- Robot body -->
                <rect x="6" y="16" width="12" height="6" rx="2" stroke="currentColor" stroke-width="2" fill="none"/>
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
                            <!-- Robot Icon matching your image -->
                            <circle cx="12" cy="2.5" r="1.5" stroke="currentColor" stroke-width="1.5" fill="none"/>
                            <circle cx="12" cy="2.5" r="0.3" fill="currentColor"/>
                            <line x1="12" y1="4" x2="12" y2="6" stroke="currentColor" stroke-width="1.5"/>
                            <rect x="5" y="6" width="14" height="10" rx="3" stroke="currentColor" stroke-width="1.5" fill="none"/>
                            <circle cx="9" cy="10" r="1" fill="currentColor"/>
                            <circle cx="9" cy="10" r="0.3" fill="white"/>
                            <circle cx="15" cy="10" r="1" fill="currentColor"/>
                            <circle cx="15" cy="10" r="0.3" fill="white"/>
                            <path d="M10 13 Q12 14.5 14 13" stroke="currentColor" stroke-width="1.5" fill="none" stroke-linecap="round"/>
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
                                    <svg width="40" height="40" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <!-- Question mark circle background -->
                                        <circle cx="12" cy="12" r="10" fill="currentColor" opacity="0.1"/>
                                        <!-- Question mark shape -->
                                        <path d="M9.5 9C9.5 7.5 10.5 6.5 12 6.5C13.5 6.5 14.5 7.5 14.5 9C14.5 10.5 12 11 12 13" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                                        <!-- Question mark dot -->
                                        <circle cx="12" cy="16" r="1" fill="currentColor"/>
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
                                placeholder="Ask me about compilers"
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
                        {#if currentPhase}
                            <!-- Show current active phase -->
                            <div class="generate-welcome">
                                <div class="generate-icon">
                                    <svg width="48" height="48" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d={phaseConfig[currentPhase].icon} stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                        {#if phaseConfig[currentPhase].iconExtra}
                                            <path d={phaseConfig[currentPhase].iconExtra} stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                        {/if}
                                    </svg>
                                </div>
                                <h4>{phaseConfig[currentPhase].name}</h4>
                                <p class="phase-description">{phaseConfig[currentPhase].description}</p>
                            </div>
                            
                            <div class="generate-action-section">
                                <div class="action-header">
                                    <h5>Ready to Generate Input</h5>
                                    <p>Click the button below to generate sample input for this phase</p>
                                </div>
                                
                                <button 
                                    class="generate-main-btn"
                                    on:click={() => generatePhaseInput(currentPhase)}
                                >
                                    <div class="btn-icon">
                                        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                            <path d="M12 2L15.09 8.26L22 9L17 14L18.18 21L12 17.77L5.82 21L7 14L2 9L8.91 8.26L12 2Z" fill="currentColor"/>
                                        </svg>
                                    </div>
                                    <div class="btn-content">
                                        <span class="btn-title">Generate {phaseConfig[currentPhase].name} Input</span>
                                        <span class="btn-subtitle">Create sample content for this phase</span>
                                    </div>
                                    <div class="btn-arrow">
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                            <path d="M5 12H19M19 12L12 5M19 12L12 19" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                                        </svg>
                                    </div>
                                </button>
                                
                                <div class="help-text">
                                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
                                        <path d="M9.09 9A3 3 0 0 1 12 6C13.66 6 15 7.34 15 9C15 10.66 13.66 12 12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                                        <circle cx="12" cy="16" r="1" fill="currentColor"/>
                                    </svg>
                                    <span>The AI will generate appropriate sample input based on your current phase requirements</span>
                                </div>
                            </div>
                        {:else}
                            <!-- No active phase message -->
                            <div class="no-phase-message">
                                <div class="no-phase-icon">
                                    <svg width="40" height="40" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
                                        <path d="M8 12L16 12" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
                                    </svg>
                                </div>
                                <h4>No Active Phase</h4>
                                <p>Please navigate to a compiler phase in the main workspace to generate input for that specific phase.</p>
                                <div class="available-phases-info">
                                    <h6>Available phases:</h6>
                                    <ul>
                                        <li>Source Code</li>
                                        <li>Lexer Rules</li>
                                        <li>Parser Grammar</li>
                                        <li>Analyzer Config</li>
                                        <li>Translator</li>
                                        <li>Optimiser</li>
                                    </ul>
                                </div>
                            </div>
                        {/if}
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
        z-index: 2100;
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

    /* Adjust container when no active phase to prevent scrolling */
    .generate-container:has(.no-phase-message) {
        padding: 1rem;
        gap: 1rem;
        overflow-y: hidden;
        justify-content: center;
    }

    .generate-welcome {
        text-align: center;
        color: #6B7280;
        padding: 2rem 1rem;
        border-bottom: 1px solid #E5E7EB;
        background: linear-gradient(135deg, #F8FAFC, #F3F4F6);
        border-radius: 8px 8px 0 0;
        margin: -1.5rem -1.5rem 2rem -1.5rem;
    }

    .generate-icon {
        margin: 0 auto 1.5rem;
        width: 80px;
        height: 80px;
        background: linear-gradient(135deg, #F3E8FF, #E0E7FF);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #8451C7;
        box-shadow: 0 4px 12px rgba(132, 81, 199, 0.15);
    }

    .generate-welcome h4 {
        margin: 0 0 1rem 0;
        color: #374151;
        font-size: 1.5rem;
        font-weight: 700;
        letter-spacing: -0.025em;
    }

    .phase-description {
        margin: 0;
        font-size: 1rem;
        line-height: 1.6;
        color: #6B7280;
        max-width: 280px;
        margin: 0 auto;
    }

    /* Generate Action Section */
    .generate-action-section {
        display: flex;
        flex-direction: column;
        gap: 2rem;
    }

    .action-header {
        text-align: center;
    }

    .action-header h5 {
        margin: 0 0 0.5rem 0;
        color: #374151;
        font-size: 1.25rem;
        font-weight: 600;
    }

    .action-header p {
        margin: 0;
        color: #6B7280;
        font-size: 0.95rem;
        line-height: 1.5;
    }

    .generate-main-btn {
        display: flex;
        align-items: center;
        gap: 1rem;
        width: 100%;
        padding: 1.5rem;
        background: linear-gradient(135deg, #8451C7, #AFA2D7);
        color: white;
        border: none;
        border-radius: 12px;
        cursor: pointer;
        transition: all 0.3s ease;
        box-shadow: 0 4px 12px rgba(132, 81, 199, 0.25);
        position: relative;
        overflow: hidden;
    }

    .generate-main-btn::before {
        content: '';
        position: absolute;
        top: 0;
        left: -100%;
        width: 100%;
        height: 100%;
        background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
        transition: left 0.5s ease;
    }

    .generate-main-btn:hover::before {
        left: 100%;
    }

    .generate-main-btn:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 25px rgba(132, 81, 199, 0.35);
    }

    .btn-icon {
        width: 48px;
        height: 48px;
        background: rgba(255, 255, 255, 0.2);
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
    }

    .btn-content {
        flex: 1;
        text-align: left;
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .btn-title {
        font-size: 1.1rem;
        font-weight: 600;
        line-height: 1.2;
    }

    .btn-subtitle {
        font-size: 0.9rem;
        opacity: 0.9;
        font-weight: 400;
        line-height: 1.3;
    }

    .btn-arrow {
        flex-shrink: 0;
        transition: transform 0.3s ease;
    }

    .generate-main-btn:hover .btn-arrow {
        transform: translateX(4px);
    }

    .help-text {
        display: flex;
        align-items: flex-start;
        gap: 0.75rem;
        background: #F8FAFC;
        border: 1px solid #E5E7EB;
        border-radius: 8px;
        padding: 1rem;
        color: #6B7280;
        font-size: 0.9rem;
        line-height: 1.5;
    }

    .help-text svg {
        flex-shrink: 0;
        margin-top: 0.1rem;
        color: #8451C7;
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

    /* Current Phase Section Styles - REMOVED (no longer used) */

    /* No Phase Message Styles */
    .no-phase-message {
        text-align: center;
        padding: 1.5rem;
        color: #6B7280;
        max-height: fit-content;
    }

    .no-phase-icon {
        margin: 0 auto 1rem;
        width: 50px;
        height: 50px;
        background: #F3F4F6;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #9CA3AF;
    }

    .no-phase-message h4 {
        margin: 0 0 0.5rem 0;
        color: #374151;
        font-size: 1rem;
        font-weight: 600;
    }

    .no-phase-message p {
        margin: 0 0 1rem 0;
        font-size: 0.85rem;
        line-height: 1.4;
    }

    .available-phases-info {
        background: #F9FAFB;
        border-radius: 6px;
        padding: 0.75rem;
        text-align: left;
        margin-top: 0.5rem;
    }

    .available-phases-info h6 {
        margin: 0 0 0.3rem 0;
        color: #374151;
        font-size: 0.85rem;
        font-weight: 600;
    }

    .available-phases-info ul {
        margin: 0;
        padding-left: 1rem;
        color: #6B7280;
        font-size: 0.8rem;
    }

    .available-phases-info li {
        margin-bottom: 0.1rem;
    }

    /* Removed test buttons CSS - no longer needed */



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
        background: linear-gradient(135deg, #2d3748, #1a202c);
        border-bottom-color: #4A5568;
    }

    :global(html.dark-mode) .generate-welcome h4 {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .generate-icon {
        background: linear-gradient(135deg, #2D1B69, #1a1b3a);
    }

    :global(html.dark-mode) .phase-description {
        color: #A0AEC0;
    }

    :global(html.dark-mode) .action-header h5 {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .action-header p {
        color: #A0AEC0;
    }

    :global(html.dark-mode) .help-text {
        background: #2d3748;
        border-color: #4A5568;
        color: #A0AEC0;
    }

    :global(html.dark-mode) .help-text svg {
        color: #C4B5FD;
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

    /* Dark mode - Current Phase Section - REMOVED (no longer used) */

    /* Dark mode - No Phase Message */
    :global(html.dark-mode) .no-phase-message {
        color: #9CA3AF;
    }

    :global(html.dark-mode) .no-phase-message h4 {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .no-phase-icon {
        background: #2d3748;
        color: #6B7280;
    }

    :global(html.dark-mode) .available-phases-info {
        background: #2d3748;
        border-color: #4A5568;
    }

    :global(html.dark-mode) .available-phases-info h6 {
        color: #E5E7EB;
    }

    :global(html.dark-mode) .available-phases-info ul {
        color: #A0AEC0;
    }

    /* Removed dark mode test buttons CSS - no longer needed */

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