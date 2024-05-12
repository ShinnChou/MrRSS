<script lang="ts" setup>
import { Icon } from '@iconify/vue';

const props = defineProps({
  isRefreshing: Boolean
})
</script>

<template>
    <nav>
      <div class="util">
        <button class="btn" title="Filter">
          <Icon icon="material-symbols:filter-list" />
        </button>
      </div>
      <div class="function">
        <button class="btn" @click="$emit('handle-click-refresh')" title="Refresh">
          <Icon v-if="!isRefreshing" icon="material-symbols:refresh" />
          <Icon v-else icon="material-symbols:error" disabled="true" />
        </button>
        <button class="btn" title="Translate">
          <Icon icon="material-symbols:g-translate" />
        </button>
        <button class="btn" @click="$emit('delete-history-content')" title="Delete history">
          <Icon icon="material-symbols:auto-delete-outline" />
        </button>
        <button class="btn" title="Settings">
          <Icon icon="material-symbols:settings" />
        </button>
      </div>
    </nav>
</template>

<style scoped>
nav {
  position: sticky;
  top: 0;
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #e0e0e0;
}

.util, .function {
  display: flex;
  align-items: center;
}

.btn {
  width: 32px;
  height: 32px;
  font-size: large;
  background-color: #e0e0e0;
  border: none;
  cursor: pointer;

  display: flex;
  justify-content: center;
  align-items: center;
}

.btn:hover {
  background-color: #d0d0d0;
}

p {
  height: 32px;
  display: flex;
  align-items: center;
}
</style>