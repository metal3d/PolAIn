<script setup>
import { watch, nextTick, useTemplateRef, computed } from 'vue';
import MarkdownIt from 'markdown-it';
import MarkdownItHighlight from 'markdown-it-highlightjs';
import 'highlight.js/styles/github-dark.css';

const props = defineProps(['message', "onContent"]);
const markdown = new MarkdownIt()
  .use(MarkdownItHighlight);
const container = useTemplateRef('message');
const cssClasses = computed(() => {
  return {
    'message-content': true,
    'user': props.message.role === 'user',
    'assistant': props.message.role === 'assistant',
  };
});


const enhanceImages = async () => {
  await nextTick(); // wait for DOM updates 
  const images = container.value?.querySelectorAll('img') || [];
  console.log('Images found:', images.length);
  images.forEach(img => {
    // Skip if already wrapped
    if (img.parentElement?.classList.contains('image-wrapper')) return;

    const wrapper = document.createElement('div');
    wrapper.className = 'image-wrapper';

    const spinner = document.createElement('div');
    spinner.className = 'spinner';

    img.parentNode.insertBefore(wrapper, img);
    wrapper.appendChild(img);
    wrapper.appendChild(spinner);

    img.style.opacity = 0;
    img.style.visibility = 'hidden';

    const onLoaded = () => {
      spinner.remove();
      img.style.opacity = 1;
      img.style.visibility = 'visible';
      props.onContent();
    };

    if (img.complete && img.naturalWidth !== 0) {
      onLoaded();
    } else {
      img.addEventListener('load', onLoaded, { once: true });
      img.addEventListener('error', onLoaded, { once: true });
    }
  });
};

watch(() => props.message.content, enhanceImages, { immediate: true });
</script>

<template>
  <div class="message-container">
    <div :class="cssClasses">
      <div ref="message" v-html="markdown.render(props.message.content)"></div>
    </div>
  </div>
</template>

<style>
.message-container {
  display: flex;
  flex-direction: column;
  margin: 1rem 0;
}

.message-content {
  padding: 1rem;
  border-radius: 1rem;
  word-wrap: break-word;
  color: #121212;
  box-shadow: 0 0 24px rgba(0, 0, 0, 0.3);
}

.message-content.assistant {
  background-color: #383838;
  color: #fafafa;
}

@media (prefers-color-scheme: dark) {
  .message-content {
    color: #fafafa;
  }
  .message-content.assistant {
    background-color: #303030;
  }

  .message-content.user {
    background-color: #242424;
  }
}

@media (min-width: 974px) {
  .message-content {
    max-width: 80%;
  }

  .message-content.user {
    margin-right: auto;
  }

  .message-content.assistant {
    margin-left: auto;
  }
}

@media (min-width: 1280px) {
  .message-content {
    max-width: 90%;
  }

  .message-container {
    margin: 1rem 15%;
  }
}

.message-content p {
  display: flex;
  flex-direction: column;
}

.message-content p img {
  max-width: 80%;
  height: auto;
  align-self: center;
  justify-self: center;
  opacity: 0;
  transition: opacity 0.5s ease-in-out;
  visibility: hidden;
}


.image-wrapper {
  position: relative;
  display: inline-block;
  text-align: center;
}

.image-wrapper img {
  box-shadow: 0 0 24px rgba(0, 0, 0, 0.3);
}

.spinner {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 24px;
  height: 24px;
  border: 3px solid #ccc;
  border-top-color: #333;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  z-index: 1;
}

@keyframes spin {
  to {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}
</style>
