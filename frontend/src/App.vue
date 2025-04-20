<script setup>
import { ref, onMounted, useTemplateRef } from 'vue';
import { Ask } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";
import Prompt from "./components/Prompt.vue";
import Message from "./components/Message.vue";

const history = ref([]);
const messageHistory = useTemplateRef('messageHistory');
const waitingResponse = ref(false);
//
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
});


</script>

<template>
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
</template>

<style>
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
