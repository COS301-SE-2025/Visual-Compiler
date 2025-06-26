<script lang="ts">
  import { AddToast } from '$lib/stores/toast';

  let code_text = '';
  let isDefaultInput = false;
  let previous_code_text = '';

  export let onCodeSubmitted: (code: string) => void = () => {};

  function handleDefaultInput() {
    if (!isDefaultInput) {
      previous_code_text = code_text;
      code_text = 'int blue = 13 + 22;';
      isDefaultInput = true;
    } else {
      code_text = previous_code_text;
      isDefaultInput = false;
    }
  }

  // handleFileChange
  // Return type: void
  // Parameter type(s): Event
  // Handles the file upload event, validates the file type, and reads its content.
  function handleFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    if (!file.name.toLowerCase().endsWith('.txt')) {
      AddToast('Only .txt files are allowed. Please upload a valid plain text file.', 'error');
      input.value = ''; // Reset the input
      return;
    }

    const reader = new FileReader();
    reader.onload = () => {
      code_text = reader.result as string;
      AddToast('File uploaded successfully!', 'success');
    };
    reader.onerror = () => {
      AddToast('Failed to read file.', 'error');
    };
    reader.readAsText(file);
  }

  // submitCode
  // Return type: void
  // Parameter type(s): none
  // Dispatches the current code text to the parent component.
  function submitCode() {
    if (!code_text.trim()) return;
    const user_id = localStorage.getItem('user_id');
    if (!user_id) {
      AddToast('User not logged in.', 'error');
      return;
    }
    fetch('http://localhost:8080/api/lexing/code', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        users_id: user_id,
        source_code: code_text
      })
    })
      .then(res => {
        if (!res.ok) throw new Error('Failed to save source code');
        AddToast('Code confirmed and saved!', 'success');
        onCodeSubmitted(code_text);
      })
      .catch(() => AddToast('Failed to save source code', 'error'));
  }
</script>

<div class="code-input-container">
  <div class="code-input-header-row">
    <h2 class="code-input-header">Enter Source Code</h2>
    <button
      type="button"
      class="default-source-btn"
      title={isDefaultInput ? "Restore your input" : "Insert default source code"}
      aria-label={isDefaultInput ? "Restore your input" : "Insert default source code"}
      on:click={handleDefaultInput}
    >
      {#if isDefaultInput}
        🧹
      {:else}
        🪄
      {/if}
    </button>
  </div>
  <textarea
    bind:value={code_text}
    placeholder="Paste or type your source code here…"
    rows="10"
  ></textarea>

  <div class="controls">
    <label class="upload-btn" title="📝 Hi there! Please make sure to upload a .txt file. Other file types won't work here!">
      Upload File
      <input
        type="file"
        accept=".txt"
        on:change={handleFileChange}
        title="📝 Please upload a '.txt' file"
      />
    </label>

    <button
      type="button"
      class="confirm-btn"
      on:click={submitCode}
      disabled={!code_text.trim()}
    >
      Confirm Code
    </button>
  </div>
</div>

<style>
  .code-input-container {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
    margin-top: 1rem;
    margin-left: 1rem;
    margin-right: 1rem;
  }

  textarea {
    resize: vertical;
    padding: 0.4rem;
    font-family: Menlo, Monaco, 'Courier New', monospace;
    font-size: 0.9rem;
    line-height: 1.2;
    border: 1px solid #ccc;
    border-radius: 4px;
    box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.1);
    height: 100px;
  }

  .controls {
    display: flex;
    gap: 0.5rem;
    margin-top: 0.3rem;
    margin-bottom: 0.7rem;
    justify-content: center;
  }

  .upload-btn {
    position: relative;
    overflow: hidden;
    display: inline-block;
    padding: 0.5rem 1rem;
    background: #646468;
    color: white;
    border-radius: 4px;
    font-size: 0.95rem;
    cursor: pointer;
    font: inherit;
  }
  .upload-btn input[type='file'] {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%;
    opacity: 0;
    cursor: pointer;
  }

  .confirm-btn {
    padding: 0.5rem 1.5rem;
    background: #001a6e;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 0.95rem;
    cursor: pointer;
    transition: background 0.2s ease;
  }
  .confirm-btn:disabled {
    background: #ccc;
    cursor: not-allowed;
  }
  .confirm-btn:not(:disabled):hover {
    background: #074799;
  }
  .upload-btn:hover {
    background: #838386;
  }

  .default-source-btn {
    background: #e0e7ff;
    color: #1e40af;
    border: none;
    border-radius: 50%;
    width: 2.2rem;
    height: 2.2rem;
    font-size: 1.3rem;
    cursor: pointer;
    transition: background 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .default-source-btn:hover,
  .default-source-btn:focus {
    background: #d0d9ff;
  }

  .code-input-header-row {
    display: flex;
    align-items: center;
    gap: 0.7rem;
    margin-bottom: 0.3rem;
    justify-content: center;
  }

  .code-input-header {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 700;
  }
</style>