<script setup lang="ts">
import { AlertCircle, CheckCircle2, Info, TriangleAlert, X } from '@lucide/vue'
import { useToast, type ToastTone } from '@/composables/useToast'

// Keep the collection as a top-level ref so Vue unwraps it in the template.
// Iterating `toast.items` would iterate the readonly ref object itself and
// render its internal fields as empty toast cards.
const { items, dismiss } = useToast()

const icons = {
  success: CheckCircle2,
  warning: TriangleAlert,
  error: AlertCircle,
  info: Info,
} satisfies Record<ToastTone, typeof Info>
</script>

<template>
  <div class="toast-host" aria-live="polite" aria-atomic="false">
    <TransitionGroup name="toast">
      <article v-for="item in items" :key="item.id" class="toast-card" :class="`toast-card--${item.tone}`" role="status">
        <component :is="icons[item.tone]" :size="20" :stroke-width="3" aria-hidden="true" />
        <p>{{ item.message }}</p>
        <button type="button" class="toast-close" aria-label="Close" @click="dismiss(item.id)">
          <X :size="18" :stroke-width="3" aria-hidden="true" />
        </button>
      </article>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-host {
  position: fixed;
  z-index: 1000;
  top: 20px;
  right: 20px;
  display: grid;
  width: min(360px, calc(100vw - 32px));
  gap: 12px;
  pointer-events: none;
}

.toast-card {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border: 3px solid var(--brutal-border-color);
  border-radius: var(--brutal-radius);
  background: var(--brutal-bg);
  color: var(--brutal-fg);
  box-shadow: 5px 5px 0 var(--brutal-shadow-color);
  pointer-events: auto;
}

.toast-card p { margin: 0; font-weight: 800; line-height: 1.4; }
.toast-card--success { background: #dbe9b8; }
.toast-card--warning { background: var(--brutal-accent); }
.toast-card--error { background: #f2b3aa; }
.toast-card--info { background: #b9dce8; }
.toast-close { display: grid; place-items: center; border: 0; background: transparent; color: inherit; cursor: pointer; }
.toast-close:focus-visible { outline: 3px solid var(--brutal-ring); outline-offset: 2px; }
.toast-enter-active, .toast-leave-active { transition: transform 160ms ease, opacity 160ms ease; }
.toast-enter-from, .toast-leave-to { transform: translateX(16px); opacity: 0; }

@media (max-width: 600px) {
  .toast-host { top: 12px; right: 16px; }
}

@media (prefers-reduced-motion: reduce) {
  .toast-enter-active, .toast-leave-active { transition: none; }
}
</style>
