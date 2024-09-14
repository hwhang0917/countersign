<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useWebSocket, useIntervalFn } from '@vueuse/core'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle
} from '@/components/ui/card'
import { HoverCard, HoverCardContent, HoverCardTrigger } from '@/components/ui/hover-card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { Skeleton } from '@/components/ui/skeleton'

const { data, send } = useWebSocket('ws://localhost:8080/ws/otp')

const textInput = ref()
const askText = ref()
const isDialogOpen = ref(false)
useIntervalFn(() => {
  if (!askText.value) return
  send(
    JSON.stringify({
      ask_text: askText.value
    })
  )
}, 1000)
const resetAskText = () => {
  askText.value = ''
  textInput.value = ''
}
const setAskText = () => {
  askText.value = textInput.value
  isDialogOpen.value = true
}
const parsedData = computed(() => JSON.parse(data.value))
const progress = computed(() => {
  if (parsedData.value && parsedData.value.success) {
    return (parsedData.value.time_left / 30) * 100
  }
  return 0
})
const isLoading = computed(() => !parsedData.value)

watch(
  () => isDialogOpen.value,
  (value) => {
    if (!value) {
      resetAskText()
    }
  }
)
</script>

<template>
  <main class="flex h-screen w-screen items-center justify-center">
    <Card class="w-2/3 p-4">
      <CardHeader>
        <CardTitle>Countersign</CardTitle>
      </CardHeader>
      <CardContent>
        <CardDescription>
          <div class="mb-2 flex gap-1">
            <p>인증받을 문어를 입력해주세요.</p>
            <HoverCard>
              <HoverCardTrigger class="cursor-pointer underline">(?)</HoverCardTrigger>
              <HoverCardContent class="w-120">
                <p class="font-bold">문어는 구두로 말하기 쉬운 단어일 수록 좋습니다.</p>
                <small>예) 바나나, 바람, 종달새</small>
              </HoverCardContent>
            </HoverCard>
          </div>
          <Input v-model="textInput" @keyup.enter="setAskText" />
        </CardDescription>
      </CardContent>
      <CardFooter class="flex justify-center px-6 pb-6">
        <Button class="w-full" @click="setAskText" :disabled="!textInput">인증 단어 조회</Button>
      </CardFooter>
    </Card>
    <Dialog v-model:open="isDialogOpen">
      <DialogContent>
        <DialogHeader v-if="isLoading">
          <DialogTitle class="mb-4">
            <Skeleton class="h-7 w-40 bg-gray-200" />
          </DialogTitle>
          <DialogDescription>
            <Skeleton class="h-5 w-full bg-gray-200" />
          </DialogDescription>
          <DialogDescription>
            <Skeleton class="h-5 w-full bg-gray-200" />
          </DialogDescription>
        </DialogHeader>
        <DialogHeader v-else>
          <DialogTitle class="mb-4">문어: {{ askText }} / 답어: {{ parsedData.otp }}</DialogTitle>
          <DialogDescription>
            상대방에게 문어를 말해주세요. 그리고 제한시간 내에 답어를 확인해주세요.
          </DialogDescription>
          <Progress :model-value="progress" />
          <DialogDescription> {{ parsedData.time_left }}초 동안 유효합니다. </DialogDescription>
        </DialogHeader>
        <DialogFooter v-if="isLoading">
          <Skeleton class="h-10 w-20 bg-gray-200" />
        </DialogFooter>
        <DialogFooter v-else>
          <Button @click="isDialogOpen = false">초기화</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </main>
</template>
