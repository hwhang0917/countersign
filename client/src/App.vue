<script setup lang="ts">
import { ref } from "vue";
const x = ref(0);
</script>

<template>
    <div class="flex w-screen justify-center text-red-500">
        <p>{{ x }}</p>
    </div>
</template>
