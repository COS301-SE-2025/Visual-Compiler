<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { addToast } from '$lib/stores/toast';

  // bound source code text
  let codeText = '';

  // event to let parent know the code was confirmed
  const dispatch = createEventDispatcher<{ codeSubmitted: string }>();

  // read .txt file into the textarea
  function handleFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;

    // Check file extension
    if (!file.name.toLowerCase().endsWith('.txt')) {
      addToast('Only .txt files are allowed. Please upload a valid plain text file.', 'error');
      input.value = ''; // Reset the input
      return;
    }

    const reader = new FileReader();
    reader.onload = () => {
      codeText = reader.result as string;
      addToast('File uploaded successfully!', 'success');
    };
    reader.onerror = () => {
      addToast('Failed to read file.', 'error');
    };
    reader.readAsText(file);
  }

  // notify parent that code is ready
  function submitCode() {
    if (!codeText.trim()) return;
    dispatch('codeSubmitted', codeText);
    addToast('Code confirmed!', 'success');
  }
</script>

<div class="code-input-container">
  <textarea bind:value={codeText}
    placeholder="Paste or type your source code here…"
    rows="10">

  </textarea>
    


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
      disabled={!codeText.trim()}
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
    box-shadow: inset 0 1px 2px rgba(0,0,0,0.1);
    height: 86px;
  }

  .controls {
    display: flex;
    gap: 0.5rem;
    
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
  .upload-btn input[type="file"] {
    position: absolute;
    left: 0; top: 0;
    width: 100%; height: 100%;
    opacity: 0;
    cursor: pointer;
   
  }

  .confirm-btn {
    padding: 0.5rem 1.5rem;
    background: #041a47;
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
    background: #36486c;
  }
  .upload-btn:hover {
    background: #838386;
  }
</style>
