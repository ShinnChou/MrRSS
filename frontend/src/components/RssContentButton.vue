<script lang="ts" setup>
import { Icon } from '@iconify/vue';

const props = defineProps({
    selectedFeed: {
      type: Object,
      default: null
    }
})
</script>

<template>
    <div class="ContentFunction">
      <button class="btn" title="Back to top" :disabled="!selectedFeed" :style="!selectedFeed ? { color: 'gray' } : {}">
        <Icon icon="material-symbols:vertical-align-top" />
      </button>
      <button class="btn" title="Bookmark" :disabled="!selectedFeed" :style="!selectedFeed ? { color: 'gray' } : {}">
        <Icon icon="material-symbols:collections-bookmark" />
      </button>
      <button class="btn" @click="selectedFeed && $emit('modifyFeedContentReaded', selectedFeed, false)" title="Unread" :disabled="!selectedFeed" :style="!selectedFeed ? { color: 'gray' } : {}">
        <Icon icon="material-symbols:thread-unread" />
      </button>
      <button class="btn" title="Open in browser" :disabled="!selectedFeed"
        :style="!selectedFeed ? { color: 'gray' } : {}">
        <Icon icon="material-symbols:open-in-browser" />
      </button>
      <button class="btn" title="Share" :disabled="!selectedFeed" :style="!selectedFeed ? { color: 'gray' } : {}">
        <Icon icon="material-symbols:share" />
      </button>
      <button class="btn" title="Translate" :disabled="!selectedFeed" :style="!selectedFeed ? { color: 'gray' } : {}">
        <Icon icon="material-symbols:g-translate" />
      </button>
      <button class="btn" title="Chat with AI" :disabled="!selectedFeed"
        :style="!selectedFeed ? { color: 'gray' } : {}">
        <Icon icon="material-symbols:robot-2" />
      </button>
    </div>
</template>

<style scoped>
.ContentFunction {
  position: sticky;
  top: 0;
  width: 100%;
  display: flex;
  justify-content: end;
  align-items: center;
  background-color: #e0e0e0;
}

.ContentFunction .btn {
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

.ContentFunction .btn:hover {
  background-color: #d0d0d0;
}
</style>