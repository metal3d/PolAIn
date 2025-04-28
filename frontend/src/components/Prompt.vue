<script setup>
import { useTemplateRef, onMounted, ref } from "vue";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import _ from "../i18n.js"

const answering = ref(false)
const props = defineProps(['sendPrompt']);
const userInput = useTemplateRef('userInput');

const translations = ref({
  placeholder: "",
  promptSend: "",
})

async function updateTranslation() {
  translations.value = {
    placeholder: await _("prompt.placeholder"),
    promptSend: await _("prompt.send"),
  }
}

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

onMounted(() => {
  EventsOn("ask-start", () => {
    answering.value = true;
  });
  EventsOn("ask-done", () => {
    answering.value = false;
    userInput.value.focus();
  });
  updateTranslation()
});

</script>
<template>
  <div class="prompt-container">
    <div class="prompt-wrapper">
      <div class="upload-buttons">
        <button class="upload image" title="Upload image">🖼️</button>
        <button class="upload audio" title="Upload audio">🎧</button>
      </div>
      <textarea :placeholder="translations.placeholder" ref="userInput" :disabled="answering"
        @keyup="handleTextareaKeys" />
    </div>
    <button @click="sendPrompt()" ref="sendButton" :disabled="answering">{{ translations.promptSend }}</button>
  </div>
</template>
<style>
.prompt-wrapper {
  position: relative;
  flex-grow: 1;
  display: flex;
  align-items: center;
}

.prompt-container {
  display: flex;
  padding: 10px;
}

.prompt-container textarea {
  flex-grow: 1;
  padding: 10px;
  border: 1px solid color-mix(in srgb, var(--view-fg-color), #000 15%);
  border-radius: 5px;
  margin-right: 5px;
}

.prompt-container button {
  padding: 10px 20px;
  background-color: var(--success-bg-color);
  color: var(--success-fg-color);
  border: none;
  border-radius: 5px;
  cursor: pointer;
  transition: .2s opacity ease-in-out;
}


button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

button:hover:not(:disabled) {
  background-color: var(--success-bg-color);
  opacity: .8;
}

.upload-buttons {
  position: absolute;
  right: 0;
  display: flex;
}

button.upload {
  font-size: 1rem;
  background-color: transparent;
}
</style>
