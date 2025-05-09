<script setup>
import { ref, onMounted, useTemplateRef } from 'vue';
import { Ask, GetSelectedModel } from "../wailsjs/go/main/App";
import { EventsOn, OnFileDrop } from "../wailsjs/runtime/runtime";
import Prompt from "./components/Prompt.vue";
import Message from "./components/Message.vue";
import Files from "./components/Files.vue";
import _ from "./i18n.js"


const messageHistory = useTemplateRef('messageHistory');
const history = ref([]);
const waitingResponse = ref(false);
const currentModel = ref({ name: "" });
const showHelp = ref(false);
const toastMessage = ref({
  hidden: true,
  type: "",
  title: "",
  message: "",
})

const translations = ref({
  currentModelLabel: "",
  helpText: "",
  closeLabel: "",
});

// Update translations
async function updateTranslation() {
  translations.value = {
    currentModelLabel: await _("model.current"),
    helpText: await _("about.help", true),
    closeLabel: await _("close"),
  }
}

function showToast(type, title, message) {
  toastMessage.value = {
    hidden: false,
    type: type,
    title: title,
    message: message,
  };
  setTimeout(() => {
    toastMessage.value.hidden = true;
  }, 5000);
}

// when content changes
function onContent() {
  messageHistory.value.scrollTop = messageHistory.value.scrollHeight;
}

function upsertMessage(origChunk) {
  const chunk = origChunk.chunk;
  const html = origChunk.html;
  const thinking = origChunk.thinkingHtml;
  const message = history.value.find((msg) => msg.id === chunk.id);
  if (!message) {
    const newMessage = {
      id: chunk.id,
      role: chunk.role,
      content: html,
      thinking: thinking,
    };
    history.value.push(newMessage);
  } else {
    message.thinking = thinking;
    message.content = html
  }
  return message;
}

// send the prompt to the App
function sendPrompt(prompt) {
  const message = {
    id: Date.now(),
    role: "user",
    content: prompt,
    thinking: "",
  };
  history.value.push(message);
  waitingResponse.value = true;
  onContent();
  Ask(prompt)
    .then(() => {
      waitingResponse.value = false;
    })
    .catch((error) => {
      showToast("error", "Error", error);
    });
}

function setCurrentModel() {
  GetSelectedModel()
    .then((model) => {
      currentModel.value = model;
    })
    .catch((error) => {
      showToast("error", "", "Error fetching current model: " + error);
    });
}

onMounted(() => {
  // Listen for events from the backend
  EventsOn("chunk", (chunk) => {
    upsertMessage(chunk);
    onContent();
  });
  EventsOn("new-conversation", () => {
    history.value = [];
    onContent();
  })
  EventsOn("selected-model", (model) => {
    currentModel.value = model;
  });
  EventsOn("show-help", () => {
    showHelp.value = true;
  });
  OnFileDrop((x, y, paths) => {
    console.log("File dropped at", x, y, paths);
  })
  updateTranslation();
  setCurrentModel();
  // hide the help popup when clicking outside of it
  document.addEventListener('keyup', (event) => {
    if (event.key === "Escape" && showHelp.value) {
      showHelp.value = false;
    }
  });
});


</script>

<template>
  <div class="on-top">
    <p>{{ translations.currentModelLabel }} :
      <strong>{{ currentModel.name }}</strong>
      <span v-if="currentModel.uncensored"> 🔞</span>
      <small> :: {{ currentModel.description }}</small>
      <span v-if="currentModel.vision"> 👁️</span>
    </p>
  </div>
  <div class="message-history" ref="messageHistory">
    <Message v-for="message in history" :key="message.id" :message="message" :onContent="onContent"
      :model="currentModel" />
    <div v-if="waitingResponse" class="thinking">
      <span>🧠</span>
      <span>🧠</span>
      <span>🧠</span>
    </div>
  </div>

  <Files :class="{ 'hidden': !currentModel.vision }" />
  <Prompt :sendPrompt="sendPrompt" :model="currentModel" />
  <div class="popup" v-if="showHelp" tabindex="-1">
    <article v-html="translations.helpText">
    </article>
    <button @click="showHelp = false">{{ translations.closeLabel }}</button>
  </div>
  <div :class="['toast', toastMessage.type]" v-if="!toastMessage.hidden">
    <strong>{{ toastMessage.title }}</strong>
    <p>{{ toastMessage.message }}</p>
  </div>
</template>

<style>
.hidden {
  visibility: hidden;
}

.popup {
  position: fixed;
  display: flex;
  flex-direction: column;
  background-color: var(--view-bg-color);
  color: var(--view-fg-color);
  padding: 10px;
  border-radius: .5rem;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  height: 80vh;
  width: 80vw;
  margin-top: 1rem;
  margin-left: 10vw;
  z-index: 1001;
}

.popup h1 {
  position: relative;
  font-size: 2rem;
  margin: 0;
  padding: 0;
  height: calc(255px + 2rem);
}

.popup h1:after {
  position: absolute;
  top: 2.2rem;
  left: 0;
  background-image: url(./assets/images/appicon.png);
  background-size: contain;
  width: 100%;
  height: 255px;
  background-repeat: no-repeat;
  background-position: calc(50%);
  content: " ";
}

.popup article {
  overflow-y: auto;
  height: 90%;
  margin: 1rem;
}

.popup button {
  padding: 10px;
  background-color: var(--success-bg-color);
  color: var(--success-fg-color);
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.popup blockquote {
  font-style: italic;
  background-color: var(--slate-bg-color);
  color: var(--slate-fg-color);
  padding: 1rem;
  position: relative;
  border-radius: 1rem;
  box-shadow: 0 0 24px rgba(0, 0, 0, 0.3);
}

.popup blockquote::before {
  content: "«";
  font-size: 2rem;
  color: var(--slate-fg-color);
  position: absolute;
  top: 0;
  left: -1rem;
}

.popup blockquote::after {
  content: "»";
  font-size: 2rem;
  color: var(--slate-fg-color);
  position: absolute;
  bottom: 0;
  right: -1rem;
}

small {
  font-size: 0.8rem;
}


.on-top {
  z-index: 1000;
  border-radius: .5rem;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  text-align: center;
  background-color: var(--view-bg-color);
  color: var(--view-fg-color);
}

.on-top p {
  margin: .25em;
}

.message-history {
  flex-grow: 1;
  padding: 20px;
  overflow-y: auto;
}


.thinking {
  display: flex;
  justify-content: center;
  gap: 5px;
}

.thinking>span {
  opacity: 0;
  animation: thinkingAnimation 3s infinite;
}

.thinking>span:nth-child(1) {
  animation-delay: 0s;
}

.thinking>span:nth-child(2) {
  animation-delay: 0.5s;
}

.thinking>span:nth-child(3) {
  animation-delay: 1s;
}


@keyframes thinkingAnimation {
  0% {
    opacity: 0;
  }

  20% {
    opacity: 1;
  }

  40% {
    opacity: 0;
  }
}

.toast {
  position: fixed;
  bottom: 55px;
  left: 80%;
  transform: translateX(-50%);
  background-color: var(--window-bg-color);
  color: var(--view-fg-color);
  padding: 10px;
  border-radius: .5rem;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  z-index: 1001;
  min-width: 200px;
}

.toast.error {
  background-color: var(--error-bg-color);
  color: var(--error-fg-color);
}
</style>
