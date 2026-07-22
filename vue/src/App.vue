<script setup lang="ts">
import { onMounted, ref } from 'vue';
const load = ref(false)
const baseurl = "https://api.3mua.cn/api/bmoji/v1/emote/"
type listobj = {
  "id": number,
  "text": string,
  "type": number,
  "url": string,
  "UrlType": number
}
const emoteList = ref<listobj[]>([])

function GetList() {
  fetch(baseurl + 'list').then((d) => d.json()).then((d) => {
    emoteList.value.push(...d.data.all_packages)
    load.value = true
  })
}
onMounted(() => {
  GetList()
})
</script>

<template>
  <RouterView v-slot="{ Component }">
    <component :is="Component" :list="emoteList"></component>
  </RouterView>
</template>

<style>
body {
  margin: 0;
  background-color: #F1F2F3;
}

a {
  color: #81ff81;
  text-decoration: none;
}
h1,
p {
  color: #18191C;
}
.bts {
  display: flex;
  gap: 2rem;

  span {
    background-color: #00AEEC;
    padding: 0.5rem;
    box-sizing: border-box;
    border-radius: 4px;
    transition: 0.3s;
    cursor: pointer;
    a{
      text-decoration: none;
      color: #000
    }
  }
  span:hover{
    transition: 0.3s;
    box-shadow: #353535 0 0 8px;
    background-color: #00bbff;
  }

}
</style>
