<script setup>
import { watch, nextTick, useTemplateRef, computed, onMounted, ref } from 'vue';
import _ from "../i18n.js"
import 'mathjax/es5/tex-mml-svg.js';
import hljs from 'highlight.js';
import 'highlight.js/styles/github-dark.css';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime.js';


const props = defineProps(['message', "onContent", "model"]);
const container = useTemplateRef('message');

const cssClasses = computed(() => {
  return {
    'message-content': true,
    'user': props.message.role === 'user',
    'assistant': props.message.role === 'assistant',
  };
});

const translations = ref({
  thinkingLabel: "",
});

async function updateTranslation() {
  translations.value.thinkingLabel = await _("thinking.label")
}

async function enhanceHighlight() {
  return new Promise((resolve) => {
    const pre = container.value?.querySelectorAll('pre') || [];
    pre.forEach((block) => {
      if (block.querySelectorAll('code').length === 0) {
        return;
      }
      hljs.highlightElement(block);
    });
    resolve()
  })
}

// reformat to use highlightjs, image loading, links, and so on
async function formatMessage() {
  await nextTick(); // wait for DOM updates

  MathJax.typesetPromise()
    .then(fixLinks)
    .then(enhanceHighlight)
    .then(enhanceImages)
    .then(props.onContent);

}

// make linls open in a new browser
function fixLinks() {
  const links = container.value?.querySelectorAll('a') || [];
  links.forEach(link => {
    const href = link.getAttribute('href');
    // make all links open in a new browser
    if (href.length > 0) {
      link.addEventListener('click', (e) => {
        e.preventDefault();
        BrowserOpenURL(href)
      });
    }
  });
}

// add a spinner while loading images
const enhanceImages = async () => {
  const images = container.value?.querySelectorAll('img') || [];
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

watch(() => props.message.content, formatMessage, { immediate: true, });
onMounted(() => {
  MathJax.svgStylesheet();
  updateTranslation();
});
</script>

<template>
  <div class="message-container">
    <div class="reasoning" v-if="props.model?.reasoning && props.message.role == 'assistant'">
      <details>
        <summary>{{ translations.thinkingLabel }}</summary>
        <div v-html="props.message.thinking"></div>
      </details>
    </div>
    <div :class="cssClasses">
      <div ref="message" v-html="props.message.content"></div>
    </div>
  </div>
</template>

<style>
.hljs {
  padding: 1rem;
}

.message-container {
  display: flex;
  flex-direction: column;
  margin: 1rem 0;
}

.message-content {
  padding: 1rem;
  border-radius: 1rem;
  word-wrap: break-word;
  color: var(--adw-color-fg);
  box-shadow: 0 0 24px rgba(0, 0, 0, 0.3);
  background-color: var(--view-bg-color);
  position: rlative;
}

.message-content.assistant {
  background-color: var(--slate-bg-color);
  color: var(--slate-fg-color);
}

.reasoning {
  margin-left: auto;
  width: 100%;
  margin-top: 1rem;
  margin-bottom: 1rem;
}

.reasoning summary {
  text-align: center;
  font-weight: bold;
  cursor: pointer;
  padding: 0.5rem;
}

.reasoning details div {
  border-radius: 1rem;
  background-color: color-mix(in srgb, var(--slate-bg-color) 25%, transparent);
  padding: 1rem;
}

@media (min-width: 974px) {

  .message-content,
  .reasoning {
    max-width: 80%;
  }

  .message-content.user {
    margin-right: auto;
  }

  .message-content.assistant,
  .reasoning {
    margin-left: auto;
  }
}

@media (min-width: 1280px) {

  .message-content,
  .reasoning {
    max-width: 90%;
  }

  .message-container {
    margin: 1rem 15%;
  }
}

.message-content p {
  display: block;
  position: relative;
}

.message-content p img {
  max-width: 80%;
  height: auto;
  align-self: center;
  justify-self: center;
  opacity: 0;
  transition: opacity 0.5s ease-in-out;
  visibility: hidden;
  --default-contextmenu: show;
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
  margin: auto;
  width: 24px;
  height: 24px;
  border: 3px solid #ccc;
  border-top-color: #333;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  z-index: 1;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}

.image-wrapper {
  position: relative;
  text-align: center;
  width: 100%;
}

@keyframes spin {
  to {
    transform: translate(-50%, -50%) rotate(360deg);
  }
}
</style>
