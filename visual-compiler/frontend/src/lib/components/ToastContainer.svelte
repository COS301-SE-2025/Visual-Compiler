<script lang="ts">
  import { toasts } from '$lib/stores/toast';
  import { fly } from 'svelte/transition';

  export let position = 'top-right'; // future-proofing for placement if needed
</script>

<style>
.toast-container {
  position: fixed;
  top: 1.5rem;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  z-index: 9999;
  align-items: center;
}

.toast {
  min-width: 260px;
  max-width: 90vw;
  padding: 1rem 1.25rem;
  border-radius: 8px;
  font-size: 0.9rem;
  font-weight: 500;
  color: white;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  animation: fadeIn 0.3s ease;
}

.success {
  background-color: #36486c;
}

.error {
  background-color: #36486c;
}

.info {
  background-color: #36486c;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>

<div class="toast-container">
  {#each $toasts as toast (toast.id)}
    <div class="toast {toast.type}" transition:fly={{ y: -10, duration: 200 }}>
      {toast.message}
    </div>
  {/each}
</div>
