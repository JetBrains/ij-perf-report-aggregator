
 * [Vite](https://vitejs.dev) is not used because it doesn't build and verify sources in a correct way. Speed or correct compilation? 
   Of course, speed is not so important if compilation errors are not checked — no need to run not compilable code to catch errors in runtime.
    * No official plugin for [typescript](https://github.com/vitejs/vite/issues/245).
    * No official plugin for [eslint](https://github.com/vitejs/vite/issues/818). 
   
    Still, 
 * [Vue-ECharts](https://github.com/ecomfe/vue-echarts) is not used in addition to ECharts even if auto resize is not implemented (`width: 100%` is supported, but subsequent window resize will not lead to chart resizing). To reduce dependency set, better to use [small solution](https://stackoverflow.com/a/27801087) instead of depending on yet another library.  