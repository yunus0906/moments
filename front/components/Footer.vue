<template>
  <div class="footer p-4 flex justify-center items-center text-xs ">
    <div class="flex flex-col gap-1 items-center">
      <div class="text-xs text-center text-gray-500 py-2">本站不以盈利为目的，所有内容收集整理于网络，本站不参与制作。</div>
      <a class="my-2 text-gray-500" v-if="sysConfig.beiAnNo" href="https://beian.miit.gov.cn/" target="_blank">{{ sysConfig.beiAnNo }}</a>
    </div>
  </div>
</template>

<script setup lang="ts">
import type {SysConfigVO} from "~/types";

const sysConfig = useState<SysConfigVO>('sysConfig')

</script>

<style scoped>

</style>