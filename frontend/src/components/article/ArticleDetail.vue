<script setup lang="ts">
import { PhNewspaper } from '@phosphor-icons/vue';
import { useArticleDetail } from '@/composables/article/useArticleDetail';
import ArticleToolbar from './ArticleToolbar.vue';
import ArticleContent from './ArticleContent.vue';
import ImageViewer from '../common/ImageViewer.vue';
import ContextMenu from '../common/ContextMenu.vue';
import type { ContextMenuItem } from '../common/ContextMenu.vue';

const {
  article,
  showContent,
  articleContent,
  isLoadingContent,
  imageViewerSrc,
  imageViewerAlt,
  imageContextMenu,
  close,
  toggleRead,
  toggleFavorite,
  openOriginal,
  toggleContentView,
  closeImageViewer,
  closeImageContextMenu,
  downloadImage,
  t,
} = useArticleDetail();

// Context menu items for image
const imageContextMenuItems: ContextMenuItem[] = [
  {
    label: t('viewImage'),
    action: 'view',
    icon: 'PhMagnifyingGlassPlus',
  },
  {
    label: t('downloadImage'),
    action: 'download',
    icon: 'PhDownloadSimple',
  },
];

function handleImageContextMenuAction(action: string) {
  if (action === 'view') {
    imageViewerSrc.value = imageContextMenu.value.src;
    imageViewerAlt.value = '';
  } else if (action === 'download') {
    downloadImage(imageContextMenu.value.src);
  }
  closeImageContextMenu();
}
</script>

<template>
  <main
    :class="[
      'flex-1 bg-bg-primary flex flex-col h-full absolute w-full md:static md:w-auto z-30 transition-transform duration-300',
      article ? 'translate-x-0' : 'translate-x-full md:translate-x-0',
    ]"
  >
    <div
      v-if="!article"
      class="hidden md:flex flex-col items-center justify-center h-full text-text-secondary text-center px-4"
    >
      <PhNewspaper :size="48" class="mb-4 sm:mb-5 opacity-50 sm:w-16 sm:h-16" />
      <p class="text-sm sm:text-base">{{ t('selectArticle') }}</p>
    </div>

    <div v-else class="flex flex-col h-full bg-bg-primary">
      <ArticleToolbar
        :article="article"
        :show-content="showContent"
        @close="close"
        @toggle-content-view="toggleContentView"
        @toggle-read="toggleRead"
        @toggle-favorite="toggleFavorite"
        @open-original="openOriginal"
      />

      <!-- Original webpage view -->
      <div v-if="!showContent" class="flex-1 bg-white w-full">
        <iframe
          :src="article.url"
          class="w-full h-full border-none"
          sandbox="allow-scripts allow-same-origin allow-popups"
        ></iframe>
      </div>

      <!-- RSS content view -->
      <ArticleContent
        v-else
        :article="article"
        :article-content="articleContent"
        :is-loading-content="isLoadingContent"
      />
    </div>

    <!-- Image Viewer Modal -->
    <ImageViewer
      v-if="imageViewerSrc"
      :src="imageViewerSrc"
      :alt="imageViewerAlt"
      @close="closeImageViewer"
    />

    <!-- Image Context Menu -->
    <ContextMenu
      v-if="imageContextMenu.show"
      :items="imageContextMenuItems"
      :x="imageContextMenu.x"
      :y="imageContextMenu.y"
      @close="closeImageContextMenu"
      @action="handleImageContextMenuAction"
    />
  </main>
</template>
