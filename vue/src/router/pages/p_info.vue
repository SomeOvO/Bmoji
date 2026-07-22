<script lang="ts" setup>
  import JSZip from 'jszip'
  import { saveAs } from 'file-saver'
  import { onMounted, ref } from 'vue'
  type infoobj = {
    "id": number,
    "text": string,
    "emote": [
      {
        "text": string,
        "url": string,
        "type": number,
        "gif_url": string,
        "meta": {
          "alias": string
        }
      }]
  }
  const id = window.location.search.split("?")[1]
  const info = ref<infoobj>()
  const baseurl = "https://api.3mua.cn/api/bmoji/v1/emote/"
  onMounted(() => {
    fetch(baseurl + "info?id=" + id).then((d) => d.json()).then((d) => {
      info.value = d
      info.value!.emote.map((e: { [x: string]: string }) => {
        if (e["gif_url"] !== "") {
          hasgif.value = true
        }
      })
    })
  })
  const view = ref({
    "name": "",
    "url": "",
    "ot": "",
    "type": -1,
    "display": false
  })
  function Info(num: number) {
    const a = info.value?.emote[num]
    view.value.name = a?.text ?? ''
    view.value.ot = a?.meta.alias ?? ''
    view.value.type = a?.type ?? -1
    if (!viewgif.value || a?.gif_url == "") {
      view.value.url = a?.url ?? ""
    } else {
      view.value.url = a?.gif_url ?? ""
    }
    view.value.display = true
  }
  function copy(data: string) {
    navigator.clipboard.writeText(data)
    alert("复制成功\n" + data)
  }
  function click(e: MouseEvent, type: number) {
    const doc = e.target as HTMLElement
    if (doc.className === "ei" || doc.className === "download") {
      switch (type) {
        case 1:

          view.value.display = false

          break;
        case 2:
          todownload.value = false
          break;
      }
    }

  }

  async function download(type: number) {
    const zip = new JSZip
    switch (type) {
      case 1:
        await Promise.all(info.value!.emote.map(async (e: { url: string | URL | Request; text: string }) => {
          const resp = await fetch(e.url, {
            method: "GET",
            referrerPolicy: "no-referrer"
          })
          const blob = await resp.blob()
          const type = resp.headers.get("content-type")!.split("/")[1]
          zip.file(e.text + "." + type, blob)
        }))
        const content = await zip.generateAsync({ type: "blob" })
        saveAs(content, info.value?.text)
        break;
      case 2:
        await Promise.all(info.value!.emote.map(async (e: { gif_url: string | URL | Request; url: string | URL | Request; text: string }) => {
          const resp = await fetch((() => {
            if (e.gif_url !== "") {
              return e.gif_url
            } else {
              return e.url
            }
          })(), {
            method: "GET",
            referrerPolicy: "no-referrer"
          })
          const blob = await resp.blob()
          const type = resp.headers.get("content-type")!.split("/")[1]
          zip.file(e.text + "." + type, blob)
        }))
        const content2 = await zip.generateAsync({ type: "blob" })
        saveAs(content2, info.value?.text)
        break;
      case 3:
        const js = info.value
        const blob = new Blob([JSON.stringify(js)], { type: 'application/json' });
        saveAs(blob, js?.text + ".json")
        break;
      case 4:
        const json = {
          "name":"",
          "type":"",
          "prefix":"bmoji_",
          "icon":"",
          "items":([] as string[])
        }
        await Promise.all(info.value!.emote.map(async (e: { gif_url: string | URL | Request; url: string | URL | Request; text: string }) => {
          const resp = await fetch((() => {
            if (e.gif_url !== "") {
              return e.gif_url
            } else {
              return e.url
            }
          })(), {
            method: "GET",
            referrerPolicy: "no-referrer"
          })
          const blob = await resp.blob()
          const type = resp.headers.get("content-type")!.split("/")[1]
          zip.file(json.prefix + e.text + "." + type, blob)
          json.type = type ?? "png"
          json.items.push(e.text)
        }))
        json.icon = json.items[0] ?? ""
        json.name = info.value?.text ?? ""
        const blob4 = new Blob([JSON.stringify(json)],{type : "applocation/json"})
        zip.file("info.json",blob4)
        const content3 = await zip.generateAsync({ type: "blob" })
        saveAs(content3, info.value?.text)
        break
    }

  }
  const hasgif = ref(false)
  const viewgif = ref(false)
  function selgif() {
    viewgif.value = !viewgif.value
  }
  const todownload = ref(false)
  function open(url: string) {
    const link = document.createElement('a');
    link.href = url;
    link.target = '_blank';
    link.rel = 'noreferrer'; // 关键：告诉浏览器不要发送 Referer
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
</script>
<template>
<div class="info" v-if="info?.emote">
  <h1>{{ info.text }}</h1>
  <p>共{{ info.emote.length }}张<span> - </span><span @click="todownload = true" style="cursor: pointer;">下载</span></p>
  <p @click="selgif()" v-if="hasgif && !viewgif" style="cursor: pointer;">存在Gif，点击切换</p>
  <p @click="selgif()" v-if="hasgif && viewgif" style="cursor: pointer;">返回至静态图片</p>
  <p v-if="info.emote[0].type !== 4">点击图片可查看详情 </p>
  <p v-if="info.emote[0].type === 4">有颜文字，故无法提供更多功能</p>
  <div class="meta">
    <div class="box" v-for="(v, i) in info.emote" :key="v.text">
      <img v-if="v.type !== 4 && ((viewgif && v.gif_url == '') || !viewgif || v.type == 4)" @click="Info(i)"
        referrerpolicy="no-referrer" :src="v.url" :alt="v.meta.alias">
      <img v-if="v.type !== 4 && ((viewgif && v.gif_url !== '') || v.type == 4)" @click="Info(i)"
        referrerpolicy="no-referrer" :src="v.gif_url" :alt="v.meta.alias">
      <span v-if="v.type === 4">{{ v.text }}</span>
    </div>
  </div>
  <div class="ei" v-if="view.display" @click="click($event, 1)">
    <img referrerpolicy="no-referrer" :src="view.url" alt="">
    <div class="ev">

      <span>{{ view.name }}</span>
      <span>{{ view.ot }}</span>
      <p>Url链接</p>
      <div class="copy" @click="copy(view.url)"><img src="/img/copy.svg" alt=""><span :title="view.url">{{ view.url
      }}</span></div>
      <p>MarkDown</p>
      <div class="copy" @click="copy(view.name + '(' + view.url + ')')"><img src="/img/copy.svg" alt=""><span
          :title="view.name + '(' + view.url + ')'">{{ view.name + "(" + view.url + ")" }}</span></div>
      <p @click="open(view.url)" style="cursor: pointer;">新标签页打开</p>
    </div>
    <span class="cancel" @click="view.display = false">关闭</span>
  </div>
  <div class="download" v-if="todownload" @click="click($event, 2)">
    <h1>下载...</h1>
    <div class="ev ed">
      <span @click="download(1)" v-if="info.emote[0].type !== 4">下载压缩包</span>
      <span @click="download(2)" v-if="hasgif">下载压缩包(动图)</span>
      <span @click="download(3)">下载json格式</span>
      <span @click="download(4)">保存为Waline表情包</span>
    </div>
    <span @click="todownload = false" class="cancel">关闭</span>
  </div>
</div>
</template>

<style scoped>
.copy {
  margin: 0.4rem 0;
  display: flex;
  overflow: hidden;
  width: fit-content;
  justify-self: start;
  align-self: start;
  height: 2rem;
  justify-content: start;
  align-items: center;
  width: 25rem;

  img {
    cursor: pointer;
    height: 1.5rem !important;
    width: 1.5rem !important;
  }

  span {
    background-color: #adadad;
    border-radius: 2px;
    box-shadow: #000 0 0 2px;
    text-wrap-mode: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
  }
}

.ev {
  display: flex;
  flex-direction: column;
  width: 25rem;
  height: 15rem;
  align-items: center;
  border-radius: 12px;
  background-color: #ffffff;
  overflow: hidden;
  padding: 0 1rem;

  p {

    align-self: start;
    text-align: start;
  }

  img {
    width: 10rem;
  }
}

.ei,
.download {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  position: fixed;
  background-color: #000000e8;
  margin: 0;
  width: 100dvw;
  height: 100dvh;

  img {
    width: 13rem;
  }
}

.cancel {
  cursor: pointer;
  font-size: 1.5rem;
  font-weight: 600;
  color: #ff5252;
}

.download {
  h1 {
    color: #fff;
  }

  .ev {
    background-color: transparent;
  }


}

.ed {
  gap: 2rem;
  padding: 1rem;

  width: fit-content;
  height: fit-content;
  justify-content: center;
  align-items: start;
  text-align: start;

  span {
    cursor: pointer;
    border-radius: 12px;
    padding: 1rem;
    background-color: #000;
    color: #fff;
    font-size: 1.2rem;
  }
}

.box {
  justify-content: center;
  display: flex;
  width: 8rem;
  cursor: pointer;

  img {
    width: 100%;
  }
}

.meta {
  justify-content: center;
  overflow-x: hidden;
  width: 100%;
  display: flex;
  flex-wrap: wrap;
}

.info {
  justify-content: center;
  align-items: center;
  flex-direction: column;
  display: flex;
  width: 100dvw;
  height: 100dvh;

  h1 {
    margin: 0;
  }

  p {
    margin: 0;
  }
}

@media (width < 450px) {
  .ev {
    width: 20rem;
  }

  .copy {
    width: 20rem;
  }
}
</style>
