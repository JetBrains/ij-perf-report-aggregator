 * [Vite](https://vitejs.dev) is used instead of Vue CLI because it is more modern tool. The only concern — speed or correct compilation?
    * No official plugin for [typescript](https://github.com/vitejs/vite/issues/245).
    * No official plugin for [eslint](https://github.com/vitejs/vite/issues/818). 
 * [ECharts](https://echarts.apache.org/en/index.html) is used instead of [amCharts](https://www.amcharts.com) because it is more modern charting library. 
    * Declarative approach allows to describe desired data visualization using less code and in a more simple way.
    * Performance. ECharts is able to render thousands points on line chart. No need to use granularity to aggregate data on server side.
 * [Vue-ECharts](https://github.com/ecomfe/vue-echarts) is not used in addition to ECharts even if auto resize is not implemented (`width: 100%` is supported, but subsequent window resize will not lead to chart resizing). To reduce dependency set, better to use [small solution](https://stackoverflow.com/a/27801087) instead of depending on yet another library.
  * [pnpm](https://pnpm.js.org) is used instead of Yarn, because Yarn 2 PnP is not supported yet by IntelliJ IDEA. And in general pnpm does the job perfectly. NPM is out of comparison since at least it is slower and no hope that is finally reliable.
    * Strict — "pnpm creates a non-flat node_modules, so code has no access to arbitrary packages". Yarn 2 PnP also supports it, but PnP cannot be used yet.