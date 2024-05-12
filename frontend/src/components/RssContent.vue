<script lang="ts" setup>
const props = defineProps({
  selectedFeed: Object
})
</script>

<template>
  <article v-if="selectedFeed">
    <h1>{{ selectedFeed.Title }}</h1>
    <div class="RssInfo">
      <div class="RssTitle">
        <img v-if="selectedFeed.FeedImage" :src="selectedFeed.FeedImage"
          :alt="`Image of feed: ${selectedFeed.FeedTitle}`" class="RssImage" />
        <p>{{ selectedFeed.FeedTitle }}</p>
      </div>
      <time>{{ selectedFeed.Time }}</time>
    </div>
    <div v-html="selectedFeed.Content" class="Content"></div>
  </article>
  <article v-else>
  </article>
</template>

<style scoped>
article {
  padding: 32px;
  background-color: #f9f9f9;

  overflow-x: hidden;
  overflow-y: scroll;
}

article::-webkit-scrollbar {
  width: 8px;
}
article::-webkit-scrollbar-thumb {
  background-color: #ccc;
  cursor: pointer;
}
article::-webkit-scrollbar-track {
  background-color: #f1f1f1;
}
article::-webkit-scrollbar-thumb:hover {
  background-color: #999;
}

h1 {
  font-size: 32px;
  margin-bottom: 16px;
  max-width: 720px;
  margin: auto;
}

.RssInfo {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 8px;
  padding-top: 32px;
  padding-bottom: 64px;
  max-width: 720px;
  margin: auto;
}

.RssTitle {
  display: flex;
  align-items: center;
}

.RssImage {
  min-width: 16px;
  max-width: 16px;
  min-height: 16px;
  max-height: 16px;
  margin-right: 4px;
}

time {
  color: #666;
}

p {
  margin: 0;
}

.Content {
  max-width: 720px;
  margin: auto;
}

.Content :deep(h1, h2, h3, h4, h5, h6) {
  margin-top: 32px;
  margin-bottom: 16px;
}

.Content :deep(p) {
  text-align: left;
  margin-top: 16px;
  margin-bottom: 16px;
  line-height: 24px;
}

.Content :deep(ul, ol) {
  text-align: left;
}

.Content :deep(li) {
  text-align: left;
  margin-top: 8px;
  margin-bottom: 8px;
}

.Content :deep(a) {
  text-decoration: underline;
  color: #666;
}

.Content :deep(a:hover) {
  color: #333;
}

.Content :deep(blockquote) {
  margin-left: 8px;
  margin-right: 8px;
  padding: 8px;
  border-left: 2px solid #666;
  color: #666;
  background-color: #f0f0f0;
}

.Content :deep(img) {
  max-width: 100%;
  height: auto;
}

.Content :deep(img:only-child) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: auto;
}

.Content :deep(pre) {
  text-align: left;
  margin-top: 16px;
  margin-bottom: 16px;
  padding: 8px;
  border: 1px solid #666;
  background-color: #f0f0f0;
  white-space: pre-wrap;
  font-family: 'Fira Code', monospace;
}

.Content :deep(code) {
  text-align: left;
  background-color: #f0f0f0;
  padding: 2px;
  border-radius: 4px;
  font-family: 'Fira Code', monospace;
}

.Content :deep(table) {
  width: 100%;
  border-collapse: collapse;
}

.Content :deep(th, td) {
  border: 1px solid #666;
  padding: 8px;
}

.Content :deep(th) {
  background-color: #f0f0f0;
}

.Content :deep(tr:nth-child(even)) {
  background-color: #f9f9f9;
}

.Content :deep(tr:nth-child(odd)) {
  background-color: #f0f0f0;
}

.Content :deep(iframe) {
  width: 100%;
  height: 400px;
}

.Content :deep(hr) {
  margin-top: 32px;
  margin-bottom: 32px;
  border: none;
  border-top: 1px solid #666;
}

.Content :deep(figure) {
  margin-top: 16px;
  margin-bottom: 16px;
}

.Content :deep(figcaption) {
  text-align: center;
  font-size: small;
  color: #666;
}

.Content :deep(button) {
  display: none;
}
</style>