<template>
  <div class="card flex justify-content-center">
    <Button label="Copy Link" @click="copylink" />
  </div>
</template>

<script setup lang="ts">
function copylink() {
  let url = window.location.href
  if (url.indexOf("customRage=") > 0) {
    navigator.clipboard.writeText(url);
  } else {
    // const sql = `BETWEEN toDate('${ago.getFullYear()}-${ago.getMonth() + 1}-${ago.getDate()}') AND toDate('${now.getFullYear()}-${now.getMonth() + 1}-${now.getDate()}')`
    const now  = new Date()
    const ago = getDateAgoByDuration("3M")
    const f = `${ago.getFullYear()}-${ago.getMonth() + 1}-${ago.getDate()}:${now.getFullYear()}-${now.getMonth() + 1}-${now.getDate()}`
    url = url
      .replace(new RegExp("&timeRange=.+&"), "")
      .replace(new RegExp("&customRange=.+&"), "")
      .replace(new RegExp("timeRange=.+&"), "")
      .replace(new RegExp("customRange=.+&"), "")
      .replace(new RegExp("timeRange=.+$"), "")
      .replace(new RegExp("customRange=.+$"), "")
    navigator.clipboard.writeText(url + "&timeRange=custom&customRange=" + f);
  }

}
const duration = /(-?\d*\.?\d+(?:e[+-]?\d+)?)\s*([a-zμ]*)/gi

function getDateAgoByDuration(s: string): Date {
  s = s.replaceAll(/(\d),(\d)/g, "$1$2")
  console.log(s)
  let days = 0
  s.replaceAll(duration, (_, ...args: string[]) => {
    console.log(args)
    const [n, unit] = args
    const count = Number.parseInt(n, 10)
    switch (unit) {
      case "d":
        days = count
        break
      case "w":
        days = count * 7
        break
      case "m":
      case "M":
      case "month":
        days = count * 31
        break
      case "y":
        days = count*365
        break
      default:
        console.error("Fail to get ")
        days = 0
    }
    return ""
  })
  const date = new Date()
  date.setDate(date.getDate() - days)
  return date
}
</script>
<style scoped></style>
