<script setup>
import { useTemplateRef, onMounted, ref } from "vue";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import _ from "../i18n.js"

const answering = ref(false)
const props = defineProps(['sendPrompt']);
const userInput = useTemplateRef('userInput');

const placeholder = ref('');
const promptSend = ref('')

async function updateTranslation() {
  placeholder.value = await _("prompt.placeholder")
  promptSend.value = await _("prompt.send")
}

onMounted(() => {
  EventsOn("ask-start", () => {
    answering.value = true;
  });
  EventsOn("ask-done", () => {
    answering.value = false;
  });
  updateTranslation()
});



function sendPrompt() {
  const prompt = userInput.value.value;
  if (prompt) {
    props.sendPrompt(prompt);
    userInput.value.value = ""; // Clear the input after sending
    userInput.value.focus();
  }
}

function handleTextareaKeys(event) {
  if (event.key === "Enter" && !(event.metaKey || event.ctrlKey || event.shiftKey)) {
    event.preventDefault(); // Prevent default behavior of Enter key
    sendPrompt();
  }
}

</script>
<template>
  <div class="prompt-container">
    <textarea :placeholder="placeholder" ref="userInput" :disabled="answering" @keyup="handleTextareaKeys" />
    <button @click="sendPrompt()" ref="sendButton" :disabled="answering">{{ promptSend }}</button>
  </div>
</template>
<style>
.prompt-container {
  display: flex;
  padding: 10px;
  background-color: #eee;
}

.prompt-container textarea {
  flex-grow: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
  margin-right: 5px;
}

.prompt-container button {
  padding: 10px 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

button:hover {
  background-color: #45a049;
}

button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

@media (prefers-color-scheme: dark) {
  .prompt-container {
    background-color: #121212;
    color: white;
  }

  .prompt-container textarea {
    background-color: #121212;
    color: white;
  }
}
</style>
