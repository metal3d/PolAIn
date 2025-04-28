<script setup>
import { onMounted, ref } from 'vue';
import { EventsOn } from '../../wailsjs/runtime/runtime';
import { RemoveFile } from '../../wailsjs/go/main/App';

const files = ref([]);

function drop(file) {
  // find the index of the file
  const index = files.value.findIndex((f) => f === file);
  if (index > -1) {
    // remove the file
    RemoveFile(index).then(() => {
      files.value.splice(index, 1);
    });
  } else {
    console.error('File not found:', file);
  }
}

onMounted(() => {
  EventsOn("register-files", (content) => {
    files.value.push(content);
  });
  EventsOn("new-conversation", () => {
    files.value = [];
  });
})
</script>
<template>

  <div v-if="files.length" class="file-container">
    <span v-for="file in files" :key="file">
      <button @click="drop(file)">❌</button>
      <img :key="file" :src="file" class="file" />
    </span>
  </div>

</template>
<style scoped>
.file-container {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 10px;
}

span {
  position: relative;
}

button {
  position: absolute;
  border: 0;
  right: 0;
}

button:hover {
  background-color: transparent;
  color: var(--adw-color-fg);
}

img {
  max-width: 80px;
}
</style>
