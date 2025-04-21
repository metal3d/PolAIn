<script setup>
import { ref, onMounted, useTemplateRef } from 'vue';
import { Ask, GetSelectedModel } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";
import MarkdownIt from 'markdown-it';
import Prompt from "./components/Prompt.vue";
import Message from "./components/Message.vue";
import _ from "./i18n.js"

const history = ref([]);
const messageHistory = useTemplateRef('messageHistory');
const waitingResponse = ref(false);
const currentModel = ref("");
const showHelp = ref(false);
const markdown = new MarkdownIt();

const currentModelLabel = ref("")
const helpText = ref("")
const closeLabel = ref("")

// Update translations
async function updateTranslation() {
  currentModelLabel.value = await _("model.current")
  helpText.value = await _("about.help")
  closeLabel.value = await _("close")
}

// when content changes
function onContent() {
  messageHistory.value.scrollTop = messageHistory.value.scrollHeight;
}

function upsertMessage(chunk) {
  const message = history.value.find((msg) => msg.id === chunk.id);
  if (!message) {
    const newMessage = {
      id: chunk.id,
      role: chunk.role,
      content: chunk.choices[0].delta.content,
    };
    history.value.push(newMessage);
    return newMessage;
  }
  message.content += chunk.choices[0].delta.content;
  return message;
}

// send the prompt to the App
function sendPrompt(prompt) {
  const message = {
    id: Date.now(),
    role: "user",
    content: prompt,
  };
  history.value.push(message);
  waitingResponse.value = true;
  onContent();
  Ask(prompt)
    .then((response) => {
      waitingResponse.value = false;
      console.log("Response received.", response);
      message.content = message.content
    });
}

function setCurrentModel() {
  GetSelectedModel()
    .then((model) => {
      currentModel.value = model;
    })
    .catch((error) => {
      console.error("Error fetching current model:", error);
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
    showHelp.value = !showHelp.value;
  });
  updateTranslation();
  setCurrentModel();

});

// Open all links externally
document.body.addEventListener('click', function (e) {
  if (e.target && e.target.nodeName == 'A' && e.target.href) {
    const url = e.target.href;
    if (
      !url.startsWith('http://#') &&
      !url.startsWith('file://') &&
      !url.startsWith('http://wails.localhost:')
    ) {
      e.preventDefault();
      window.runtime.BrowserOpenURL(url);
    }
  }
});

</script>

<template>
  <div class="on-top">
    <p>{{ currentModelLabel }} :
      <strong>{{ currentModel.name }}</strong>
      <span v-if="currentModel.uncensored"> 🔞</span>
    </p>
    <small>{{ currentModel.description }}</small>
  </div>
  <div class="message-history" ref="messageHistory">
    <Message v-for="message in history" :key="message.id" :message="message" :onContent="onContent" />
    <div v-if="waitingResponse" class="thinking">
      <span>🧠</span>
      <span>🧠</span>
      <span>🧠</span>
    </div>
    <!--div v-if="waitingResponse" class="thinking">Loading</div-->
  </div>

  <Prompt :sendPrompt="sendPrompt" />
  <div class="popup" v-if="showHelp">
    <article v-html="markdown.render(helpText)">
    </article>
    <button @click="showHelp = false">{{ closeLabel }}</button>
  </div>
</template>

<style>
.popup {
  position: fixed;
  display: flex;
  flex-direction: column;
  background-color: rgba(255, 255, 255, 0.8);
  padding: 10px;
  border-radius: .5rem;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  height: 80vh;
  width: 80vw;
  margin-top: 1rem;
  margin-left: 10vw;
  z-index: 1001;
}

.popup article {
  overflow-y: auto;
  height: 90%;
  margin: 1rem;
}

.popup button {
  padding: 10px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.popup blockquote {
  font-style: italic;
  color: #555;
}

small {
  font-size: 0.8rem;
}


@media (prefers-color-scheme: dark) {
  .popup {
    background-color: #242424;
    color: white;
  }

  .popup blockquote {
    color: #aaa;
  }

  .popup button {
    background-color: #4CAF50;
    color: white;
  }
}


.on-top {
  position: absolute;
  top: 0;
  right: 0;
  padding: 10px;
  z-index: 1000;
  padding: 1rem;
  border-radius: .5rem;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  text-align: center;
  background-color: rgba(255, 255, 255, 0.1);
  color: black;
}

.on-top p {
  margin: .25em;
}

@media (prefers-color-scheme: dark) {
  .on-top {
    background-color: rgba(0, 0, 0, 0.8);
    color: white;
  }
}

.message-history {
  flex-grow: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: #f0f0f0;
}


@media (prefers-color-scheme: dark) {
  .message-history {
    background-color: #121212;
    color: white;
  }
}

.thinking {
  display: flex;
  justify-content: center;
  gap: 5px;
  /* Espacement entre les span */
}

.thinking>span {
  opacity: 0;
  /* Initialement caché */
  animation: thinkingAnimation 3s infinite;
  /* Animation */
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
</style>
