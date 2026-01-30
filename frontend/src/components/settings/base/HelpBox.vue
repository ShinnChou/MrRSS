<script setup lang="ts">
import { computed } from 'vue';
import type { Component } from 'vue';
import { PhInfo, PhQuestion, PhLightbulb } from '@phosphor-icons/vue';

interface Props {
  type?: 'info' | 'help' | 'tip';
  icon?: Component;
  title?: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'help',
  icon: undefined,
  title: undefined,
});

const defaultIcon = computed(() => {
  switch (props.type) {
    case 'info':
      return PhInfo;
    case 'tip':
      return PhLightbulb;
    default:
      return PhQuestion;
  }
});

const displayIcon = computed(() => props.icon || defaultIcon.value);

const boxClass = computed(() => {
  switch (props.type) {
    case 'info':
      return 'help-box-info';
    case 'tip':
      return 'help-box-tip';
    default:
      return 'help-box-help';
  }
});
</script>

<template>
  <div class="help-box" :class="boxClass">
    <div v-if="title || displayIcon" class="help-box-header">
      <component :is="displayIcon" v-if="displayIcon" :size="18" class="help-box-icon" />
      <span v-if="title" class="help-box-title">{{ title }}</span>
    </div>
    <div class="help-box-content">
      <slot />
    </div>
  </div>
</template>

<style scoped>
@reference "../../../../style.css";

.help-box {
  @apply py-2.5 px-3 rounded-lg border;
}

.help-box-header {
  @apply flex items-center gap-2 mb-2;
}

.help-box-icon {
  @apply shrink-0;
}

.help-box-title {
  @apply text-sm font-medium;
}

.help-box-content {
  @apply text-xs sm:text-sm text-text-secondary;
}

.help-box-content :deep(p) {
  @apply mb-2 last:mb-0;
}

.help-box-content :deep(ol),
.help-box-content :deep(ul) {
  @apply list-inside space-y-1 my-2;
}

.help-box-content :deep(ol) {
  @apply list-decimal;
}

.help-box-content :deep(ul) {
  @apply list-disc;
}

.help-box-content :deep(a) {
  @apply text-accent underline underline-offset-2 hover:no-underline;
}

/* Info type - blue */
.help-box-info {
  background-color: rgba(59, 130, 246, 0.05);
  border-color: rgba(59, 130, 246, 0.3);
}

.help-box-info .help-box-icon,
.help-box-info .help-box-title {
  @apply text-blue-500;
}

/* Help type - purple (default) */
.help-box-help {
  background-color: rgba(139, 92, 246, 0.05);
  border-color: rgba(139, 92, 246, 0.3);
}

.help-box-help .help-box-icon,
.help-box-help .help-box-title {
  @apply text-purple-500;
}

/* Tip type - amber */
.help-box-tip {
  background-color: rgba(245, 158, 11, 0.05);
  border-color: rgba(245, 158, 11, 0.3);
}

.help-box-tip .help-box-icon,
.help-box-tip .help-box-title {
  @apply text-amber-500;
}
</style>
